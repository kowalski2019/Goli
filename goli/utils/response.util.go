package utils

import (
	"encoding/json"

	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func SendOctetStreamResponse(w http.ResponseWriter, status int, octet []byte, mediaName string) {
	w.Header().Set("Content-Type", "multipart/form-data")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")
	w.Header().Set("Content-Disposition", "attachment; filename="+mediaName)
	w.WriteHeader(status)
	w.Write(octet)
}

func SendJsonResponse(w http.ResponseWriter, status int, json []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)
}

func SendResponse(w http.ResponseWriter, statusCode int, statusName string, description string) {
	var res bson.M
	if description != "" {
		res = bson.M{"statusCode": statusCode, "statusName": statusName, "description": description}
	} else {
		res = bson.M{"statusCode": statusCode, "statusName": statusName}
	}

	jsonAsBytes, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	SendJsonResponse(w, statusCode, jsonAsBytes)
}

func SendOkResponse(w http.ResponseWriter, res string) {
	SendResponse(w, http.StatusOK, "OK", res)
}
func SendBadRequestResponse(w http.ResponseWriter, res string) {
	SendResponse(w, http.StatusBadRequest, "BadRequest", res)
}
func SendUnauthorizedResponse(w http.ResponseWriter, res string) {
	SendResponse(w, http.StatusUnauthorized, "Unauthorized", res)
}
func SendForbiddenResponse(w http.ResponseWriter, res string) {
	SendResponse(w, http.StatusForbidden, "Forbidden", res)
}
func SendNotFoundResponse(w http.ResponseWriter, res string) {
	SendResponse(w, http.StatusNotFound, "NotFound", res)
}
func SendConflictResponse(w http.ResponseWriter, res string) {
	SendResponse(w, http.StatusConflict, "Conflict", res)
}
func SendInternalServerErrorResponse(w http.ResponseWriter, res string) {
	SendResponse(w, http.StatusInternalServerError, "InternalServerError", res)
}
