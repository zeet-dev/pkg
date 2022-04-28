package utils

func MergeMap(m1, m2 map[string]string) map[string]string {
	m3 := make(map[string]string)
	for k, v := range m1 {
		m3[k] = v
	}
	for k, v := range m2 {
		m3[k] = v
	}
	return m3
}
