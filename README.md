# trigger

[![GoDoc](https://godoc.org/github.com/graphikDB/trigger?status.svg)](https://godoc.org/github.com/graphikDB/trigger)

a decision & trigger framework backed by Google's Common Expression Language used in [graphikDB](https://graphikdb.github.io/graphik/)

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
