package trigger

import (
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/graphikDB/generic"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"os"
	"time"
)

func init() {
	var err error
	var declarations = []*expr.Decl{
		decls.NewVar("this", decls.NewMapType(decls.String, decls.Any)),
	}
	for _, function := range Functions {
		declarations = append(declarations, function.decl)
	}
	env, err := cel.NewEnv(cel.Declarations(declarations...))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	globalEnv = &environment{
		env:   env,
		cache: generic.NewCache(1 * time.Minute),
	}
}

type environment struct {
	env   *cel.Env
	cache *generic.Cache
}

func (e *environment) Program(expression string) (cel.Program, error) {
	if val, ok := e.cache.Get(expression); ok {
		if program, ok := val.(cel.Program); ok {
			return program, nil
		}
	}
	ast, iss := globalEnv.env.Compile(expression)
	if iss.Err() != nil {
		return nil, iss.Err()
	}
	var overloads []*functions.Overload
	for _, function := range Functions {
		overloads = append(overloads, function.overload)
	}
	program, err := globalEnv.env.Program(
		ast,
		cel.Functions(overloads...),
	)
	if err != nil {
		return nil, err
	}
	e.cache.Set(expression, program, 5*time.Minute)
	return program, nil
}

var globalEnv *environment
