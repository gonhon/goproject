package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "gee.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	//删除表
	db.Exec("Drop table IF EXISTS User;")
	//建表
	db.Exec("create table User(Name text);")

	//插入数据
	res, err := db.Exec("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam")
	if err == nil {
		affecte, _ := res.RowsAffected()
		log.Printf("insert %v", affecte)
	}

	//查询数据
	row := db.QueryRow("select * from User limit 1")

	var name string
	if err := row.Scan(&name); err == nil {
		log.Println(name)
	} else {
		log.Fatal(err)
	}

}
