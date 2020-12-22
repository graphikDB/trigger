package eval

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types/ref"
	"github.com/pkg/errors"
	"sort"
	"sync"
	"time"
)

// Trigger creates values as map[string]interface{} if it's decisider returns no errors against a Mapper
type Trigger struct {
	e        *cel.Env
	decision *Decision
	programs map[string]cel.Program
	mu       sync.RWMutex
}

// NewTrigger creates a new trigger instance from the decision & trigger expressions
func NewTrigger(decision *Decision, triggerExpressions []string) (*Trigger, error) {
	if len(triggerExpressions) == 0 {
		return nil, ErrNoExpressions
	}
	e, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("this", decls.NewMapType(decls.String, decls.Any)),
			decls.NewVar("now", decls.Int),
		),
	)
	if err != nil {
		return nil, err
	}
	var programs = map[string]cel.Program{}
	for _, expression := range triggerExpressions {
		if expression == "" {
			return nil, errors.New("eval: empty trigger expression")
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
	return &Trigger{
		e:        e,
		decision: decision,
		programs: programs,
		mu:       sync.RWMutex{},
	}, nil
}

// Trigger executes it's decision against the Mapper and then overwrites the
func (t *Trigger) Trigger(data map[string]interface{}) error {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if len(t.programs) == 0 {
		return ErrNoExpressions
	}
	if err := t.decision.Eval(data); err == nil {
		for exp, program := range t.programs {
			out, _, err := program.Eval(map[string]interface{}{
				"this": data,
				"now":  time.Now().Unix(),
			})
			if err != nil {
				return errors.Wrapf(err, "eval: failed to evaluate trigger (%s)", exp)
			}
			if patchFields, ok := out.Value().(map[ref.Val]ref.Val); ok {
				for k, v := range patchFields {
					data[k.Value().(string)] = v.Value()
				}
			}
			if patchFields, ok := out.Value().(map[string]interface{}); ok {
				for k, v := range patchFields {
					data[k] = v
				}
			}
			if patchFields, ok := out.Value().(map[string]string); ok {
				for k, v := range patchFields {
					data[k] = v
				}
			}
		}
	}
	return nil
}

// AddExpression adds an expression to the decision tree
func (n *Trigger) AddExpression(expression string) error {
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

// Expressions returns the triggers raw CEL expressions
func (e *Trigger) Expressions() []string {
	var exp []string
	for ex, _ := range e.programs {
		exp = append(exp, ex)
	}
	sort.Strings(exp)
	return exp
}
