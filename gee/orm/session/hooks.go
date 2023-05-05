package session

import (
	"reflect"

	"github.com/limerence-code/goproject/gee/orm/log"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

func (s *Session) CallMethod(method string, val interface{}) {
	fm := reflect.ValueOf(s.refTable.Model).MethodByName(method)
	if val != nil {
		fm = reflect.ValueOf(val).MethodByName(method)
	}
	if fm.IsValid() {
		params := []reflect.Value{reflect.ValueOf(s)}
		if v := fm.Call(params); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
}
