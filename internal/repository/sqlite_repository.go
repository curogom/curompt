package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"github.com/curogom/curo-prompt/internal/model"
	"github.com/curogom/curo-prompt/internal/parser"
)

// sqliteRepository implements PromptRepository using SQLite
type sqliteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new SQLite repository
func NewSQLiteRepository(dbPath string) (PromptRepository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	repo := &sqliteRepository{db: db}

	if err := repo.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

// initSchema creates the database schema
func (r *sqliteRepository) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS prompts (
		id TEXT PRIMARY KEY,
		tool TEXT NOT NULL,
		raw_prompt TEXT NOT NULL,
		role TEXT,
		inputs TEXT,
		invariants TEXT,
		output_format TEXT,
		timestamp INTEGER NOT NULL,
		command TEXT,
		working_dir TEXT,
		metadata TEXT,
		created_at INTEGER NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_tool ON prompts(tool);
	CREATE INDEX IF NOT EXISTS idx_timestamp ON prompts(timestamp DESC);
	`

	_, err := r.db.Exec(schema)
	return err
}

// Save saves a collected prompt
func (r *sqliteRepository) Save(ctx context.Context, prompt *model.CollectedPrompt) error {
	var role, inputs, invariants, outputFormat string
	var inputsJSON, invariantsJSON, outputFormatJSON []byte
	var metadataJSON []byte

	if prompt.Prompt != nil {
		role = prompt.Prompt.Role

		var err error
		if len(prompt.Prompt.Inputs) > 0 {
			inputsJSON, err = json.Marshal(prompt.Prompt.Inputs)
			if err == nil {
				inputs = string(inputsJSON)
			}
		}

		if len(prompt.Prompt.Invariants) > 0 {
			invariantsJSON, err = json.Marshal(prompt.Prompt.Invariants)
			if err == nil {
				invariants = string(invariantsJSON)
			}
		}

		if len(prompt.Prompt.OutputFormat) > 0 {
			outputFormatJSON, err = json.Marshal(prompt.Prompt.OutputFormat)
			if err == nil {
				outputFormat = string(outputFormatJSON)
			}
		}
	}

	if len(prompt.Metadata) > 0 {
		var err error
		metadataJSON, err = json.Marshal(prompt.Metadata)
		if err == nil {
			// metadata는 JSON 문자열로 저장
		} else {
			metadataJSON = []byte("{}")
		}
	} else {
		metadataJSON = []byte("{}")
	}

	query := `
		INSERT OR REPLACE INTO prompts (
			id, tool, raw_prompt, role, inputs, invariants, output_format,
			timestamp, command, working_dir, metadata, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now().Unix()
	_, err := r.db.ExecContext(ctx, query,
		prompt.ID,
		prompt.Tool,
		prompt.RawPrompt,
		role,
		inputs,
		invariants,
		outputFormat,
		prompt.Timestamp,
		prompt.Command,
		prompt.WorkingDir,
		string(metadataJSON),
		now,
	)

	return err
}

// FindByID finds a prompt by ID
func (r *sqliteRepository) FindByID(ctx context.Context, id string) (*model.CollectedPrompt, error) {
	query := `SELECT id, tool, raw_prompt, role, inputs, invariants, output_format,
	          timestamp, command, working_dir, metadata
	          FROM prompts WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)

	var collectedPrompt model.CollectedPrompt
	var role, inputs, invariants, outputFormat, metadata string

	err := row.Scan(
		&collectedPrompt.ID,
		&collectedPrompt.Tool,
		&collectedPrompt.RawPrompt,
		&role,
		&inputs,
		&invariants,
		&outputFormat,
		&collectedPrompt.Timestamp,
		&collectedPrompt.Command,
		&collectedPrompt.WorkingDir,
		&metadata,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("prompt not found: %s", id)
		}
		return nil, err
	}

	// Parsed prompt 재구성
	if role != "" || inputs != "" || invariants != "" || outputFormat != "" {
		collectedPrompt.Prompt = &parser.Prompt{
			Raw: collectedPrompt.RawPrompt,
		}
		collectedPrompt.Prompt.Role = role

		if inputs != "" {
			if err := json.Unmarshal([]byte(inputs), &collectedPrompt.Prompt.Inputs); err != nil {
				// If invalid, Inputs will remain empty - this is acceptable
				_ = err
			}
		}
		if invariants != "" {
			if err := json.Unmarshal([]byte(invariants), &collectedPrompt.Prompt.Invariants); err != nil {
				// If invalid, Invariants will remain empty - this is acceptable
				_ = err
			}
		}
		if outputFormat != "" {
			if err := json.Unmarshal([]byte(outputFormat), &collectedPrompt.Prompt.OutputFormat); err != nil {
				// If invalid, OutputFormat will remain empty - this is acceptable
				_ = err
			}
		}
	}

	// Metadata 재구성
	if metadata != "" {
		if err := json.Unmarshal([]byte(metadata), &collectedPrompt.Metadata); err != nil {
			// If invalid, Metadata will remain empty - this is acceptable
			_ = err
		}
	}

	return &collectedPrompt, nil
}

// FindByTool finds all prompts collected by a specific tool
func (r *sqliteRepository) FindByTool(ctx context.Context, tool string) ([]*model.CollectedPrompt, error) {
	query := `SELECT id, tool, raw_prompt, role, inputs, invariants, output_format,
	          timestamp, command, working_dir, metadata
	          FROM prompts WHERE tool = ? ORDER BY timestamp DESC`

	rows, err := r.db.QueryContext(ctx, query, tool)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prompts []*model.CollectedPrompt

	for rows.Next() {
		var p model.CollectedPrompt
		var role, inputs, invariants, outputFormat, metadata string

		err := rows.Scan(
			&p.ID,
			&p.Tool,
			&p.RawPrompt,
			&role,
			&inputs,
			&invariants,
			&outputFormat,
			&p.Timestamp,
			&p.Command,
			&p.WorkingDir,
			&metadata,
		)
		if err != nil {
			return nil, err
		}

		// Parsed prompt 재구성
		if role != "" || inputs != "" {
			p.Prompt = &parser.Prompt{Raw: p.RawPrompt, Role: role}
			if inputs != "" {
				if err := json.Unmarshal([]byte(inputs), &p.Prompt.Inputs); err != nil {
					// If invalid, Inputs will remain empty - this is acceptable
					_ = err
				}
			}
			if invariants != "" {
				if err := json.Unmarshal([]byte(invariants), &p.Prompt.Invariants); err != nil {
					// If invalid, Invariants will remain empty - this is acceptable
					_ = err
				}
			}
			if outputFormat != "" {
				if err := json.Unmarshal([]byte(outputFormat), &p.Prompt.OutputFormat); err != nil {
					// If invalid, OutputFormat will remain empty - this is acceptable
					_ = err
				}
			}
		}

		if metadata != "" {
			if err := json.Unmarshal([]byte(metadata), &p.Metadata); err != nil {
				// If invalid, Metadata will remain empty - this is acceptable
				_ = err
			}
		}

		prompts = append(prompts, &p)
	}

	return prompts, rows.Err()
}

// FindRecent finds recent prompts (limit)
func (r *sqliteRepository) FindRecent(ctx context.Context, limit int) ([]*model.CollectedPrompt, error) {
	query := `SELECT id, tool, raw_prompt, role, inputs, invariants, output_format,
	          timestamp, command, working_dir, metadata
	          FROM prompts ORDER BY timestamp DESC LIMIT ?`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prompts []*model.CollectedPrompt

	for rows.Next() {
		var p model.CollectedPrompt
		var role, inputs, invariants, outputFormat, metadata string

		err := rows.Scan(
			&p.ID,
			&p.Tool,
			&p.RawPrompt,
			&role,
			&inputs,
			&invariants,
			&outputFormat,
			&p.Timestamp,
			&p.Command,
			&p.WorkingDir,
			&metadata,
		)
		if err != nil {
			return nil, err
		}

		// Parsed prompt 재구성
		if role != "" || inputs != "" {
			p.Prompt = &parser.Prompt{Raw: p.RawPrompt, Role: role}
			if inputs != "" {
				if err := json.Unmarshal([]byte(inputs), &p.Prompt.Inputs); err != nil {
					// If invalid, Inputs will remain empty - this is acceptable
					_ = err
				}
			}
			if invariants != "" {
				if err := json.Unmarshal([]byte(invariants), &p.Prompt.Invariants); err != nil {
					// If invalid, Invariants will remain empty - this is acceptable
					_ = err
				}
			}
			if outputFormat != "" {
				if err := json.Unmarshal([]byte(outputFormat), &p.Prompt.OutputFormat); err != nil {
					// If invalid, OutputFormat will remain empty - this is acceptable
					_ = err
				}
			}
		}

		if metadata != "" {
			if err := json.Unmarshal([]byte(metadata), &p.Metadata); err != nil {
				// If invalid, Metadata will remain empty - this is acceptable
				_ = err
			}
		}

		prompts = append(prompts, &p)
	}

	return prompts, rows.Err()
}

// Close closes the repository connection
func (r *sqliteRepository) Close() error {
	return r.db.Close()
}
