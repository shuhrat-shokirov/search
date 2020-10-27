package search

import (
	"context"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

type Result struct {
	Phrase  string
	Line    string
	LineNum int64
	ColNum  int64
}

func All(ctx context.Context, phrase string, files []string) <-chan []Result {
	resultChan := make(chan []Result)
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(ctx)
	for i := 0; i < len(files); i++ {
		wg.Add(1)
		go func(ctx context.Context, filename string, i int, resultChan chan<- []Result) {
			defer wg.Done()
			results := FindAll(phrase, filename)
			if len(results) > 0 {
				resultChan <- results
			}
		}(ctx, files[i], i, resultChan)
	}

	go func() {
		defer close(resultChan)
		wg.Wait()

	}()

	cancel()
	return resultChan
}

func FindAll(phrase, fileName string) (results []Result) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println("error while open file: ", err)
		return results
	}

	split := strings.Split(string(data), "\n")
	for i, line := range split {
		if strings.Contains(line, phrase) {
			result := Result{
				Phrase:  phrase,
				Line:    line,
				LineNum: int64(i + 1),
				ColNum:  int64(strings.Index(line, phrase)) + 1,
			}
			results = append(results, result)
		}
	}

	return results
}

func Any(ctx context.Context, phrase string, files []string) <-chan Result {
	resultChan := make(chan Result)
	wg := sync.WaitGroup{}
	result := Result{}

	for i := 0; i < len(files); i++ {
		data, err := ioutil.ReadFile(files[i])
		if err != nil {
			log.Println("error while open file: ", err)
		}

		if strings.Contains(string(data), phrase) {
			res := FindAny(phrase, string(data))
			if (Result{}) != res {
				result = res
				break
			}
		}
	}
	log.Println("find result: ", result)

	wg.Add(1)
	go func(ctx context.Context, ch chan<- Result) {
		defer wg.Done()
		if (Result{}) != result {
			ch <- result
		}
	}(ctx, resultChan)

	go func() {
		defer close(resultChan)
		wg.Wait()
	}()
	return resultChan
}

func FindAny(phrase, search string) (result Result) {
	for i, line := range strings.Split(search, "\n") {
		if strings.Contains(line, phrase) {
			return Result{
				Phrase:  phrase,
				Line:    line,
				LineNum: int64(i + 1),
				ColNum:  int64(strings.Index(line, phrase)) + 1,
			}
		}
	}
	return result
}