package webrunner

import (
	"testing"

	"github.com/gosom/google-maps-scraper/grid"
)

func TestCoveragePointTarget(t *testing.T) {
	tests := []struct {
		coverage string
		want     int
		wantErr  bool
	}{
		{coverage: "economical", want: 5},
		{coverage: "balanced", want: 10},
		{coverage: "", want: 10},
		{coverage: "complete", want: 20},
		{coverage: "unknown", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.coverage, func(t *testing.T) {
			got, err := coveragePointTarget(tt.coverage)
			if (err != nil) != tt.wantErr {
				t.Fatalf("coveragePointTarget(%q) error = %v, wantErr %v", tt.coverage, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("coveragePointTarget(%q) = %d, want %d", tt.coverage, got, tt.want)
			}
		})
	}
}

func TestValidateCityBounds(t *testing.T) {
	tests := []struct {
		name    string
		bbox    grid.BoundingBox
		wantErr bool
	}{
		{
			name: "city sized area",
			bbox: grid.BoundingBox{
				MinLat: -20.0, MinLon: -44.1,
				MaxLat: -19.7, MaxLon: -43.8,
			},
		},
		{
			name: "country sized area",
			bbox: grid.BoundingBox{
				MinLat: -33.75, MinLon: -73.99,
				MaxLat: 5.27, MaxLon: -34.79,
			},
			wantErr: true,
		},
		{
			name: "oversized north south extent",
			bbox: grid.BoundingBox{
				MinLat: -2.0, MinLon: -44.0,
				MaxLat: 2.0, MaxLon: -43.9,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCityBounds(tt.bbox)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateCityBounds() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
