package search

import (
	"context"
	"log"
	"testing"
)

func TestAll_user(t *testing.T) {
	ch := All(context.Background(), "Need", []string{"test.txt"})
	results, ok := <-ch
	if !ok {
		t.Errorf("error: %v", ok)
	}
	log.Println("result: ", results)
}

func TestAny_user(t *testing.T) {
	results := Any(context.Background(), "abs1", []string{"test.txt", "test.txt"})
	result, err := <-results
	if !err {
		log.Println("error: ", err)
	}
	log.Println("result.Phrase: ", result.Phrase)
	log.Println("result.Line: ", result.Line)
	log.Println("result.LineNum: ", result.LineNum)
	log.Println("result.ColNum: ", result.ColNum)
}