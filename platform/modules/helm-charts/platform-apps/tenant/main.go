package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	gitea "code.gitea.io/sdk/gitea"
)

type Config struct {
	GiteaURL   string
	GiteaToken string
	Port       string
}

// type (
// 	argoCDApplication struct {
// 		APIVersion string   `yaml:"apiVersion"`
// 		Kind       string   `yaml:"kind"`
// 		Metadata   metadata `yaml:"metadata"`
// 		Spec       spec     `yaml:"spec"`
// 	}
// 	metadata struct {
// 		Name       string   `yaml:"name"`
// 		Namespace  string   `yaml:"namespace"`
// 		Finalizers []string `yaml:"finalizers,omitempty"`
// 	}
// 	spec struct {
// 		Project     string      `yaml:"project"`
// 		Source      source      `yaml:"source"`
// 		Destination destination `yaml:"destination"`
// 	}
// 	source struct {
// 		RepoURL        string `yaml:"repoURL"`
// 		TargetRevision string `yaml:"targetRevision"`
// 		Path           string `yaml:"path"`
// 	}
// 	destination struct {
// 		Server    string `yaml:"server"`
// 		Namespace string `yaml:"namespace"`
// 	}
// )

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

type TenantPageData struct {
	Title   string
	Tenants []*gitea.ContentsResponse
}

type tenantHandler struct {
	cfg    *Config
	logger *log.Logger
	client *gitea.Client
}

func (h *tenantHandler) handler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		tenants, resp, err := h.client.ListContents("gitea_admin", "gitea_admin/platform-ops", "tenants", "main")
		if err != nil {
			// http.Error(w, "Failed to fetch tenants", http.StatusInternalServerError)
			h.logger.Print("Error fetching tenants: ", err)
		}
		if tenants == nil || resp.StatusCode == http.StatusNotFound {
			tenants = []*gitea.ContentsResponse{}
		}

		data := TenantPageData{
			Title:   "Platform Tenants",
			Tenants: tenants,
		}

		tmpl, _ := template.ParseFiles("templates/tenant.html")
		tmpl.Execute(w, data)

	case http.MethodPost:
		tenant := r.FormValue("name")
		h.logger.Printf("Received request to create tenant: %s", tenant)

		yaml := argoAppYAML(tenant, "https://gitea.example.com/company/platform-ops.git", "helm", tenant)

		file := gitea.CreateFileOptions{
			Content: yaml,
		}
		h.client.CreateFile("company", "platform-ops", fmt.Sprintf("tenants/%s", tenant), file)
		http.Error(w, "POST method not implemented", http.StatusNotImplemented)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func argoAppYAML(name, repoURL, path, namespace string) string {
	return fmt.Sprintf(
		`apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: %s
  namespace: argocd
spec:
  project: default
  source:
    repoURL: %s
    targetRevision: HEAD
    path: %s
  destination:
    server: https://kubernetes.default.svc
    namespace: %s
`, name, repoURL, path, namespace)
}

func main() {
	cfg := &Config{
		GiteaURL:   getEnvOrDefault("GITEA_URL", "http://localhost:3000"),
		GiteaToken: getEnvOrDefault("GITEA_TOKEN", "your_gitea_token_here"),
		Port:       getEnvOrDefault("PORT", "8000"),
	}

	logger := log.Default()
	httpClient := http.DefaultClient

	giteaClient, err := gitea.NewClient(cfg.GiteaURL, gitea.SetHTTPClient(httpClient))
	if err != nil {
		logger.Fatal("Failed to create Gitea client:", err)
	}

	tenant := &tenantHandler{cfg: cfg, logger: logger, client: giteaClient}
	http.HandleFunc("/tenant", tenant.handler)

	fmt.Println("Server started at :" + cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		logger.Fatal("Failed to start server:", err)
	}
}
