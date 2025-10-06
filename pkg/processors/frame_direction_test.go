package processors

import (
	"testing"
)

func TestFrameDirectionString(t *testing.T) {
	tests := []struct {
		direction    FrameDirection
		expectedString string
	}{
		{FrameDirectionUpstream, "Upstream"},
		{FrameDirectionDownstream, "Downstream"},
	}

	for _, test := range tests {
		if test.direction.String() != test.expectedString {
			t.Errorf("Expected %s, got %s", test.expectedString, test.direction.String())
		}
	}
}