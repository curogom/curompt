package collector

import (
	"context"

	"github.com/curogom/curompt/internal/model"
)

// Collector collects prompts from various sources
type Collector interface {
	// Collect collects prompts and returns collected prompts
	Collect(ctx context.Context) ([]*model.CollectedPrompt, error)

	// Name returns the collector name
	Name() string
}
