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
	err = decision.Eval(map[string]interface{}{
		"email": "bob@gmail.com",
	})
	fmt.Println(err == trigger.ErrDecisionDenied)
	// Output: true
}

func ExampleNewTrigger() {
	// create a decision that passes if the event equals signup
	decision, err := trigger.NewDecision("this.event == 'signup' && has(this.email)")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// create a trigger based on the new decision that hashes a password and creates an updated_at timestamp
	// this would in theory be applied to a newly created user after signup
	trigg, err := trigger.NewTrigger(decision, `
	{
		'updated_at': now(),
		'password': sha1(this.password)
	}
`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	user := map[string]interface{}{
		"event":    "signup",
		"name":     "bob",
		"email":    "bob@acme.com",
		"password": "123456",
	}
	data, err := trigg.Trigger(user)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(data["updated_at"].(int64) > 0, data["password"])
	// Output: true 7c4a8d09ca3762af61e59520943dc26494f8941b
}

func TestDecision_Eval(t *testing.T) {
	type fields struct {
		expression string
	}
	type args struct {
		data map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "hello world equality",
			fields: fields{
				expression: "this.text == 'hello world'",
			},
			args: args{
				data: map[string]interface{}{
					"text": "hello world",
				},
			},
			wantErr: false,
		},
		{
			name: "hello world inequlity",
			fields: fields{
				expression: "this.text != 'hello world'",
			},
			args: args{
				data: map[string]interface{}{
					"text": "hello world",
				},
			},
			wantErr: true,
		},
		{
			name: "hello world sha1",
			fields: fields{
				expression: "sha1(this.text) != 'hello world'",
			},
			args: args{
				data: map[string]interface{}{
					"text": "hello world",
				},
			},
			wantErr: false,
		},
		{
			name: "hello world sha3",
			fields: fields{
				expression: "sha3(this.text) != 'hello world'",
			},
			args: args{
				data: map[string]interface{}{
					"text": "hello world",
				},
			},
			wantErr: false,
		},
		{
			name: "hello world sha256",
			fields: fields{
				expression: "sha256(this.text) != 'hello world'",
			},
			args: args{
				data: map[string]interface{}{
					"text": "hello world",
				},
			},
			wantErr: false,
		},
		{
			name: "hello world base64Encode",
			fields: fields{
				expression: "base64Encode(this.text) == 'aGVsbG8gd29ybGQ='",
			},
			args: args{
				data: map[string]interface{}{
					"text": "hello world",
				},
			},
			wantErr: false,
		},
		{
			name: "hello world base64Decode",
			fields: fields{
				expression: "base64Decode(this.text) == 'hello world'",
			},
			args: args{
				data: map[string]interface{}{
					"text": "aGVsbG8gd29ybGQ=",
				},
			},
			wantErr: false,
		},
		{
			name: "hello world jsonDecode",
			fields: fields{
				expression: "jsonDecode(this.text).text == 'hello world'",
			},
			args: args{
				data: map[string]interface{}{
					"text": `{ "text": "hello world"}`,
				},
			},
			wantErr: false,
		},
		{
			name: "hello world includes",
			fields: fields{
				expression: "includes(this.text, 'hello world')",
			},
			args: args{
				data: map[string]interface{}{
					"text": []string{"hello world"},
				},
			},
			wantErr: false,
		},
		{
			name: "1993 includes",
			fields: fields{
				expression: "includes(this.dob, 1993)",
			},
			args: args{
				data: map[string]interface{}{
					"dob": []int64{1993},
				},
			},
			wantErr: false,
		},
		{
			name: "hello world replace",
			fields: fields{
				expression: "replace(this.text, ' ', '') == 'helloworld'",
			},
			args: args{
				data: map[string]interface{}{
					"text": "hello world",
				},
			},
			wantErr: false,
		},
		{
			name: "hello world join",
			fields: fields{
				expression: "join(this.text, ' ') == 'hello world'",
			},
			args: args{
				data: map[string]interface{}{
					"text": []string{"hello", "world"},
				},
			},
			wantErr: false,
		},
		{
			name: "hello world titleCase",
			fields: fields{
				expression: "titleCase(this.text) == 'Hello World'",
			},
			args: args{
				data: map[string]interface{}{
					"text": "hello world",
				},
			},
			wantErr: false,
		},
		{
			name: "hello world split",
			fields: fields{
				expression: "includes(split(this.text, ' '), 'hello')",
			},
			args: args{
				data: map[string]interface{}{
					"text": "hello world",
				},
			},
			wantErr: false,
		},
		{
			name: "denver to la",
			fields: fields{
				expression: "int(geoDistance(this.denver, this.los_angelas)) > 1336367 && int(geoDistance(this.denver, this.los_angelas)) < 1536367",
			},
			args: args{
				data: map[string]interface{}{
					"denver":      []float64{39.739235, -104.990250},
					"los_angelas": []float64{34.052235, -118.243683},
				},
			},
			wantErr: false,
		},
		{
			name: "render hello world",
			fields: fields{
				expression: "render(this.text, this.data) == 'hello world'",
			},
			args: args{
				data: map[string]interface{}{
					"text": "{{ .text }}",
					"data": map[string]interface{}{
						"text": "hello world",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "parseClaims",
			fields: fields{
				expression: "parseClaims(this.jwt).name == 'John Doe'",
			},
			args: args{
				data: map[string]interface{}{
					"jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				},
			},
			wantErr: false,
		},
		{
			name: "typeOf",
			fields: fields{
				expression: "typeOf(this.jwt) == 'string'",
			},
			args: args{
				data: map[string]interface{}{
					"jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				},
			},
			wantErr: false,
		},
		{
			name: "hello world encrypt",
			fields: fields{
				expression: "decrypt(this.secret, encrypt(this.secret, 'hello world')) == 'hello world'",
			},
			args: args{
				data: map[string]interface{}{
					"secret": "this is a secret",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision, err := trigger.NewDecision(tt.fields.expression)
			if err != nil && !tt.wantErr {
				t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if decision != nil {
				if err := decision.Eval(tt.args.data); (err != nil) != tt.wantErr {
					if (err != nil) != tt.wantErr {
						t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
				}
			}
		})
	}
}

//
//func TestTrigger_Trigger(t1 *testing.T) {
//	type fields struct {
//		decision   *trigger.Decision
//		expression string
//	}
//	type args struct {
//		data map[string]interface{}
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    map[string]interface{}
//		wantErr bool
//	}{
//		{
//			name:    "",
//			fields:  fields{
//				decision:   nil,
//				expression: "",
//			},
//			args:    args{
//				data: map[string]interface{}{},
//			},
//			want:    nil,
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t1.Run(tt.name, func(t1 *testing.T) {
//			t, err := trigger.NewTrigger(tt.fields.decision, tt.fields.expression)
//			if (err != nil) != tt.wantErr {
//				t1.Errorf("Trigger() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if t != nil {
//				got, err := t.Trigger(tt.args.data)
//				if (err != nil) != tt.wantErr {
//					t1.Errorf("Trigger() error = %v, wantErr %v", err, tt.wantErr)
//					return
//				}
//				if !reflect.DeepEqual(got, tt.want) {
//					t1.Errorf("Trigger() got = %v, want %v", got, tt.want)
//				}
//			}
//		})
//	}
//}
