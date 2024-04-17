package util

// 对指针进行安全解引用
func Deref[T any](ptr *T) T {
	if ptr != nil {
		return *ptr
	}
	return *new(T)
}
