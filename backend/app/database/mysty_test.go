package database

import "fmt"

type MyCoolStruct struct {
	test uint16
}

func (s *MyCoolStruct) Thing() {
	// So you can do stuff with your struct here...
	// Thing is a method of this struct ya...
	// s here refers to this struct...
	fmt.Println(s.test)
}

func Test() {
	cool := MyCoolStruct{1}
	cool.Thing();
}

// THat's a method which is good syntax to know...
// So we make the struct here and call it's method
// I just put it here so you remember lol