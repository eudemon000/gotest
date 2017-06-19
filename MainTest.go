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
	disPack "gotest/src/dis"
	"regexp"
	msgLog "gotest/src/logPackage"
)

type DataFormat interface {
	//格式化字符串
	Format(str string) (result string, ok bool)
}

//用来实现DataFormat接口
type InitImp struct {

}


var initDataFormat *InitImp

var urls = make([]string, 0)

var currentUrl = ""

var manage *queen.MsgQueenManager

var initData *disPack.InitData

/*var handler = queen.MessageHandler(func(data interface{}) {
	//fmt.Println("这是消息：", data)
	start(data.(string))
})*/

var handler = disPack.MessageHandler(func(data interface{}) {
	//fmt.Println("这是消息：", data)
	start(data.(string))
})

var mutexLock sync.Mutex

var cUrl string

func main() {
	accept := runtime.NumCPU()
	runtime.GOMAXPROCS(accept)
	initDataFormat = new(InitImp)
	initData = new(disPack.InitData)
	//initData.Handler = handler
	initData.H = handler
	initData.ConUrl = new(sqlConn.DataBasePip)
	disPack.NewDispathcer(initData)
	cUrl = "http://www.99.com.cn"
	initData.Push("http://www.99.com.cn")
	/*manage = queen.NewmsgQueenManager(10, 10, handler)
	manage.PushData("http://www.99.com.cn")*/
	var c chan int
	c <- 1
	/*urls = append(urls, "http://www.99.com.cn")
	a := make(chan int)
	for _, item := range urls {
		go start(a, item)
	}*/

}

//实现DataFormat的方法
func (initImp *InitImp)Format(str string) (result string, ok bool) {
	//fmt.Println("接口方法调用", str)
	//首先判断是不是是不是javascript，#或*开头的,如果是代表不是合法URL
	ok, err := regexp.MatchString("^javascript|^#|^\\*", str)
	if err !=nil {
		msgLog.Msg(msgLog.Error, err)
		return "", false
	}
	if ok {
		return "", false
	}

	//判断是不是http开头的，http和https均可判断
	ok, err = regexp.MatchString("^http", str)
	if err != nil {
		msgLog.Msg(msgLog.Error, err)
		return "", false
	}
	if ok {
		return str, true
	}

	//还要一种是相对路径，分两种情况，1、"/"开头；2、非"/"开头
	ok, err = regexp.MatchString("^/{1}[a-zA-Z0-9]{1,}?", str)
	if ok {
		//需要找路径根
		strs := strings.Split(cUrl, "/")
		re := strs[0] + "//" + strs[2] + str
		return re, false
	}

	ok, err = regexp.MatchString("[a-zA-Z0-9]{1,}?", str)
	if err != nil {
		msgLog.Msg(msgLog.Error, err)
		return "", false
	}
	if ok {
		postion := strings.LastIndex(cUrl, "/")
		postion += 1
		a := cUrl[0:postion]
		re := a + str
		return re, true
	}
	return "", false
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
		msgLog.Msg(msgLog.Error, err)
		return
	}
	cUrl = doc.Url.String()
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
		msgLog.Msg(msgLog.Error, err)
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
			//manage.PushData(item)
			r, ok := initDataFormat.Format(item)
			if ok {
				//msgLog.Msg(msgLog.Info, r)
				initData.Push(r)
			}
		}
	}

}

func formatStr(str, setCharset string) string {
	setCharset = strings.ToLower(setCharset)
	if strings.Contains(setCharset, "gbk") {
		de := mahonia.NewDecoder("gbk")
		result := de.ConvertString(str)
		//result := Decode(str, "gbk")
		return result
	} else if strings.Contains(setCharset, "gb2312") {
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
			//fmt.Println(err)
			msgLog.Msg(msgLog.Error, err)
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
			/*if strings.Index(tempUrl, "http") != -1 {
				array = append(array, tempUrl)
			}*/
			result, ok := initDataFormat.Format(tempUrl)
			if ok {
				array = append(array, result)
			}
			//manage.PushData(tempUrl)
		}
	})
	return array
}

func sendSearch(data myHttp.Body) {
	myHttp.AddElsearch(data)
}



