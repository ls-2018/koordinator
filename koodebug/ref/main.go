package main

import (
	"fmt"
	"reflect"
)

type dog struct {
}
type animal interface {
	speak()
}

func (a *dog) speak() {}

func main() {
	var a []dog
	plugins := reflect.ValueOf(&a).Elem()
	pluginType := plugins.Type().Elem()
	fmt.Println(pluginType)
	typAnimal := reflect.TypeOf((*animal)(nil)).Elem()
	fmt.Println(reflect.TypeOf(&dog{}).Implements(typAnimal))

}
