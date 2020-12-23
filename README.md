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

## Standard Definitions/Library

<table style="width=100%" border="1">
  <col width="15%">
  <col width="40%">
  <col width="45%">
  <tr>
    <th>Symbol</th>
    <th>Type</th>
    <th>Description</th>
  </tr>
  <tr>
    <th rowspan="1">
      !_
    </th>
    <td>
      (bool) -> bool
    </td>
    <td>
      logical not
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      -_
    </th>
    <td>
      (int) -> int
    </td>
    <td>
      negation
    </td>
  </tr>
  <tr>
    <td>
      (double) -> double
    </td>
    <td>
      negation
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      _!=_
    </th>
    <td>
      (A, A) -> bool
    </td>
    <td>
      inequality
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      _%_
    </th>
    <td>
      (int, int) -> int
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <td>
      (uint, uint) -> uint
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      _&&_
    </th>
    <td>
      (bool, bool) -> bool
    </td>
    <td>
      logical and
    </td>
  </tr>
  <tr>
    <td>
      (bool, ...) -> bool
    </td>
    <td>
      logical and (variadic)
    </td>
  </tr>
  <tr>
    <th rowspan="3">
      _*_
    </th>
    <td>
      (int, int) -> int
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <td>
      (uint, uint) -> uint
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <td>
      (double, double) -> double
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <th rowspan="9">
      _+_
    </th>
    <td>
      (int, int) -> int
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <td>
      (uint, uint) -> uint
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <td>
      (double, double) -> double
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <td>
      (string, string) -> string
    </td>
    <td>
      String concatenation. Space and time cost proportional to the sum of the
      input sizes.
    </td>
  </tr>
  <tr>
    <td>
      (bytes, bytes) -> bytes
    </td>
    <td>
      bytes concatenation
    </td>
  </tr>
  <tr>
    <td>
      (list(A), list(A)) -> list(A)
    </td>
    <td>
      List concatenation. Space and time cost proportional to the sum of the
      input sizes.
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Timestamp, google.protobuf.Duration) -> google.protobuf.Timestamp
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Duration, google.protobuf.Timestamp) -> google.protobuf.Timestamp
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Duration, google.protobuf.Duration) -> google.protobuf.Duration
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <th rowspan="6">
      _-_
    </th>
    <td>
      (int, int) -> int
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <td>
      (uint, uint) -> uint
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <td>
      (double, double) -> double
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Timestamp, google.protobuf.Timestamp) -> google.protobuf.Duration
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Timestamp, google.protobuf.Duration) -> google.protobuf.Timestamp
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Duration, google.protobuf.Duration) -> google.protobuf.Duration
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <th rowspan="3">
      _/_
    </th>
    <td>
      (int, int) -> int
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <td>
      (uint, uint) -> uint
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <td>
      (double, double) -> double
    </td>
    <td>
      arithmetic
    </td>
  </tr>
  <tr>
    <th rowspan="8">
      _<=_
    </th>
    <td>
      (bool, bool) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (int, int) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (uint, uint) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (double, double) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (string, string) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (bytes, bytes) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Timestamp, google.protobuf.Timestamp) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Duration, google.protobuf.Duration) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <th rowspan="8">
      _<_
    </th>
    <td>
      (bool, bool) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (int, int) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (uint, uint) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (double, double) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (string, string) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (bytes, bytes) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Timestamp, google.protobuf.Timestamp) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Duration, google.protobuf.Duration) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      _==_
    </th>
    <td>
      (A, A) -> bool
    </td>
    <td>
      equality
    </td>
  </tr>
  <tr>
    <th rowspan="8">
      _>=_
    </th>
    <td>
      (bool, bool) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (int, int) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (uint, uint) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (double, double) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (string, string) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (bytes, bytes) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Timestamp, google.protobuf.Timestamp) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Duration, google.protobuf.Duration) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <th rowspan="8">
      _>_
    </th>
    <td>
      (bool, bool) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (int, int) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (uint, uint) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (double, double) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (string, string) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (bytes, bytes) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Timestamp, google.protobuf.Timestamp) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Duration, google.protobuf.Duration) -> bool
    </td>
    <td>
      ordering
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      _?_:_
    </th>
    <td>
      (bool, A, A) -> A
    </td>
    <td>
      The conditional operator. See above for evaluation semantics. Will
      evaluate the test and only one of the remaining sub-expressions.
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      _[_]
    </th>
    <td>
      (list(A), int) -> A
    </td>
    <td>
      list indexing. Constant time cost.
    </td>
  </tr>
  <tr>
    <td>
      (map(A, B), A) -> B
    </td>
    <td>
      map indexing.  For string keys, cost is proportional to the size of the
      map keys times the size of the index string.
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      in
    </th>
    <td>
      (A, list(A)) -> bool
    </td>
    <td>
      list membership. Time cost proportional to the product of the size of
      both arguments.
    </td>
  </tr>
  <tr>
    <td>
      (A, map(A, B)) -> bool
    </td>
    <td>
      map key membership. Time cost proportional to the product of the size of
      both arguments.
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      _||_
    </th>
    <td>
      (bool, bool) -> bool
    </td>
    <td>
      logical or
    </td>
  </tr>
  <tr>
    <td>
      (bool, ...) -> bool
    </td>
    <td>
      logical or (variadic)
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      bool
    </th>
    <td>
      type(bool)
    </td>
    <td>
      type denotation
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      bytes
    </th>
    <td>
      type(bytes)
    </td>
    <td>
      type denotation
    </td>
  </tr>
  <tr>
    <td>
      (string) -> bytes
    </td>
    <td>
      type conversion
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      contains
    </th>
    <td>
      string.(string) -> bool
    </td>
    <td>
      Tests whether the string operand contains the substring. Time cost
      proportional to the product of sizes of the arguments.
    </td>
  </tr>
  <tr>
    <th rowspan="4">
      double
    </th>
    <td>
      type(double)
    </td>
    <td>
      type denotation
    </td>
  </tr>
  <tr>
    <td>
      (int) -> double
    </td>
    <td>
      type conversion
    </td>
  </tr>
  <tr>
    <td>
      (uint) -> double
    </td>
    <td>
      type conversion
    </td>
  </tr>
  <tr>
    <td>
      (string) -> double
    </td>
    <td>
      type conversion
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      duration
    </th>
    <td>
      (string) -> google.protobuf.Duration
    </td>
    <td>
      type conversion, duration should end with "s", which stands for seconds
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      dyn
    </th>
    <td>
      type(dyn)
    </td>
    <td>
      type denotation
    </td>
  </tr>
  <tr>
    <td>
      (A) -> dyn
    </td>
    <td>
      type conversion
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      endsWith
    </th>
    <td>
      string.(string) -> bool
    </td>
    <td>
      Tests whether the string operand ends with the suffix argument. Time cost
      proportional to the product of the sizes of the arguments.
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      getDate
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      get day of month from the date in UTC, one-based indexing
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      get day of month from the date with timezone, one-based indexing
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      getDayOfMonth
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      get day of month from the date in UTC, zero-based indexing
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      get day of month from the date with timezone, zero-based indexing
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      getDayOfWeek
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      get day of week from the date in UTC, zero-based, zero for Sunday
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      get day of week from the date with timezone, zero-based, zero for Sunday
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      getDayOfYear
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      get day of year from the date in UTC, zero-based indexing
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      get day of year from the date with timezone, zero-based indexing
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      getFullYear
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      get year from the date in UTC
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      get year from the date with timezone
    </td>
  </tr>
  <tr>
    <th rowspan="3">
      getHours
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      get hours from the date in UTC, 0-23
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      get hours from the date with timezone, 0-23
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Duration.() -> int
    </td>
    <td>
      get hours from duration
    </td>
  </tr>
  <tr>
    <th rowspan="3">
      getMilliseconds
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      get milliseconds from the date in UTC, 0-999
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      get milliseconds from the date with timezone, 0-999
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Duration.() -> int
    </td>
    <td>
      milliseconds from duration, 0-999
    </td>
  </tr>
  <tr>
    <th rowspan="3">
      getMinutes
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      get minutes from the date in UTC, 0-59
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      get minutes from the date with timezone, 0-59
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Duration.() -> int
    </td>
    <td>
      get minutes from duration
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      getMonth
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      get month from the date in UTC, 0-11
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      get month from the date with timezone, 0-11
    </td>
  </tr>
  <tr>
    <th rowspan="3">
      getSeconds
    </th>
    <td>
      google.protobuf.Timestamp.() -> int
    </td>
    <td>
      get seconds from the date in UTC, 0-59
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Timestamp.(string) -> int
    </td>
    <td>
      get seconds from the date with timezone, 0-59
    </td>
  </tr>
  <tr>
    <td>
      google.protobuf.Duration.() -> int
    </td>
    <td>
      get seconds from duration
    </td>
  </tr>
  <tr>
    <th rowspan="6">
      int
    </th>
    <td>
      type(int)
    </td>
    <td>
      type denotation
    </td>
  </tr>
  <tr>
    <td>
      (uint) -> int
    </td>
    <td>
      type conversion
    </td>
  </tr>
  <tr>
    <td>
      (double) -> int
    </td>
    <td>
      type conversion
    </td>
  </tr>
  <tr>
    <td>
      (string) -> int
    </td>
    <td>
      type conversion
    </td>
  </tr>
  <tr>
    <td>
      (enum E) -> int
    </td>
    <td>
      type conversion
    </td>
  </tr>
  <tr>
    <td>
      (google.protobuf.Timestamp) -> int
    </td>
    <td>
      Convert timestamp to int64 in seconds since Unix epoch.
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      list
    </th>
    <td>
      type(list(dyn))
    </td>
    <td>
      type denotation
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      map
    </th>
    <td>
      type(map(dyn, dyn))
    </td>
    <td>
      type denotation
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      matches
    </th>
    <td>
      (string, string) -> bool
    </td>
    <td>
      Matches first argument against regular expression in second argument.
      Time cost proportional to the product of the sizes of the arguments.
    </td>
  </tr>
  <tr>
    <td>
      string.(string) -> bool
    </td>
    <td>
      Matches the self argument against regular expression in first argument.
      Time cost proportional to the product of the sizes of the arguments.
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      null_type
    </th>
    <td>
      type(null)
    </td>
    <td>
      type denotation
    </td>
  </tr>
  <tr>
    <th rowspan="4">
      size
    </th>
    <td>
      (string) -> int
    </td>
    <td>
      string length
    </td>
  </tr>
  <tr>
    <td>
      (bytes) -> int
    </td>
    <td>
      bytes length
    </td>
  </tr>
  <tr>
    <td>
      (list(A)) -> int
    </td>
    <td>
      list size. Time cost proportional to the length of the list.
    </td>
  </tr>
  <tr>
    <td>
      (map(A, B)) -> int
    </td>
    <td>
      map size. Time cost proportional to the number of entries.
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      startsWith
    </th>
    <td>
      string.(string) -> bool
    </td>
    <td>
      Tests whether the string operand starts with the prefix argument. Time
      cost proportional to the product of the sizes of the arguments.
    </td>
  </tr>
  <tr>
    <th rowspan="5">
      string
    </th>
    <td>
      type(string)
    </td>
    <td>
      type denotation
    </td>
  </tr>
  <tr>
    <td>
      (int) -> string
    </td>
    <td>
      type conversion
    </td>
  </tr>
  <tr>
    <td>
      (uint) -> string
    </td>
    <td>
      type conversion
    </td>
  </tr>
  <tr>
    <td>
      (double) -> string
    </td>
    <td>
      type conversion
    </td>
  </tr>
  <tr>
    <td>
      (bytes) -> string
    </td>
    <td>
      type conversion
    </td>
  </tr>
  <tr>
    <th rowspan="1">
      timestamp
    </th>
    <td>
      (string) -> google.protobuf.Timestamp
    </td>
    <td>
      Type conversion of strings to timestamps according to RFC3339. Example: "1972-01-01T10:00:20.021-05:00"
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      type
    </th>
    <td>
      type(dyn)
    </td>
    <td>
      type denotation
    </td>
  </tr>
  <tr>
    <td>
      (A) -> type(dyn)
    </td>
    <td>
      returns type of value
    </td>
  </tr>
  <tr>
    <th rowspan="4">
      uint
    </th>
    <td>
      type(uint)
    </td>
    <td>
      type denotation
    </td>
  </tr>
  <tr>
    <td>
      (int) -> uint
    </td>
    <td>
      type conversion
    </td>
  </tr>
  <tr>
    <td>
      (double) -> uint
    </td>
    <td>
      type conversion
    </td>
  </tr>
  <tr>
    <td>
      (string) -> uint
    </td>
    <td>
      type conversion
    </td>
  </tr>
  <tr>
    <th rowspan="2">
      E (for fully-qualified enumeration E)
    </th>
    <td>
      (int) -> enum E
    </td>
    <td>
      type conversion when in int32 range, otherwise error
    </td>
  </tr>
  <tr>
    <td>
      (string) -> enum E
    </td>
    <td>
      type conversion for unqualified symbolic name, otherwise error
    </td>
  </tr>
  <tr>
    <td>
      now()
    </td>
    <td>
      current time in unix seconds
    </td>
  </tr>
  <tr>
    <td>
      uuid()
    </td>
    <td>
      randomly generated uuidV4 string
    </td>
  </tr>
</table>
