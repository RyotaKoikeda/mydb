package main

import (
	"database/sql"
	"fmt"
	"os"

	"utils/log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// MySQLへの接続設定
	dsn := "root:root@tcp(127.0.0.1:3306)/sakila"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 接続確認
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Info("Connected to MySQL successfully!")

	queryBytes, err := os.ReadFile("query.sql")
	if err != nil {
		log.Error(err.Error())
		return
	}
	query := string(queryBytes)

	rows, err := db.Query(query)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer rows.Close()

	// カラム情報取得
	cols, err := rows.Columns()
	if err != nil {
		log.Error(err.Error())
		return
	}

	// 各行の値を取得して表示
	values := make([]sql.RawBytes, len(cols))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Error(err.Error())
			return
		}
		for i, col := range values {
			fmt.Printf("%s: %s\t", cols[i], string(col))
		}
		fmt.Println()
	}
}
