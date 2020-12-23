package trigger

import (
	"github.com/google/cel-go/cel"
	"github.com/pkg/errors"
)

var (
	ErrDecisionDenied   = errors.New("trigger: evaluation = false")
	ErrEmptyExpressions = errors.New("trigger: empty expressions")
)

// Decision is used to evaluate boolean expressions
type Decision struct {
	program    cel.Program
	expression string
}

// NewDecision creates a new Decision with the given boolean CEL expressions
func NewDecision(expression string) (*Decision, error) {
	if expression == "" {
		return nil, ErrEmptyExpressions
	}
	program, err := globalEnv.Program(expression)
	if err != nil {
		return nil, err
	}
	return &Decision{
		program:    program,
		expression: expression,
	}, nil
}

// Eval evaluates the boolean CEL expressions against the Mapper
func (n *Decision) Eval(data map[string]interface{}) error {
	out, _, err := n.program.Eval(map[string]interface{}{
		"this": data,
	})
	if err != nil {
		return errors.Wrapf(err, "trigger: failed to evaluate decision (%s)", n.expression)
	}
	if val, ok := out.Value().(bool); !ok || !val {
		return ErrDecisionDenied
	}
	return nil
}

// Expressions returns the decsions raw expression
func (e *Decision) Expression() string {
	return e.expression
}
