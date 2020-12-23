# trigger

[![GoDoc](https://godoc.org/github.com/graphikDB/trigger?status.svg)](https://godoc.org/github.com/graphikDB/trigger)

a decision & trigger framework backed by Google's Common Expression Language

```go
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

```