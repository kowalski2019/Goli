package queue

import (
	"goli/database"
	"goli/websocket"
	"log"
	"sync"
	"time"
)

// LogStreamer periodically broadcasts logs for long-running jobs
type LogStreamer struct {
	hub          *websocket.Hub
	interval     time.Duration
	stopChan     chan struct{}
	lastLogs     map[int64]string // Track last sent logs to avoid duplicates
	lastStepLogs map[int64]string // Track last sent step logs
	mu           sync.RWMutex
}

var (
	globalLogStreamer *LogStreamer
	streamerOnce      sync.Once
)

// GetLogStreamer returns the global log streamer instance
func GetLogStreamer() *LogStreamer {
	streamerOnce.Do(func() {
		globalLogStreamer = NewLogStreamer(5 * time.Minute) // Default: 5 minutes
	})
	return globalLogStreamer
}

// NewLogStreamer creates a new log streamer
func NewLogStreamer(interval time.Duration) *LogStreamer {
	return &LogStreamer{
		hub:          nil,
		interval:     interval,
		stopChan:     make(chan struct{}),
		lastLogs:     make(map[int64]string),
		lastStepLogs: make(map[int64]string),
	}
}

// SetHub sets the WebSocket hub for broadcasting
func (ls *LogStreamer) SetHub(hub *websocket.Hub) {
	ls.hub = hub
}

// Start starts the log streaming service
func (ls *LogStreamer) Start() {
	if ls.hub == nil {
		log.Println("LogStreamer: No hub set, cannot start")
		return
	}
	log.Printf("LogStreamer: Starting with interval %v", ls.interval)
	go ls.run()
}

// Stop stops the log streaming service
func (ls *LogStreamer) Stop() {
	close(ls.stopChan)
	log.Println("LogStreamer: Stopped")
}

// run periodically checks for running jobs and broadcasts their logs
func (ls *LogStreamer) run() {
	ticker := time.NewTicker(ls.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ls.broadcastLogs()
		case <-ls.stopChan:
			return
		}
	}
}

// broadcastLogs checks all running jobs and broadcasts their logs
func (ls *LogStreamer) broadcastLogs() {
	if ls.hub == nil {
		return
	}

	runningJobs, err := database.GetRunningJobs()
	if err != nil {
		log.Printf("LogStreamer: Error getting running jobs: %v", err)
		return
	}

	for _, job := range runningJobs {
		// Check if job has been running for at least the interval duration
		if job.StartedAt == nil {
			continue
		}
		elapsed := time.Since(*job.StartedAt)
		if elapsed < ls.interval {
			continue // Job hasn't been running long enough
		}

		// Get current job with logs
		currentJob, err := database.GetJob(job.ID)
		if err != nil {
			log.Printf("LogStreamer: Error getting job %d: %v", job.ID, err)
			continue
		}

		// Check if logs have changed
		ls.mu.RLock()
		lastJobLogs := ls.lastLogs[job.ID]
		ls.mu.RUnlock()

		if currentJob.Logs != lastJobLogs && currentJob.Logs != "" {
			// Get running step if any
			var stepID *int64
			var stepLogs string
			runningStep, err := database.GetRunningStep(job.ID)
			if err == nil && runningStep != nil {
				stepID = &runningStep.ID
				stepLogs = runningStep.Logs

				// Check if step logs have changed
				ls.mu.RLock()
				lastStepLogs := ls.lastStepLogs[runningStep.ID]
				ls.mu.RUnlock()

				if stepLogs != lastStepLogs {
					ls.hub.BroadcastLogUpdate(job.ID, currentJob.Logs, stepID, stepLogs)
					ls.mu.Lock()
					ls.lastStepLogs[runningStep.ID] = stepLogs
					ls.mu.Unlock()
				}
			} else {
				// Broadcast job-level logs only
				ls.hub.BroadcastLogUpdate(job.ID, currentJob.Logs, nil, "")
			}

			// Update last logs
			ls.mu.Lock()
			ls.lastLogs[job.ID] = currentJob.Logs
			ls.mu.Unlock()
		}
	}
}

// ClearJob clears tracking for a completed/cancelled job
func (ls *LogStreamer) ClearJob(jobID int64) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	delete(ls.lastLogs, jobID)
	// Note: We don't clear step logs as they're keyed by step ID, not job ID
}
