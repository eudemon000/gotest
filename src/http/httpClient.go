package http

import (
	_"fmt"
	_"net/http"
	"net/http"
	"encoding/json"
	"strings"
	"fmt"
	"io/ioutil"
	msgLog "gotest/src/logPackage"
)

type Body struct {
	Url		string
	Md5 		string
	Tag 		string
}

func AddElsearch(data Body) {

	b, err := json.Marshal(data)
	if err != nil {
		//fmt.Println(err)
		msgLog.Msg(msgLog.Error, err)
	}
	sendData := string(b)
	fmt.Println(sendData)
	resp, err := http.Post("http://192.168.0.215:9200/sanpotel_search/cancer/" + data.Md5, "application/x-www-form-urlencoded", strings.NewReader(sendData))
	if err != nil {
		//fmt.Println(err)
		msgLog.Msg(msgLog.Error, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msgLog.Msg(msgLog.Error, err)
	}
	//fmt.Println(string(body))
	msgLog.Msg(msgLog.Info, body)
}


