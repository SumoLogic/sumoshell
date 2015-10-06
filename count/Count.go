package count

import (
	"github.com/SumoLogic/sumoshell/group"
	"github.com/SumoLogic/sumoshell/util"
)

type count struct {
	ct     *int
	base   map[string]interface{}
	output func(map[string]interface{})
}

func makeCount() count {
	ct := 0
	return count{&ct, make(map[string]interface{}), util.NewJsonWriter().Write}
}

func aggregateCount(output grouper.Merger, key string, base map[string]interface{}) util.SumoAggOperator {
	ct := 0
	count := count{&ct, base, output.Write}
	return count
}

func Build(args []string) (util.SumoAggOperator, error) {
	if len(args) == 0 {
		return makeCount(), nil
	} else {
		keyFields := args
		// key is meaningless for count
		return grouper.NewAggregate(aggregateCount, keyFields, ""), nil
	}
}

func (ct count) Flush() {
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
