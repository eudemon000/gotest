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

type DataBasePip struct {

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

	/*if rows.Next() {
		err := rows.Scan(scanArgs...)
		msgLog.CheckErr(err)
		record := make(map[string]interface{})
		for i, col := range values {
			if col != nil {
				//record[columns[i]] = string(col.([]byte))
				record[columns[i]] = col
				//record[columns[i]] = "aaa"
			}
		}
		fmt.Println(string(record["url"]))
		isExist = true
	} else {
		isExist = false
	}*/
	if rows.Next() {
		return true
	} else {
		return false
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

//插入待爬取表
func InsertContinue(url string) (ok bool) {
	result, err := Db.Exec("insert into tbl_continue_url(url) values(?)", url)
	if err != nil {
		msgLog.Msg(msgLog.Error, err)
		return false
	}
	id, insertErr := result.LastInsertId()
	if insertErr != nil {
		msgLog.Msg(msgLog.Error, insertErr)
		return false
	}
	if id > 0 {
		return true
	}
	return false
}

//检查是否存在
func (pip *DataBasePip)CheckData(data interface{}) bool {
	rows, err := Db.Query("select * from tbl_tbl_continue_url as c where c.url like ?", data.(string))
	if err != nil {
		msgLog.Msg(msgLog.Error, err)
		return true
	}
	columns, err := rows.Columns()
	if err != nil {
		msgLog.Msg(msgLog.Error, err)
		return true
	}

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
		return true
	} else {
		return false
	}
	defer rows.Close()
	return false
}

//添加到待爬取文件
func (pip *DataBasePip)PushContinue(data interface{}) bool {
	result, err := Db.Exec("insert into tbl_tbl_continue_url(url) values()?", data.(string))
	if err != nil {
		msgLog.Msg(msgLog.Error, err)
		return false
	}
	lastId, lastErr := result.LastInsertId()
	if lastErr != nil {
		msgLog.Msg(msgLog.Error, lastErr)
		return false
	}
	if lastId > 0 {
		return true
	}
	return false
}
//读取待爬取URL
func (pip *DataBasePip)PullContinue() interface{} {
	rows, err := Db.Query("select * from tbl_tbl_continue_url limit 0, 20")
	defer rows.Close()
	if err != nil {
		msgLog.Msg(msgLog.Error, err)
		return nil
	}
	columns, err := rows.Columns()
	if err != nil {
		msgLog.Msg(msgLog.Error, err)
		return nil
	}
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs)
		if err != nil {
			msgLog.Msg(msgLog.Error, err)
			return nil
		}
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		return record
	}

	return nil
}


