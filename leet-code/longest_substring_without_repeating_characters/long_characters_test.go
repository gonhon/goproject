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
