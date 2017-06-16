package pipeline

import (
	"fmt"
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	msgLog "gotest/src/logPackage"
)

//var db = &sql.DB{}
var Db *sql.DB

type Urls struct {
	Url		string
	Md5		string
	Is_crawl	string
	Layer		int
	Content		string
}

func init() {
	initDB()
}

func initDB() {
	var err error
	Db, err = sql.Open("mysql", "root:Sanpotel9958!@tcp(192.168.0.215:3306)/sanpotel_search?charset=utf8")
	//db, err = sql.Open("mysql", "debian-sys-maint:98lq22Jdd0SosgmM@tcp(localhost:3306)/sanpotel_search?charset=utf8")
	Db.SetMaxOpenConns(20)
	msgLog.CheckErr(err)
}

func QueryURL(url string) (isExist bool) {
	rows, queryErr := Db.Query("SELECT * FROM tbl_urls as u where u.url like ?", url)
	msgLog.CheckErr(queryErr)
	//strs, _ := rows.Columns()
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	if rows.Next() {
		err := rows.Scan(scanArgs...)
		msgLog.CheckErr(err)
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		fmt.Println(record)
		isExist = true
	} else {
		isExist = false
	}
	defer rows.Close()
	return
}

func Insert(u *Urls) (e error) {
	//initDB()
	//db.SetMaxOpenConns(50)
	//result, err := db.Exec("insert into tbl_urls(url, md5, is_crawl, layer) values(?, ?, ?, ?)", u.Url, u.Md5, u.Is_crawl, u.Layer)
	var result, err = Db.Exec("insert into tbl_urls(url, md5, is_crawl, layer, tag) values(?, ?, ?, ?, ?)", u.Url, u.Md5, u.Is_crawl, u.Layer, u.Content)
	if err != nil {
		//fmt.Println(err)
		msgLog.CheckErr(err)
	}
	e = err
	fmt.Println(result)
	return
}

func InsertTag(str, url string) (e error) {
	//initDB()
	//db.SetMaxOpenConns(50)
	//result, err := db.Exec("insert into tbl_urls(url, md5, is_crawl, layer) values(?, ?, ?, ?)", u.Url, u.Md5, u.Is_crawl, u.Layer)
	var result, err = Db.Exec("insert into tbl_tag(tag, url) values(?, ?)", str, url)
	if err != nil {
		msgLog.Msg(msgLog.Error, err)
	}
	e = err
	fmt.Println(result)
	return
}


