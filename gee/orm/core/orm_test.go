package core

import (
	"errors"
	"fmt"
	"testing"

	"github.com/limerence-code/goproject/gee/orm/session"
	_ "github.com/mattn/go-sqlite3"
)

func TestEngineRun(t *testing.T) {
	engine, _ := NewEngine("sqlite3", "gee.db")
	defer engine.Close()

	session, _ := engine.NewSession()

	session.Raw("Drop table IF EXISTS User;").Exec()
	session.Raw("create table User(Name text);").Exec()
	session.Raw("create table User(Name text);").Exec()

	result, _ := session.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)

}

//----------------------事务----------------------

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "gee.db")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestTransactionRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s, _ := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return nil, errors.New("Error")
	})
	if err == nil || s.HasTable() {
		t.Fatal("failed to rollback")
	}
}

func TestTransactionCommit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s, _ := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return
	})
	u := &User{}
	_ = s.First(u)
	if err != nil || u.Name != "Tom" {
		t.Fatal("failed to commit")
	}
}
