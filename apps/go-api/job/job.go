package job

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type (
	Job struct {
		Id     string
		Status int
	}
)

const (
	Succeeded = iota
	Running
	Failed
)

func NewJob() *Job {
	return &Job{
		Id:     generateID(),
		Status: Running,
	}
}

func generateID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// fallback to timestamp if random fails
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}
