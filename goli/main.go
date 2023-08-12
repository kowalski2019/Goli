package main

import (
	aux "deployer/auxiliary"
	"deployer/handler"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var port = aux.GetFromConfig("constants.port")

var host = "127.0.0.1:" + port

func main() {
	r := mux.NewRouter()
	encr := r.PathPrefix("/api/v1/").Subrouter()

	//encr.HandleFunc("/docker/compose/prepare", handler.StartAwakeDockers).Methods(http.MethodPost)
	encr.HandleFunc("/docker/compose/up", handler.StartADockerOrchestra).Methods(http.MethodPost)
	encr.HandleFunc("/docker/compose/down", handler.StopADockerOrchestra).Methods(http.MethodPost)

	encr.HandleFunc("/docker/container/start", handler.StartADocker).Methods(http.MethodPost)
	encr.HandleFunc("/docker/container/stop", handler.StopADocker).Methods(http.MethodPost)
	encr.HandleFunc("/docker/container/rm", handler.RemoveADocker).Methods(http.MethodPost)
	encr.HandleFunc("/docker/container/run", handler.RunDockerContainer).Methods(http.MethodPost)

	encr.HandleFunc("/docker/container/pause", handler.PauseADocker).Methods(http.MethodPost)
	encr.HandleFunc("/docker/container/unpause", handler.UnPauseADocker).Methods(http.MethodPost)
	encr.HandleFunc("/docker/container/inspect", handler.InspectADocker).Methods(http.MethodPost)
	encr.HandleFunc("/docker/container/logs", handler.GetADockerLogs).Methods(http.MethodPost)
	encr.HandleFunc("/docker/image/rm", handler.RemoveAnDockerImage).Methods(http.MethodPost)
	encr.HandleFunc("/docker/image/pull", handler.PullAnDockerImage).Methods(http.MethodPost)

	encr.HandleFunc("/docker/ps", handler.GetDockerPS).Methods(http.MethodPost)
	encr.HandleFunc("/docker/images", handler.GetDockerImages).Methods(http.MethodPost)

	log.Println("Action Helper Backend is running on " + host)
	log.Fatal(http.ListenAndServe(host, handlers.CORS(
		handlers.AllowedHeaders(
			[]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}
