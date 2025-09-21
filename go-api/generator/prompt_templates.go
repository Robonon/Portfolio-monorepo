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
    - Design a modular, idiomatic Go software architecture for the requested functionality.
    - Define clear package boundaries and responsibilities.
    - Use idiomatic Go package naming conventions (short, lowercase, descriptive).
    - Organize code into logical packages (e.g., cmd/, internal/, pkg/, api/).
    - Ensure separation of concerns between business logic, API handlers, configuration, and utilities.
    - Structure the module to be stateless and container-ready.
    - Use environment variables for configuration.
    - Follow 12-factor app principles.
    - Make the architecture scalable, maintainable, and extensible.
    - Ensure components are loosely coupled and highly cohesive.
    - Facilitate testability by designing for dependency injection and clear interfaces.
    - Include health check endpoints and proper logging.
    - Only use standard library and well-known, maintained libraries.
    - Document the architecture with GoDoc comments and a README.md. and add necessary files to the output schema
    - Do not include implementation details—focus on high-level structure and interactions.
	- Do not create any Go code in this stage`,
		Task:           "Design a Go software architecture for the following scope:",
		Scope:          "architecture",
		ResponseSchema: "architectureSchema",
	}

	architectureEvaluation = promptTemplateValues{
		Role:               "Go Software Architect",
		Context:            "",
		EvaluationFeedback: "",
		Instructions: `
    - Evaluate the architecture for the following criteria:
    - Clarity: Are package and module responsibilities well-defined and separated?
    - Scalability: Can the architecture handle growth in users, data, or features?
    - Maintainability: Is the codebase easy to update, refactor, and extend?
    - Modularity: Are components loosely coupled and highly cohesive?
    - Testability: Is it easy to write automated tests for components?
    - Extensibility: Can new features be added without major changes?
    - Performance: Are there any bottlenecks or inefficiencies?
    - Security: Are there clear boundaries and controls for sensitive operations?
    - Idiomatic Go: Does the architecture follow Go conventions and best practices?
    - Documentation: Is the architecture well-documented and understandable?
      `,
		Task:           "Evaluate the Go software architecture for the following scope:",
		Scope:          "architecture",
		ResponseSchema: "evaluationSchema",
	}

	interfacesPrompt = promptTemplateValues{
		Role:               "Go Software Engineer",
		Context:            "",
		EvaluationFeedback: "",
		Instructions: `
    - Define a public api interface and types for the following scope.
    - Use idiomatic Go naming conventions.
    - Include comments for all exported items.
    - Ensure interfaces are small and focused.
    - Include examples for all public methods.
    - Document any non-obvious behavior or edge cases.
    - Follow best practices for Go code organization.
    - Ensure interfaces align with the overall architecture and design principles.
    `,
		Task:           "Define interfaces for the Go application in the following scope:",
		Scope:          "interfaces",
		ResponseSchema: "interfacesSchema",
	}

	interfacesEvaluation = promptTemplateValues{
		Role:               "Go Software Engineer",
		Context:            "",
		EvaluationFeedback: "",
		Instructions: `
    - Evaluate the Go interfaces for the following criteria:
    - Clarity and purpose: Is each interface name descriptive and does it clearly define its intended behavior?
    - Size and focus: Are interfaces small and focused (prefer single-method interfaces when possible)?
    - Idiomatic Go: Do interfaces follow Go naming conventions and best practices?
    - Documentation: Are all exported interfaces and methods documented with clear comments?
    - Extensibility: Can interfaces be easily extended or implemented by other types?
    - Alignment with architecture: Do interfaces fit well with the overall package/module design?
    - Testability: Do interfaces make it easy to write tests and mock implementations?
    - No stuttering: Are interface names free from repeating the package name (e.g., io.Reader, not io.IoReader)?
    - No unnecessary methods: Do interfaces avoid including methods that not all implementers would need?
	- Pass the evaluation if you find no issues or you deem it good enough.
    `,
		Task:           "Evaluate the Go interfaces for the following scope:",
		Scope:          "interfaces",
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
Do not return any formats that could interfere with parsing a JSON response, or handle them properly. Return your response as a JSON object with the following structure:

{{.ResponseSchema}}

Each field should contain the respective code or documentation. Do not include explanations or comments outside the code.
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
