package dispathcer

import (
	queen "gotest/src/msgqueen"
)


type InitData struct {
	Handler queen.MessageHandler
}

var manage *queen.MsgQueenManager

var cUrl string

func NewDispathcer(init *InitData) {
	manage = queen.NewmsgQueenManager(10, 10, init.Handler)
}

func (init *InitData)Push(str string) {
	manage.PushData(str)
}

