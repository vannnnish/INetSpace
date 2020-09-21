package db

import (
	"fmt"
	"netspace/db/mysql"
)

// UserSignUp  通过用户及密码创建
func UserSignUp(username string, passwd string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_user(`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		fmt.Println("Failed to insert :", err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("Failed to insert err:", err)
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}
	return false
}

func UserSignIn(username string, encpwd string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"select user_name,user_pwd from tbl_user where user_name = ? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found:", username)
		return false
	}
	pRows := mysql.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].(string)) == encpwd {
		return true
	}
	return false
}

func UpdateToken(userame string, token string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"replace into tbl_user_token(`user_name`,`user_token`) values(?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(userame, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

func GetUserInfo(username string) (User, error) {
	user := User{}
	stmt, err := mysql.DBConn().Prepare(
		"select user_name,signup_at from tbl_user where user_name = ? limit 1")
	if err != nil {
		fmt.Println(err)
		return user, err
	}
	defer stmt.Close()
	// 执行查询操作
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil

}
