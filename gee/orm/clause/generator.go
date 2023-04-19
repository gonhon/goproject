package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[OERDERBY] = _orderBy
}

func genBindVars(num int) string {
	var vars []string

	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

// INSERT INTO $tableName ($fields)
func _insert(vals ...interface{}) (string, []interface{}) {
	//表名
	tableName := vals[0]
	//将字符串数组按照,连接
	fields := strings.Join(vals[1].([]string), ",")
	return fmt.Sprintf("INSERT INFO %s (%v)", tableName, fields), []interface{}{}
}

// VALUES ($v1), ($v2), ...
func _values(vals ...interface{}) (string, []interface{}) {
	var bindStr string
	var sql strings.Builder
	var vs []interface{}
	sql.WriteString("VALUES ")

	for i, val := range vals {
		v := val.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		//非最后一个添加,
		if i+1 != len(vals) {
			sql.WriteString(",")
		}
		vs = append(vs, v...)
	}
	return sql.String(), vs

}

// SELECT $fields FROM $tableName
func _select(vals ...interface{}) (string, []interface{}) {
	tableName := vals[0]
	fields := strings.Join(vals[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}

}

// LIMIT $num
func _limit(vals ...interface{}) (string, []interface{}) {
	return "LIMIT ?", vals
}

// WHERE $desc
func _where(vals ...interface{}) (string, []interface{}) {
	desc, vs := vals[0], vals[1:]
	return fmt.Sprintf("WHERE %s", desc), vs
}
func _orderBy(vals ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", vals[0]), []interface{}{}
}
