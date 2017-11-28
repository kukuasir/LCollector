package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteData(w http.ResponseWriter, res interface{}) {
	data, err := json.Marshal(res)
	if err != nil {
		log.Fatal("json marshal error: ", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Server", "BTK")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
	defer func() {
		w.Write(data)
	}()
}
