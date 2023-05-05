package session

import (
	"errors"
	"reflect"

	"github.com/limerence-code/goproject/gee/orm/clause"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordVars := make([]interface{}, 0)
	for _, val := range values {
		s.CallMethod(BeforeInsert, val)
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
	s.CallMethod(AfterInsert, nil)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *Session) Find(values interface{}) error {
	s.CallMethod(BeforeQuery, nil)

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

		err := rows.Scan(values...)
		s.CallMethod(AfterQuery, dest.Addr().Interface())
		if err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

func (s *Session) Update(kv ...interface{}) (int64, error) {
	s.CallMethod(BeforeUpdate, nil)
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.refTable.Name, m)
	sql, vals := s.clause.Build(clause.UPDATE, clause.WHERE)
	res, err := s.Raw(sql, vals...).Exec()
	s.CallMethod(AfterUpdate, nil)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	s.CallMethod(BeforeDelete, nil)
	s.clause.Set(clause.DELETE, s.refTable.Name)
	sql, vals := s.clause.Build(clause.UPDATE, clause.WHERE)
	res, err := s.Raw(sql, vals...).Exec()
	s.CallMethod(AfterDelete, nil)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.refTable.Name)
	sql, vals := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vals...).QueryRow()
	var count int64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.OERDERBY, desc)
	return s
}

func (s *Session) First(vals interface{}) error {

	dest := reflect.Indirect(reflect.ValueOf(vals))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("Not Found")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
