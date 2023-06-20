package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ApiCall struct {
	url     string
	ch chan string
}

func main() {
	var apiCalls []ApiCall
	ch := make(chan string)
	apiCalls = append(apiCalls, ApiCall{"https://cdn.apicep.com/file/apicep/01001-001.json", ch})
	apiCalls = append(apiCalls, ApiCall{"http://viacep.com.br/ws/01001001/json", ch})

	for _, apiCall := range apiCalls {
		go CallApi(apiCall)
	}

	select {
	case res := <- ch:
		fmt.Println(res)
		os.Exit(0)
	case <- time.After(time.Second * 1):
		fmt.Println("Application timed out!")
		os.Exit(1)
	}
}

func CallApi(a ApiCall) error {
	funcStart := time.Now()
	req, err := http.Get(a.url)
	if err != nil { return err }
	defer req.Body.Close()
	
	body, err := io.ReadAll(req.Body)
	if err != nil {  return err }
	
	funcEnd := time.Since(funcStart)
	fmt.Printf("Called api %s in %s\n", a.url, funcEnd)
	a.ch <- string(body)
	return nil
}