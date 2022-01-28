package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/getevo/evo/lib/gpath"
	"go/parser"
	"go/token"
	"log"
)


func main()  {
	/*var file,_ = os.Open("./test.go")*/
	src,_ := gpath.ReadFile("./test.go")
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "", src, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	for _,item := range f.Comments{
		spew.Dump(item.Text())
	}
	fmt.Println(  )
}
