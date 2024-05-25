package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundFloat64(t *testing.T) {
	tests := map[string]struct {
		input      float64
		precesiont uint
		want       float64
	}{
		"data 1": {
			input:      3.14,
			precesiont: 1,
			want:       3.1,
		},
		"data 2": {
			input:      1.6666666,
			precesiont: 2,
			want:       1.67,
		},
		"data 3": {
			input:      1.9090909,
			precesiont: 5,
			want:       1.90909,
		},
		"pembulatan diatas 0.5": {
			input:      1.06,
			precesiont: 1,
			want:       1.1,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := RoundFloat64(tc.input, tc.precesiont)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAlmostEqual(t *testing.T) {
	tests := map[string]struct {
		input1 float64
		input2 float64
		want   bool
	}{
		"data 1": {
			input1: 3.13,
			input2: 3.15,
			want:   true,
		},
		"data 2": {
			input1: 3.11112,
			input2: 3.20000,
			want:   true,
		},
		"data 3": {
			input1: 3.1010,
			input2: 3.3000,
			want:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := AlmostEqual(tc.input1, tc.input2)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStrToInt(t *testing.T) {
	tests := map[string]struct {
		input      string
		defaultVal int
		want       int
	}{
		"data 1": {
			input:      "200",
			defaultVal: 200,
			want:       200,
		},
		"data 2": {
			input:      "xxx500",
			defaultVal: 500,
			want:       500,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := StrToInt(tc.input, tc.defaultVal)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStrToInt64(t *testing.T) {
	tests := map[string]struct {
		input      string
		defaultVal int64
		want       int64
	}{
		"data 1": {
			input:      "200",
			defaultVal: 200,
			want:       200,
		},
		"data 2": {
			input:      "xxx500",
			defaultVal: 500,
			want:       500,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := StrToInt64(tc.input, tc.defaultVal)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStrToUint64(t *testing.T) {
	tests := map[string]struct {
		input      string
		defaultVal uint64
		want       uint64
	}{
		"data 1": {
			input:      "200",
			defaultVal: 200,
			want:       200,
		},
		"data 2": {
			input:      "xxx500",
			defaultVal: 500,
			want:       500,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := StrToUint64(tc.input, tc.defaultVal)
			assert.Equal(t, tc.want, got)
		})
	}
}
