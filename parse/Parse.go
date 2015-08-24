package main
import "log"
import "os"
import "strings"
import "github.com/SumoLogic/sumoshell/util"
import "regexp"

type Parser struct {
	pattern string
	extractions []string
	regex *regexp.Regexp
	output *util.JsonWriter
}

const Wildcard = '*'
const genericError = "parse takes arguments like: `parse \"[key=*]\" as key`\n"
func main() {

	if len(os.Args) < 2 {
		log.Printf("Error! No arguments provided.")
		log.Printf(genericError)
		return
	}
	parseExpression := os.Args[1]
	numExtractions := findNumExtractions(parseExpression)
	//         (parse, str)	(as)  (foo, bar, baz)			
	expectedArgs := 2 +       1  +  numExtractions
	if (len(os.Args) < expectedArgs) {
		log.Printf(genericError)
		log.Printf("Expected more arguments\n")
		return
	}
	as := os.Args[2]
	if (as != "as") {
		log.Printf(genericError)
		log.Printf("Expected `as`\n")
		return
	}
	extractions := make([]string, len(os.Args) - 3)
	for i, arg := range os.Args[3:] {
		extractions[i] = strings.Trim(arg, ",")
	}
	util.ConnectToStdIn(Parser{parseExpression, extractions, 
		regexFromPat(parseExpression), util.NewJsonWriter()})
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
	if (util.IsPlus(inp)) {
		matches := p.regex.FindStringSubmatch(util.ExtractRaw(inp))
		if len(matches) == 1 + len(p.extractions) {
			for i, match := range matches[1:] {
				inp[p.extractions[i]] = match
			}
			p.output.Write(inp)
		}
	}
}
