package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vaughan0/go-ini"
	"strconv"
)

var conn_str string

func load_config(filepath string) {
	file, _ := ini.LoadFile(filepath)
	db_user, _ := file.Get("bugspad", "user")
	db_pass, _ := file.Get("bugspad", "password")
	db_host, _ := file.Get("bugspad", "host")
	db_name, _ := file.Get("bugspad", "database")
	conn_str = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", db_user, db_pass, db_host, db_name)

}

func main() {
	load_config("config/bugspad.ini")

	db, err := sql.Open("mysql", conn_str)
	if err != nil {
		// handle error
		fmt.Print(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, status, summary FROM bugs")
	if err == nil {
		var id int64
		var status, summary string
		for rows.Next() {
			err = rows.Scan(&id, &status, &summary)
			//fmt.Println(id, status, summary)
			m := make(Bug)
			m["id"] = id
			m["status"] = status
			m["summary"] = summary
			data, _ := json.Marshal(m)
			sdata := string(data)
			sid := strconv.FormatInt(id, 10)
			redis_hset("bugs", sid, sdata)
			update_redis_bug_status(sid, status)
		}
		fmt.Println("All bug indexes loaded in Redis.")
	} else {
		fmt.Println("err in loading data")
		fmt.Println(err.Error())
	}
	defer rows.Close()

}
