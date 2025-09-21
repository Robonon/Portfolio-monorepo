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

	// JSON ArchitectureSchema in schemas.json
	ArchitectureSchema struct {
		Packages  []ArchitecturePackage `json:"packages"`
		Diagram   string                `json:"diagram"`
		Summary   string                `json:"summary"`
		Filepaths []string              `json:"filepaths"`
	}
)

var (
	interfacesSchema     = map[string]any{"type": "string"}
	unitTestsSchema      = map[string]any{"type": "string"}
	implementationSchema = map[string]any{"type": "string"}
	integrationSchema    = map[string]any{"type": "string"}
	documentationSchema  = map[string]any{"type": "string"}
)
