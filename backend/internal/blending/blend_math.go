package blend

import (
	"maps"
	"math"
	"slices"
)

func FindIntersectKeys(m1, m2 map[string]int) []string {
	//Smaller map on outside loop leads to better stack order

	if len(m1) > len(m2) {
		m1, m2 = m2, m1
	}
	result := make([]string, 0)

	for k := range m1 {
		if _, ok := m2[k]; ok {
			result = append(result, k)
		}
	}

	return result
}

// Calculate Log Weighted Cosine Similarity
func CalculateLWCS(lambda float64, user1Map, user2Map map[string]int) int {
	common_keys := FindIntersectKeys(user1Map, user2Map)

	dot_product := 0.0
	for _, v := range common_keys {
		v1 := user1Map[v]
		v2 := user2Map[v]
		dot_product += float64(v1 * v2)
	}

	magA := math.Sqrt(getMagnitude(slices.Collect(maps.Values(user1Map))))
	magB := math.Sqrt(getMagnitude(slices.Collect(maps.Values(user2Map))))
	if magA == 0 || magB == 0 {
		return 0
	}

	num_factor := math.Log10(float64(dot_product))
	denom_factor := math.Log10(float64(magA * magB))

	logWeightedValue := (num_factor / denom_factor)
	if logWeightedValue < 0 {
		logWeightedValue = 0
	}

	directCosineValue := dot_product / (magA * magB)
	finalValue := lambda*logWeightedValue + (1-lambda)*directCosineValue

	return int(finalValue * 100) //0.X float -> XX int for percentage value
}

func getMagnitude(arr []int) float64 {
	sum := 0.0
	for _, v := range arr {
		sum += float64(v * v)
	}
	return sum
}

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}
