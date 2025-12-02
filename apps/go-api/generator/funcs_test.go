package generator

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
func TestGenerateID(t *testing.T) {
	id1 := generateID()
	id2 := generateID()
	if id1 == "" || id2 == "" {
		t.Error("generateID returned empty string")
	}
	if id1 == id2 {
		t.Error("generateID returned duplicate IDs")
	}
}

func TestNewEvaluationResult(t *testing.T) {
	e := newEvaluationResult()
	if e == nil {
		t.Error("newEvaluationResult returned nil")
		return
	}
	if e.Pass {
		t.Error("newEvaluationResult Pass should be false")
	}
	if e.Feedback == nil {
		t.Error("newEvaluationResult Feedback should not be nil")
	}
	if len(e.Feedback) != 0 {
		t.Error("newEvaluationResult Feedback should be empty")
	}
}

func TestZipDirectory(t *testing.T) {
	dir := t.TempDir()
	file1 := filepath.Join(dir, "file1.txt")
	file2 := filepath.Join(dir, "file2.txt")
	os.WriteFile(file1, []byte("hello"), 0644)
	os.WriteFile(file2, []byte("world"), 0644)

	zipPath := filepath.Join(dir, "test.zip")
	err := ZipDirectory(dir, zipPath)
	if err != nil {
		t.Fatalf("ZipDirectory failed: %v", err)
	}
	info, err := os.Stat(zipPath)
	if err != nil || info.IsDir() {
		t.Error("ZipDirectory did not create a zip file")
	}
}

func TestLoadSchemas(t *testing.T) {
	dir := t.TempDir()
	schema := map[string]json.RawMessage{
		"TestSchema": json.RawMessage(`{"type":"object"}`),
	}
	data, _ := json.Marshal(schema)
	schemaPath := filepath.Join(dir, "schemas.json")
	os.WriteFile(schemaPath, data, 0644)

	schemas, err := loadSchemas(schemaPath)
	if err != nil {
		t.Fatalf("loadSchemas failed: %v", err)
	}
	if _, ok := schemas["TestSchema"]; !ok {
		t.Error("TestSchema not found in loaded schemas")
	}
}

func TestCreateFolders(t *testing.T) {
	logger := newTestLogger()
	dir := t.TempDir()
	paths := &[]string{"a/b/c/", "d/e/"}
	err := createFolders(logger, paths, dir)
	if err != nil {
		t.Fatalf("createFolders failed: %v", err)
	}
}
func TestCreateFileWithContent(t *testing.T) {
	logger := newTestLogger()
	dir := t.TempDir()
	path := "a/b/c.txt"
	err := createFileWithContent(logger, path, "content", dir)
	if err != nil {
		t.Fatalf("createFileWithContent failed: %v", err)
	}
}

func TestLogCompletedJobs(t *testing.T) {
	logger := newTestLogger()
	jobs := &sync.Map{}
	job1 := &Job{ID: "1", Status: StatusDone}
	job2 := &Job{ID: "2", Status: StatusFailed}
	job3 := &Job{ID: "3", Status: StatusRunning}
	jobs.Store("1", job1)
	jobs.Store("2", job2)
	jobs.Store("3", job3)
	logCompletedJobs(logger, jobs)
	if _, ok := jobs.Load("1"); ok {
		t.Error("Completed job 1 was not deleted")
	}
	if _, ok := jobs.Load("2"); ok {
		t.Error("Completed job 2 was not deleted")
	}
	if _, ok := jobs.Load("3"); !ok {
		t.Error("Running job 3 was deleted")
	}
}
