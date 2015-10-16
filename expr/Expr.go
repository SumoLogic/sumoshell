package expr 

import (
	"github.com/SumoLogic/sumoshell/util"
	"strings"
	"os/exec"
	"fmt"
)

type ExprOperator struct {
	lsh string
	expr string
	output *util.JsonWriter
}

const genericError = "expr takes arguments like: newV = 5 + 5 * x"

func Build(args []string) (util.SumoOperator, error) {

	lsh := args[1]
	eq := args[2]
	if eq != "=" {
		return nil, util.ParseError("Expected `=` found `" + eq + "`\n" + genericError)
	}
	value := strings.Join(args[3:], " ")

	return &ExprOperator{lsh, value, util.NewJsonWriter()}, nil
}

func (w ExprOperator) Process(inp map[string]interface{}) {
	if util.IsPlus(inp) {
		vars := []string{}
		for k, v := range inp {
			_, isNumber := util.CoerceNumber(v)
			var s string
			if isNumber == nil {
				s = fmt.Sprintf("%v=%v", k, v)
			} else {
				s = fmt.Sprintf("%v=\"%v\"", k, v)
			}
	        	vars = append(vars, s)	
		}
		varStr := strings.Join(vars, ";")
		pythonCmd := varStr + "; print " + w.expr
		cmdStr := []string{"-c", pythonCmd}
		cmd := exec.Command("python", cmdStr...)
		out, err := cmd.Output()
		if err != nil {
			fmt.Println("ERROR", pythonCmd)
			return
		}
		strOut := strings.Trim(string(out), "\n")
		numRep, isNum := util.CoerceNumber(strOut)
		if isNum == nil {
			inp[w.lsh] = numRep
		} else {
			inp[w.lsh] = strOut
		}	
		w.output.Write(inp)
	}
}
