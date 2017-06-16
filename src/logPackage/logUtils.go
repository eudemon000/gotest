package log

import (
	"log"
	"os"
	"time"
	"fmt"
	"strconv"
)

const (
	Info		string = "info"
	Error		string = "error"
	Warning		string = "warning"
)

func init() {
	os.MkdirAll("log", os.ModePerm)
}


func Msg(msgType string, msg ...interface{}) {
	if !(msgType == Info || msgType == Error || msgType == Warning) {
		return
	}
	fName := fileName(msgType)
	f, err := os.OpenFile(fName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0777)
	if err != nil {
		fmt.Println(err)
		/*os.MkdirAll("log", os.ModePerm)
		f, _ = os.Create(fName)*/
	}

	defer f.Close()

	switch msgType {
	case Info:
		logger := log.New(f, "[info]", log.Ldate | log.Ltime | log.Llongfile)
		logger.Println(msg)
	case Error:
		logger := log.New(f, "[error]", log.Ldate | log.Ltime | log.Llongfile)
		logger.Println(msg)
	case Warning:
		logger := log.New(f, "[warring]", log.Ldate | log.Ltime | log.Llongfile)
		logger.Println(msg)
	default:
		fmt.Println("default")
	}

}

func fileName(t string) string {
	y, m, d := time.Now().Date()
	return "log/" + strconv.Itoa(y) + strconv.Itoa(int(m)) + strconv.Itoa(d) + "_" + t + ".log"
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
		Msg(Error, err)
	}
}





