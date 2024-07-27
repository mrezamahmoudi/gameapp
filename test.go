package main

import "fmt"

type person struct {
	name string
}

func rename(p *person) {
	p.name = "test"
}

func main() {

	p := person{name: "Ali"}

	rename(&p)

	fmt.Println(p.name)
}
