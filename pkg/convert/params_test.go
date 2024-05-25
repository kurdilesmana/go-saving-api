package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortToMap(t *testing.T) {
	tests := map[string]struct {
		input string
		want  map[string]string
	}{
		"data 1": {
			input: "nama,-posisi",
			want: map[string]string{
				"nama":   "ASC",
				"posisi": "DESC",
			},
		},
		"data 2": {
			input: "-nama,-posisi",
			want: map[string]string{
				"nama":   "DESC",
				"posisi": "DESC",
			},
		},
		"data 3": {
			input: "nama,posisi",
			want: map[string]string{
				"nama":   "ASC",
				"posisi": "ASC",
			},
		},
		"empty": {
			input: "",
			want:  map[string]string{},
		},
		"one data": {
			input: "nama",
			want: map[string]string{
				"nama": "ASC",
			},
		},
		"one data desc": {
			input: "-nama",
			want: map[string]string{
				"nama": "DESC",
			},
		},
		"minus in the middle": {
			input: "nama,created-at",
			want: map[string]string{
				"nama":       "ASC",
				"created-at": "ASC",
			},
		},
		"space in the middle": {
			input: "nama,-created at",
			want: map[string]string{
				"nama":       "ASC",
				"created at": "DESC",
			},
		},
		"space after comma (trimming)": {
			input: "nama,  -created_at",
			want: map[string]string{
				"nama":       "ASC",
				"created_at": "DESC",
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := SortToMap(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}
