package ideamart

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type StatusCode string

type Response struct {
	StatusCode   string `json:"statusCode"`
	StatusDetail string `json:"statusDetail"`
}

const (
	contentType string = "application/json"
)

var (
	SuccessResponse Response = Response{"S1000", "Success"}
)

func (r Response) Success() bool {
	return strings.HasPrefix(r.StatusCode, "S")
}

func (r Response) Error() *Error {
	if strings.HasPrefix(r.StatusCode, "E") {
		e := apiErrorFromCode(r.StatusCode)
		return &e
	}
	return nil
}

func isErrorCode(code string) bool {
	return strings.HasPrefix(code, "E")
}

func isSuccessCode(code string) bool {
	return (strings.HasPrefix(code, "S") || strings.HasPrefix(code, "P"))
}

func sendSuccessResponse(res http.ResponseWriter) {
	resBody, err := json.Marshal(SuccessResponse)
	if err != nil {
		res.WriteHeader(500)
		return
	}
	res.Header().Add("Content-Type", contentType)
	res.Write(resBody)
	return
}

func sendErrorResponse(res http.ResponseWriter) {
	res.WriteHeader(500)
}

func doRequest(endpoint string, request interface{}, response interface{}) error {
	reqBody, err := json.Marshal(request)
	if err != nil {
		return err
	}
	resp, err := http.Post(endpoint, contentType, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resBody, response)
	if err != nil {
		log.Print("Error parsing request response: ", err, string(resBody), response)
		return ErrInvalidJSON
	}
	//log.Print(string(resBody))
	return nil
}

func unmarshalRequest(req *http.Request, data interface{}) error {
	reqBody, err := ioutil.ReadAll(req.Body)
	log.Print(string(reqBody))
	err = json.Unmarshal(reqBody, data)
	return err
}
