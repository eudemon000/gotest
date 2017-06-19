package dispathcer

import (
	queen "gotest/src/msgqueen"
	dataPip "gotest/src/pipeline"
	"fmt"
	"container/list"
)


type InitData struct {
	//Handler		queen.MessageHandler
	ConUrl		ContinueUrl
	H		MessageHandler
}

type MessageHandler func(data interface{})

//待爬取URL
type ContinueUrl interface {
	//检查是否存在
	CheckData(data interface{}) bool
	//添加到待爬取文件
	PushContinue(data interface{}) bool
	//读取待爬取URL
	PullContinue() list.List
}

var mainCh chan int = make(chan int)


var manage *queen.MsgQueenManager

var continueUrl ContinueUrl

var handler MessageHandler

func NewDispathcer(init *InitData) {
	//manage = queen.NewmsgQueenManager(10, 10, init.Handler)
	continueUrl = init.ConUrl
	handler = init.H
	go queryUrl()

}

//待爬取的URL放入表里
func (init *InitData)Push(str string) {
	<- mainCh
	//首先检查该URL是否爬取过
	isExist := dataPip.QueryURL(str)
	if isExist {
		return
	}

	//检查待爬取的表里是否有该URL
	exist := continueUrl.CheckData(str)
	if !exist {
		//放入待爬取的表
		continueUrl.PushContinue(str)
		//mainCh <- 1
	}

	//manage.PushData(str)

}

//从待爬取列表读取URL
func queryUrl() {
	for {
		//a := <- mainCh
		result := continueUrl.PullContinue()
		for e := result.Front(); e != nil; e = e.Next() {
			tempMap := e.Value.(map[string]string)
			fmt.Println("tempMap=", tempMap["url"])
			//manage.PushData(tempMap["url"])
			handler(tempMap["url"])
		}
		mainCh <- 1
		fmt.Println("queryUrl")

	}
}

