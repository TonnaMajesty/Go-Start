package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

func main() {
	// Create a Resty Client
	client := resty.New()

	file, err := os.Open("/Users/tonnamajesty/www/go/tzh/go-start/output.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Multipart/Form-Data request with byte[]
	//data := []byte("Hello World")
	resp, err := client.R().
		SetMultipartField("file", "hello.txt", "multipart/form-data", file).
		SetFormData(map[string]string{
			"key":      "value",
			"username": "testuser",
		}).
		Post("http://localhost:4600/api/charge-station/debug/multipart")

	// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("Error      :", err)
	fmt.Println("Status Code:", resp.StatusCode())
	fmt.Println("Status     :", resp.Status())
	fmt.Println("Body       :\n", resp)
}
