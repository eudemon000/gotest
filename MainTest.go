package main

import (
	"fmt"
	"runtime"
	_"time"
	queen "gotest/src/msgqueen"
	"github.com/PuerkitoBio/goquery"
	sqlConn "gotest/src/pipeline"
	hex "gotest/src/utils"
	"strings"
	"sync"
	/*"bytes"*/
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"github.com/henrylee2cn/pholcus/common/mahonia"
	"bytes"
	myHttp "gotest/src/http"
)


var urls = make([]string, 0)

var currentUrl = ""

var manage *queen.MsgQueenManager

var handler = queen.MessageHandler(func(data interface{}) {
	fmt.Println("这是消息：", data)
	start(data.(string))
})

var mutexLock sync.Mutex

func main() {
	accept := runtime.NumCPU()
	runtime.GOMAXPROCS(accept)

	manage = queen.NewmsgQueenManager(10, 10, handler)
	manage.PushData("http://www.99.com.cn")
	//manage.PushData("http://www.qq.com")
	var c chan int
	c <- 1
	/*urls = append(urls, "http://www.99.com.cn")
	a := make(chan int)
	for _, item := range urls {
		go start(a, item)
	}*/

}

func start(url string) {
	if currentUrl == url {
		return
	}
	/**
	1、首先爬去全文取出所有待爬的URL放到数组里
	2、获取页面编码格式，根据编码格式格式化内容
	3、
	 */

	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	//首先获取页面的编码格式，编码格式一般在head的meta中
	var webCharset string
	headTag := doc.Find("head")
	metaTag := headTag.Find("meta")
	webCharset = checkCharset(metaTag)

	//获取页面的关键词，根据编码进行编码转换，并保存到数据库
	keyword := checkTag(metaTag, webCharset)
	fmt.Println("转码后===》", keyword)
	urls := new(sqlConn.Urls)
	urls.Url = url
	urls.Md5, _ = hex.Md5(url)
	urls.Content = keyword
	urls.Layer = 1
	urls.Is_crawl = "YES"
	sqlErr := sqlConn.Insert(urls)
	if sqlErr != nil {
		fmt.Println("数据插入失败===》", sqlErr)
	} else {
		if strings.Contains(urls.Content, "肿瘤") {
			b := myHttp.Body{urls.Url, urls.Md5, urls.Content}
			fmt.Println(b)
			sendSearch(b)
		}
	}


	//获取接下来需要爬取的URL，放放入队列中
	bodyTag := doc.Find("body")
	resultUrl := make([]string, 50)
	bodyTag.Each(func(i int, bodySelect *goquery.Selection) {
		resultUrl = findUrls(bodySelect)
	})
	for _, item := range resultUrl {
		//fmt.Println(index, item)
		if len(item) > 0 {
			manage.PushData(item)
		}
	}



	/*



	//先获取所有需要爬取的URL放入到数组中，待爬取







	document := doc.Find("html")
	var webCharset string

	urls := new(sqlConn.Urls)

	*//*defer func() {
		if rErr := recover(); rErr != nil {
			var b myHttp.Body
			b = myHttp.Body{urls.Url, urls.Md5, urls.Content}
			fmt.Println(b)
			//sendSearch(b)
		}
	}()*//*

	document.Each(func(i int, s *goquery.Selection) {
		aTag := s.Find("a")
		urlArr = findUrls(aTag)

		metaTag := s.Find("meta")
		if len(webCharset) <= 0 {
			webCharset = checkCharset(metaTag)
		}
		fmt.Println(webCharset)

		var content string
		content = checkTag(metaTag, webCharset)

		urls.Url = url
		urls.Md5, _ = hex.Md5(url)
		urls.Is_crawl = "NO"
		urls.Layer = 1
		urls.Content = content
		//panic(urls)
	})

	b := myHttp.Body{urls.Url, urls.Md5, urls.Content}
	fmt.Println(b)
	sendSearch(b)

	err1 := sqlConn.Insert(urls)
	if err1 != nil {
		fmt.Println(err1)
	}
	for index, item := range urlArr {
		fmt.Println("index=", index, "item=", item)
		if len(item) > 0 {
			manage.PushData(item)
		}
	}*/
	//fmt.Println(len(urls))
}

func formatStr(str, setCharset string) string {
	setCharset = strings.ToLower(setCharset)
	if strings.Contains(setCharset, "gbk") {
		de := mahonia.NewDecoder("gbk")
		result := de.ConvertString(str)
		//result := Decode(str, "gbk")
		return result
	} else if strings.Contains(setCharset, "gb2312") {
		fmt.Println("这是gb2312页面")
		de := mahonia.NewDecoder("gb2312")
		result := de.ConvertString(str)
		//result := Decode(str, "gb2312")
		return result

	}
	return str

}

//检查页面的编码类型
func checkCharset(sele *goquery.Selection) (webCharset string) {
	//var webCharset string
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	sele.Each(func(i int, m *goquery.Selection) {
		var wOK bool
		webCharset, wOK = m.Attr("charset")
		if !wOK {
			httpEquiv, hOk := m.Attr("http-equiv")
			if hOk && httpEquiv == "Content-Type" {
				content, _ := m.Attr("content")
				webCharset = content
				//fmt.Println(webCharset)
				panic(webCharset)
			}
		} else {
			panic(webCharset)
		}


	})
	//fmt.Println("charset===>", webCharset)
	return
}

//检查meta信息
func checkTag(sele *goquery.Selection, webCharset string) string {
	var tag string
	sele.Each(func(i int, m *goquery.Selection) {

		result, ok := m.Attr("name")
		if ok {
			if result == "keywords" || result =="Keywords" || result == "description" || result == "Description" {
				content, _ := m.Attr("content")
				//fmt.Println(content)
				tag = formatStr(content, webCharset)
				//tag = content
				fmt.Println("content===>", content)
				/*if content != "" {
					content = formatStr(content, webCharset)
					err := sqlConn.InsertTag(content, url)
					if err != nil {
						fmt.Println(err)
					}
					//fmt.Println(content, err)
				}*/
				//return tag
			}
		}

	})

	return tag
}

func formatString(str, setCharset string) string {
	de := mahonia.NewDecoder("gb2312")
	result := de.ConvertString(str)
	return result
	/*a := []byte(str)
	r := bytes.NewReader(a)
	d, _ := charset.NewReader(r, "gb2312")

	fmt.Println("de===>", aaa)

	result, _ := ioutil.ReadAll(d)
	fmt.Println("result===>", string(result))
	return "aaa"*/
}

func Decode(str, setCharset string) string {
	b := []byte(str)
	i := bytes.NewReader(b)
	r, _:= charset.NewReader(i, setCharset)
	d, _ := ioutil.ReadAll(r)
	fmt.Println("d===>", string(d))
	return string(d)

}

//获取页面上所有的URL
func findUrls(bodySelect *goquery.Selection) []string {
	var array []string = make([]string, 50)
	aTag := bodySelect.Find("a")
	aTag.Each(func(index int, node *goquery.Selection) {
		tempUrl, ok := node.Attr("href")
		if ok {
			//此处暂时判断链接以http开头，未来需要判断相对地址，暂时不做处理
			if strings.Index(tempUrl, "http") != -1 {
				array = append(array, tempUrl)
			}
			//manage.PushData(tempUrl)
		}
	})
	return array
}

func sendSearch(data myHttp.Body) {
	myHttp.AddElsearch(data)
}

