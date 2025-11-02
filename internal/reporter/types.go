package reporter

import "github.com/curogom/curo-prompt/internal/evaluator"

// Reporter generates reports from evaluation results
type Reporter interface {
	// Generate generates a report from evaluation results
	Generate(result *evaluator.EvaluationResult) (string, error)

	// Name returns the reporter format name
	Name() string
}
