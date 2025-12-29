package blend

import (
	"fmt"
	"maps"
	"math"
	"slices"
)

func GetLogCombinedScore(a, b int) float64 {
	left, right := float64(a), float64(b)
	logSum := math.Log(left) + math.Log(right)
	var diffFactor float64
	if left > right {
		diffFactor = math.Log(right) / math.Log(left)
	} else {
		diffFactor = math.Log(left) / math.Log(right)
	}

	return logSum * diffFactor
}

func FindIntersectKeys[T any](m1, m2 map[string]T) []string {
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

func GetTopDuoArtists(lambda float64, user1Map, user2Map map[string]int) {

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

// func calcOverModality(typeBlend TypeBlend, durationWeights ...int) int {
// 	return (typeBlend.OneMonth*durationWeights[0] + typeBlend.ThreeMonth*durationWeights[1] + typeBlend.OneYear*durationWeights[2]) / len(durationWeights)
// }

// Takes in equal number of input numbers and weights as slice of ints
// Weights need to be between 0 and 10
func combineNumbersWithWeights(inputsAndWeights ...int) (int, error) {
	if len(inputsAndWeights)%2 != 0 {
		return 0, fmt.Errorf(" need equal number of inputs and weights")
	}

	numInputs := len(inputsAndWeights) / 2
	runningSum := 0.0
	// fmt.Println("Num inputs ", numInputs)
	// fmt.Println("---------------------")
	max := inputsAndWeights[numInputs]
	for _, v := range inputsAndWeights[numInputs:] {
		if v > max {
			max = v
		}
	}
	for i := 0; i <= numInputs-1; i++ {
		weight := inputsAndWeights[i+numInputs]

		if weight < 0 || weight > 10 {
			return 0, fmt.Errorf(" abnormal weight values, need to be between 0 and 10:%d", weight)
		}
		weightFloat := (float64(weight) / float64(max))
		runningSum += float64(inputsAndWeights[i]) * weightFloat
		// fmt.Println("i:", i)
		// fmt.Println("input:", inputsAndWeights[i])
		// fmt.Println("weight:", weightFloat)
		// fmt.Println("runningSum:", runningSum)
		// fmt.Println("---------------------")
		// runningSum = runningSum
	}
	finalSum := int(runningSum / float64(numInputs))
	if finalSum > 100.0 || finalSum < 0.0 {
		return 0, fmt.Errorf(" abnormal sum of modalities:%d", finalSum)
	}
	// fmt.Println("Final sum: ", finalSum)
	return finalSum, nil
}
