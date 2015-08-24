package main
import (
	"os"
	"github.com/SumoLogic/sumoshell/util"
	"github.com/SumoLogic/sumoshell/group"
	"time"
)

type count struct {
	ct *int
	base map[string]interface{}
	output func(map[string]interface{})
}

func makeCount() count {
	ct := 0
	return count{&ct, make(map[string]interface{}), util.NewJsonWriter().Write}
}

func aggregateCount(output grouper.Merger, key string, base map[string]interface{}) util.SumoAggOperator {
	ct := 0
	ticker := time.NewTicker(100 * time.Millisecond)
	count := count{&ct, base, output.Write}
	go flush(count, ticker)
	return count
}

func main() {
	relevantArgs := os.Args[1:]

	if len(relevantArgs) == 0 {
		ticker := time.NewTicker(1 * time.Second)
		ct := makeCount()
		go flush(ct, ticker)
		util.ConnectToStdIn(ct)
		ct.Flush()
	} else if len(relevantArgs) > 0 {
		keyFields := relevantArgs
		// key is meaningless for count
		agg := grouper.NewAggregate(aggregateCount, keyFields, "")
		util.ConnectToStdIn(agg)
		agg.Flush()
	}
}

func flush(ct count, ticker *time.Ticker) {
	for {
		select {
			case <- ticker.C:	
				ct.Flush()
		}
	}
}

func (ct count)Flush() {
	ct.output(util.CreateStartRelation())
	ct.output(util.CreateRelation(currentState(ct)))
	ct.output(util.CreateEndRelation())
}

func currentState(ct count) map[string]interface{} {
	ret := make(map[string]interface{})
	for key, val := range ct.base {
		ret[key] = val
	}
	ret["_count"] = *ct.ct
	return ret
}

func (ct count) Process(inp map[string]interface{}) {
	if util.IsPlus(inp) {
		*ct.ct += 1
	}
}
