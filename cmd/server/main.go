package main

import (
	"flag"
	"fmt"
	"godrop/internal/common/colored"
	"os"
)

var (
	usage = `All Commands`
)

func main() {
	commands := []Command{&ServerCommand{}, &GetCommand{}}
	buildUsage(commands)
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	for _, cmd := range commands {
		if os.Args[1] == cmd.Name() {
			cmd.Execute()
		}
	}
	flag.Usage()

}

func buildUsage(cmds []Command) {

	usage += "\n"

	for _, cmd := range cmds {
		usage += colored.BuildColoredString("- <RED>%s<RED> : <GREEN>%s<GREEN>\n", cmd.Name(), cmd.Describe())
	}

	usage += colored.BuildColoredString("<BLUE>if you don't know how to use it use -h<BLUE>")
	flag.Usage = func() {
		fmt.Printf("%s\n", usage)
	}
}
