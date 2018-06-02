package urls_checker

import (
	"net/http"
	"fmt"
	"io/ioutil"
)

type UrlRequestBody struct {
	Domain string `json:"domain"`
	Path   string `json:"path"`
}

type UrlResponseBody struct {
	Location string `json:"location"`
}

var urlsChecker = NewUrlsChecker()

func StartServer() {
	http.HandleFunc("/check_url", checkURL)
	fmt.Println("server started")

	http.ListenAndServe(":8080", nil)
}

func checkURL(writer http.ResponseWriter, req *http.Request) {
	defer catchErrors(writer)

	bodyBytes    := getBodyBytes(req)

	reqBody      := deserializeRequestBody(bodyBytes)

	resBody      := responseFor(reqBody)

	jsonResponse := serializeToJSON(resBody)

	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonResponse)
}

func catchErrors(writer http.ResponseWriter) {
	if r := recover(); r != nil {
		http.Error(writer, r.(error).Error(), http.StatusInternalServerError)
	}
}

func handleError(err error) {
	if err == nil { return }
	panic(err)
}

func getBodyBytes(req *http.Request) []byte {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	handleError(err)

	return bodyBytes
}

func deserializeRequestBody(bodyBytes []byte) *UrlRequestBody {
	reqBody := &UrlRequestBody{}
	err := reqBody.UnmarshalJSON(bodyBytes)
	handleError(err)

	return reqBody
}

func serializeToJSON(resBody UrlResponseBody) []byte {
	jsonResponse, err := resBody.MarshalJSON()
	handleError(err)

	return jsonResponse
}

func responseFor(reqBody *UrlRequestBody) UrlResponseBody {
	url 		 := reqBody.Domain + reqBody.Path
	exists 	 := urlsChecker.Exists(url)
	location := (map[bool]string{true: url, false: ""})[exists]

	return UrlResponseBody { Location: location }
}



