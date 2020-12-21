package eval_test

import (
	"fmt"
	"github.com/graphikDB/eval"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	decision, err := eval.NewDecision(eval.AllTrue, []string{"this.name == 'bob'"})
	if err != nil {
		t.Fatal(err.Error())
	}
	if err := decision.AddExpression("this.email != ''"); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := decision.Eval(map[string]interface{}{
			"name":  "bob",
			"email": "bob@acme.com",
		}); err != nil {
		t.Fatal(err.Error())
	}
	if err := decision.Eval(map[string]interface{}{
			"name":  "bob3",
			"email": "bob@acme.com",
		}); err == nil {
		t.Fatal("expected an error since bob3 != bob")
	}
	if len(decision.Expressions()) != 2 {
		t.Fatal("expected 2 expressions")
	}
	trigg, err := eval.NewTrigger(decision, []string{"{'name': 'coleman'}"})
	if err != nil {
		t.Fatal(err.Error())
	}
	person := map[string]interface{}{
		"name":  "bob",
		"email": "bob@acme.com",
	}
	if err := trigg.Trigger(person); err != nil {
		t.Fatal(err.Error())
	}
	if person["name"] != "coleman" {
		t.Fatal("failed to trigger")
	}
	fmt.Println("trigger expressions: ", strings.Join(trigg.Expressions(), ","))
}

func ExampleNewDecision() {
	decision, err := eval.NewDecision(eval.AllTrue, []string{"this.email.endsWith('acme.com')"})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := decision.AddExpression("this.name != ''"); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := decision.Eval(map[string]interface{}{
			"name":  "bob",
			"email": "bob@acme.com",
		}); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(strings.Join(decision.Expressions(), ","))
	// Output: this.email.endsWith('acme.com'),this.name != ''
}

func ExampleNewTrigger() {
	decision, err := eval.NewDecision(eval.AllTrue, []string{"this.email.endsWith('acme.com')"})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	trigg, err := eval.NewTrigger(decision, []string{"{'admin': true}"})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	person := map[string]interface{}{
		"name":  "bob",
		"email": "bob@acme.com",
	}
	if err := trigg.Trigger(person); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(person["admin"])
	// Output: true
}
