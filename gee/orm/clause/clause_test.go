package clause

import (
	"reflect"
	"testing"
)

func TestSelect(t *testing.T) {
	var clause Clause
	clause.Set(LIMIT, 3)
	clause.Set(SELECT, "User", []string{"*"})
	clause.Set(WHERE, "Name = ? ", "Tom")
	clause.Set(OERDERBY, "Age ASC")
	sql, vars := clause.Build(SELECT, WHERE, OERDERBY, LIMIT)
	t.Log(sql, vars)
	if sql != "SELECT * FROM User WHERE Name = ?  ORDER BY Age ASC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}

}
