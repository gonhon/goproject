package session

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/limerence-code/goproject/gee/orm/log"
	"github.com/limerence-code/goproject/gee/orm/schema"
)

func (s *Session) Model(val interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(val) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(val, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is no set ...")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	tab := s.refTable
	var columns []string
	for _, field := range s.refTable.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", tab.Name, desc)).Exec()
	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, val := s.dialect.TableExistSql(s.refTable.Name)
	row := s.Raw(sql, val...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.refTable.Name
}
