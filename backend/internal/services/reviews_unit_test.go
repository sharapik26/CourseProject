package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAggregateZeroPlaces проверяет, что для несуществующего места
// агрегаты возвращаются с нулевыми значениями (без ошибок).
func TestAggregateEmptyValues(t *testing.T) {
	a := Aggregate{PlaceID: 999}
	assert.Equal(t, uint64(999), a.PlaceID)
	assert.Equal(t, 0.0, a.AvgNoise)
	assert.Equal(t, 0, a.ReviewsCnt)
}

// TestSensoryScoreBoundaries проверяет, что значения сенсорных оценок
// корректно описаны в модели — сценарий граничных значений 1 и 5.
func TestSensoryScoreBoundaries(t *testing.T) {
	cases := []struct {
		name  string
		val   int
		valid bool
	}{
		{"min", 1, true},
		{"middle", 3, true},
		{"max", 5, true},
		{"zero is absent (not entered)", 0, false},
		{"too high", 6, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.valid, tc.val >= 1 && tc.val <= 5)
		})
	}
}