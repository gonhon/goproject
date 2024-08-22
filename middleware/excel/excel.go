/*
 * @Author: gaoh
 * @Date: 2024-08-14 14:24:11
 * @LastEditTime: 2024-08-22 19:00:44
 */
package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tealeg/xlsx"

	"github.com/spf13/viper"
)

var (
	config Config
)

// 定义config结构体
type Config struct {
	FileName  string
	TableName string
	//跳过前面行数
	SkipRows int
	Mysql    MysqlConfig
}

// json中的嵌套对应结构体的嵌套
type MysqlConfig struct {
	UserName string
	Password string
	Ip       string
	Port     int
	Database string
}

// 初始化配置
func init() {
	viper := viper.New()
	viper.AddConfigPath("./config")
	viper.SetConfigName("excel")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println(err)
	}
}

func main() {

	// 解析Excel文件
	file := config.FileName
	xlsxFile, err := xlsx.OpenFile(file)
	if err != nil {
		panic(err)
	}
	// 连接数据库
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.Mysql.UserName, config.Mysql.Password, config.Mysql.Ip, config.Mysql.Port, config.Mysql.Database))
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	name := config.TableName
	// 遍历所有Sheet
	for i, sheet := range xlsxFile.Sheets {
		// 检查数据库是否存在该表
		tableName := fmt.Sprintf("%s_%d", name, i+1)
		var count int
		db.QueryRow("SELECT count(*) as count FROM information_schema.tables WHERE table_schema = 'db_name' AND table_name = ?", tableName).Scan(&count)
		//// 如果存在,使用name_n作为表名
		if count <= 0 {
			// 获取字段定义
			fields := getFields(sheet)

			// 创建表
			stmt := fmt.Sprintf("CREATE TABLE %s (%s)", tableName, strings.Join(fields, ","))
			db.Exec(stmt)
		}
		wg.Add(1)
		go func() {
			// 插入数据
			num := insertData(db, tableName, sheet)
			fmt.Printf("%s-插入{%d}条数据-%s\n", sheet.Name, num, tableName)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("数据导入完毕!")
}

// 获取字段定义
func getFields(sheet *xlsx.Sheet) []string {
	var fields []string

	fields = append(fields, fmt.Sprintf("`%s` VARCHAR(255)", sheet.Name))
	for _, row := range sheet.Rows {
		for _, cell := range row.Cells {
			// 获取值类型
			switch cell.Type() {
			case xlsx.CellTypeString:
				fields = append(fields, fmt.Sprintf("`%s` VARCHAR(255)", cell.String()))
			case xlsx.CellTypeNumeric:
				fields = append(fields, fmt.Sprintf("`%s` INT", cell.String()))
			case xlsx.CellTypeDate:
				fields = append(fields, fmt.Sprintf("`%s` DATE", cell.String()))
			case xlsx.CellTypeBool:
				fields = append(fields, fmt.Sprintf("`%s` TINYINT(1)", cell.String()))
			}
		}
		break
	}
	return fields
}

// 插入数据
func insertData(db *sql.DB, name string, sheet *xlsx.Sheet) int64 {
	// 获取列名
	var colNames []string
	colNames = append(colNames, sheet.Name)
	for _, cell := range sheet.Rows[0].Cells {
		colNames = append(colNames, cell.String())
	}
	var count int64
	var values []string
	// 遍历行数据
	for idx, row := range sheet.Rows {
		if idx < config.SkipRows {
			continue
		}
		// 拼接值
		var colValues []string

		for _, cell := range row.Cells {
			switch cell.Type() {
			case xlsx.CellTypeString:
				colValues = append(colValues, "'"+cell.String()+"'")
			case xlsx.CellTypeNumeric:
				i, _ := cell.Int()
				colValues = append(colValues, strconv.Itoa(i))
			case xlsx.CellTypeBool:
				colValues = append(colValues, strconv.FormatBool(cell.Bool()))
			}
		}
		values = append(values, fmt.Sprintf("('%s',%s)", sheet.Name, strings.Join(colValues, ",")))
		if len(values) == 50 {
			count += saveData(name, db, colNames, values)
			values = values[:0]
		}
	}
	if len(values) > 0 {
		count += saveData(name, db, colNames, values)
	}
	return count
}

func saveData(name string, db *sql.DB, colNames, values []string) int64 {
	// 拼接SQL并执行
	stmt := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		name,
		strings.Join(colNames, ","),
		strings.Join(values, ","),
	)
	// fmt.Printf("%s<------>%s\n", sheet.Name, stmt)
	res, err := db.Exec(stmt)
	if err != nil {
		fmt.Printf("insert data error: %v\n", err)
		return 0
	}
	//获取影响行数
	count, _ := res.RowsAffected()
	return count
}
