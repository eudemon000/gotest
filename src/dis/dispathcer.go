package dispathcer

import (
	queen "gotest/src/msgqueen"
	dataPip "gotest/src/pipeline"
	"fmt"
)


type InitData struct {
	Handler		queen.MessageHandler
	ConUrl		ContinueUrl
}

//待爬取URL
type ContinueUrl interface {
	//检查是否存在
	CheckData(data interface{}) bool
	//添加到待爬取文件
	PushContinue(data interface{}) bool
	//读取待爬取URL
	PullContinue() interface{}
}


var manage *queen.MsgQueenManager

var continueUrl ContinueUrl

func NewDispathcer(init *InitData) {
	manage = queen.NewmsgQueenManager(10, 10, init.Handler)
	continueUrl = init.ConUrl
	fmt.Println("aaa")
}

//待爬取的URL放入表里
func (init *InitData)Push(str string) {
	//首先检查该URL是否爬取过
	isExist := dataPip.QueryURL(str)
	if isExist {
		return
	}
	//放入待爬取的表
	continueUrl.PushContinue(str)

	manage.PushData(str)
}

//从待爬取列表读取URL
func queryUrl() {
	results := continueUrl.PullContinue()
	fmt.Println(results)
}

