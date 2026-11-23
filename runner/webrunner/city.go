package webrunner

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gosom/google-maps-scraper/grid"
)

const (
	cityGeocoderURL = "https://nominatim.openstreetmap.org/search"
	maxCityPoints   = 20
	maxCitySpanKm   = 300
)

type geocodingResult struct {
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	BoundingBox []string `json:"boundingbox"`
}

// resolveCityPoints geocodes a city and creates a bounded grid of search points.
// Nominatim's bounding box is deliberately used as a conservative rectangular
// approximation; it is not an administrative-boundary polygon.
func resolveCityPoints(ctx context.Context, city, coverage string) (string, []string, error) {
	city = strings.TrimSpace(city)
	if city == "" {
		return "", nil, fmt.Errorf("city is required")
	}

	target, err := coveragePointTarget(coverage)
	if err != nil {
		return "", nil, err
	}

	u, err := url.Parse(cityGeocoderURL)
	if err != nil {
		return "", nil, err
	}
	q := u.Query()
	q.Set("q", city)
	q.Set("format", "jsonv2")
	q.Set("limit", "1")
	q.Set("addressdetails", "1")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("User-Agent", "google-maps-scraper/1.0 (city coverage geocoding)")

	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("could not locate city %q: %w", city, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("city lookup failed with status %s", resp.Status)
	}

	var results []geocodingResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return "", nil, fmt.Errorf("decode city lookup: %w", err)
	}
	if len(results) == 0 {
		return "", nil, fmt.Errorf("city %q was not found; include state and country, e.g. Belo Horizonte, MG, Brazil", city)
	}

	result := results[0]
	lat, err := strconv.ParseFloat(result.Lat, 64)
	if err != nil {
		return "", nil, fmt.Errorf("invalid city latitude: %w", err)
	}
	lon, err := strconv.ParseFloat(result.Lon, 64)
	if err != nil {
		return "", nil, fmt.Errorf("invalid city longitude: %w", err)
	}
	if len(result.BoundingBox) != 4 {
		return "", nil, fmt.Errorf("city lookup did not return a usable bounding box")
	}

	// Nominatim order: south, north, west, east.
	south, err := strconv.ParseFloat(result.BoundingBox[0], 64)
	if err != nil {
		return "", nil, fmt.Errorf("invalid city bounding box: %w", err)
	}
	north, err := strconv.ParseFloat(result.BoundingBox[1], 64)
	if err != nil {
		return "", nil, fmt.Errorf("invalid city bounding box: %w", err)
	}
	west, err := strconv.ParseFloat(result.BoundingBox[2], 64)
	if err != nil {
		return "", nil, fmt.Errorf("invalid city bounding box: %w", err)
	}
	east, err := strconv.ParseFloat(result.BoundingBox[3], 64)
	if err != nil {
		return "", nil, fmt.Errorf("invalid city bounding box: %w", err)
	}

	bbox := grid.BoundingBox{MinLat: south, MinLon: west, MaxLat: north, MaxLon: east}
	if _, err := grid.ParseBoundingBox(fmt.Sprintf("%f,%f,%f,%f", south, west, north, east)); err != nil {
		return "", nil, fmt.Errorf("invalid city boundary: %w", err)
	}
	if err := validateCityBounds(bbox); err != nil {
		return "", nil, fmt.Errorf("could not use geocoding result for %q: %w", city, err)
	}

	// Increase cell size until the generated grid fits the selected point budget.
	cellSizeKm := 2.0
	for grid.EstimateCellCount(bbox, cellSizeKm) > target {
		cellSizeKm *= 1.25
		if cellSizeKm > 2000 {
			return "", nil, fmt.Errorf("city area is too large to cover safely")
		}
	}

	cells := grid.GenerateCells(bbox, cellSizeKm)
	if len(cells) == 0 || len(cells) > maxCityPoints {
		return "", nil, fmt.Errorf("could not generate a safe search grid for %q", city)
	}

	points := make([]string, 0, len(cells))
	for _, cell := range cells {
		points = append(points, cell.GeoCoordinates())
	}

	return fmt.Sprintf("%f,%f", lat, lon), points, nil
}

// validateCityBounds prevents a country or state returned by an ambiguous
// geocoding query from being treated as a city and sampled with only 5-20 points.
func validateCityBounds(bbox grid.BoundingBox) error {
	midLat := (bbox.MinLat + bbox.MaxLat) / 2
	latSpanKm := (bbox.MaxLat - bbox.MinLat) * 111.32
	lonSpanKm := (bbox.MaxLon - bbox.MinLon) * 111.32 * math.Abs(math.Cos(midLat*math.Pi/180))

	if latSpanKm > maxCitySpanKm || lonSpanKm > maxCitySpanKm {
		return fmt.Errorf(
			"the matched area is approximately %.0f km north-south by %.0f km east-west; this is too large for city coverage. Enter the municipality, state, and country (for example, Iturama, MG, Brazil)",
			latSpanKm,
			lonSpanKm,
		)
	}

	return nil
}

func coveragePointTarget(coverage string) (int, error) {
	switch strings.ToLower(strings.TrimSpace(coverage)) {
	case "", "balanced":
		return 10, nil
	case "economical":
		return 5, nil
	case "complete":
		return maxCityPoints, nil
	default:
		return 0, fmt.Errorf("invalid city coverage %q: use economical, balanced, or complete", coverage)
	}
}
