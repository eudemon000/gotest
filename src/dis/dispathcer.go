package dispathcer

import (
	queen "gotest/src/msgqueen"
	dataPip "gotest/src/pipeline"
	"fmt"
	"container/list"
	"time"
	"sync"
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

type myQueue struct {
	l 	*list.List
	lock 	*sync.Mutex
	mux 	*sync.Mutex
	isLock	bool
}

var mq *myQueue


func NewDispathcer(init *InitData) {
	//manage = queen.NewmsgQueenManager(10, 10, init.Handler)
	mq = new(myQueue)
	mq.l = list.New()
	mq.lock = new(sync.Mutex)
	mq.mux = new(sync.Mutex)
	continueUrl = init.ConUrl
	handler = init.H
	go queryUrl()
	time.Sleep(time.Second)
}

//待爬取的URL放入表里
func (init *InitData)Push(listData list.List) {
	//首先检查该URL是否爬取过
	mq.mux.Lock()
	defer mq.mux.Unlock()

	if mq.isLock {
		mq.lock.Unlock()
		mq.isLock = false
	}

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
	//<- mainCh
	//manage.PushData(str)

}

//从待爬取列表读取URL
func queryUrl() {
	for {
		fmt.Println("loop")
		//首先判断队列里是否有
		if mq.l.Len() > 0 {
			handler(*mq.l)
		} else {	//如果没有数据，则解除队列的锁，让其他函数可以向队列写入数据
			mq.lock.Lock()
			mq.isLock = true
			result := continueUrl.PullContinue()
			mq.l.PushBackList(&result)
		}




		/*if result.Len() > 0 {
			handler(result)
			//manage.PushData(result)
		}*/
	}
}

