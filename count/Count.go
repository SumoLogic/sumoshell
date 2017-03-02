package count

import (
	"github.com/SumoLogic/sumoshell/group"
	"github.com/SumoLogic/sumoshell/util"
	"sync"
)

const Count = "_count"

type count struct {
	ct     *int
	base   map[string]interface{}
	output func(map[string]interface{})
	mu     *sync.Mutex
}

func makeCount() count {
	ct := 0
	return count{&ct, make(map[string]interface{}), util.NewJsonWriter().Write, &sync.Mutex{}}
}

func aggregateCount(output grouper.Merger, key string, base map[string]interface{}) util.SumoAggOperator {
	ct := 0
	count := count{&ct, base, output.Write, &sync.Mutex{}}
	return count
}

func Build(args []string) (util.SumoAggOperator, error) {
	args = args[1:]
	if len(args) == 0 {
		return makeCount(), nil
	} else {
		keyFields := args
		// key is meaningless for count
		return grouper.NewAggregate(aggregateCount, keyFields, "", Count), nil
	}
}

func (ct count) Flush() {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	ct.output(util.CreateStartRelation())
	ct.output(util.CreateRelation(currentState(ct)))
	ct.output(util.CreateEndRelation())
}

func currentState(ct count) map[string]interface{} {
	ret := make(map[string]interface{})
	for key, val := range ct.base {
		ret[key] = val
	}
	ret[Count] = *ct.ct
	return ret
}

func (ct count) Process(inp map[string]interface{}) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	if util.IsStartRelation(inp) {
		*ct.ct = 0
	} else if util.IsEndRelation(inp) {
		//ct.Flush()
	} else if util.IsPlus(inp) || util.IsRelation(inp) {
		*ct.ct += 1
	}
}
