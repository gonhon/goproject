package long_characters

import (
	"fmt"
	"testing"
)

func TestLongestSubstring(t *testing.T) {
	lengthOfLongestSubstring("abcabcbb z")

}

func Test1(t *testing.T) {
	a := lengthOfLongestSubstring2("abcabcbb")
	fmt.Println(a)
}

func Test2(t *testing.T) {
	a := "abcabcbb"
	for i := 0; i < len(a); i++ {
		fmt.Println(a[i])
	}

}
