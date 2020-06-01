package goBully

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

// RequestGET - send http get request
func RequestGET(url string, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Fatalf("[request.RequestGET] Error instantiating request with error %s", err)
		return nil, err
	}
	if token != "" {
		req.Header.Add("Authorization", token)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Fatalf("[request.RequestGET] Error sending request with error %s", err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Fatalf("[request.RequestGET] Error parsing response with error %s", err)
		return nil, err
	}
	return body, nil
}

// RequestPOST - send http post request
func RequestPOST(url string, payloadString string, token string) ([]byte, error) {
	payload := strings.NewReader(payloadString)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		logrus.Fatalf("[request.RequestPOST] Error instantiating request with error %s", err)
		return nil, err
	}
	if token != "" {
		req.Header.Add("Authorization", token)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Fatalf("[request.RequestPOST] Error sending request with error %s", err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Fatalf("[request.RequestPOST] Error parsing response with error %s", err)
		return nil, err
	}
	return body, nil
}
