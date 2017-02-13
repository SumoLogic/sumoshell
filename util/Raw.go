package util

import (
	"io"
	"os"
	"fmt"
	"bufio"
	"encoding/json"
	"log"
	"bytes"
	"sort"
	"strconv"
	"strings"
	"sync"
	"github.com/reiver/go-whitespace"
)

type RawInputHandler struct {
	output io.Writer
	buff   []rune
}

type SumoOperator interface {
	Process(map[string]interface{})
}

type SumoAggOperator interface {
	Process(map[string]interface{})
	Flush()
}

type ParseError string

func (e ParseError) Error() string {
	return string(e)
}

func ConnectToStdIn(operator SumoOperator) {
	fi, _ := os.Stdin.Stat() // get the FileInfo struct describing the standard input.
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		ConnectToReader(operator, os.Stdin)
	} else {
		fmt.Println("No input")
		return
	}
}

func ConnectToReader(operator SumoOperator, reader io.Reader) {
	bio := bufio.NewReader(reader)
	var line, hasMoreInLine, err = bio.ReadLine()
	var buf []byte
	for err != io.EOF {
		buf = append(buf, line...)
		if !hasMoreInLine && len(buf) > 0 {
			var rawMsg interface{}
			err := json.Unmarshal(buf, &rawMsg)
			buf = []byte{}
			if err != nil {
				log.Println("Error parsing json", err)
				log.Println(string(line))
			} else {
				mapMessage, ok := rawMsg.(map[string]interface{})
				if ok {
					operator.Process(mapMessage)
				} else {
					log.Println("Unexpected JSON")
				}
			}
		}
		line, hasMoreInLine, err = bio.ReadLine()
	}
}

const Plus = "PLUS"
const StartRelation = "StartRelation"
const EndRelation = "EndRelation"
const Raw = "_raw"
const Type = "_type"
const Meta = "Meta"
const Relation = "Relation"

func IsPlus(inp map[string]interface{}) bool {
	tpe, ok := inp[Type].(string)
	return ok && tpe == Plus
}

func IsStartRelation(inp map[string]interface{}) bool {
	tpe, ok := inp[Type].(string)
	return ok && tpe == StartRelation
}

func IsEndRelation(inp map[string]interface{}) bool {
	tpe, ok := inp[Type].(string)
	return ok && tpe == EndRelation
}

func IsRelation(inp map[string]interface{}) bool {
	tpe, ok := inp[Type].(string)
	return ok && tpe == Relation
}

func IsMeta(inp map[string]interface{}) bool {
	tpe, ok := inp[Type].(string)
	return ok && tpe == Meta
}

func CreateStartRelation() map[string]interface{} {
	return map[string]interface{}{Type: StartRelation}
}

func CreateStartRelationMeta(origin string) map[string]interface{} {
	return map[string]interface{}{Type: StartRelation, Meta: origin}
}

func CreateEndRelation() map[string]interface{} {
	return map[string]interface{}{Type: EndRelation}
}

func CreateRelation(inp map[string]interface{}) map[string]interface{} {
	inp[Type] = Relation
	return inp
}

func CreateMeta(inp map[string]interface{}) map[string]interface{} {
	inp[Type] = Meta
	return inp
}

func ExtractRaw(inp map[string]interface{}) string {
	raw, ok := inp[Raw].(string)
	if ok {
		return strings.TrimRight(raw, "\n")
	} else {
		return ""
	}
}

func (handler *RawInputHandler) Process(inp []byte) {
	runes := bytes.Runes(inp)
	// If not whitespace, flush, append
	if len(runes) > 0 && !whitespace.IsWhitespace(runes[0]) {
		handler.Flush()
		handler.buff = append(handler.buff, runes...)
	} else {
		// If it is whitespace, just append with a newline
		handler.buff = append(handler.buff, '\n')
		handler.buff = append(handler.buff, runes...)
	}
}

func NewRawInputHandler(inp io.Writer) *RawInputHandler {
	return &RawInputHandler{inp, []rune{}}
}

func (handler *RawInputHandler) Flush() {
	m := make(map[string]interface{})
	m[Raw] = string(handler.buff)
	m[Type] = Plus
	handler.buff = []rune{}
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("ERROR!", err)
	} else {
		handler.output.Write(b)
		handler.output.Write([]byte{'\n'})
	}
}

type JsonWriter struct {
	writer io.Writer
	mu     *sync.Mutex
}

func NewJsonWriter() *JsonWriter {
	return &JsonWriter{os.Stdout, &sync.Mutex{}}
}

func (writer *JsonWriter) Write(inp map[string]interface{}) {
	jsonBytes, err := json.Marshal(inp)
	//fmt.Printf(b)
	if err != nil {
		fmt.Printf("ERROR!", err)
	} else {
		writer.mu.Lock()
		writer.writer.Write(jsonBytes)
		writer.writer.Write([]byte{'\n'})
		writer.mu.Unlock()
	}
}

func CoerceNumber(v interface{}) (float64, error) {
	return strconv.ParseFloat(fmt.Sprint(v), 64)
}

type Datum []map[string]interface{}
type By func(p1, p2 *map[string]interface{}) bool

func (a datumSorter) Len() int      { return len(a.data) }
func (a datumSorter) Swap(i, j int) { a.data[i], a.data[j] = a.data[j], a.data[i] }

// planetSorter joins a By function and a slice of Planets to be sorted.
type datumSorter struct {
	data Datum
	by   func(p1, p2 map[string]interface{}) bool // Closure used in the Less method.
}

func (a datumSorter) Less(i, j int) bool {
	return a.by(a.data[i], a.data[j])
}

func SortByField(field string, data Datum) {
	by := func(p1, p2 map[string]interface{}) bool {
		v1, err1 := CoerceNumber(p1[field])
		v2, err2 := CoerceNumber(p2[field])

		if err1 != nil || err2 != nil {
			fmt.Print(data)
			panic(err1)
		}

		return v1 < v2
	}
	sort.Sort(sort.Reverse(&datumSorter{data, by}))
}
