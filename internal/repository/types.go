package repository

import (
	"context"

	"github.com/curogom/curompt/internal/model"
)

// PromptRepository handles prompt data persistence
type PromptRepository interface {
	// Save saves a collected prompt
	Save(ctx context.Context, prompt *model.CollectedPrompt) error

	// FindByID finds a prompt by ID
	FindByID(ctx context.Context, id string) (*model.CollectedPrompt, error)

	// FindByTool finds all prompts collected by a specific tool
	FindByTool(ctx context.Context, tool string) ([]*model.CollectedPrompt, error)

	// FindRecent finds recent prompts (limit)
	FindRecent(ctx context.Context, limit int) ([]*model.CollectedPrompt, error)

	// FindAll returns every stored prompt
	FindAll(ctx context.Context) ([]*model.CollectedPrompt, error)

	// Close closes the repository connection
	Close() error
}
