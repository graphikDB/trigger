# trigger

https://graphikdb.github.io/trigger/

[![GoDoc](https://godoc.org/github.com/graphikDB/trigger?status.svg)](https://godoc.org/github.com/graphikDB/trigger)

a decision & trigger framework backed by Google's Common Expression Language used in [graphikDB](https://graphikdb.github.io/graphik/)

- [x] Full Text Search Expression Macros/Functions(`startsWith, endsWith, contains`)
- [x] RegularExp Expression Macros/Functions(`matches`)
- [x] Geographic Expression Macros/Functions(`geoDistance`)
- [x] Cryptographic Expression Macros/Functions(`encrypt, decrypt, sha1, sha256, sha3`)
- [x] JWT Expression Macros/Functions(`parseClaims, parseHeader, parseSignature`)
- [x] Collection Expression Macros/Functions(`in, map, filter, exists`)
- [x] String Manipulation Expression Macros/Functions(`replace, join, titleCase, lowerCase, upperCase, trimSpace, trimPrefix, trimSuffix, render`)
- [x] URL Introspection Expression Macros/Functions(`parseHost, parseScheme, parseQuery, parsePath`)

Use Case:

Since this expression language requires just input data(map[string]interface) and an expression string, Go programs may use it to embed flexible logic that may be changed at runtime without having to recompile.

- Authorization Middleware/Policy Evaluation/Rule Engine

- Database or API "triggers" for mutating data before its commited

- Search Engine(filter something based on a decision)

## Examples

#### restrict access based on domain

```go
	decision, err := trigger.NewDecision("this.email.endsWith('acme.com')")
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return
	}
	// create a trigger based on the new decision that hashes a password and creates an updated_at timestamp
	// this would in theory be applied to a newly created user after signup
	trigg, err := trigger.NewTrigger(decision, `
	{
		'updated_at': now(),
		'password': this.password.sha1()
	}
`)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
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
|replace     |replace(text string, old string, new string) string        |full string replacement of the old value with the new value                                                 |
|join        |join(arr list(string), sep string) string                  |joins the array into a single string with the given separator                                               |
|titleCase   |titleCase(string) string                                   |converts the input into title case string                                                                   |
|lowerCase   |lowerCase(string) string                                   |converts the input into lower case string                                                                   |
|upperCase   |upperCase(string) string                                   |converts the input into upper case string                                                                   |
|trimSpace   |trimSpace(string) string                                   |removes white spaces from the input string                                                                  |
|trimPrefix  |trimPrefix(string) string                                  |removes prefix from the input string                                                                        |
|trimSuffix  |trimSuffix(string) string                                  |removes suffix from the input string                                                                        |
|geoDistance |geoDistance(this list(float64), that list(float64)) float64|haversine distance between two coordinates [lat,lng]                                                        |
|render      |render(tmplate string, data map[string]interface) string   |renders the input template with the provided data map                                                       |
|parseClaims |parseClaims(jwt string) map[string]interface) | returns the payload of the jwt as a map
|parseHeader| parseHeader(jwt string) map[string]interface | returns the header of the jwt as a map
|parseSignature| parseSignature(jwt string) string | returns the signature of the jwt as a string
|typeOf |typeOf(any) string | returns the go type of the input
|encrypt|encrypt(secret string, msg string) string| aes encrypt a message with a given secret
|decrypt|decrypt(secret string, msg string) string| aes decrypt a message with a given secret
