package decorator

import (
	"log"
	"net/http"
	"testing"
)

func TestHttp(t *testing.T) {
	http.HandleFunc("/v1/hello", WithDebugLog(WithServerHeader(hello)))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func TestPipeline(t *testing.T) {
	http.HandleFunc("/v4/hello", Handler(hello,
		WithServerHeader, WithBasicAuth, WithDebugLog))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
