package util

// 对指针进行安全解引用
func Deref[T any](ptr *T) T {
	if ptr != nil {
		return *ptr
	}
	return *new(T)
}

// 检查一个元素是否存在于一个切片里
func Contain[T comparable](slice []T, val T) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// 通过循环让一个函数先后传入多个参数并执行
func Loop[T any](callback func(T), argu ...T) {
	for _, value := range argu {
		callback(value)
	}
}
