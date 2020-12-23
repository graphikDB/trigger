package trigger_test

import (
	"fmt"
	"github.com/graphikDB/trigger"
	"testing"
)

func Test(t *testing.T) {
	decision, err := trigger.NewDecision("this.name == 'bob'")
	if err != nil {
		t.Fatal(err.Error())
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
	trigg, err := trigger.NewTrigger(decision, "{'name': 'coleman'}")
	if err != nil {
		t.Fatal(err.Error())
	}
	person := map[string]interface{}{
		"name":  "bob",
		"email": "bob@acme.com",
	}
	data, err := trigg.Trigger(person)
	if err != nil {
		t.Fatal(err.Error())
	}
	if data["name"] != "coleman" {
		t.Fatal("failed to trigger")
	}
	fmt.Println("trigger expressions: ", trigg.Expression())
}

func ExampleNewDecision() {
	decision, err := trigger.NewDecision("this.email.endsWith('acme.com')")
	if err != nil {
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
	fmt.Println(decision.Expression())
	// Output: this.email.endsWith('acme.com')
}

func ExampleNewTrigger() {
	decision, err := trigger.NewDecision("this.email.endsWith('acme.com')")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	trigg, err := trigger.NewTrigger(decision, `
	{
		'admin': true,
		'updated_at': now(),
		'email_hash': sha1(this.email)
	}
`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	person := map[string]interface{}{
		"name":  "bob",
		"email": "bob@acme.com",
	}
	data, err := trigg.Trigger(person)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(data["admin"], data["updated_at"].(int64) > 0, data["email_hash"])
	// Output: true true 6fd706dd2d151c2bf79218a2acd764a7d3eed7e3
}
