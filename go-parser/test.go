package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

// comment goes here


// @comment -> goes here!

type MyStruct struct {
	Field1 string
}

func (MyStruct)Func()  {

}


// another
func Func() string {
	fmt.Println("test")
	spew.Dump("test")
	return "hi"
}