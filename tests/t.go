package main

import (
	"fmt"
	"reflect"
)

type A struct {
	Name string
}

func SliceResource2SliceRuntimeObj(resources interface{}) []*A {
	res := make([]*A, 0)
	if reflect.TypeOf(resources).Kind() == reflect.Slice {
		s := reflect.ValueOf(resources)
		for i := 0; i < s.Len(); i++ {
			ele := s.Index(i)
			res = append(res, ele.Interface().(*A))
		}
	}
	return res

}

func main() {
	rs := []*A{&A{
		Name: "a",
	}}
	rs2 := SliceResource2SliceRuntimeObj(rs)
	fmt.Printf("%T\n", rs2[0])

}
