package job

import "testing"

func TestCreateJob(t *testing.T) {
	job := NewJob()

	if job != nil && job.Id == "" {
		t.Error("Job id should not be empty")
	}

	if job != nil && job.Status != Running {
		t.Errorf("Expected status Running, got %v", job.Status)
	}
}

func TestFailedJob(t *testing.T) {

}

func TestSucceedJob(t *testing.T) {

}

func TestDeleteJob(t *testing.T) {}
