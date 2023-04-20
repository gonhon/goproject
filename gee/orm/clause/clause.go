package clause

import "strings"

type Type int

type Clause struct {
	//存储sql
	sql map[Type]string
	//存储sql与编译后代替换的值
	sqlVars map[Type][]interface{}
}

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	OERDERBY
)

func (c *Clause) Set(name Type, values ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](values...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

func (c *Clause) Build(types ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}

	for _, tp := range types {
		if sql, ok := c.sql[tp]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[tp]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
