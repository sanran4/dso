package util

import (
	"encoding/json"
	"reflect"
	"strconv"
)

// PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

// RangeStructer takes the first argument, which must be a struct, and
// returns the value of each field in a slice. It will return nil
// if there are no arguments or first argument is not a struct
func RangeStructer(args ...interface{}) []interface{} {
	if len(args) == 0 {
		return nil
	}

	v := reflect.ValueOf(args[0])
	if v.Kind() != reflect.Struct {
		return nil
	}

	out := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		out[i] = v.Field(i).Interface()
	}

	return out
}

func StructConv(args ...interface{}) interface{} {
	if len(args) < 2 {
		return nil
	}
	arg1 := reflect.ValueOf(args[0])
	if arg1.Kind() != reflect.Struct {
		return nil
	}

	type ret struct {
		key string
		val string
	}

	fields := reflect.TypeOf(arg1)
	values := reflect.ValueOf(arg1)
	num := fields.NumField()

	rets := []ret{}
	//out := make([]interface{}, num)
	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)
		var r ret
		var v string
		switch value.Kind() {
		case reflect.String:
			v = value.String()
		case reflect.Int:
			v = (strconv.FormatInt(value.Int(), 10))
		case reflect.Int32:
			v = strconv.FormatInt(value.Int(), 10)
		case reflect.Int64:
			v = strconv.FormatInt(value.Int(), 10)
		default:
			v = value.String()
		}
		r.key = field.Name
		r.val = v
		//fmt.Print("Type:", field.Type, ",", field.Name, "=", value, "\n")
		rets = append(rets, r)
	}
	return rets
}
