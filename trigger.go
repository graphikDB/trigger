package eval

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/pkg/errors"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"time"
)

// Trigger creates values as map[string]interface{} if it's decisider returns no errors against a Mapper
type Trigger struct {
	e          *cel.Env
	decision   *Decision
	program    cel.Program
	expression string
}

// NewTrigger creates a new trigger instance from the decision & trigger expressions
func NewTrigger(decision *Decision, triggerExpression string) (*Trigger, error) {
	if triggerExpression == "" {
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
		),
	)
	if err != nil {
		return nil, err
	}
	ast, iss := e.Compile(triggerExpression)
	if iss.Err() != nil {
		return nil, iss.Err()
	}
	program, err := e.Program(
		ast,
		cel.Functions(
			&functions.Overload{
				Operator: "now",
				Function: defaultFuncMap["now"],
			}),
	)
	if err != nil {
		return nil, err
	}
	return &Trigger{
		e:          e,
		decision:   decision,
		program:    program,
		expression: triggerExpression,
	}, nil
}

// Trigger executes it's decision against the Mapper and then overwrites the
func (t *Trigger) Trigger(data map[string]interface{}) (map[string]interface{}, error) {
	if err := t.decision.Eval(data); err == nil {
		out, _, err := t.program.Eval(map[string]interface{}{
			"this": data,
			"now":  time.Now().Unix(),
		})
		if err != nil {
			return nil, errors.Wrapf(err, "eval: failed to evaluate trigger (%s)", t.expression)
		}
		if patchFields, ok := out.Value().(map[ref.Val]ref.Val); ok {
			newData := map[string]interface{}{}
			for k, v := range patchFields {
				newData[k.Value().(string)] = v.Value()
			}
			return newData, nil
		}
		if patchFields, ok := out.Value().(map[string]interface{}); ok {
			return patchFields, nil
		}
		if patchFields, ok := out.Value().(map[string]string); ok {
			newData := map[string]interface{}{}
			for k, v := range patchFields {
				newData[k] = v
			}
			return newData, nil
		}
	}
	return map[string]interface{}{}, nil
}

// Expression returns the triggers raw CEL expressions
func (e *Trigger) Expression() string {
	return e.expression
}
