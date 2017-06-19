package msgqueen

import (
	"os"
	"sync"
	_"time"
)

type msgQueen struct {
	index   	int
	data    	interface{}
	pre     	*msgQueen
	next    	*msgQueen
	hasData 	bool
}


type MessageHandler func(data interface{})

type MsgQueenManager struct {
	cur            *msgQueen //当前消息即将插入位置
	size           int
	readThreadSize int
	lock           sync.Mutex
	mux            sync.Mutex
	isLock         bool
	CallBack       MessageHandler
}

var queen *msgQueen

func productMsgQueen(size int) *msgQueen {

	data := make([]msgQueen, size)
	tmpMsg := &(data[0])
	tmpMsg.index = 0
	tmpMsg.next = &(data[1])
	for i := 1; i < size-1; i++ {
		tmpMsg = &(data[i])
		tmpMsg.index = i
		tmpMsg.next = &(data[i+1])
		tmpMsg.pre = &(data[i-1])

	}
	tmpMsg = &(data[size-1])
	tmpMsg.index = size - 1
	tmpMsg.pre = &(data[size-2])
	// first item
	tmpMsg.next = &(data[0])
	tmpMsg.next.pre = tmpMsg
	return tmpMsg.next
}

func NewmsgQueenManager(size, readThreadSize int, handler MessageHandler) *MsgQueenManager {
	if size <= 0 && size%2 != 0 && handler == nil {
		os.Exit(1)
		return nil
	}
	mQm := &MsgQueenManager{}
	mQm.cur = productMsgQueen(size << 1)
	mQm.readThreadSize = readThreadSize
	mQm.size = size << 1
	head := mQm.cur
	queen = head
	mQm.CallBack = handler
	go mQm.readMsg(head, mQm.readThreadSize)

	return mQm
}

func (mQm *MsgQueenManager) PushData(data interface{}) {
	mQm.lock.Lock()
	defer mQm.lock.Unlock()

	/*isExist := continueUrl.CheckData(data)
	if !isExist {
		return
	}*/



	if mQm.isLock {
		mQm.mux.Unlock()
		mQm.isLock = false
	}
	if mQm.cur.hasData {
		var mq *msgQueen = mQm.cur
		// 准备消息空间
		mQm.size = mQm.size << 1
		queen := productMsgQueen(mQm.size)
		head := mq
		tail := mq.next
		newTail := queen.pre
		head.next = queen
		queen.pre = head
		tail.pre = newTail
		newTail.next = tail
		mQm.cur = mQm.cur.next
	}

	mQm.cur.data = data
	mQm.cur.hasData = true
	mQm.cur = mQm.cur.next

	return
}

func (mQm *MsgQueenManager) readMsg(mq *msgQueen, size int) {
	var ds int
	ds = size
	for {
		//time.Sleep(0)
		if mq.hasData {
			ds = size
			mq.hasData = false
			mQm.CallBack(mq.data)
		} else {
			mQm.mux.Lock()
			mQm.isLock = true
			ds = 0
		}
		for i := 0; i < ds; i++ {
			mq = mq.next
		}
	}
}
