package main

import (
	aux "goli/auxiliary"
	"goli/database"
	"goli/handler"
	"goli/queue"
	"goli/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var port = aux.GetFromConfig("constants.port")

var host = "127.0.0.1:" + port

func main() {
	// Initialize database
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDatabase()

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Initialize and start job queue
	jobQueue := queue.GetQueue()
	queue.SetWebSocketHub(wsHub) // Pass hub to queue for broadcasting
	jobQueue.Start()
	defer jobQueue.Stop()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
		jobQueue.Stop()
		os.Exit(0)
	}()

	r := mux.NewRouter()
	encr := r.PathPrefix("/api/v1/").Subrouter()

	// Job management endpoints
	encr.HandleFunc("/jobs", handler.ListJobsHandler).Methods(http.MethodGet)
	encr.HandleFunc("/jobs", handler.CreateJobHandler).Methods(http.MethodPost)
	encr.HandleFunc("/jobs/{id}", handler.GetJobHandler).Methods(http.MethodGet)

	// Pipeline management endpoints
	encr.HandleFunc("/pipelines", handler.ListPipelinesHandler).Methods(http.MethodGet)
	encr.HandleFunc("/pipelines", handler.CreatePipelineHandler).Methods(http.MethodPost)
	encr.HandleFunc("/pipelines/upload", handler.UploadPipelineHandler).Methods(http.MethodPost)
	encr.HandleFunc("/pipelines/{id}", handler.GetPipelineHandler).Methods(http.MethodGet)
	encr.HandleFunc("/pipelines/{id}/run", handler.RunPipelineHandler).Methods(http.MethodPost)

	// WebSocket endpoint
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeWebSocket(wsHub, w, r)
	})

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

	// Serve static files (UI)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))

	log.Println("Action Helper Backend is running on " + host)
	log.Fatal(http.ListenAndServe(host, handlers.CORS(
		handlers.AllowedHeaders(
			[]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}
