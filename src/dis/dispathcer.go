package dispathcer

import (
	queen "gotest/src/msgqueen"
)

//格式化数据接口
type DataFormat interface {
	Format(str string) (result string, ok bool)
}

type InitData struct {
	DataFormat DataFormat
	Handler queen.MessageHandler
}

var df *DataFormat

var manage *queen.MsgQueenManager

func NewDispathcer(init *InitData) {
	manage = queen.NewmsgQueenManager(10, 10, init.Handler)
}

func (init *InitData)Push(str string) {
	result, ok := init.DataFormat.Format(str)
	if !ok {
		panic("格式不合法")
		return
	}
	manage.PushData(result)
}
