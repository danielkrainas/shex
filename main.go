package main

import (
	//"log"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

func main() {
	help := &helpCommand{}
	println("")
	_, err := commandParser.Parse()
	if err != nil {
		if aerr, ok := err.(appError); ok {
			fmt.Println(aerr.reason)
		} else if _, ok := err.(usageError); ok {
			help.Execute([]string{})
		} else if ferr, ok := err.(*flags.Error); ok {
			if ferr.Type == flags.ErrCommandRequired || ferr.Type == flags.ErrRequired {
				var args []string
				if commandParser.Active != nil {
					args = []string{commandParser.Active.Name}
				}

				help.Execute(args)
			}
		} else {
			println(err.Error())
			os.Exit(1)
		}
	}

	println("")
}
