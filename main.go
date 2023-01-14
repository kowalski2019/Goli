package main

import (
	"log"
	"net/http"

	"deployer/handler"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const port = "127.0.0.1:8125"

func main() {
	r := mux.NewRouter()
	encr := r.PathPrefix("/api/v1/").Subrouter()
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

	log.Println("Action Helper Backend is running on " + port)
	log.Fatal(http.ListenAndServe(port, handlers.CORS(
		handlers.AllowedHeaders(
			[]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}
