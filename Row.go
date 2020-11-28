package eDB

import (
	"encoding/json"
	"reflect"
	"strings"
)

type Row struct {
	columnValues []interface{}
}

func NewRow() *Row {

	return &Row{[]interface{}{}}
}

func (this *Row) SetColumn(index int, value interface{}) bool {
	if strings.Contains(reflect.TypeOf(value).String(), "*") {
		panic("set value err, Value cannot be a ptr ")
	}
	if this.GetSize() > index && index >= 0 {
		this.columnValues[index] = value
		return true
	} else if index == this.GetSize() {
		if value == nil {
			this.columnValues = append(this.columnValues, "NULL")
			return true
		}
		this.columnValues = append(this.columnValues, value)
		return true
	}
	panic("out of range by Row ")
}

func (this *Row) UpdateColumn(index int, value interface{}) bool {
	if strings.Contains(reflect.TypeOf(value).String(), "*") {
		panic("set value err, Value cannot be a ptr ")
	}
	if index >= 0 && index < this.GetSize() {
		if value == nil {
			this.columnValues[index] = "NULL"
			return true
		} else {
			this.columnValues[index] = value
			return true
		}
	}

	panic("out of range by Row ")

}

func (this *Row) GetColumnValues(index int) interface{} {
	return this.columnValues[index]
}

func (this *Row) GetSize() int {
	return len(this.columnValues)
}

func (this *Row) String() string {
	b, _ := json.Marshal(this.columnValues)
	return string(b)
}
