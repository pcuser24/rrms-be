package utils

func Ternary[T any](condition bool, a, b T) T {
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
