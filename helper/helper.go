package helper

import (
	"encoding/csv"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func ReadCSV(file string, sep rune) [][]string {
	n, _ := os.ReadFile(file)

	w := strings.NewReader(string(n))
	r := csv.NewReader(w)

	r.Comma = sep

	res, err := r.ReadAll()
	if err != nil {
		panic(err.Error())
	}

	return res
}

func ParseFloat(value string) float64 {
	f, err := strconv.ParseFloat(value, 64)

	if err != nil {
		panic(err.Error())
	}

	return f
}

func IsStruct(value interface{}) reflect.Value {
	v := reflect.ValueOf(value)

	kind := v.Kind().String()

	if kind != "struct" {
		panic("value is not struct")
	}

	return v
}
