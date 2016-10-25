package main

import (
	"fmt"
	"strings"
	"go/parser"
	"go/token"
	"os"
	"github.com/codegangsta/cli"
	"github.com/miyo/go2ir/synthesijer"
)

func parse(src string) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, src, nil, 0)
	if err != nil {
		panic(err)
	}

	target_name := src[:strings.LastIndex(src, ".")]
	target := synthesijer.Module{Name: target_name}

	synthesijer.Parse(file, &target)

	synthesijer.GenerateModule(&target)

}

func main() {
	app := cli.NewApp()
	app.Name = "go2ir"
	app.Usage = "Generating Synthesijer-IR from Go programming"
	app.Version = "0.1.1"
	app.Action = func(c *cli.Context) {
		if len(c.Args()) == 0 {
			fmt.Println("usage: go2ir sources")
			return
		}
		for _, src := range c.Args(){
			parse(src)
		}
	}
	app.Run(os.Args)
}

