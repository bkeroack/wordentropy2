package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type json_error_msg struct {
	error_message string
}

func emit_json_error(w http.ResponseWriter, msg string, n int) {
	err_json := json_error_msg{
		error_message: msg,
	}
	err_str, _ := json.Marshal(err_json)
	log.Printf(err_json.error_message)
	http.Error(w, string(err_str), n)
}
