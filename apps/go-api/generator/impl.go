package generator

import (
	"api/config"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
)

type (
	generatorImpl struct {
		logger     *slog.Logger
		httpClient *http.Client
		config     *config.Config
	}
)

const (
	StatusRunning string = "running"
	StatusDone    string = "done"
	StatusFailed  string = "failed"
)

func (g *generatorImpl) GenerateModule(input GeneratorInput) error {
	var (
		jobs sync.Map
	)

	// Create and store new job
	job, err := newJob()
	if err != nil {
		return err
	}
	jobs.Store(job.ID, job)
	g.logger.Info("Created new job", "jobID", job.ID)

	// Init git repo
	err = initGitRepo(g.logger, g.config.OutputDir)
	if err != nil {
		job.Status = StatusFailed
		job.Error = err.Error()
		return err
	}

	job.Status = StatusDone
	jobs.Store(job.ID, job)
	g.logger.Info("Job completed", "jobID", job.ID, "status", job.Status)

	// Start a goroutine to log completed jobs every 30 seconds
	// go func() {
	// 	for {
	// 		logCompletedJobs(g.logger, &jobs)
	// 		time.Sleep(30 * time.Second)
	// 	}
	// }()
	return nil
}

func (g *generatorImpl) ZipOutput() error {
	zipPath := g.config.OutputDir + ".zip"
	err := ZipDirectory(g.config.OutputDir, zipPath)
	if err != nil {
		return err
	}
	return nil
}

func evaluate(
	logger *slog.Logger,
	client *http.Client,
	stage *stage,
	url string,
	model string,
	evaluation *EvaluationResult) error {

	logger.Debug("Create evaluation prompt")
	prompt, err := createPrompt(stage.evaluationPromptValues)
	if err != nil {
		return err
	}

	logger.Debug("Prompting ollama for evaluation...")
	rawEvalResult, err := promptOllama(
		logger,
		client,
		url,
		model,
		prompt,
		json.RawMessage(stage.evaluationPromptValues.ResponseSchema))
	if err != nil {
		return err
	}
	logger.Debug("Raw evaluation response", "evalRes", rawEvalResult)

	// parse evaluation
	err = json.Unmarshal([]byte(rawEvalResult), evaluation)
	if err != nil {
		return err
	}
	logger.Debug("Parsed evaluation", "pass", evaluation.Pass, "feedback", evaluation.Feedback)

	FeedbackJSON, err := json.Marshal(evaluation.Feedback)
	if err != nil {
		return err
	}
	stage.promptValues.EvaluationFeedback = string(FeedbackJSON)
	return nil
}

func architectureStage(
	logger *slog.Logger,
	httpClient *http.Client,
	archRes *ArchitectureResponse,
	stage *stage,
	input *GeneratorInput,
	url, outputDir string) error {

	parsedEvalResult := newEvaluationResult()
	parsedEvalResult.Pass = false
	stage.promptValues.Scope = input.Scope
	stage.evaluationPromptValues.Scope = input.Scope

	for !parsedEvalResult.Pass {

		prompt, err := createPrompt(stage.promptValues)
		if err != nil {
			return err
		}

		rawGenResult, err := promptOllama(
			logger,
			httpClient,
			url,
			"gemma3",
			prompt,
			stage.responseSchema)
		if err != nil {
			return err
		}
		// parse result
		err = json.Unmarshal([]byte(rawGenResult), &archRes)
		if err != nil {
			return err
		}

		stage.promptValues.Context = rawGenResult
		stage.evaluationPromptValues.Context = rawGenResult

		if input.SkipEval {
			logger.Info("Skipping evaluation...")
			parsedEvalResult.Pass = true
			continue
		}

		logger.Info("Evaluating result...")
		err = evaluate(logger, httpClient, stage, url, "gemma3", parsedEvalResult)
		if err != nil {
			return err
		}

		feedbackJSON, err := json.Marshal(parsedEvalResult.Feedback)
		if err != nil {
			return err
		}

		stage.promptValues.EvaluationFeedback = string(feedbackJSON)
	}

	if len(archRes.DirectoryPaths) > 0 {
		// creates folders and returns updated context
		logger.Info("Creating folders for architecture...")
		err := createFolders(logger, &archRes.DirectoryPaths, outputDir)
		if err != nil {
			return err
		}
	}
	// commit files
	logger.Info("Committing changes to git...")
	err := commitEverything(logger, stage.name, outputDir)
	if err != nil {
		return err
	}

	return nil
}

func interfaceStage(
	logger *slog.Logger,
	httpClient *http.Client,
	archRes *ArchitectureResponse,
	ifaceRes *InterfaceResponse,
	stage *stage,
	input *GeneratorInput,
	url, outputDir string) error {

	parsedEvalResult := newEvaluationResult()
	parsedEvalResult.Pass = false
	stage.promptValues.Scope = input.Scope
	stage.evaluationPromptValues.Scope = input.Scope

	for i := range archRes.DirectoryPaths {

		parsedEvalResult.Pass = false

		for !parsedEvalResult.Pass {

			dir := archRes.DirectoryPaths[i]

			stage.promptValues.Context = dir
			// Create prompt
			prompt, err := createPrompt(stage.promptValues)
			if err != nil {
				return err
			}
			logger.Debug("prompt", "prompt", prompt)

			// Send prompt
			rawGenResult, err := promptOllama(
				logger,
				httpClient,
				url,
				"gemma3",
				prompt,
				stage.responseSchema)
			if err != nil {
				return err
			}

			// Parse result
			err = json.Unmarshal([]byte(rawGenResult), &ifaceRes)
			if err != nil {
				return err
			}

			// Update context
			stage.promptValues.Context = rawGenResult
			stage.evaluationPromptValues.Context = rawGenResult

			if input.SkipEval {
				logger.Info("Skipping evaluation...")
				parsedEvalResult.Pass = true
				continue
			}

			logger.Info("Evaluating result...")
			err = evaluate(logger, httpClient, stage, url, "gemma3", parsedEvalResult)
			if err != nil {
				return err
			}
			logger.Info("Evaluation", "pass", parsedEvalResult.Pass)

			logger.Info("Updating prompt feedback...")
			FeedbackJSON, err := json.Marshal(parsedEvalResult.Feedback)
			if err != nil {
				return err
			}

			stage.promptValues.EvaluationFeedback = string(FeedbackJSON)
		}

		createFileWithContent(logger, ifaceRes.FilePath, ifaceRes.Code, outputDir)
		// commit files
		logger.Info("Committing changes to git...")
		err := commitEverything(logger, stage.name, outputDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func unitTestStage(
	logger *slog.Logger,
	httpClient *http.Client,
	archRes *ArchitectureResponse,
	ifaceRes *InterfaceResponse,
	stage *stage,
	input *GeneratorInput,
	url, outputDir string) {

}
