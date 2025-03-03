package response

import (
	"fmt"
	"github.com/codecrafters-io/http-server-starter-go/app/http/response/compressions"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Response struct {
	body       string
	status     int
	statusText string
	headers    map[string]string
}

func New() *Response {
	return &Response{}
}

func (res *Response) Compress(compressionType string, payload string) *Response {

	availableCompressionTypes := []string{"gzip"}
	requestCompressionsTypes := make([]string, 10)
	if strings.Contains(compressionType, ", ") {
		requestCompressionsTypes = strings.Split(compressionType, ", ")
	} else {
		requestCompressionsTypes = append(requestCompressionsTypes, compressionType)
	}

	res.headers = make(map[string]string)

	for _, requestCompressionsType := range requestCompressionsTypes {
		if slices.Contains(availableCompressionTypes, requestCompressionsType) {
			if requestCompressionsType == "gzip" {
				compressedPayload, err := compressions.GzipCompression(payload)
				if err != nil {
					continue
				}

				res.headers["Content-Type"] = "text/plain"
				res.headers["Content-Encoding"] = requestCompressionsType
				res.headers["Content-Length"] = strconv.Itoa(len(compressedPayload))
				res.body = compressedPayload
				break
			}
		}
	}

	return res
}

func (res *Response) Text(payload string) *Response {
	res.headers = make(map[string]string)
	res.headers["Content-Type"] = "text/plain"
	res.headers["Content-Length"] = strconv.Itoa(len(payload))

	res.body = payload
	return res
}

func (res *Response) Success() *Response {
	res.status = http.StatusOK
	res.statusText = "OK"

	return res
}

func (res *Response) InternalServerError() *Response {
	res.status = http.StatusInternalServerError
	res.statusText = "Internal Server Error"

	return res
}

func (res *Response) Created() *Response {
	res.status = http.StatusCreated
	res.statusText = "Created"

	return res
}

func (res *Response) Download(filePath string, fileName string) *Response {

	fmt.Println(filePath + fileName)
	// Validate File Existence
	absolutePath := fmt.Sprintf("%v%v", filePath, fileName)
	fileInfo, fileStatus := os.ReadFile(absolutePath)
	if fileStatus != nil {
		res.body = ""
		return res.NotFound()
	}

	// Populate Response
	res.headers = make(map[string]string)
	res.headers["Content-Type"] = "application/octet-stream"
	res.headers["Content-Length"] = strconv.Itoa(len(fileInfo))
	res.body = string(fileInfo)

	return res.Success()
}

func (res *Response) NotFound() *Response {
	res.status = http.StatusNotFound
	res.statusText = "Not Found"

	return res
}

func (res *Response) Build() string {
	var response string

	response += fmt.Sprintf("%v %v %v\r\n", "HTTP/1.1", res.status, res.statusText)

	for header, value := range res.headers {
		response += fmt.Sprintf("%v: %v\r\n", header, value)
	}
	// End of Headers
	response += "\r\n"
	if len(res.body) > 0 {
		response += res.body
	}

	fmt.Printf("Payload: %v\n", response)
	return response
}
