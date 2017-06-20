package dispathcer

import (
	queen "gotest/src/msgqueen"
	dataPip "gotest/src/pipeline"
	"fmt"
	"container/list"
	"time"
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

var manage *queen.MsgQueenManager

var continueUrl ContinueUrl

var handler MessageHandler

var ch chan int = make(chan int)


func NewDispathcer(init *InitData) {
	//manage = queen.NewmsgQueenManager(10, 10, init.Handler)
	continueUrl = init.ConUrl
	handler = init.H
	go queryUrl()
	time.Sleep(time.Second)
}

//待爬取的URL放入表里
func (init *InitData)Push(listData list.List) {
	//首先检查该URL是否爬取过
	var n *list.Element
	for e := listData.Front(); e != nil; e = n {
		n = e.Next()
		listData.Remove(e)

		var str string
		switch e.Value.(type) {
		case string:
			str = e.Value.(string)
		}

		isExist := dataPip.QueryURL(str)
		if isExist {
			//return
			//如果存在，跳出当前循环
			continue
		}

		//检查待爬取的表里是否有该URL
		exist := continueUrl.CheckData(str)
		if !exist {
			//放入待爬取的表
			continueUrl.PushContinue(str)
		}

	}
	r := <- ch
	fmt.Println("ch=", r)
	//<- mainCh
	//manage.PushData(str)

}

//从待爬取列表读取URL
func queryUrl() {
	for {
		fmt.Println("loop")
		//首先阻塞一下，等push被调用后，解除阻塞
		ch <- 1
		result := continueUrl.PullContinue()
		if result.Len() > 0 {
			go handler(result)
		}
	}
}

