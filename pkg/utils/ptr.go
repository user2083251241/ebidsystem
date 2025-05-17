package utils

// Ptr 返回任意类型的指针：
func Ptr[T any](v T) *T {
	return &v
}
