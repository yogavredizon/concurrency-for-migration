package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"test.com/helper"
)

func SqlConnect(dsn, driver string) *sql.DB {
	db, err := sql.Open(driver, dsn)

	if err != nil {
		return nil
	}

	return db
}

func CreateTable(db *sql.DB, table interface{}) {
	value := helper.IsStruct(table)

	n := value.NumField()

	if n < 1 {
		panic("undefined field in table")
	}

	str := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (", strings.ToLower(value.Type().Name()))
	columnType := ""

	for i := 0; i < n; i++ {
		field := value.Type().Field(i)

		if tag, ok := field.Tag.Lookup("bson"); !ok {
			errText := fmt.Sprintf("Field %v doesn't have tag", field.Name)
			panic(errText)
		} else {
			switch field.Type.String() {
			case "int", "int8", "int16", "int32", "int64":
				columnType = "int"
			case "float", "float32", "float64":
				columnType = "float"
			case "string":
				columnType = "varchar(255)"
			case "bool":
				columnType = "boolean"
			}

			column := strings.ToLower(tag)

			if column == "_id" {
				column = column[1:]
				columnType += " primary key"
			}

			if field.Tag.Get("type") == "desc" {
				columnType = "TEXT"
			}
			str = str + column + " " + columnType + ","
		}

	}
	str = str[:len(str)-1] + ")"

	_, err := db.Exec(str)

	if err != nil {
		panic(err)
	}
}

func CreateRawInsert(value interface{}) (string, []interface{}) {
	table := helper.IsStruct(value)

	n := table.Type().NumField()

	if n < 1 {
		panic("Value nof valid")
	}

	paramMarks := ""
	fields := []interface{}{}

	for i := 0; i < n; i++ {
		paramMarks += "$" + strconv.Itoa(i+1) + ","
		fields = append(fields, table.Field(i).Interface())
	}

	paramMarks = paramMarks[:len(paramMarks)-1]
	raw := fmt.Sprintf("INSERT INTO %v VALUES(%v)", strings.ToLower(table.Type().Name()), paramMarks)

	return raw, fields
}
