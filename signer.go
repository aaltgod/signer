package main

import (
	"fmt"
	"sync"
)

var combinedResult string

func SingleHash(dataCh chan string) {

	data := <-dataCh
	mu := &sync.Mutex{}
	go func(data string, mu *sync.Mutex){
		mu.Lock()
		md5Data := DataSignerMd5(data)
		crc32Md5Data := DataSignerCrc32(md5Data)
		crc32Data := DataSignerCrc32(data)
		dataCh <- crc32Md5Data + "~" + crc32Data
		mu.Unlock()
	}(data, mu)
}

func MultiHash(dataCh chan string){

	singleHashedData := <-dataCh
	dataCh <- DataSignerCrc32(singleHashedData)
}

func CombineResults(dataCh chan string){

	finallyHashedData := <-dataCh
	combinedResult += finallyHashedData + "_"
}

func ExecutePipeline(data string) {

	dataCh := make(chan string, 1)
	dataCh <- data

	SingleHash(dataCh)
	MultiHash(dataCh)
	CombineResults(dataCh)

}

func main() {

	data := []string{"0", "1", "2", "3", "4"}
	for _, d := range data {
		ExecutePipeline(d)
	}
	fmt.Println(combinedResult)
}
