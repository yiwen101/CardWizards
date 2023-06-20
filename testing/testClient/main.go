package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func realURL(u string) string {
	encode := url.Values{}
	encode.Add("query", "query")
	encode.Add("q", "q1")
	encode.Add("q", "q2")
	encode.Add("vd", "1")

	return fmt.Sprintf("%s?%s", u, encode.Encode())
}

func main() {
	client, err := client.NewClient()
	if err != nil {
		return
	}
	req := &protocol.Request{}
	res := &protocol.Response{}

	// Construct your own request URL using the "net/url" library
	req.Header.SetMethod(consts.MethodPost)
	req.SetRequestURI(realURL("http://127.0.0.1:8080/arithmatic/Add"))
	err = client.Do(context.Background(), req, res)
	if err != nil {
		return
	}

	// Send "Json" request
	req.Reset()
	req.Header.SetMethod(consts.MethodPost)
	req.Header.SetContentTypeBytes([]byte("application/json"))
	req.SetRequestURI("http://127.0.0.1:8080/arithmatic/Add")
	data := struct {
		firstArguement  int
		secondArguement int
	}{
		1, 2,
	}
	jsonByte, _ := json.Marshal(data)
	req.SetBody(jsonByte)
	err = client.Do(context.Background(), req, res)
	if err != nil {

		return
	}
	fmt.Println(res)
	fmt.Printf("%v", string(res.Body()))
}
