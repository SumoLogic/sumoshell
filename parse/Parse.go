package parse

import "log"
import "strings"
import "github.com/SumoLogic/sumoshell/util"
import "regexp"

type Parser struct {
	pattern     string
	extractions []string
	regex       *regexp.Regexp
	output      *util.JsonWriter
}

const Wildcard = '*'
const genericError = "parse takes arguments like: `parse \"[key=*]\" as key`\n"

func Build(args []string) (util.SumoOperator, error) {
	// [parse x as y, z, w]
	if len(args) < 2 {
		log.Printf("Error! No arguments provided.")
		log.Printf(genericError)
		return nil, util.ParseError("Error! No arguments provided\n" + genericError)
	}
	parseExpression := args[1]
	numExtractions := findNumExtractions(parseExpression)
	//         (parse pattern)	(as)  (foo, bar, baz)
	expectedArgs := 2 + 1 + numExtractions
	if len(args) < expectedArgs {
		return nil, util.ParseError("Expected more arguments\n" + genericError)
	}
	as := args[2]
	if as != "as" {
		return nil, util.ParseError("Expacted `as` got " + as + "\n" + genericError)
	}
	extractions := make([]string, len(args)-3)
	for i, arg := range args[3:] {
		extractions[i] = strings.Trim(arg, ",")
	}
	ret := Parser{parseExpression, extractions, regexFromPat(parseExpression), util.NewJsonWriter()}
	return &ret, nil
}

func findNumExtractions(parseExpression string) int {
	var count = 0
	for _, element := range parseExpression {
		if element == Wildcard {
			count += 1
		}
	}
	return count
}

func regexFromPat(pat string) *regexp.Regexp {
	regex := ".*?" + strings.Replace(regexp.QuoteMeta(pat), "\\*", "(.*?)", -1) + ".*"
	return regexp.MustCompile(regex)
}

func (p Parser) Process(inp map[string]interface{}) {
	if util.IsPlus(inp) {
		matches := p.regex.FindStringSubmatch(util.ExtractRaw(inp))
		if len(matches) == 1+len(p.extractions) {
			for i, match := range matches[1:] {
				inp[p.extractions[i]] = match
			}
			p.output.Write(inp)
		}
	}
}
