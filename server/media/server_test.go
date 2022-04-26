package media_test

import (
	"os"

	//"io/ioutil"
	"net/http"
	"testing"
)

//type Request struct {
//	Data []byte `json:"data"`
//}
//
//func TestRequestMedia(t *testing.T) {
//	data, _ := ioutil.ReadFile("../../public/test.jpg")
//
//	req := Request{
//		Data: data,
//	}
//	t.Log(len(data))
//	encoded, err := json.Marshal(req)
//	if err != nil {
//		t.Log("failed http post : ", err)
//	}
//
//	buf := bytes.NewBuffer(encoded)
//	resp, err := http.Post("http://localhost:10424/request/media", "application/json", buf)
//	if err != nil {
//		t.Log("failed http post : ", err)
//	}
//
//	t.Log(resp)
//
//}

func TestPostFile(t *testing.T) {

	client := &http.Client{}

	data, err := os.Open("../../public/IU.mp4")
	if err != nil {
		t.Log("err :", err)
	}

	req ,err := http.NewRequest("POST","http://localhost:10424/request/file",data)
	if err != nil {
		t.Log("err :", err)
	}

	resp, err := client.Do(req)
	if err != nil{
		t.Log("err : ",err)
	}

	t.Log("resp : ",resp)
}
