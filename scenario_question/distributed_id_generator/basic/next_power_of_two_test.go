package basic

import "testing"

func TestNextPowerOfTwo(t *testing.T) {
	tests := []struct {
		n    uint32
		want uint32
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 4},
		{4, 4},
		{5, 8},
		{6, 8},
		{7, 8},
		{8, 8},
		{9, 16},
		{10, 16},
		{11, 16},
		{12, 16},
		{13, 16},
		{14, 16},
		{15, 16},
		{16, 16},
		{17, 32},
		{18, 32},
		{19, 32},
		{20, 32},
		{21, 32},
		{22, 32},
		{23, 32},
		{24, 32},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := NextPowerOfTwo(tt.n); got != tt.want {
				t.Errorf("NextPowerOfTwo() = %v, want %v", got, tt.want)
			}
		})
	}
}
