package utils

func GetPointValue[T any](value T) *T {
	return &value
}
