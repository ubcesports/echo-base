package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/ubcesports/echo-base/internal/util"
)

func Login(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var user User
	json.NewDecoder(request.Body).Decode(&user)
}