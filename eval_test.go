package eval_test

import (
	"fmt"
	"github.com/graphikDB/eval"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	e, err := eval.New([]string{"this.name == 'bob'"})
	if err != nil {
		t.Fatal(err.Error())
	}
	if err := e.AddExpression("this.email != ''"); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := e.Eval(eval.MapperFunc(func() map[string]interface{} {
		return map[string]interface{}{
			"name":  "bob",
			"email": "bob@acme.com",
		}
	}), eval.AllTrue); err != nil {
		t.Fatal(err.Error())
	}
	if err := e.Eval(eval.MapperFunc(func() map[string]interface{} {
		return map[string]interface{}{
			"name":  "bob3",
			"email": "bob@acme.com",
		}
	}), eval.AllTrue); err == nil {
		t.Fatal("expected an error since bob3 != bob")
	}
	if len(e.Expressions()) != 2 {
		t.Fatal("expected 2 expressions")
	}
}

func ExampleNew() {
	e, err := eval.New([]string{"this.email.endsWith('acme.com')"})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := e.AddExpression("this.name != ''"); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := e.Eval(eval.MapperFunc(func() map[string]interface{} {
		return map[string]interface{}{
			"name":  "bob",
			"email": "bob@acme.com",
		}
	}), eval.AllTrue); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(strings.Join(e.Expressions(), ","))
	// Output: this.email.endsWith('acme.com'),this.name != ''
}
