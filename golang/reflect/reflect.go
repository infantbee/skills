package main

import (
	"fmt"
	"os/exec"
	"reflect"
)

type order struct {
	id         int
	customerid int
	comment    string
}

type employee struct {
	id      int
	name    string
	address string
	salary  int
	country string
}

// reflect.Value
func parseValue(v reflect.Value) {
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			viv := v.Field(i)
			fmt.Printf("\n3: field[%d], value:[%v], type:[%T], Type:[%v], Name:[%v]", i, viv, viv, viv.Type(), viv.Type().Name())
			parseValue(viv)
		}

	case reflect.Array:

	case reflect.Map:

	case reflect.String:

	default:
		fmt.Println("Unsupported type")
		return
	}
}

func parseInterface(in, out interface{}) {
	t := reflect.TypeOf(in)
	v := reflect.ValueOf(in)
	fmt.Println("1: type: ", t, ", value: ", v)
	fmt.Println("2: t.kind: ", t.Kind(), ", v.kind: ", v.Kind())
	//reflect.Type: 代表interface{}实际类型main.order; 而reflect.Kind: 代表具体类型struct
	parseValue(v)

}

func main() {
	o := order{id: 1, customerid: 1, comment: "order table"}
	parseInterface(o, o)

	e := employee{id: 10, name: "xiaoming", address: "wuhan", salary: 100, country: "china"}
	parseInterface(e, e)

	cpath := "/home/leizi/Desktop/cpath"
	apath := "/home/leizi/Desktop/apath"

	var cmd *exec.Cmd
	cmd = exec.Command("mv", cpath, apath)
	ob, err1 := cmd.Output()
	fmt.Printf("\nmove file error:[%#v], resp:[%s]", err1, string(ob))

}
