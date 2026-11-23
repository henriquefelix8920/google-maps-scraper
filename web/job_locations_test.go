package web

import "testing"

func TestValidateCoordinatePair(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid coordinates", value: "-19.9167,-43.9345"},
		{name: "spaces are trimmed", value: " -19.9167, -43.9345 "},
		{name: "missing longitude", value: "-19.9167", wantErr: true},
		{name: "invalid latitude", value: "north,-43.9345", wantErr: true},
		{name: "latitude out of range", value: "91,-43.9345", wantErr: true},
		{name: "longitude out of range", value: "-19.9167,-181", wantErr: true},
		{name: "not a coordinate pair", value: "curitiba", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCoordinatePair(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateCoordinatePair(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestJobDataValidateAdditionalLocations(t *testing.T) {
	data := JobData{
		Keywords: []string{"dentists in Belo Horizonte"},
		Lang:     "pt",
		Depth:    10,
		MaxTime:  10,
		FastMode: true,
		Lat:      "-19.9167",
		Lon:      "-43.9345",
		Locations: []string{
			"-19.9000,-43.9500",
			"-19.9400,-43.9200",
		},
	}

	if err := data.Validate(); err != nil {
		t.Fatalf("valid job with multiple locations was rejected: %v", err)
	}

	data.Locations = []string{"-19.9000,west"}
	if err := data.Validate(); err == nil {
		t.Fatal("invalid additional location was accepted")
	}
}
