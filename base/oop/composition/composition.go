package composition

import "fmt"

type auther struct {
	firstName string
	lastName  string
	bio       string
}

func (a *auther) fullName() string {
	return fmt.Sprintf("%s %s", a.firstName, a.lastName)
}

type post struct {
	title   string
	content string
	*auther
}

func (p *post) desc() {
	fmt.Println("Title: ", p.title)
	fmt.Println("Content: ", p.content)
	fmt.Println("Author: ", p.fullName())
	fmt.Println("Bio: ", p.bio)
}

type posts struct {
	ps []post
}

func (ps *posts) showPs() {
	for _, v := range ps.ps {
		v.desc()
		fmt.Println()
	}
}

func show() {
	auther1 := &auther{
		"Naveen",
		"Ramanathan",
		"Golang Enthusiast",
	}
	post1 := post{
		"Inheritance in Go",
		"Go supports composition instead of inheritance",
		auther1,
	}
	post1.desc()

	fmt.Println("===============")

	post2 := post{
		"Inheritance in Go",
		"Go supports composition instead of inheritance",
		auther1,
	}
	post3 := post{
		"Inheritance in Go",
		"Go supports composition instead of inheritance",
		auther1,
	}
	ps := posts{[]post{post1, post2, post3}}
	ps.showPs()
}
