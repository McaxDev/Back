package util

// 从数组里找到target的函数
func FindTarget(nums []int, target int) (int, int, bool) {

	// 第一遍遍历数组
	for i := 0; i < len(nums); i++ {

		// 第二遍遍历数组
		for j := i + 1; j < len(nums); j++ {

			// 检查是否满足条件并返回
			if nums[i]+nums[j] == target {

				// 如果找到了，就返回索引和true
				return i, j, true
			}
		}
	}

	// 如果找不到，就返回0和false
	return 0, 0, false
}

// 优化版
func FindTargetOptimized(nums []int, target int) (int, int, bool) {

	// 创建哈希表
	hash := make(map[int]int)

	// 执行遍历操作
	for i := 0; i < len(nums); i++ {

		// 检查哈希表的键是否存在
		if value, exists := hash[nums[i]]; exists {

			// 如果存在，就返回索引
			return i, value, true
		}

		// 计算差值
		left := target - nums[i]
		// 将差值存入哈希表
		hash[left] = i
	}
	// 如果不存在，就返回false
	return 0, 0, false
}
