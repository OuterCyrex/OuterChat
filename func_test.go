package main

import (
	"OuterChat/model"
	"fmt"
	"reflect"
	"testing"
)

func TestValidator(t *testing.T) {
	Field := "Name"
	data := "Outer"

	user := model.UserBasic{}
	val := reflect.ValueOf(&user).Elem()
	typeOf := reflect.TypeOf(&user).Elem()

	if _, ok := typeOf.FieldByName(Field); ok {
		val.FieldByName(Field).Set(reflect.ValueOf(data))
	}

	fmt.Println("Name:", user.Name)
}
