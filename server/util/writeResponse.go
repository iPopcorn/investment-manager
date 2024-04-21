package util

import (
	"log"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, response []byte, err error) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Response: \n%s", string(response))

	_, err = w.Write(response)

	if err != nil {
		log.Println("Failed to write response to writer")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
