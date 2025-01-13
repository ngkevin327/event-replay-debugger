package queue

import (
	"encoding/json"
	"errors"
	"sync"
)

// ErrEmpty is returned when no jobs are available.
var ErrEmpty = errors.New("queue empty")

// CollectionJob identifies an incident collection task.
type CollectionJob struct {
	IncidentID string `json:"incident_id"`
}

// RedisQueue stores collection jobs (in-memory when Redis is unavailable).
type RedisQueue struct {
	mu   sync.Mutex
	jobs []CollectionJob
}

// NewRedisQueue creates a queue backed by process memory for MVP.
func NewRedisQueue() *RedisQueue {
	return &RedisQueue{}
}

// PushJob enqueues a collection job.
func (q *RedisQueue) PushJob(job CollectionJob) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.jobs = append(q.jobs, job)
	return nil
}

// PopJob removes and returns the oldest job.
func (q *RedisQueue) PopJob() (CollectionJob, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.jobs) == 0 {
		return CollectionJob{}, ErrEmpty
	}
	job := q.jobs[0]
	q.jobs = q.jobs[1:]
	return job, nil
}

// MarshalJob serializes a job for Redis list storage.
func MarshalJob(job CollectionJob) ([]byte, error) {
	return json.Marshal(job)
}
