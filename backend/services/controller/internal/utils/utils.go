package utils

import (
	"encoding/json"
	"io"
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

func MarshallEncoder(data any, w io.Writer) {
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Error to encode message into json: %q", err)
	}
}

func MarshallDecoder(data any, r io.Reader) error {
	err := json.NewDecoder(r).Decode(data)
	if err != nil {
		log.Printf("Error to decode message into json: %q", err)
	}

	return err
}
