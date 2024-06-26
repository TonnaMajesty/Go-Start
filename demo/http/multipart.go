package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	// Create a buffer to write the form data
	var buffer bytes.Buffer

	// Create a multipart writer with a random boundary
	writer := multipart.NewWriter(&buffer)

	// Add a form field with key "name" and value "Alice"
	if err := writer.WriteField("name", "Alice"); err != nil {
		log.Fatal(err)
	}

	// Add a form file with key "photo" and the content of "photo.jpg"
	file, err := os.Open("/Users/tonnamajesty/www/go/tzh/go-start/output.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", "photo.jpg")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(part, file); err != nil {
		log.Fatal(err)
	}

	// Close the writer to finish the form data
	if err := writer.Close(); err != nil {
		log.Fatal(err)
	}

	// Create a new HTTP request with the buffer as the body
	request, err := http.NewRequest("POST", "http://localhost:4600/api/charge-station/debug/multipart", &buffer)
	if err != nil {
		log.Fatal(err)
	}

	// Set the Content-Type header to multipart/form-data with the boundary
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a new HTTP client and send the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Print the status code and the response body
	fmt.Println(response.StatusCode)
	fmt.Println(response.Body)
}
