package core

import (
	"fmt"
	"testing"

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
