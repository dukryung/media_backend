package media_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
	"testing"
)

type Request struct {
	Data []byte `json:"data"`
}

func TestRequestMedia(t *testing.T) {
	data, _ := ioutil.ReadFile("../../public/test.jpg")

	req := Request{
		Data: data,
	}
	t.Log(len(data))
	encoded, err := json.Marshal(req)
	if err != nil {
		t.Log("failed http post : ", err)
	}

	buf := bytes.NewBuffer(encoded)
	resp, err := http.Post("http://localhost:10424/request/media", "application/json", buf)
	if err != nil {
		t.Log("failed http post : ", err)
	}

	t.Log(resp)

}

func TestPostFile(t *testing.T) {
	var wg sync.WaitGroup


	for i := 0 ; i < 10000 ; i++ {
		wg.Add(1)
		go func () {
			defer wg.Done()
			fileName := "test.jpg"
			filePath := path.Join("../../public/", fileName)

			file, err := os.Open(filePath)
			if err != nil {
				t.Log("err : ", err)
				return
			}
			defer file.Close()

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("attachment", filepath.Base(file.Name()))
			if err != nil {
				t.Log("err : ", err)
				return
			}

			io.Copy(part, file)

			writer.Close()

			r, _ := http.NewRequest("POST", "http://localhost:10424/request/file", body)
			r.Header.Add("Content-Type", writer.FormDataContentType())
			client := &http.Client{}
			client.Do(r)

		}()
	}
	wg.Wait()

}
