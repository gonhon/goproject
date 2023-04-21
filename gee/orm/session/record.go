package session

import (
	"reflect"

	"github.com/limerence-code/goproject/gee/orm/clause"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordVars := make([]interface{}, 0)
	for _, val := range values {
		tab := s.Model(val).RefTable()
		//insert
		s.clause.Set(clause.INSERT, tab.Name, tab.FieldNames)
		recordVars = append(recordVars, tab.RecordValues(val))
	}
	//value绑定
	s.clause.Set(clause.VALUES, recordVars...)
	//获取sql和参数
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	//执行sql
	res, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	tab := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	s.clause.Set(clause.SELECT, tab.Name, tab.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.OERDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range tab.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}

		if err := rows.Scan(values...); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))

	}
	return rows.Close()
}
