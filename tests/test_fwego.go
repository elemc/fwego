package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	resp, err := http.Get("http://localhost:4000")
	panicOnError(err)
	//fmt.Printf("[%v]\n", resp)

	if resp.StatusCode == 200 {
		//fmt.Printf("Status code is 200\n")
	} else {
		panic(fmt.Errorf("Status code is not 200, as expected. Status code: %d", resp.StatusCode))
	}

	var allData []byte
	data := make([]byte, 1024)
	count, err := resp.Body.Read(data)
	if err != io.EOF {
		panicOnError(err)
	}
	//fmt.Printf("Count is %d\n", count)
	allData = append(allData, data...)
	if err == io.EOF {
		count = 0
	}

	for count != 0 {
		data := make([]byte, 1024)
		count, err = resp.Body.Read(data)
		if err != io.EOF {
			panicOnError(err)
		}
		//fmt.Printf("Count is %d\n", count)
		allData = append(allData, data...)

		if err == io.EOF {
			count = 0
		}
	}

	body := string(allData)
	testStr := "<td><a href=\"/fwego\">fwego</a></td>"
	if !strings.Contains(body, testStr) {
		panic(fmt.Errorf("In body not found fwego record!\n%s", body))
	}
}
