package main

import "fmt"

func SingleHash(data string) chan string{
	if len(data) > MaxInputDataLen {
		fmt.Errorf("len of DATA > 100")
	}

	ch := make(chan string, 1)

	go func(ch chan<- string, data string) {
		md5Data := DataSignerMd5(data)
		crc32Md5Data := DataSignerCrc32(md5Data)
		crc32Data := DataSignerCrc32(data)
		ch <- crc32Md5Data + "~" + crc32Data
	}(ch, data)

	return ch
}

func Multihash(data string) string{

	ch := make(chan string, 1)
	ch <- DataSignerCrc32(data)

	return <-ch

}

func CombineResults(data chan string) string{
	result :
}

func ExecutePipeline(data string) {

	inputCh := make(chan string, 1)
	hashCh := make(chan string, 1)
	combineCh := make(chan string,1)

	inputCh <-data
	hashCh <-SingleHash(<-inputCh)
	combineCh <-Multihash(<-hashCh)


}

func main() {

	data := "sss"

	ExecutePipeline(data)
}
