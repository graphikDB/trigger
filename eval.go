package eval

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/pkg/errors"
	"sync"
)

type EvalType int

const (
	AllTrue EvalType = 0
	AnyTrue EvalType = 1
)

var (
	ErrEvalDenied    = errors.New("eval: evaluation = false")
	ErrNoExpressions = errors.New("eval: no expressions")
)

type Mapper interface {
	AsMap() map[string]interface{}
}

type MapperFunc func() map[string]interface{}

func (m MapperFunc) AsMap() map[string]interface{} {
	return m()
}

type Eval struct {
	e        *cel.Env
	programs map[string]cel.Program
	mu       sync.RWMutex
}

func New(expressions []string) (*Eval, error) {
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
	return &Eval{
		e:        e,
		programs: programs,
		mu:       sync.RWMutex{},
	}, nil
}

func (n *Eval) AddExpression(expression string) error {
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

func (n *Eval) Eval(mapper Mapper, typ EvalType) error {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if len(n.programs) == 0 {
		return ErrNoExpressions
	}
	if typ == AllTrue {
		for exp, program := range n.programs {
			out, _, err := program.Eval(map[string]interface{}{
				"this": mapper.AsMap(),
			})
			if err != nil {
				return errors.Wrapf(err, "eval: failed to evaluate expression (%s)", exp)
			}
			if val, ok := out.Value().(bool); !ok || !val {
				return ErrEvalDenied
			}
		}
	} else if typ == AnyTrue {
		for exp, program := range n.programs {
			out, _, err := program.Eval(map[string]interface{}{
				"this": mapper.AsMap(),
			})
			if err != nil {
				return errors.Wrapf(err, "eval: failed to evaluate expression (%s)", exp)
			}
			if val, ok := out.Value().(bool); ok && val {
				return nil
			}
		}
		return ErrEvalDenied
	}
	return nil
}

func (e *Eval) Expressions() []string {
	var exp []string
	for ex, _ := range e.programs {
		exp = append(exp, ex)
	}
	return exp
}
