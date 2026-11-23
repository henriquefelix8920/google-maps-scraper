package web

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var jobs []Job

const (
	StatusPending = "pending"
	StatusWorking = "working"
	StatusOK      = "ok"
	StatusFailed  = "failed"
)

type SelectParams struct {
	Status string
	Limit  int
	Offset int
}

// JobPage contains one page of jobs and the metadata needed to navigate it.
type JobPage struct {
	Jobs        []Job
	CurrentPage int
	TotalPages  int
	Total       int
	HasPrev     bool
	HasNext     bool
	PrevPage    int
	NextPage    int
	HasPages    bool
}

type JobRepository interface {
	Get(context.Context, string) (Job, error)
	Create(context.Context, *Job) error
	Delete(context.Context, string) error
	Select(context.Context, SelectParams) ([]Job, error)
	Update(context.Context, *Job) error
}

type Job struct {
	ID     string
	Name   string
	Date   time.Time
	Status string
	Data   JobData
}

func (j *Job) Validate() error {
	if j.ID == "" {
		return errors.New("missing id")
	}

	if j.Name == "" {
		return errors.New("missing name")
	}

	if j.Status == "" {
		return errors.New("missing status")
	}

	if j.Date.IsZero() {
		return errors.New("missing date")
	}

	if err := j.Data.Validate(); err != nil {
		return err
	}

	return nil
}

type JobData struct {
	Keywords     []string      `json:"keywords"`
	Lang         string        `json:"lang"`
	Zoom         int           `json:"zoom"`
	Lat          string        `json:"lat"`
	Lon          string        `json:"lon"`
	FastMode     bool          `json:"fast_mode"`
	Radius       int           `json:"radius"`
	Depth        int           `json:"depth"`
	Email        bool          `json:"email"`
	ExtraReviews bool          `json:"extra_reviews"`
	MaxTime      time.Duration `json:"max_time"`
	Proxies      []string      `json:"proxies"`
	Locations    []string      `json:"locations,omitempty"`
}

func (d *JobData) Validate() error {
	if len(d.Keywords) == 0 {
		return errors.New("missing keywords")
	}

	if d.Lang == "" {
		return errors.New("missing lang")
	}

	if len(d.Lang) != 2 {
		return errors.New("invalid lang")
	}

	if d.Depth == 0 {
		return errors.New("missing depth")
	}

	if d.MaxTime == 0 {
		return errors.New("missing max time")
	}

	if d.FastMode {
		if d.Lat == "" || d.Lon == "" {
			return errors.New("missing geo coordinates")
		}
		if err := validateCoordinatePair(d.Lat + "," + d.Lon); err != nil {
			return fmt.Errorf("invalid primary geo coordinates: %w", err)
		}
	}

	if len(d.Locations) > 20 {
		return errors.New("a maximum of 20 additional search points is allowed")
	}
	for i, location := range d.Locations {
		if err := validateCoordinatePair(location); err != nil {
			return fmt.Errorf("invalid additional search point %d: %w", i+1, err)
		}
	}

	return nil
}

func validateCoordinatePair(value string) error {
	parts := strings.Split(strings.TrimSpace(value), ",")
	if len(parts) != 2 {
		return errors.New("use latitude,longitude")
	}

	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return fmt.Errorf("invalid latitude: %w", err)
	}
	lon, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return fmt.Errorf("invalid longitude: %w", err)
	}
	if !(lat >= -90 && lat <= 90) {
		return errors.New("latitude must be between -90 and 90")
	}
	if !(lon >= -180 && lon <= 180) {
		return errors.New("longitude must be between -180 and 180")
	}

	return nil
}
