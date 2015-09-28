package main

import (
	"fmt"
	"github.com/SumoLogic/sumoshell/parse"
	"github.com/SumoLogic/sumoshell/filter"
	"github.com/SumoLogic/sumoshell/util"
	"os"
)

type Builder func([]string) (util.SumoOperator, error)

func main() {
	operators := map[string]Builder{
		"parse": parse.Build,
		"filter": filter.Build,
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

		operator, err := builder(os.Args[1:])
		if err != nil {
			fmt.Println(err)
		} else {
			util.ConnectToStdIn(operator)
		}
	}
}
