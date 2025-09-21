package generator

import (
	"api/config"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type (
	generatorImpl struct {
		logger     *slog.Logger
		httpClient *http.Client
		config     *config.Config
	}
	stage struct {
		name                   string
		promptValues           promptTemplateValues
		evaluationPromptValues promptTemplateValues
		schema                 string
	}
)

const (
	StatusPending  string = "pending"
	StatusRunning  string = "running"
	StatusDone     string = "done"
	StatusFailed   string = "failed"
	StageCompleted string = ""
)

var (
	schemas, _ = loadSchemas("go-api/generator/schemas.json")
	Jobs       sync.Map
	stages     = []stage{
		{
			name:                   "architecture",
			promptValues:           architecturePrompt,
			evaluationPromptValues: architectureEvaluation,
			schema:                 schemas["ArchitectureSchema"],
		},
		{
			name:                   "interfaces",
			promptValues:           interfacesPrompt,
			evaluationPromptValues: interfacesEvaluation,
			schema:                 schemas["InterfacesSchema"],
		},
		{
			name:                   "unit_tests",
			promptValues:           unitTestsPrompt,
			evaluationPromptValues: unitTestsEvaluation,
			schema:                 schemas["UnitTestsSchema"],
		},
		// {
		// 	name:                   "implementation",
		// 	promptValues:           implementationPrompt,
		// 	evaluationPromptValues: implementationEvaluation,
		// },
		{
			name:                   "integration",
			promptValues:           integrationPrompt,
			evaluationPromptValues: integrationEvaluation,
			schema:                 schemas["IntegrationSchema"],
		},
		{
			name:                   "documentation",
			promptValues:           documentationPrompt,
			evaluationPromptValues: documentationEvaluation,
			schema:                 schemas["DocumentationSchema"],
		},
	}
)

func (g *generatorImpl) GenerateModule(input GeneratorInput) error {
	var (
		Url              = g.config.LLMUrl
		model            = "gemma3"
		job              = &Job{ID: generateID(), Status: StatusRunning}
		context          = make(map[string]string)
		rawGenResult     = ""
		rawEvalResult    = ""
		parsedEvalResult = newEvaluationResult()
	)

	Jobs.Store(job.ID, job)
	g.logger.Info("Job created", "jobID", job.ID)

	// Init git repo
	err := initGitRepo(g.logger, g.config.OutputDir)
	if err != nil {
		job.Status = StatusFailed
		job.Error = err.Error()
		return err
	}
	g.logger.Info("Git repository initialized")

	for i := range stages {
		stage := &stages[i]
		stage.promptValues.Scope = input.Scope
		stage.promptValues.ResponseSchema = stage.schema
		stage.evaluationPromptValues.ResponseSchema = schemas["EvaluationSchema"]
		g.logger.Info("Starting stage: " + stage.name)
		g.logger.Debug("Starting stage", "jobID", job.ID, "stage", stage.name)
		parsedEvalResult.Pass = false
		job.Stage = stage.name

		// Loop until evaluation passes
		for !parsedEvalResult.Pass {

			prompt, err := createPrompt(stage.promptValues)
			if err != nil {
				job.Status = StatusFailed
				job.Error = err.Error()
				return err
			}
			g.logger.Debug(stage.name+" prompt", "prompt", prompt)
			g.logger.Info("Prompting generation...")
			rawGenResult, err = promptLLM(
				g.logger,
				g.httpClient,
				Url,
				model,
				prompt,
				stage.schema)
			if err != nil {
				job.Status = StatusFailed
				job.Error = err.Error()
				return err
			}
			g.logger.Debug(stage.name+" result", "result", rawGenResult)

			// execute loop based on stage
			switch stage.name {
			case "architecture":
				// parse result
				var archRes ArchitectureSchema
				err = json.Unmarshal([]byte(rawGenResult), &archRes)
				if err != nil {
					job.Status = StatusFailed
					job.Error = err.Error()
					return err
				}

				if len(archRes.Filepaths) > 0 {
					// reads or creates files and returns updated context
					g.logger.Info("Creating context from architecture files")
					context, err = createContextFromFiles(g.logger, &archRes.Filepaths, g.config.OutputDir)
					if err != nil {
						job.Status = StatusFailed
						job.Error = err.Error()
						return err
					}
					// update context in stage prompts
					jsonContext, err := json.Marshal(context)
					if err != nil {
						job.Status = StatusFailed
						job.Error = err.Error()
						return err
					}

					stage.promptValues.Context = string(jsonContext)
					stage.evaluationPromptValues.Context = string(jsonContext)

				}

				// case "interfaces":
				// 	var ifaceRes interfacesSchema
				// 	err = json.Unmarshal([]byte(result), &ifaceRes)
				// 	context, err = readFiles(g.logger, &ifaceRes.Filepaths)

				// case "unit_tests":
				// 	var utRes unitTestsSchema
				// 	err = json.Unmarshal([]byte(result), &utRes)
				// 	context, err = readFiles(g.logger, &utRes.Filepaths)
				// case "integration":
				// 	var intRes integrationSchema
				// 	err = json.Unmarshal([]byte(result), &intRes)
				// 	context, err = readFiles(g.logger, &intRes.Filepaths)
				// case "documentation":
				// 	var docRes documentationSchema
				// 	err = json.Unmarshal([]byte(result), &docRes)
				// 	context, err = readFiles(g.logger, &docRes.Filepaths)
			}

			// commit files
			g.logger.Info("Committing changes to git...")
			err = commitEverything(g.logger, stage.name, g.config.OutputDir)
			if err != nil {
				job.Status = StatusFailed
				job.Error = err.Error()
				return err
			}

			// prepare evaluation prompt
			prompt, err = createPrompt(stage.evaluationPromptValues)
			if err != nil {
				job.Status = StatusFailed
				job.Error = err.Error()
				return err
			}

			// evaluate result with context
			g.logger.Info("Prompt evaluation...")
			rawEvalResult, err = promptLLM(
				g.logger,
				g.httpClient,
				Url,
				model,
				prompt,
				schemas["EvaluationSchema"])
			if err != nil {
				job.Status = StatusFailed
				job.Error = err.Error()
				return err
			}
			g.logger.Debug("Raw evaluation response", "evalRes", rawEvalResult)

			// parse evaluation
			err = json.Unmarshal([]byte(rawEvalResult), parsedEvalResult)
			if err != nil {
				job.Status = StatusFailed
				job.Error = err.Error()
				return err
			}
			g.logger.Debug("Parsed evaluation", "pass", parsedEvalResult.Pass, "feedback", parsedEvalResult.Feedback)

			FeedbackJSON, err := json.Marshal(parsedEvalResult.Feedback)
			if err != nil {
				job.Status = StatusFailed
				job.Error = err.Error()
				return err
			}
			stage.promptValues.EvaluationFeedback = string(FeedbackJSON)

		}
	}
	job.Status = StatusDone
	job.Result = rawGenResult

	g.logger.Debug("Job completed", "jobID", job.ID, "status", job.Status)
	Jobs.Store(job.ID, job)

	// Start a goroutine to log completed jobs every 30 seconds
	go func() {
		for {
			logCompletedJobs(g.logger, &Jobs)
			time.Sleep(30 * time.Second)
		}
	}()
	return nil
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
