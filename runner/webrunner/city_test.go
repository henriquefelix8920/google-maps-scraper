package webrunner

import "testing"

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
