package eval

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/pkg/errors"
	"sort"
	"sync"
)

type DecisionType int

const (
	AllTrue DecisionType = 0
	AnyTrue DecisionType = 1
)

var (
	ErrDecisionDenied = errors.New("eval: evaluation = false")
	ErrNoExpressions  = errors.New("eval: no expressions")
)

// Decision is used to evaluate boolean expressions
type Decision struct {
	e        *cel.Env
	programs map[string]cel.Program
	mu       sync.RWMutex
	dtype    DecisionType
}

// NewDecision creates a new Decision with the given boolean CEL expressions
func NewDecision(dtype DecisionType, expressions []string) (*Decision, error) {
	if len(expressions) == 0 {
		return nil, ErrNoExpressions
	}
	e, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("this", decls.NewMapType(decls.String, decls.Any)),
		),
	)
	if err != nil {
		return nil, err
	}
	var programs = map[string]cel.Program{}
	for _, expression := range expressions {
		if expression == "" {
			return nil, errors.New("empty expression")
		}
		ast, iss := e.Compile(expression)
		if iss.Err() != nil {
			return nil, iss.Err()
		}
		program, err := e.Program(ast)
		if err != nil {
			return nil, err
		}
		programs[expression] = program
	}
	return &Decision{
		e:        e,
		programs: programs,
		mu:       sync.RWMutex{},
		dtype:    dtype,
	}, nil
}

// AddExpression adds an expression to the decision tree
func (n *Decision) AddExpression(expression string) error {
	if expression == "" {
		return errors.New("eval: empty expression")
	}
	ast, iss := n.e.Compile(expression)
	if iss.Err() != nil {
		return iss.Err()
	}
	program, err := n.e.Program(ast)
	if err != nil {
		return err
	}
	n.mu.Lock()
	n.programs[expression] = program
	n.mu.Unlock()
	return nil
}

// Eval evaluates the boolean CEL expressions against the Mapper
func (n *Decision) Eval(data map[string]interface{}) error {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if len(n.programs) == 0 {
		return ErrNoExpressions
	}
	if n.dtype == AllTrue {
		for exp, program := range n.programs {
			out, _, err := program.Eval(map[string]interface{}{
				"this": data,
			})
			if err != nil {
				return errors.Wrapf(err, "eval: failed to evaluate expression (%s)", exp)
			}
			if val, ok := out.Value().(bool); !ok || !val {
				return ErrDecisionDenied
			}
		}
	} else if n.dtype == AnyTrue {
		for exp, program := range n.programs {
			out, _, err := program.Eval(map[string]interface{}{
				"this": data,
			})
			if err != nil {
				return errors.Wrapf(err, "eval: failed to evaluate expression (%s)", exp)
			}
			if val, ok := out.Value().(bool); ok && val {
				return nil
			}
		}
		return ErrDecisionDenied
	}
	return nil
}

// Expressions returns the decsions raw expressions
func (e *Decision) Expressions() []string {
	var exp []string
	for ex, _ := range e.programs {
		exp = append(exp, ex)
	}
	sort.Strings(exp)
	return exp
}
