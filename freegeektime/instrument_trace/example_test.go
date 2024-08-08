package instrument_trace

func a() {
	defer Trance()()
	b()
}

func b() {
	defer Trance()()
	c()
}
func c() {
	defer Trance()()
	d()
}

func d() {
	defer Trance()()
}

func ExampleTrance() {
	a()
}
