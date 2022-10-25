package long_characters

import "testing"

func TestLengthOfLongestSubstringBak(t *testing.T) {
	s := "asdasdasd"
	res := lengthOfLongestSubstringBak(s)
	res = lengthOfLongestSubstring2(s)
	t.Logf("res:%d", res)
}
