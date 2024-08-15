/*
 * @Author: gaoh
 * @Date: 2024-08-14 14:24:11
 * @LastEditTime: 2024-08-14 17:49:53
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
)

func main() {
	// 解析Excel文件
	file := "./赣榆乡镇增值税号和水表类别完善.xlsx"
	xlsxFile, err := xlsx.OpenFile(file)
	if err != nil {
		panic(err)
	}

	// 连接数据库
	db, err := sql.Open("mysql", "wpg:DmJme(ZFl9txW@2P@tcp(10.10.102.205:3306)/waterdb_ys_gy")

	var wg sync.WaitGroup
	name := "client_water_info_all_6"
	// 遍历所有Sheet
	for _, sheet := range xlsxFile.Sheets {
		// 获取Sheet名作为表名
		//name := sheet.Name

		// 检查数据库是否存在该表
		var count int
		db.QueryRow("SELECT count(*) as count FROM information_schema.tables WHERE table_schema = 'db_name' AND table_name = ?", name).Scan(&count)
		//
		//// 如果存在,使用name_n作为表名
		if count <= 0 {
			//name = fmt.Sprintf("%s_%d", name, count)
			//}else {
			// 获取字段定义
			fields := getFields(sheet)

			// 创建表
			stmt := fmt.Sprintf("CREATE TABLE %s (%s)", name, strings.Join(fields, ","))
			db.Exec(stmt)
		}
		wg.Add(1)
		go func() {
			// 插入数据
			num := insertData(db, name, sheet)
			fmt.Printf("sheet:%s插入%d条数据\n", sheet.Name, num)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("数据导入完毕!")
}

// 获取字段定义
func getFields(sheet *xlsx.Sheet) []string {
	var fields []string

	fields = append(fields, fmt.Sprintf("`%s` VARCHAR(255)", "sheetName"))
	for _, row := range sheet.Rows {
		for _, cell := range row.Cells {
			// 获取值类型
			switch cell.Type() {
			case xlsx.CellTypeString:
				fields = append(fields, fmt.Sprintf("`%s` VARCHAR(255)", cell.String()))
			case xlsx.CellTypeNumeric:
				fields = append(fields, fmt.Sprintf("`%s` INT", cell.String()))
			// case xlsx.CellTypeInt:
			// 	fields = append(fields, fmt.Sprintf("`%s` INT", cell.String()))
			// case xlsx.CellTypeFloat:
			// 	fields = append(fields, fmt.Sprintf("`%s` FLOAT", cell.String()))
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
	colNames = append(colNames, "sheetName")
	for _, cell := range sheet.Rows[0].Cells {
		colNames = append(colNames, cell.String())
	}
	var count int64
	var values []string
	// 遍历行数据
	for idx, row := range sheet.Rows {
		if idx == 0 {
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
			// case xlsx.CellTypeFloat:
			// 	colValues = append(colValues, strconv.FormatFloat(cell.Num(), 'f', -1, 64))
			case xlsx.CellTypeBool:
				colValues = append(colValues, strconv.FormatBool(cell.Bool()))
				/* case xlsx.CellTypeDate:
				colValues = append(colValues, cell.Date().Format("2006-01-02"))
				*/
			}
		}
		values = append(values, fmt.Sprintf("('%s',%s)", sheet.Name, strings.Join(colValues, ",")))
		if len(values) == 50 {
			count += saveData(name, sheet, db, &colNames, &values, count)
		}
	}
	if len(values) > 0 {
		count += saveData(name, sheet, db, &colNames, &values, count)
	}
	return count
}

func saveData(name string, sheet *xlsx.Sheet, db *sql.DB, colNames, values *[]string, count int64) int64 {
	// 拼接SQL并执行
	stmt := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		name,
		strings.Join(*colNames, ","),
		strings.Join(*values, ","),
	)
	db.Exec(stmt)
	// INSERT INTO client_water_info_all_2 (户号,税号,水表类别) VALUES ('79001196','','')
	// fmt.Printf("%s<------>%s\n", sheet.Name, stmt)
	count += int64(len(*values))
	*values = (*values)[:0]
	return count
}
