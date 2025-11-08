package middlewares

import (
	aux "goli/auxiliary"
	response_util "goli/utils"
	"net/http"

	"strings"
)

var auth_key = aux.GetFromConfig("constants.auth_key")

func VerifyAuth(w http.ResponseWriter, r *http.Request) bool {
	if GetAuthKeyFromRequest(r) == auth_key {
		return true
	} else {
		response_util.SendUnauthorizedResponse(w, "Wrong auth key provided")
		return false
	}
}

func ExtractAuthKey(token string) string {
	// array[0] = Goli-Auth-Key
	return strings.Split(token, " ")[1]
}

func GetAuthKeyFromRequest(r *http.Request) string {
	return ExtractAuthKey(r.Header.Get("Authorization"))
}
