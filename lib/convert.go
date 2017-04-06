package lib

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

func InterfaceToFloat64(v interface{}) (f float64, err error) {
	kinds := []string{reflect.String.String(), reflect.Float32.String(), reflect.Float64.String(),
		reflect.Int.String(), reflect.Int8.String(), reflect.Int32.String(), reflect.Int64.String()}

	vv := reflect.ValueOf(v)
	vk := reflect.ValueOf(v).Type().Kind().String()
	if !InArray(kinds, vk) {
		err = errors.New("v kind must in :" + strings.Join(kinds, ","))
		return
	}

	switch vk {
	case reflect.String.String():
		f, err = strconv.ParseFloat(vv.String(), 64)
		return
	case reflect.Float32.String(), reflect.Float64.String():
		f = vv.Float()
		return
	}

	err = errors.New("not support other kind")
	return
}
