package llm

import (
	"api/config"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type (
	PromptInput struct {
		Prompt string `json:"prompt"`
	}

	evaluationResult struct {
		Issues      string `json:"issues"`
		Suggestions string `json:"suggestions"`
		Pass        bool   `json:"pass"`
	}

	Job struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Result string `json:"result,omitempty"`
		Error  string `json:"error,omitempty"`
		Stage  string `json:"stage,omitempty"`
	}
)

const (
	StatusPending  string = "pending"
	StatusRunning  string = "running"
	StatusDone     string = "done"
	StatusFailed   string = "failed"
	StageCompleted string = ""
	evalSchema            = `{
  "type": "object",
  "properties": {
    "issues": { "type": "string" },
    "suggestions": { "type": "string" },
    "pass": { "type": "boolean" }
  },
  "required": ["issues", "suggestions", "pass"]
}`
	outputSchema = `{
	"type": "object",
	"properties": {
		"architecture": { "type": "string" },
		"interfaces": { "type": "string" },
		"unit_tests": { "type": "string" },
		"implementation": { "type": "string" },
		"integration": { "type": "string" },
		"documentation": { "type": "string" }
		},
		"required": ["architecture", "interfaces", "unit_tests", "implementation", "integration", "documentation"]
}`
)

var Jobs sync.Map // map[string]*Job
var stages = [6]string{"architecture", "interfaces", "unit_tests", "implementation", "integration", "documentation"}

func GenerateModuleHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input PromptInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			logger.Error("Failed to decode JSON", "error", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		logger.Info("Received prompt", "prompt", input.Prompt)

		jobID := generateID()
		job := &Job{ID: jobID, Status: StatusPending}
		Jobs.Store(jobID, job)

		logger.Info("Job created", "jobID", jobID)
		// Respond immediately
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "Module generation started, output will be saved to output.go",
			"jobId":  jobID,
		})

		go func() {
			var Url = config.GetConfig().LLMUrl + "/api/generate"
			var model = "gemma3"
			var context = ""
			var evaluation = newEvaluationResult()
			var result = ""
			var err error
			var evalRes string
			logger.Info("Job started", "jobID", jobID)
			job.Status = StatusRunning

			for _, stage := range stages {
				job.Stage = stage
				logger.Info("Starting stage", "jobID", jobID, "stage", stage)
				for !evaluation.Pass {
					// Get result from LLM
					result, err = StreamPromptToFile(Url, model, promptTemplate(input.Prompt, context, evaluation.Suggestions), "output.go", outputSchema)
					if err != nil {
						job.Status = StatusFailed
						job.Error = err.Error()
						return
					}
					evalRes, err = StreamPromptToFile(Url, model, evaluationTemplate(stage, context, input.Prompt), "feedback.go", evalSchema)
					if err != nil {
						job.Status = StatusFailed
						job.Error = err.Error()
						return
					}
					err = json.Unmarshal([]byte(evalRes), &evaluation)
					if err != nil {
						job.Status = StatusFailed
						job.Error = err.Error()
						return
					}
				}
			}
			job.Status = StatusDone
			job.Result = result

			logger.Info("Job completed", "jobID", jobID, "status", job.Status)
			Jobs.Store(jobID, job)
		}()

		// Start a goroutine to log completed jobs every 30 seconds
		go func() {
			for {
				logCompletedJobs(logger, &Jobs)
				time.Sleep(30 * time.Second)
			}
		}()
	}
}

func newEvaluationResult() *evaluationResult {
	return &evaluationResult{
		Issues:      "",
		Suggestions: "",
		Pass:        false,
	}
}

func promptTemplate(userPrompt, context, evaluationFeedback string) string {
	return fmt.Sprintf(`
You are an expert Go developer.

## Context
%s

## Evaluation Feedback
%s

## Instructions
- Use cloud native and idiomatic Go best practices.
- The code must:
  - Be stateless and container-ready.
  - Use environment variables for configuration.
  - Organize code using idiomatic Go structure (e.g., cmd/, internal/, pkg/, etc.).
  - Include GoDoc comments for all exported items.
  - Provide a README.md with usage instructions.
  - Use dependency injection and context where appropriate.
  - Include health checks and proper logging.
  - Follow 12-factor app principles.
  - Only use standard library and well-known, maintained libraries.
  - Ensure proper error handling and security best practices.

## Task
Generate a Go package that accomplishes the following:

---
%s
---

## Output Format
Return your response as a JSON object with the following structure:
{
  "architecture": "...",
  "interfaces": "...",
  "unit_tests": "...",
  "implementation": "...",
  "integration": "...",
  "documentation": "..."
}
Each field should contain the respective code or documentation. Do not include explanations or comments outside the code.
`, context, evaluationFeedback, userPrompt)
}

func evaluationTemplate(stage, context, prompt string) string {
	return fmt.Sprintf(`
You are an expert Go developer.

## Context
%s

## Instructions
- Use cloud native and idiomatic Go best practices.
- Evaluate the code for:
  - Quality
  - Adherence to best practices
  - Potential issues
- Provide constructive feedback.
- Identify security vulnerabilities or performance bottlenecks.
- Suggest improvements or refactoring opportunities.
- The evaluation passes if there are no suggested changes to the %s

## Task
Evaluate the following %s with the given context and provide feedback on its quality, adherence to best practices, and any potential issues given the following goal, and finally if it passes or fails the evaluation:

---
%s
---

## Output Format
Return your response as a JSON object with the following structure:
{
  "issues": [...],
  "suggestions": [...],
  "pass": true/false
}
Each field should contain the respective code or documentation. Do not include explanations or comments outside the code.
`, context, stage, stage, prompt)
}

func logCompletedJobs(logger *slog.Logger, jobs *sync.Map) {
	jobs.Range(func(key, value any) bool {
		job := value.(*Job)
		if job.Status == StatusDone || job.Status == StatusFailed {
			logger.Info("Job completed", "id", job.ID, "status", job.Status, "result", job.Result, "error", job.Error)
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

// Prompt a goal
/*
1. Architecture Design
Prompt:
Ask the LLM to design the overall architecture for your application.
Output:
List of required packages (e.g., cmd/, internal/service, internal/storage, pkg/utils, api/handlers)
Description of each package’s responsibility
High-level diagram or outline of interactions

2. Interface Definition
Prompt:
For each package, request the definition of all public interfaces and types.
Output:
Interface files for each package (e.g., internal/service/service.go with type Service interface { ... })
Types and contracts that other packages will depend on

3. Unit Test Generation
Prompt:
For each package, ask for a complete set of Go unit tests based on the interfaces.
Output:
Test files (e.g., internal/service/service_test.go)
Tests for all public methods and behaviors

4. Implementation
Prompt:
For each package, request the implementation code that passes the previously generated tests.
Output:
Implementation files (e.g., internal/service/service_impl.go)
Code that fulfills the interfaces and passes all tests

5. Integration
Prompt:
Ask for the main entrypoint (cmd/main.go) and integration glue code.
Output:
Main application file
Wiring of packages using dependency injection
Configuration via environment variables

6. Documentation & Health Checks
Prompt:
Request GoDoc comments, a README.md, and health check endpoints.
Output:
Documentation for all exported items
Usage instructions
Health check implementation (e.g., /healthz endpoint)

*/
