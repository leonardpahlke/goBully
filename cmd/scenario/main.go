package main

import (
	"github.com/sirupsen/logrus"
)

type testStruct struct {
	Message string `json:"message"`
}

func main() {
	var responses []testStruct
	var responseTestStructHE = testStruct{
		Message: "he",
	}
	var responseTestStructHEHE = testStruct{
		Message: "hehe",
	}

	responses = append(responses, responseTestStructHEHE)

	testMethod(&responseTestStructHE, &responses)

	for _, res := range responses {
		logrus.Info("res - ", res.Message)
	}
	logrus.Info("responseTestStructHE - ", responseTestStructHE.Message)
}

func testMethod(test *testStruct, testList *[]testStruct) {
	*test = testStruct{
		Message: "hello",
	}
	*testList = append(*testList, *test)
}
