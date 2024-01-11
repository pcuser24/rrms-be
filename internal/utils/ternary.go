package utils

func Ternary(condition bool, a, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func PtrDerefence[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}
