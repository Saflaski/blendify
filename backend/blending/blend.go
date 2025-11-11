package blend

// --- Functions for creating the blend number

func findIntersectKeys(m1, m2 map[string]int) ([]string) {
	//Smaller map on outside loop leads to better stack order

	if len(m1) > len(m2) {
		m1 , m2 = m2, m1
	}

	result := make([]string, 0)

	for k:= range m1 {
		if _, ok := m2[k]; ok {
			result = append(result, k)
		}
	}

	return result
}


func getBlend(userA, userB )
