package grouper

import (
	"fmt"
	"github.com/SumoLogic/sumoshell/util"
	"os"
	"strings"
	"sync"
	"time"
)

type Grouper struct {
	constructor func(map[string]interface{}) util.SumoAggOperator
	operators   map[string]util.SumoAggOperator
	merger      Merger
	by          []string
	key         string
	mu          *sync.Mutex
}

type builder func(Merger, string, map[string]interface{}) util.SumoAggOperator

func NewAggregate(
	constructor builder,
	by []string,
	key string,
	sortCol string) Grouper {
	merger := NewMerger(sortCol)
	ctor := func(base map[string]interface{}) util.SumoAggOperator {
		return constructor(merger, key, base)
	}

	mu := &sync.Mutex{}
	return Grouper{ctor, make(map[string]util.SumoAggOperator), merger, by, key, mu}
}

func (g Grouper) Flush() {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, v := range g.operators {
		v.Flush()
	}
	g.merger.Flush()
}

func (g Grouper) Process(inp map[string]interface{}) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var keys []string
	for _, key := range g.by {
		val, ok := inp[key]
		if ok {
			keys = append(keys, fmt.Sprint(val))
		} else {
			keys = append(keys, "")
		}
	}

	groupKey := strings.Join(keys, "-")

	_, ok := g.operators[groupKey]
	if !ok {
		nextId := len(g.operators)
		base := make(map[string]interface{})
		for i, key := range g.by {
			base[key] = keys[i]
		}
		base[Id] = nextId
		g.operators[groupKey] = g.constructor(base)
	}
	op, _ := g.operators[groupKey]
	op.Process(inp)
}

type Merger struct {
	// one map for each grouper
	aggregate map[int]map[string]interface{}
	output    *util.JsonWriter
	mu        *sync.Mutex
	sortCol   string
}

func NewMerger(sortCol string) Merger {
	mu := &sync.Mutex{}
	m := Merger{make(map[int]map[string]interface{}), util.NewJsonWriter(), mu, sortCol}
	ticker := time.NewTicker(100 * time.Millisecond)
	go flush(m, ticker)
	return m
}

const Id = "_Id"

func ExtractId(inp map[string]interface{}) int {
	raw, ok := inp[Id].(int)
	if ok {
		return raw
	} else {
		return -1
	}
}

func WithId(id int) map[string]interface{} {
	return map[string]interface{}{Id: id}
}

func (m Merger) Process(inp map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !util.IsStartRelation(inp) && !util.IsEndRelation(inp) {
		m.aggregate[ExtractId(inp)] = inp
	}
}

func (m Merger) Write(inp map[string]interface{}) {
	m.Process(inp)
}

func (m Merger) Flush() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.output.Write(util.CreateStartRelationMeta("merger"))
	// Output keys sorted by index so the ui is consistent
	if m.sortCol == "" {
		for i := 0; i < len(m.aggregate); i++ {
			m.output.Write(util.CreateRelation(m.aggregate[i]))
		}
	} else {
		aggs := make([]map[string]interface{}, 0, len(m.aggregate))
		for i := 0; i < len(m.aggregate); i++ {
			aggs = append(aggs, m.aggregate[i])
		}
		util.SortByField(m.sortCol, aggs)
		for i := 0; i < len(aggs); i++ {
			m.output.Write(util.CreateRelation(aggs[i]))

		}
	}
	m.output.Write(util.CreateEndRelation())
	queryString := strings.Join(os.Args[0:], " ")
	m.output.Write(util.CreateMeta(map[string]interface{}{"_queryString": queryString}))
}

func flush(m Merger, ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			m.Flush()
		}
	}
}
