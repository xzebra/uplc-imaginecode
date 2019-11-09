package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB
var (
	DBUser = "BN7SqiyTM6"
	DBName = "BN7SqiyTM6"
	DBPassword = "Fp0l0c0HBK"
)

func DBInit() {
	var err error
	DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUser, DBPassword, "remotemysql.com", "3306", DBName))
	if err != nil { panic(err) }

	err = DB.Ping()
	if err != nil { panic(err) }

	fmt.Println("connected to database")
}

