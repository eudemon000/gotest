package msgqueen

import (
	"container/list"
	"sync"
)


//用List实现的同步队列

//用来存放消息的队列
var items *list.List

type MsgHandler func(data interface{})

type MsgManager struct {
	lock		sync.Mutex
	mux		sync.Mutex
	isLock		bool
	callback	MsgHandler
}


func NewMsgManager(handler MsgHandler) *MsgManager {
	items = list.New()
	mQm := new(MsgManager)
	mQm.callback = handler
	go mQm.readMsgData()
	return mQm
}

func (mQm *MsgManager)readMsgData() {
	for {
		if items.Len() > 0 {
			var n *list.Element
			for e := items.Front(); e != nil; e = n {
				n = e.Next()
				items.Remove(e)
				mQm.callback(e.Value)
			}
		} else {
			mQm.mux.Lock()
			mQm.isLock = true
		}
	}
}

func (mQm *MsgManager) PushMsgData(data interface{}) {
	mQm.lock.Lock()
	defer mQm.lock.Unlock()

	if mQm.isLock {
		mQm.mux.Unlock()
		mQm.isLock = false
	}

	items.PushBack(data)
}
