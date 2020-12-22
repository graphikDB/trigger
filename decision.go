package eval

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/pkg/errors"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

var (
	ErrDecisionDenied   = errors.New("eval: evaluation = false")
	ErrEmptyExpressions = errors.New("eval: empty expressions")
)

// Decision is used to evaluate boolean expressions
type Decision struct {
	e          *cel.Env
	program    cel.Program
	expression string
}

// NewDecision creates a new Decision with the given boolean CEL expressions
func NewDecision(expression string) (*Decision, error) {
	if expression == "" {
		return nil, ErrEmptyExpressions
	}
	e, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("this", decls.NewMapType(decls.String, decls.Any)),
			decls.NewFunction("now",
				decls.NewOverload(
					"now",
					[]*expr.Type{},
					decls.Int,
				),
			),
			decls.NewFunction("uuid",
				decls.NewOverload(
					"uuid",
					[]*expr.Type{},
					decls.String,
				),
			),
			decls.NewFunction("sha1",
				decls.NewOverload(
					"sha1",
					[]*expr.Type{decls.String},
					decls.String,
				),
			),
			decls.NewFunction("sha256",
				decls.NewOverload(
					"sha256",
					[]*expr.Type{decls.String},
					decls.String,
				),
			),
			decls.NewFunction("base64Decode",
				decls.NewOverload(
					"base64Decode",
					[]*expr.Type{decls.String},
					decls.String,
				),
			),
			decls.NewFunction("base64Encode",
				decls.NewOverload(
					"base64Encode",
					[]*expr.Type{decls.String},
					decls.String,
				),
			),
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
			},
			&functions.Overload{
				Operator: "sha256",
				Function: defaultFuncMap["sha256"],
			},
			&functions.Overload{
				Operator: "base64Encode",
				Function: defaultFuncMap["base64Encode"],
			},
			&functions.Overload{
				Operator: "base64Decode",
				Function: defaultFuncMap["base64Decode"],
			},
			&functions.Overload{
				Operator: "uuid",
				Function: defaultFuncMap["uuid"],
			},
			&functions.Overload{
				Operator: "sha1",
				Function: defaultFuncMap["sha1"],
			}),
	)
	if err != nil {
		return nil, err
	}
	return &Decision{
		e:          e,
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
