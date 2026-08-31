package filtergraph

import (
	"fmt"
	"log/slog"
)

// Pass represents an optimization, validation, or transformation pass executed on a filtergraph.
type Pass interface {
	Name() string
	Run(g *Graph) error
}

// Pipeline orchestrates sequential execution of graph passes.
type Pipeline struct {
	passes []Pass
	logger *slog.Logger
}

// NewPipeline creates a new optimization pass pipeline.
func NewPipeline(logger *slog.Logger) *Pipeline {
	if logger == nil {
		logger = slog.Default()
	}
	return &Pipeline{
		passes: make([]Pass, 0),
		logger: logger,
	}
}

// Add appends one or more passes to the pipeline.
func (p *Pipeline) Add(passes ...Pass) *Pipeline {
	p.passes = append(p.passes, passes...)
	return p
}

// Execute runs all registered passes sequentially against the graph.
func (p *Pipeline) Execute(g *Graph) error {
	for _, pass := range p.passes {
		p.logger.Debug("running filtergraph pass", slog.String("pass", pass.Name()))
		if err := pass.Run(g); err != nil {
			return fmt.Errorf("filtergraph pass %q failed: %w", pass.Name(), err)
		}
	}
	return nil
}

// DefaultPipeline returns a standard optimization pipeline containing Validation, AutoSplit, and DCE passes.
func DefaultPipeline(logger *slog.Logger) *Pipeline {
	return NewPipeline(logger).
		Add(&ValidatePass{}).
		Add(&AutoSplitPass{}).
		Add(&DeadCodeEliminationPass{})
}
