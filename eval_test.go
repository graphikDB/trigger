package eval_test

import (
	"fmt"
	"github.com/graphikDB/eval"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	decision, err := eval.NewDecision([]string{"this.name == 'bob'"})
	if err != nil {
		t.Fatal(err.Error())
	}
	if err := decision.AddExpression("this.email != ''"); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := decision.Eval(eval.MapperFunc(func() map[string]interface{} {
		return map[string]interface{}{
			"name":  "bob",
			"email": "bob@acme.com",
		}
	}), eval.AllTrue); err != nil {
		t.Fatal(err.Error())
	}
	if err := decision.Eval(eval.MapperFunc(func() map[string]interface{} {
		return map[string]interface{}{
			"name":  "bob3",
			"email": "bob@acme.com",
		}
	}), eval.AllTrue); err == nil {
		t.Fatal("expected an error since bob3 != bob")
	}
	if len(decision.Expressions()) != 2 {
		t.Fatal("expected 2 expressions")
	}

}

func ExampleNewDecision() {
	decision, err := eval.NewDecision([]string{"this.email.endsWith('acme.com')"})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := decision.AddExpression("this.name != ''"); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := decision.Eval(eval.MapperFunc(func() map[string]interface{} {
		return map[string]interface{}{
			"name":  "bob",
			"email": "bob@acme.com",
		}
	}), eval.AllTrue); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(strings.Join(decision.Expressions(), ","))
	// Output: this.email.endsWith('acme.com'),this.name != ''
}
