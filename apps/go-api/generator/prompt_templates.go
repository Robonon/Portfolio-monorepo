package generator

import (
	"bytes"
	"html/template"
)

type (
	promptTemplateValues struct {
		Role               string
		Context            string
		EvaluationFeedback string
		Instructions       string
		Task               string
		Scope              string
		ResponseSchema     string
	}
)

var (
	architecturePrompt = promptTemplateValues{
		Role:               "Go Software Architect",
		Context:            "",
		EvaluationFeedback: "",
		Instructions: `
    - Design a modular, idiomatic Go module folder and file structure for the requested functionality.
    - Define clear package boundaries and responsibilities.
    - Use idiomatic Go package naming conventions (short, lowercase, descriptive).
    - Organize the folder structure into logical packages (e.g., cmd/, internal/, pkg/, api/).
    - Ensure separation of concerns between business logic, API handlers, configuration, and utilities.
    - Follow 12-factor app principles.
    - Make the folder structure scalable, maintainable, and extensible.
    - Do not include implementation details or Go code—focus only on folder structure and high-level interactions.
	- If there is feedback, choose one issue to address in the next iteration.
    `,
		Task:           "Design a Go module folder and file structure for the following scope:",
		Scope:          "architecture",
		ResponseSchema: "architectureSchema",
	}

	architectureEvaluation = promptTemplateValues{
		Role:               "Go Solutions Architect",
		Context:            "",
		EvaluationFeedback: "",
		Instructions: `
    - Evaluate the proposed Go module folder and file structure for the following criteria:
		- Scalability: Can the folder structure support increased features or modules without major changes?
		- Maintainability: Is the structure easy to update, refactor, and extend?
		- Modularity: Are components well-separated into distinct packages with clear responsibilities?
		- Idiomatic Go: Does the structure follow Go conventions, best practices?
		- Create one maxium of three pieces of feedback to improve the architecture or pass the evaluation.
		- You can set Pass to true if you deem the file and folder structure good enough to move on to creating the interfaces, even if there is feedback.
    `,
		Task:           "Evaluate the Go module folder and file structure is good enough for the following scope:",
		Scope:          "architecture",
		ResponseSchema: "evaluationSchema",
	}

	// implementationPrompt = promptTemplateValues{
	//}

	// implementationEvaluation = promptTemplateValues{
	// }

	unitTestsPrompt = promptTemplateValues{
		Role:               "Go Software Engineer",
		Context:            "",
		EvaluationFeedback: "",
		Instructions: `
    - Generate unit tests for the following scope.
    - Use idiomatic Go testing conventions.
    - Include table-driven tests where applicable.
    - Ensure tests cover all public methods and behaviors.
    - Document any non-obvious behavior or edge cases.
    - Follow best practices for Go code organization.
    `,
		Task:           "Generate unit tests for the Go application in the following scope:",
		Scope:          "unit_tests",
		ResponseSchema: "unitTestsSchema",
	}

	unitTestsEvaluation = promptTemplateValues{
		Role:               "Go Software Engineer",
		Context:            "",
		EvaluationFeedback: "",
		Instructions: `
    - Evaluate the Go unit tests for the following criteria:
    - Coverage: Do the tests cover all public methods and important code paths?
    - Clarity: Are the tests easy to read and understand?
    - Idiomatic Go: Do the tests follow Go naming conventions and best practices?
    - Documentation: Are all tests well-documented with clear comments?
    - Edge Cases: Do the tests cover edge cases and error conditions?
    - Performance: Do the tests run efficiently and not take excessive time?
    `,
		Task:           "Evaluate the Go unit tests for the following scope:",
		Scope:          "unit_tests",
		ResponseSchema: "evaluationSchema",
	}

	integrationPrompt = promptTemplateValues{
		Role:               "Go Software Engineer",
		Context:            "",
		EvaluationFeedback: "",
		Instructions: `
    - Define integration tests for the following scope.
    - Use idiomatic Go testing conventions.
    - Include table-driven tests where applicable.
    - Ensure tests cover all public methods and behaviors.
    - Document any non-obvious behavior or edge cases.
    - Follow best practices for Go code organization.
    `,
		Task:           "Define integration tests for the Go application in the following scope:",
		Scope:          "integration_tests",
		ResponseSchema: "integrationSchema",
	}

	integrationEvaluation = promptTemplateValues{
		Role:               "Go Software Engineer",
		Context:            "",
		EvaluationFeedback: "",
		Instructions: `
    - Evaluate the integration tests for the following criteria:
    - Coverage: Do the tests cover all public methods and important code paths?
    - Clarity: Are the tests easy to read and understand?
    - Idiomatic Go: Do the tests follow Go naming conventions and best practices?
    - Documentation: Are all tests well-documented with clear comments?
    - Edge Cases: Do the tests cover edge cases and error conditions?
    - Performance: Do the tests run efficiently and not take excessive time?
    `,
		Task:           "Evaluate the integration tests for the Go application in the following scope:",
		Scope:          "integration_tests",
		ResponseSchema: "evaluationSchema",
	}

	documentationPrompt = promptTemplateValues{
		Role:               "Go Software Engineer",
		Context:            "",
		EvaluationFeedback: "",
		Instructions: `
    - Define documentation for the following scope.
    - Use idiomatic Go documentation conventions.
    - Include examples where applicable.
    - Ensure documentation covers all public methods and behaviors.
    - Document any non-obvious behavior or edge cases.
    - Follow best practices for Go code organization.
    `,
		Task:           "Define documentation for the Go application in the following scope:",
		Scope:          "documentation",
		ResponseSchema: "documentationSchema",
	}

	documentationEvaluation = promptTemplateValues{
		Role:               "Go Software Engineer",
		Context:            "",
		EvaluationFeedback: "",
		Instructions: `
    - Evaluate the documentation for the following criteria:
    - Coverage: Does the documentation cover all public methods and important code paths?
    - Clarity: Is the documentation easy to read and understand?
    - Idiomatic Go: Does the documentation follow Go naming conventions and best practices?
    - Examples: Are there sufficient examples to illustrate usage?
    - Edge Cases: Does the documentation cover edge cases and non-obvious behavior?
    - Organization: Is the documentation well-organized and easy to navigate?
    `,
		Task:           "Evaluate the documentation for the Go application in the following scope:",
		Scope:          "documentation",
		ResponseSchema: "evaluationSchema",
	}
)

