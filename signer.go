package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)


var intermediateResult []string

func executeMd5(data string, mu *sync.Mutex, out chan string) {
	mu.Lock()
	md5Data := DataSignerMd5(data)
	mu.Unlock()
	out <-md5Data
}

func executeCrc32(data string, out chan string) {

	out <-DataSignerCrc32(data)
}

func SingleHash(in, output chan interface{}) {

	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	md5Out := make(chan string)
	crc32Out := make(chan string)
	md5Crc32Out := make(chan string)

	for d := range in {
		wg.Add(1)
		data := fmt.Sprintf("%v", d)

		go func(data string, wg *sync.WaitGroup) {
			defer wg.Done()

			go executeMd5(data, mu, md5Out)
			go executeCrc32(<-md5Out, md5Crc32Out)
			go executeCrc32(data, crc32Out)

			output <- <-crc32Out + "~" + <-md5Crc32Out
		}(data, wg)

	}

	wg.Wait()

}

func MultiHash(in, output chan interface{}){

	wg := &sync.WaitGroup{}

	for d := range in {
		singleHashedData := fmt.Sprintf("%v", d)

		wg.Add(1)

		go func(data string, wg *sync.WaitGroup) {
			defer wg.Done()

			out := make([]chan string, 6)
			for th := range out {
				out[th] = make(chan string)
				go executeCrc32(strconv.Itoa(th) + singleHashedData, out[th])

			}

			result := ""
			for th := range out {
				result += <-out[th]
			}

			output <- result
		}(singleHashedData, wg)
	}

	wg.Wait()
}

func CombineResults(in, output chan interface{}){

	for d := range in {
		finallyHashedData := fmt.Sprintf("%v", d)
		intermediateResult = append(intermediateResult, finallyHashedData)
	}

	sort.Strings(intermediateResult)
	combinedResult := strings.Join(intermediateResult, "_")
	output <- combinedResult
}

func ExecutePipeline(jobs ...job) {

	wg := &sync.WaitGroup{}
	in := make(chan interface{}, 7)

	for _, j := range jobs {
		wg.Add(1)
		out := make(chan interface{}, 7)

		go func(job job, in, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			defer close(out)

			job(in, out)
		}(j, in, out, wg)

		in = out
	}

	wg.Wait()
}

