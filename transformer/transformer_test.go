package transformer

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"
)

type TT struct {
	A int
	B string
}

func TransformA(slice []string) []interface{} {
	var res []interface{}
	log.Print(len(slice))
	for _, content := range slice {
		/**
		value, ok := content.(string)
		if !ok {
			continue
		}
		*/
		var info interface{}

		err := json.Unmarshal([]byte(content), &info)
		if err != nil {
			log.Fatal(err)
			continue
		}
		res = append(res, info)
	}
	log.Print(res)
	return res
}

func TestTransform(t *testing.T) {
	tt := TT{12, "skidoo"}
	ss := `{"A":12, "B":"skidoo", "C": {"D": "E"}}`
	var slice []string
	slice = append(slice, ss)
	TransformA(slice)
	s := reflect.ValueOf(&tt).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d %s %s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}

func TestReflect(t *testing.T) {
	var x int = 1
	fmt.Println("type: ", reflect.TypeOf(x))
	fmt.Println("value: ", reflect.ValueOf(x))
	fmt.Println("Kind:  ", reflect.ValueOf(x).Kind())

	y := reflect.ValueOf(x).Interface()
	fmt.Println(y)
}
