package eval

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/pkg/errors"
)

type DecisionType int

const (
	AllTrue DecisionType = 0
	AnyTrue DecisionType = 1
)

var (
	ErrDecisionDenied   = errors.New("eval: evaluation = false")
	ErrEmptyExpressions = errors.New("eval: empty expressions")
)

// Decision is used to evaluate boolean expressions
type Decision struct {
	e          *cel.Env
	program    cel.Program
	dtype      DecisionType
	expression string
}

// NewDecision creates a new Decision with the given boolean CEL expressions
func NewDecision(dtype DecisionType, expression string) (*Decision, error) {
	if expression == "" {
		return nil, ErrEmptyExpressions
	}
	e, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("this", decls.NewMapType(decls.String, decls.Any)),
		),
	)
	if err != nil {
		return nil, err
	}
	if expression == "" {
		return nil, errors.New("empty expression")
	}
	ast, iss := e.Compile(expression)
	if iss.Err() != nil {
		return nil, iss.Err()
	}
	program, err := e.Program(ast,
		cel.Functions(
			&functions.Overload{
				Operator: "now",
				Function: defaultFuncMap["now"],
			}),
	)
	if err != nil {
		return nil, err
	}
	return &Decision{
		e:          e,
		program:    program,
		dtype:      dtype,
		expression: expression,
	}, nil
}

// Eval evaluates the boolean CEL expressions against the Mapper
func (n *Decision) Eval(data map[string]interface{}) error {
	out, _, err := n.program.Eval(map[string]interface{}{
		"this": data,
	})
	if err != nil {
		return errors.Wrapf(err, "eval: failed to evaluate decision (%s)", n.expression)
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
