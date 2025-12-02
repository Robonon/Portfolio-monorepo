package generator

import (
	"encoding/json"
)

type (
	Job struct {
		ID     string
		Status string
		Error  string
		Stages []stage
	}
	stage struct {
		name                   string
		promptValues           promptTemplateValues
		evaluationPromptValues promptTemplateValues
		responseSchema         json.RawMessage
	}
)

func newJob() (*Job, error) {
	schemas, err := loadSchemas("./schemas.json")
	if err != nil {
		return nil, err
	}
	return &Job{
		ID:     generateID(),
		Status: StatusRunning,
		Stages: []stage{
			{
				name: "architecture",
				promptValues: NewPromptTemplateValues(
					architecturePrompt.Role,
					architecturePrompt.Context,
					architecturePrompt.EvaluationFeedback,
					architecturePrompt.Instructions,
					architecturePrompt.Task,
					architecturePrompt.Scope,
					string(schemas["ArchitectureSchema"]),
				),
				evaluationPromptValues: NewPromptTemplateValues(
					architectureEvaluation.Role,
					architectureEvaluation.Context,
					architectureEvaluation.EvaluationFeedback,
					architectureEvaluation.Instructions,
					architectureEvaluation.Task,
					architectureEvaluation.Scope,
					string(schemas["EvaluationSchema"]),
				),
				responseSchema: schemas["ArchitectureSchema"],
			},
			{
				name: "interfaces",
				promptValues: NewPromptTemplateValues(
					Prompt.Role,
					Prompt.Context,
					Prompt.EvaluationFeedback,
					Prompt.Instructions,
					Prompt.Task,
					Prompt.Scope,
					string(schemas["InterfaceSchema"]),
				),
				evaluationPromptValues: NewPromptTemplateValues(
					Prompt.Role,
					Prompt.Context,
					Prompt.EvaluationFeedback,
					Prompt.Instructions,
					Prompt.Task,
					Prompt.Scope,
					string(schemas["EvaluationSchema"]),
				),
				responseSchema: schemas["InterfaceSchema"],
			},
			{
				name: "unittests",
				promptValues: NewPromptTemplateValues(
					unitTestsPrompt.Role,
					unitTestsPrompt.Context,
					unitTestsPrompt.EvaluationFeedback,
					unitTestsPrompt.Instructions,
					unitTestsPrompt.Task,
					unitTestsPrompt.Scope,
					string(schemas["UnitTestsSchema"]),
				),
				evaluationPromptValues: NewPromptTemplateValues(
					unitTestsEvaluation.Role,
					unitTestsEvaluation.Context,
					unitTestsEvaluation.EvaluationFeedback,
					unitTestsEvaluation.Instructions,
					unitTestsEvaluation.Task,
					unitTestsEvaluation.Scope,
					string(schemas["EvaluationSchema"]),
				),
				responseSchema: schemas["UnitTestsSchema"],
			},
			// {name: "implementation", promptValues: implementationPrompt, evaluationPromptValues: implementationEvaluation},
			{
				name: "integration",
				promptValues: NewPromptTemplateValues(
					integrationPrompt.Role,
					integrationPrompt.Context,
					integrationPrompt.EvaluationFeedback,
					integrationPrompt.Instructions,
					integrationPrompt.Task,
					integrationPrompt.Scope,
					string(schemas["IntegrationSchema"]),
				),
				evaluationPromptValues: NewPromptTemplateValues(
					integrationEvaluation.Role,
					integrationEvaluation.Context,
					integrationEvaluation.EvaluationFeedback,
					integrationEvaluation.Instructions,
					integrationEvaluation.Task,
					integrationEvaluation.Scope,
					string(schemas["EvaluationSchema"]),
				),
				responseSchema: schemas["IntegrationSchema"],
			},
			{
				name: "documentation",
				promptValues: NewPromptTemplateValues(
					documentationPrompt.Role,
					documentationPrompt.Context,
					documentationPrompt.EvaluationFeedback,
					documentationPrompt.Instructions,
					documentationPrompt.Task,
					documentationPrompt.Scope,
					string(schemas["DocumentationSchema"]),
				),
				evaluationPromptValues: NewPromptTemplateValues(
					documentationEvaluation.Role,
					documentationEvaluation.Context,
					documentationEvaluation.EvaluationFeedback,
					documentationEvaluation.Instructions,
					documentationEvaluation.Task,
					documentationEvaluation.Scope,
					string(schemas["EvaluationSchema"]),
				),
				responseSchema: schemas["DocumentationSchema"],
			},
		},
	}, nil
}
