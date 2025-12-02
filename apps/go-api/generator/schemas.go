package generator

type (
	EvaluationFeedback struct {
		Issue      string `json:"issue"`
		Suggestion string `json:"suggestion"`
	}

	// JSON EvaluationSchema in schemas.json
	EvaluationResult struct {
		Feedback []EvaluationFeedback `json:"feedback"`
		Pass     bool                 `json:"pass"`
	}

	ArchitecturePackage struct {
		Name           string `json:"name"`
		Responsibility string `json:"responsibility"`
	}

	// JSON ArchitectureResponse in schemas.json
	ArchitectureResponse struct {
		Packages []ArchitecturePackage `json:"packages"`
		// Diagram   string                `json:"diagram"`
		Summary        string   `json:"summary"`
		DirectoryPaths []string `json:"directoryPaths"`
	}

	UnitTestSchema struct {
		Description     string `json:"description"`
		MethodSignature string `json:"methodSignature"`
		Code            string `json:"code"`
	}

	ImplementationSchema struct {
		Description     string `json:"description"`
		MethodSignature string `json:"methodSignature"`
		Code            string `json:"code"`
	}

	// integration

	DocumentationSchema struct {
		Summary     string                `json:"summary"`
		Packages    []ArchitecturePackage `json:"packages"`
		Interfaces  []string              `json:"interfaces"`
		DocFilePath string                `json:"docFilePath"`
	}
)

var (
	integrationSchema   = map[string]any{"type": "string"} // Will this be outside resources like databases or external APIs
	documentationSchema = map[string]any{"type": "string"}
)
