package main

import (
	"encoding/json"
	"fmt"
)

type Cat struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (c Cat) print() {
	fmt.Println(c)
}

func (c Cat) String() string {
	return fmt.Sprintf("Name==>%v Age===>%v", c.Name, c.Age)
}

func main() {
	// test01()

	var cat Cat = Cat{}
	// cat.print()
	fmt.Println(&cat)
}

func test01() {
	var cat Cat = Cat{"a",
		1}

	fmt.Println(cat)

	cat1 := Cat{"b",
		2}

	fmt.Println(cat1)

	var cat2 *Cat = new(Cat)
	cat2.Age = 10

	fmt.Println(*cat2)

	var cat3 *Cat = &Cat{}
	cat3.Age = 20
	cat3.Name = "20"

	fmt.Println(*cat3)

	var cat4 Cat /* = Cat{"a", 12} */
	fmt.Println(cat4)

	jsonStr, _ := json.Marshal(cat4)
	fmt.Println(string(jsonStr))
}
