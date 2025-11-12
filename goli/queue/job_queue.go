package queue

import (
	"fmt"
	"goli/database"
	"goli/models"
	"goli/pipeline"
	"goli/websocket"
	"log"
	"sync"
	"time"
)

// JobQueue manages the job queue and workers
type JobQueue struct {
	jobs     chan *models.Job
	workers  int
	wg       sync.WaitGroup
	stopChan chan struct{}
	mu       sync.RWMutex
	active   map[int64]*models.Job
	hub      *websocket.Hub
}

var (
	globalQueue *JobQueue
	once        sync.Once
)

// GetQueue returns the global job queue instance
func GetQueue() *JobQueue {
	once.Do(func() {
		globalQueue = NewJobQueue(3) // 3 workers by default
	})
	return globalQueue
}

// NewJobQueue creates a new job queue
func NewJobQueue(workers int) *JobQueue {
	return &JobQueue{
		jobs:     make(chan *models.Job, 100),
		workers:  workers,
		stopChan: make(chan struct{}),
		active:   make(map[int64]*models.Job),
	}
}

// SetWebSocketHub sets the WebSocket hub for broadcasting updates
func (q *JobQueue) SetWebSocketHub(hub *websocket.Hub) {
	q.hub = hub
}

// Start starts the job queue workers
func (q *JobQueue) Start() {
	log.Printf("Starting job queue with %d workers", q.workers)
	for i := 0; i < q.workers; i++ {
		q.wg.Add(1)
		go q.worker(i)
	}
}

// Stop stops all workers gracefully
func (q *JobQueue) Stop() {
	close(q.stopChan)
	close(q.jobs)
	q.wg.Wait()
	log.Println("Job queue stopped")
}

// Enqueue adds a job to the queue
func (q *JobQueue) Enqueue(job *models.Job) error {
	// Create job in database
	dbJob, err := database.CreateJob(job)
	if err != nil {
		return err
	}
	job.ID = dbJob.ID
	job.CreatedAt = dbJob.CreatedAt

	// Add to queue
	select {
	case q.jobs <- job:
		q.mu.Lock()
		q.active[job.ID] = job
		q.mu.Unlock()
		log.Printf("Job %d (%s) enqueued", job.ID, job.Name)
		return nil
	case <-time.After(30 * time.Second):
		return &QueueError{Message: "Queue is full, cannot enqueue job"}
	}
}

// EnqueueJob is an alias for Enqueue (for compatibility)
func (q *JobQueue) EnqueueJob(job *models.Job) error {
	return q.Enqueue(job)
}

// GetActiveJob returns an active job by ID
func (q *JobQueue) GetActiveJob(id int64) (*models.Job, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	job, exists := q.active[id]
	return job, exists
}

// RemoveActiveJob removes a job from active map
func (q *JobQueue) RemoveActiveJob(id int64) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.active, id)
}

// CancelJob cancels a running or pending job
func (q *JobQueue) CancelJob(id int64) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Check if job is in active map (running)
	job, exists := q.active[id]
	if exists {
		// Mark job as cancelled in database
		if err := database.UpdateJobStatus(id, models.JobStatusCancelled, "Job cancelled by user"); err != nil {
			return err
		}
		job.Status = models.JobStatusCancelled
		completedAt := time.Now()
		job.CompletedAt = &completedAt

		// Broadcast update
		if q.hub != nil {
			q.hub.BroadcastJobUpdate(job)
		}

		// Remove from active map
		delete(q.active, id)
		log.Printf("Job %d cancelled (was running)", id)
		return nil
	}

	// Job might be pending in queue or already completed
	// Check database status
	dbJob, err := database.GetJob(id)
	if err != nil {
		return err
	}

	// Only cancel if pending or running
	if dbJob.Status == models.JobStatusPending || dbJob.Status == models.JobStatusRunning {
		if err := database.UpdateJobStatus(id, models.JobStatusCancelled, "Job cancelled by user"); err != nil {
			return err
		}
		completedAt := time.Now()
		dbJob.CompletedAt = &completedAt
		dbJob.Status = models.JobStatusCancelled

		// Broadcast update
		if q.hub != nil {
			q.hub.BroadcastJobUpdate(dbJob)
		}

		log.Printf("Job %d cancelled (was %s)", id, dbJob.Status)
		return nil
	}

	return &QueueError{Message: fmt.Sprintf("Job %d cannot be cancelled (status: %s)", id, dbJob.Status)}
}

// worker processes jobs from the queue
func (q *JobQueue) worker(id int) {
	defer q.wg.Done()
	log.Printf("Worker %d started", id)

	for {
		select {
		case job, ok := <-q.jobs:
			if !ok {
				log.Printf("Worker %d: queue closed, shutting down", id)
				return
			}
			q.processJob(job)
		case <-q.stopChan:
			log.Printf("Worker %d: stop signal received, shutting down", id)
			return
		}
	}
}

// processJob processes a single job
func (q *JobQueue) processJob(job *models.Job) {
	log.Printf("Worker: Processing job %d (%s)", job.ID, job.Name)

	// Update status to running
	if err := database.UpdateJobStatus(job.ID, models.JobStatusRunning, ""); err != nil {
		log.Printf("Error updating job status: %v", err)
		return
	}
	job.Status = models.JobStatusRunning
	now := time.Now()
	job.StartedAt = &now

	// Broadcast update
	if q.hub != nil {
		q.hub.BroadcastJobUpdate(job)
	}

	// Execute pipeline if pipeline_id is provided
	if job.PipelineID != nil {
		pipelineRecord, err := database.GetPipeline(*job.PipelineID)
		if err != nil {
			log.Printf("Error loading pipeline: %v", err)
			database.UpdateJobStatus(job.ID, models.JobStatusFailed, "Failed to load pipeline: "+err.Error())
			return
		}

		// Parse pipeline definition
		pipelineDef, err := pipeline.ParsePipelineDefinition(pipelineRecord.Definition)
		if err != nil {
			log.Printf("Error parsing pipeline definition: %v", err)
			database.UpdateJobStatus(job.ID, models.JobStatusFailed, "Failed to parse pipeline: "+err.Error())
			return
		}

		// Execute the pipeline
		if err := pipeline.ExecutePipeline(job, pipelineDef); err != nil {
			log.Printf("Error executing pipeline: %v", err)
			// Status already updated by executor
			return
		}
	} else {
		// No pipeline, just mark as completed (simple job)
		time.Sleep(1 * time.Second)
		if err := database.UpdateJobStatus(job.ID, models.JobStatusCompleted, ""); err != nil {
			log.Printf("Error updating job status: %v", err)
			return
		}
		job.Status = models.JobStatusCompleted
		completedAt := time.Now()
		job.CompletedAt = &completedAt
	}

	// Broadcast final update
	if q.hub != nil {
		q.hub.BroadcastJobUpdate(job)
	}

	q.RemoveActiveJob(job.ID)
	log.Printf("Worker: Job %d (%s) completed", job.ID, job.Name)
}

// QueueError represents a queue error
type QueueError struct {
	Message string
}

func (e *QueueError) Error() string {
	return e.Message
}
