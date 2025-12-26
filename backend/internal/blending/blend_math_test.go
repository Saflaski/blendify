package blend

import (
	"math"
	"slices"
	"testing"
)

func TestLWCS(t *testing.T) {

	userA := map[string]int{
		"A": 123,
		"B": 234,
		"C": 0,
		"D": 2900,
	}

	userB := map[string]int{
		"A": 1444,
		"B": 12,
		"Z": 510,
		"Y": 3023,
	}
	t.Run("Get blend from A and B", func(t *testing.T) {
		got := FindIntersectKeys(userA, userB)
		want := []string{"A", "B"}
		if !slices.Equal(got, want) {
			t.Errorf("Intersection not equal. Got %s, want %s", got, want)
		}
	})
	t.Run("Magnitude Test", func(t *testing.T) {
		A := []int{3, 5, 10}
		got := roundTo(math.Sqrt(getMagnitude(A)), 3)

		//sqrt(3^2 + 5^2 + 10^2) = sqrt(134) = 11.5758369028
		want := roundTo(11.5758369028, 3)
		if got != want {
			t.Errorf("Incorrect magnitude. Got %f , want %f", got, want)
		}

	})
	t.Run("Get blend from A and B", func(t *testing.T) {
		blendNum := CalculateLWCS(0.8, userA, userB)
		if blendNum > 100 || blendNum <= 0 {
			t.Errorf("Number is not within acceptable range: %d", blendNum)
		}
	})

	t.Run("Get overall blend", func(t *testing.T) {
		artist, err := combineNumbersWithWeights(70, 55, 63, 5, 4, 5)
		if err != nil {
			t.Log(artist)
			t.Errorf(": %s", err)
		}
		album, _ := combineNumbersWithWeights(70, 55, 63, 5, 6, 4)
		track, _ := combineNumbersWithWeights(70, 55, 63, 5, 4, 5)
		blendNum, err := combineNumbersWithWeights(artist, album, track, 10, 10, 10)
		if err != nil {
			t.Log(blendNum)
			t.Errorf(": %s", err)
		}
	})
}

func roundTo(n float64, decimals uint) float64 {
	factor := math.Pow(10, float64(decimals))
	return math.Round(n*factor) / factor
}
