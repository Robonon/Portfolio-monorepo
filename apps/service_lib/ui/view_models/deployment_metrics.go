package viewmodels

type DeploymentMetrics struct {
	TotalPods      int
	PodsPerService map[string]int
}
