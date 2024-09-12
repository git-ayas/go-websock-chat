package main

import (
	"encoding/json"
)

func UnmarshaledJsonFromByteArray[ReturnShape comparable](jsonString []byte) (out ReturnShape, err error) {

	// log.Printf("parsing the following Json: %v", string(jsonString))
	err = json.Unmarshal(jsonString, &out)

	return

}
