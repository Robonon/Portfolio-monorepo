package generator

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func promptLLM(logger *slog.Logger, client *http.Client, llmURL, model, prompt, schema string) (string, error) {
	logger.Debug("Starting PromptLLM")
	reqBody, err := json.Marshal(ollamaRequest{Model: model, Prompt: prompt, Stream: true, Format: schema})
	if err != nil {
		return "", fmt.Errorf("failed to marshal LLM request: %w", err)
	}

	resp, err := client.Post(llmURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM returned non-200 status: %d", resp.StatusCode)
	}
	logger.Debug("Received LLM response", "status", resp.StatusCode)

	var ollamaResp ollamaResponse
	var fullResponse bytes.Buffer
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		// logger.Debug("LLM chunk", "data", scanner.Text())
		if err := json.Unmarshal(scanner.Bytes(), &ollamaResp); err != nil {
			return "", err
		}

		fullResponse.WriteString(ollamaResp.Response)

		if ollamaResp.Done {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading LLM response: %w", err)
	}

	cleaned := strings.TrimSpace(fullResponse.String())
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)
	cleaned = strings.ReplaceAll(cleaned, `\\n`, "\n")
	cleaned = strings.ReplaceAll(cleaned, `\\t`, "\t")
	cleaned = strings.ReplaceAll(cleaned, `\\\"`, `"`)
	// Unescape escaped characters
	unescaped := strings.ReplaceAll(cleaned, `\\n`, "\n")
	unescaped = strings.ReplaceAll(unescaped, `\\\"`, `"`)

	return unescaped, nil
}

func initGitRepo(logger *slog.Logger, outputDir string) error {
	os.RemoveAll(filepath.Join(outputDir, ".git"))
	logger.Debug("Initializing git repository")

	command := "git init && git config user.name 'Auto Commit' && git config user.email 'auto@example.com'"
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = outputDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to initialize git repository: %w output: %s", err, string(out))
	}
	logger.Debug("Running git command", "command", command, "dir", cmd.Dir)
	return nil
}

func commitEverything(logger *slog.Logger, message string, outputDir string) error {
	// Check for staged changes
	statusCommand := "git status --porcelain"
	cmd := exec.Command("sh", "-c", statusCommand)
	cmd.Dir = outputDir
	out, err := cmd.CombinedOutput()
	logger.Debug("git status output", "output", string(out), "error", err)

	if err != nil {
		return fmt.Errorf("failed to check git status: %w dir: %s", err, cmd.Dir)
	}

	if strings.TrimSpace(string(out)) == "" {
		logger.Info("No changes to commit")
		return nil
	}

	commitCommand := fmt.Sprintf("git add . && git commit -m '%s'", message)
	cmd2 := exec.Command("sh", "-c", commitCommand)
	cmd2.Dir = outputDir
	out2, err := cmd2.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to commit changes: %w output: %s", err, string(out2))
	}
	return nil
}

func createContextFromFiles(logger *slog.Logger, paths *[]string, outputDir string) (map[string]string, error) {
	logger.Debug("Running as user", "uid", os.Getuid(), "gid", os.Getgid())
	if paths == nil || len(*paths) == 0 {
		return map[string]string{}, nil
	}
	logger.Debug("Reading or creating files", "paths", paths)
	filesMap := make(map[string]string)
	for _, path := range *paths {
		fullPath := filepath.Join(outputDir, path)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			// If file doesn't exist, create it empty
			if os.IsNotExist(err) {
				dir := filepath.Dir(fullPath)
				if err := os.MkdirAll(dir, 0775); err != nil {
					return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
				}
				f, err := os.Create(fullPath)
				if err != nil {
					return nil, fmt.Errorf("failed to create file %s: %w", fullPath, err)
				}
				f.Close()
				filesMap[path] = ""
				continue
			}
			return nil, err
		}
		filesMap[path] = string(data)
	}
	return filesMap, nil
}

func newEvaluationResult() *EvaluationResult {
	return &EvaluationResult{
		Feedback: []EvaluationFeedback{},
		Pass:     false,
	}
}

func logCompletedJobs(logger *slog.Logger, jobs *sync.Map) {
	jobs.Range(func(key, value any) bool {
		job := value.(*Job)
		if job.Status == StatusDone || job.Status == StatusFailed {
			logger.Debug("Job completed", "id", job.ID, "status", job.Status, "result", job.Result, "error", job.Error)
			jobs.Delete(key)
		}
		return true
	})
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

func loadSchemas(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var schemas map[string]string
	if err := json.Unmarshal(data, &schemas); err != nil {
		return nil, err
	}
	return schemas, nil
}