func NewPromptTemplateValues(
	role string,
	context string,
	evaluationFeedback string,
	instructions string,
	task string,
	scope string,
	responseSchema string,
) promptTemplateValues {
	return promptTemplateValues{
		Role:               role,
		Context:            context,
		EvaluationFeedback: evaluationFeedback,
		Instructions:       instructions,
		Task:               task,
		Scope:              scope,
		ResponseSchema:     responseSchema,
	}
}

const promptTmpl = `
You are an expert {{.Role}}.

## Context
{{.Context}}

## Previous Evaluation Feedback
{{.EvaluationFeedback}}

## Instructions
{{.Instructions}}

## Task
{{.Task}}

---
{{.Scope}}
---

## Output Format

{{.ResponseSchema}}

When including code in JSON, escape all newlines as \\n and all double quotes as \\\".
Do not return any formats that could interfere with parsing a JSON response, or handle them properly. Return your response as a JSON object with the following structure:
Do not include explanations or comments outside the code.
Do not wrap your response in markdown code fences or any extra formatting. Return only raw JSON.
`

func createPrompt(templateValues promptTemplateValues) (string, error) {
	t := template.Must(template.New("prompt").Parse(promptTmpl))
	var buf bytes.Buffer
	err := t.Execute(&buf, templateValues)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
