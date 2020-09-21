package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var myDB *sql.DB

func init() {
	db, err := sql.Open("mysql", "root:123456@(192.168.123.91:3307)/fileserver?charset=utf8")
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(1000)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	myDB = db
}

func DBConn() *sql.DB {
	fmt.Println("db:", myDB)
	return myDB
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	var ret = make([]map[string]interface{}, 0)
	// TODO:
	for rows.Next() {
		var tmp = make(map[string]interface{}, 2)
		var (
			username string
			userpwd  string
		)

		err := rows.Scan(&username, &userpwd)
		if err != nil {
			fmt.Println("err:", err)
		}
		tmp["user_name"] = username
		tmp["user_pwd"] = userpwd
		ret = append(ret, tmp)
	}
	return ret
}
