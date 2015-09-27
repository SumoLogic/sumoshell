package main

import (
	"fmt"
	"github.com/SumoLogic/sumoshell/parse"
	"github.com/SumoLogic/sumoshell/util"
	"os"
)

type Builder func([]string) (error, util.SumoOperator)

func main() {
	operators := map[string]Builder{
		"parse": parse.Build,
	}

	args := os.Args
	if len(args) == 1 {
		fmt.Println("Arguments expected")
	} else {
		selectingArg := args[1]
		builder, ok := operators[selectingArg]
		if !ok {
			fmt.Println("Operator " + selectingArg + " not found")
			return
		}

		err, operator := builder(os.Args[1:])
		if err != nil {
			fmt.Println(err)
		} else {
			util.ConnectToStdIn(operator)
		}
	}
}
