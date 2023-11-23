package util_test

import (
	"testing"

	"github.com/acgtools/sub-renamer/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestSliceToSet(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		arr  []string
		want map[string]struct{}
	}{
		{
			name: "Empty slice",
			arr:  nil,
			want: map[string]struct{}{},
		},
		{
			name: "Case 01",
			arr:  []string{"A", "B", "C"},
			want: map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := util.SliceToSet(tc.arr)

			assert.Equal(t, tc.want, m)
		})
	}
}

func TestSliceFilter(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		arr  []int
		keep func(int) bool
		want []int
	}{
		{
			name: "Int slice",
			arr:  []int{1, 2, 3, 4, 5},
			keep: func(i int) bool {
				return i > 3
			},
			want: []int{4, 5},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := util.SliceFilter(tc.arr, tc.keep)

			assert.Equal(t, tc.want, res)
		})
	}
}
