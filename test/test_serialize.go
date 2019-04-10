package main

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

var (
	loop = int64(100000)
)

type A struct {
	Str  string
	Data []byte
	Int  int64
	Bool bool
}

type B struct {
	ArrS    []string
	Stru    A
	ArrStru []A
}

func testBson() {
	b := &B{
		ArrS: []string{"abc", "测试数据111", "hello world"},
		Stru: A{
			Str:  "a string",
			Data: []byte("a data"),
			Int:  10998,
			Bool: true,
		},
		ArrStru: []A{
			A{
				Str:  "a string",
				Data: []byte("a data"),
				Int:  10998,
				Bool: true,
			},
			A{
				Str:  "a string",
				Data: []byte("a data"),
				Int:  10998,
				Bool: true,
			},
		},
	}
	ret := &B{}
	data := []byte{}
	t := time.Now()
	for i := int64(0); i < loop; i++ {
		data, _ = bson.Marshal(b)
		if err := bson.Unmarshal(data, ret); err != nil {
			panic(fmt.Errorf("bson.Unmarshal failed: ", err))
		}
	}
	tused := time.Since(t)
	fmt.Printf("bson len: %v, equal: %v, %v\n", len(data), ret.ArrStru[1].Str == b.ArrStru[1].Str, ret.ArrStru[1].Int == b.ArrStru[1].Int)
	fmt.Printf("bson: %v, %v / ns\n", tused.Seconds(), tused.Nanoseconds()/loop)
}

func testJson() {
	b := &B{
		ArrS: []string{"abc", "测试数据111", "hello world"},
		Stru: A{
			Str:  "a string",
			Data: []byte("a data"),
			Int:  10998,
			Bool: true,
		},
		ArrStru: []A{
			A{
				Str:  "a string",
				Data: []byte("a data"),
				Int:  10998,
				Bool: true,
			},
			A{
				Str:  "a string",
				Data: []byte("a data"),
				Int:  10998,
				Bool: true,
			},
		},
	}
	ret := &B{}
	data := []byte{}
	t := time.Now()
	for i := int64(0); i < loop; i++ {
		data, _ = json.Marshal(b)
		if err := json.Unmarshal(data, ret); err != nil {
			panic(fmt.Errorf("json.Unmarshal failed: ", err))
		}
	}
	tused := time.Since(t)
	fmt.Printf("json len: %v, equal: %v, %v\n", len(data), ret.ArrStru[1].Str == b.ArrStru[1].Str, ret.ArrStru[1].Int == b.ArrStru[1].Int)
	fmt.Printf("json: %v, %v / ns\n", tused.Seconds(), tused.Nanoseconds()/loop)
}

func testMsgpack() {
	b := &B{
		ArrS: []string{"abc", "测试数据111", "hello world"},
		Stru: A{
			Str:  "a string",
			Data: []byte("a data"),
			Int:  10998,
			Bool: true,
		},
		ArrStru: []A{
			A{
				Str:  "a string",
				Data: []byte("a data"),
				Int:  10998,
				Bool: true,
			},
			A{
				Str:  "a string",
				Data: []byte("a data"),
				Int:  10998,
				Bool: true,
			},
		},
	}
	ret := &B{}
	data := []byte{}
	t := time.Now()
	for i := int64(0); i < loop; i++ {
		data, _ = msgpack.Marshal(b)
		if err := msgpack.Unmarshal(data, ret); err != nil {
			panic(fmt.Errorf("msgpack.Unmarshal failed: ", err))
		}

	}
	tused := time.Since(t)
	fmt.Printf("msgpack len: %v, equal: %v, %v\n", len(data), ret.ArrStru[1].Str == b.ArrStru[1].Str, ret.ArrStru[1].Int == b.ArrStru[1].Int)
	fmt.Printf("msgpack: %v, %v / ns\n", tused.Seconds(), tused.Nanoseconds()/loop)
}

func main() {
	testBson()
	testJson()
	testMsgpack()
}
