package student

type student struct {
	Age  int
	Name string
}

func New(age int, name string) student {
	return student{Age: age, Name: name}
}
