package towSum

//https://leetcode.cn/problems/two-sum/
func towSum(arrays []int, tag int) []int {
	cache := map[int]int{}
	for index, val := range arrays {
		if val, exist := cache[tag-val]; exist {
			return []int{val, index}
		}
		cache[val] = index
	}
	return nil
}
