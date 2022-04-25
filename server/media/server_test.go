package media_test

import (
	"bytes"
	"encoding/json"
	//"io/ioutil"
	"net/http"
	"testing"
)

type Request struct {
	//Data []byte `json:"data"`
	Data string `json:"data"`

}



func TestRequestMedia(t *testing.T) {
	//data, _ := ioutil.ReadFile("../../public/test.jpg")

	//mediaReq := media.MediaRequest{Data: data}

	//req := Request{
	//	Data: data,
	//}
	req := Request{
		Data: "good",
	}

	encoded , err := json.Marshal(req)
	if err != nil {
		t.Log("failed http post : ", err)
	}

	buf := bytes.NewBuffer(encoded)
	resp, err :=  http.Post("http://localhost:9094/request/media", "application/json", buf)
	if err != nil {
		t.Log("failed http post : ", err)
	}




	t.Log(resp)


}
