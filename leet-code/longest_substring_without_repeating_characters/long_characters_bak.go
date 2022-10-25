package long_characters

func lengthOfLongestSubstringBak(s string) int {
	len := len(s)
	if len == 0 {
		return 0
	}
	//求最大数
	maxFun := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	left, rigth, res := 0, 0, 0
	//存储字符的code
	var cache [256]int
	for left < len {
		//cache中不存在
		if rigth < len && cache[s[rigth]] == 0 {
			//改为1
			cache[s[rigth]]++
			rigth++
		} else {
			//已存在该字符串剔除掉
			cache[s[left]]--
			left++
		}
		res = maxFun(res, rigth-left)

	}
	return res
}
