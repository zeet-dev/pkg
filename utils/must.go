package utils

// please only use this for testing
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// please only use this for testing
func MustOk[T any](v T, ok bool) T {
	if !ok {
		panic("mustok failed")
	}
	return v
}
