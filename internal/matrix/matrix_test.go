package matrix_test

import (
	"berlin-heatmap/internal/matrix"
	"testing"
)

func TestUnreachableConstant(t *testing.T) {
	if matrix.Unreachable != 65535 {
		t.Errorf("expected Unreachable=65535, got %d", matrix.Unreachable)
	}
}

func TestMatrixInit(t *testing.T) {
	n := 3
	data := make([][]uint16, n)
	for i := range data {
		data[i] = make([]uint16, n)
		for j := range data[i] {
			if i == j {
				data[i][j] = 0
			} else {
				data[i][j] = matrix.Unreachable
			}
		}
	}

	if data[0][1] != matrix.Unreachable {
		t.Errorf("expected off-diagonal to be Unreachable, got %d", data[0][1])
	}
}
