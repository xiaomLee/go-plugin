package common

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func BenchmarkHttpAgent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		request := NewHttpAgent()
		request = request.SetHeader("X-Forwarded-For", "127.0.0.1")
		request = request.Post("https://www.baidu.com")

		data := map[string]string{
			"name":    "tom",
		}

		_, _, err := request.ContentType(TypeFormUrlencoded).SendForm(data).End()
		if err != nil {
			println(err.Error())
			return
		}
	}
}

func TestHttpAgent_Timeout(t *testing.T) {
	request := NewHttpAgent()
	request = request.SetHeader("X-Forwarded-For", "127.0.0.1").Timeout(time.Millisecond * 500)
	request = request.Post("https://www.baidu.com")

	data := map[string]string{
		"name":    "tom",
	}

	_, body, err := request.ContentType(TypeFormUrlencoded).SendForm(data).End()
	if err != nil {
		println(err.Error())
		return
	}
	println(string(body))
}

func TestHttpAgent(t *testing.T) {

	// 1. POST
	request := NewHttpAgent()
	request = request.SetHeader("SomeHead", "value")
	request = request.Post("http://127.0.0.1:80/user/login")

	data := map[string]string{
		"account":  "tom",
		"password": "123456",
	}

	_, body, err := request.ContentType(TypeFormUrlencoded).SendForm(data).End()
	if err != nil {
		println(err.Error())
		return
	}
	println(string(body))

	// GET
	_, body, err = request.Get("http://127.0.0.1:80/user/getUserInfo").End()
	if err != nil {
		println(err.Error())
		return
	}
	println(string(body))

	// UPLOAD
	file1, err := os.OpenFile("test1.txt", os.O_RDONLY, 0755)
	if err != nil {
		println(err.Error())
		return
	}
	defer file1.Close()
	file2, err := os.OpenFile("test2.txt", os.O_RDONLY, 0755)
	if err != nil {
		println(err.Error())
		return
	}
	defer file2.Close()
	data1, _ := ioutil.ReadAll(file1)
	data2, _ := ioutil.ReadAll(file2)
	request = request.SendFile(File{FileName: "test1.txt", FieldName: "a", Data: data1})
	request = request.SendFile(File{FileName: "test2.txt", FieldName: "b", Data: data2})
	_, body, err = request.Post("http://127.0.0.1:80/resource/upload").ContentType(TypeMultipartFormData).End()
	if err != nil {
		println(err.Error())
		return
	}
	println(string(body))
}
