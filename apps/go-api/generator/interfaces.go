package generator

// JSON InterfaceSchema in schemas.json
type InterfaceResponse struct {
	FilePath string `json:"filePath"`
	Code     string `json:"code"`
}

var (
	EvaluationSchema = ``

	Prompt = promptTemplateValues{
		Role:               "Go Software Engineer",
		Context:            "",
		EvaluationFeedback: "",
		Instructions: `
    - Define a single API interface for the following package.
    - Use idiomatic Go naming conventions.
    - Include comments for all exported items.
    - Ensure the interface is small and focused.
    - Include examples for all public methods.
    - Document any non-obvious behavior or edge cases.
    - Follow best practices for Go code organization.
    - Output the code as a string in the Code field in the schema, and the file path, including the file name and file ending, in the filepath field
    `,
		Task:           "Define a single API interface for the Go application in the following scope:",
		Scope:          "interfaces",
		ResponseSchema: "interface response",
	}
	EvaluationPrompt = promptTemplateValues{
		Role:               "Go Software Engineer",
		Context:            "",
		EvaluationFeedback: "",
		Instructions: `
    - Evaluate the provided Go API interface for the following criteria:
    - Clarity and purpose: Is the interface name descriptive and does it clearly define its intended behavior?
    - Size and focus: Is the interface small and focused (prefer single-method interfaces when possible)?
    - Idiomatic Go: Does the interface follow Go naming conventions and best practices?
    - Documentation: Are all exported interfaces and methods documented with clear comments?
    - Extensibility: Can the interface be easily extended or implemented by other types?
    - Alignment with architecture: Does the interface fit well with the overall package/module design?
    - Testability: Does the interface make it easy to write tests and mock implementations?
    - No stuttering: Is the interface name free from repeating the package name (e.g., io.Reader, not io.IoReader)?
    - No unnecessary methods: Does the interface avoid including methods that not all implementers would need?
    - Pass the evaluation if you find no issues or you deem it good enough.
    `,
		Task:           "Evaluate the provided Go API interface for the following scope:",
		Scope:          "interfaces",
		ResponseSchema: "evaluationSchema",
	}
)
