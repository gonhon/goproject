package schema

import (
	"go/ast"
	"reflect"

	"github.com/limerence-code/goproject/gee/orm/dialect"
)

// 字段信息
type Field struct {
	Name string
	Type string
	//约束信息
	Tag string
}

type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	FieldMap   map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.FieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		FieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		structField := modelType.Field(i)
		if !structField.Anonymous && ast.IsExported(structField.Name) {
			field := &Field{
				Name: structField.Name,
				//根据方言获取类型
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(structField.Type))),
			}
			if v, ok := structField.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, structField.Name)
			schema.FieldMap[structField.Name] = field
		}
	}
	return schema
}

func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	value := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}

	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, value.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
