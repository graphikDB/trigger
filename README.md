# trigger

[![GoDoc](https://godoc.org/github.com/graphikDB/trigger?status.svg)](https://godoc.org/github.com/graphikDB/trigger)

a decision & trigger framework backed by Google's Common Expression Language used in [graphikDB](https://graphikdb.github.io/graphik/)

## Examples

#### restrict access based on domain

```go
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
```

#### create a trigger based on signup event that adds updated_at timestamp & hashes a password

```go
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
		"event": "signup",
		"name":  "bob",
		"email": "bob@acme.com",
		"password": "123456",
	}
	data, err := trigg.Trigger(user)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(data["updated_at"].(int64) > 0, data["password"])
	// Output: true 7c4a8d09ca3762af61e59520943dc26494f8941b

```

## CEL Macro Extensions

Additional details on the standard CEL spec/library may be found [here](https://github.com/google/cel-spec/blob/master/doc/langdef.md#overview)

|function    |notation                                                   |description                                                                                                 |
|------------|-----------------------------------------------------------|------------------------------------------------------------------------------------------------------------|
|now         |now() int64                                                |current timestamp in unix secods                                                                            |
|uuid        |uuid() string                                              |random uuidv4 string                                                                                        |
|sha1        |sha1(string) string                                        |sha1 hash of the input string                                                                               |
|sha256      |sha256(string) string                                      |sha256 hash of the input string                                                                             |
|sha3        |sha3(string) string                                        |sha3 hash of the input string                                                                               |
|base64Encode|base64Encode(string) string                                |base64 encoded version of the input                                                                         |
|base64Decode|base64Decode(string) string                                |base64 decoded version of the input                                                                         |
|jsonEncode  |jsonEncode(string) string                                  |json encoded version of the input                                                                           |
|jsonDecode  |jsonDecode(string) string                                  |json decoded version of the input                                                                           |
|includes    |includes(arr list(any), element any) bool                  |returns whether the slice includes the element                                                              |
|replace     |replace(text string, old string, new string) string        |full string replacement of the old value with the new value                                                 |
|join        |join(arr list(string), sep string) string                  |joins the array into a single string with the given separator                                               |
|titleCase   |titleCase(string) string                                   |converts the input into title case string                                                                   |
|lowerCase   |lowerCase(string) string                                   |converts the input into lower case string                                                                   |
|upperCase   |upperCase(string) string                                   |converts the input into upper case string                                                                   |
|trimSpace   |trimSpace(string) string                                   |removes white spaces from the input string                                                                  |
|trimPrefix  |trimPrefix(string) string                                  |removes prefix from the input string                                                                        |
|trimSuffix  |trimSuffix(string) string                                  |removes suffix from the input string                                                                        |
|split       |split(arr list(string), sep string) string                 |slices s into all substrings separated by sep and returns a slice of the substrings between those separators|
|geoDistance |geoDistance(this list(float64), that list(float64)) float64|haversine distance between two coordinates [lat,lng]                                                        |
|render      |render(tmplate string, data map[string]interface) string   |renders the input template with the provided data map                                                       |
|parseClaims |parseClaims(jwt string) map[string]interface) string | returns the payload of the jwt as a map
|typeOf |typeOf(any) string | returns the go type of the input
|encrypt|encrypt(secret string, msg string) string| aes encrypt a message with a given secret
|decrypt|decrypt(secret string, msg string) string| aes decrypt a message with a given secret
