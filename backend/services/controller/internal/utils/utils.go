package utils

import (
	"encoding/json"
	"log"
)

func Marshall(data any) []byte {
	fmtData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error to marshall message into json: %q", err)
		return []byte(err.Error())
	}
	return fmtData
}
