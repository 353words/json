## JSON - The Fine Print
+++
title = "JSON - The Fine Print"
date = "FIXME"
tags = ["golang", "json"]
categories = ["golang"]
url = "FIXME"
author = "mikit"
+++

<!--
JSON

Serialization in general
    common mistakes: serialize inside, layers objects, validation
JSON Format
    textual, no schema
    size of 123

    types
    []byte, io.Reader/io.Writer
Exported names
    field tags
    Also -, omitempty, ...
Missing vs zero value
    ptr
    map[string]any, mapstructure
    default values
Custom serialization
    - built in support: time.Time, []byte
    - Value?
Streaming
Type pollution


-->

### Introduction

Everybody knows [JSON](https://www.json.org/), it's a simple serialization format and the default format for REST APIs.
However, as they say, the devil is in the details.
In this article, we'll explore some big picture aspects of using JSON and some low level details.

## Serialization

Before diving into JSON, I'd like to take a look at serialization in general and discuss common mistakes I've seen my customers make.

Serialization is the process of converting a value in Go to bytes on one side and converting a sequence of bytes to a value on the other side.
You might ask: Why do you need serialization?

The answer is that under the hood, computers stores everything in bytes.
When you need to transfer data between two pieces of code that don't share memory,
you first need to serialize the data and then transfer it.
You'll use serialization in network operations, saving data to disk or a database and more.

The first the common mistake is code that passes serialized data in regular function calls.
For example, say you pass time as the string `"2024-08-19T12:12:39.295144041Z"`,
which means that when you want to do some date related operations, say get the year, 
you need to call `time.Parse` to convert it back to a `time.Time` object.
You should serialize only at the "edges" of your program, when it interacts with the outside world.

If you look at the major three layers of a program, then only the API and the storage layers should serialize when dealing with the outside world.

![](layers.png)

The second mistake is using the same data structure in different layers of your code.
When you work with document oriented databases such as ElasticSearch,
it's super convenient to get a request for a document, pluck it from the database and return it "as is".
However, what you are doing is tying the storage layer to the business layer to the API layer.
Which means that is you make a change to the database schema, you have changes you API - not a recipe for happy users.

Each of these layers have a different change velocity.
The API layer has low to zero velocity, the business layer has high velocity and the storage layer has medium velocity.
Say you have a `User` type in your system.
You should have a different `User` type at each layer.

![](layers-user.png)

At the beginning, all of these types look the same.
But, over time, each layer `User` type diverge.
This way, you can make database schema changes without affecting your API.

You should also only look "down" at types, the storage layer knows only about its `User`,
while the business layer know about it's `User` and also the storage `User` so it'll be able to convert between these types.

The third mistake is not validating incoming data.
An income request that is valid JSON does not mean a valid request.
Consider the following request:

**Listing 1: A JSON Request**

```json
01 {
02   "car_id": "CAR3",
03   "lat": 0.7579787,
04   "lng": -173.9881175,
05   "passenger_count": 3,
06   "shared": true
07 }
```

Listing 1 shows a JSON request.
It's a valid JSON, but the longitude (`lng`) is above the maximal longitude value of 180.

## JSON

Let's start diving into the JSON format.

JSON is a text based format without schema.
Text based makes it readable to humans, but you pay a price in the size of the encoded data.
For example, you can encode the number `123` can in a single byte,
but in JSON it's encoded as the string `123` which is three bytes.

No schema means you can quickly develop and change the data,
however no schema means you need to work harder on validating incoming data.


Every serialization protocol defines its own set of types,
and it's up to the language to decide on the mapping between JSON types and the types in the language.

In Go, the mapping is as follows:

| JSON Type | Go Type(s)                                                  |
|-----------|-------------------------------------------------------------|
| string    | string                                                      |
| number    | float64, float32, int, int8, ... int64, uint8 ... uint64    |
| boolean   | bool                                                        |
| null      | nil                                                         |
| array     | []T, []any                                                  |
| object    | struct, map[string]any                                      |


Some points to consider:
- Not everything is nil-able in Go. In JSON you can have `null` in a string member, but in Go you can have nil in a string field.
- JSON has one number type, Go has many.
- JSON arrays can have mixed types, Go slices can't - unless you use `[]any` which is painful since you need to do type assertions.
- Using structs, you can give `encoding/json` hints on how to convert JSON types to Go types.
- JSON does not have a `timestamp` like Go's `time.Time`.
- JSON does not have binary type such as Go's `[]byte`.

## `encoding/json` API

Go has built-in support for JSON serialization and in the `encoding/json` package.
This package defines the following API:

| From | Via       |  To  | Method         |
|------|-----------|------|----------------|
| JSON | bytes     | Go   | json.Unmarshal |
| Go   | bytes     | JSON | json.Marshal   |
| JSON | io.Reader | Go   | json.NewDecoder |
| Go   | io.Writer | JSON | json.NewEncoder |

`encoding/json` can work in memory with `[]byte` or in streaming with `io.Reader` and `io.Writer`.

Pick the right method depending on the situation.
For example, if you decode an incoming HTTP request body, that implements `io.Reader`, use `json.NewDecoder`.

## Unmarshaling

Most of the time, you'll work with structs to marshal or unmarshal JSON.
`encoding/json` looks only at exported fields (ones starting with capital letter).
It ignores struct fields missing from the JSON document and any JSON members that don't appear in the struct.

**Listing 2: Redundant Fields**

```go
09     data := []byte(`
10     {
11         "login": "elliot",
12         "nick": "Mr. Robot"
13     }`)
14 
15     type User struct {
16         Login string
17         UID   int
18     }
19 
20     var u User
21     if err := json.Unmarshal(data, &u); err != nil {
22         return err
23     }
24     fmt.Printf("%+v\n", u) // {Login:elliot UID:0}
```

Listing 2 shows unmarshaling JSON document.
On lines 09-13 you define the JSON data with "login" and "nick" members.
On lines 15-18 you define the `User` type with `Login` and `U
On line 20 you define a `u` variable and on lines 21-23 you unmarshal the JSON data into it.
Finally, on line 24 you print `u`. You can see that there's no error and only `Login` is filled.

`encoding/json` works only with exported fields that start with uppercase letter.
However, in JSON the convention is to use lower case names.
`encoding/json` has a heuristic to convert names from JSON to Go,
and in Listing 2 uses the `login` from the JSON document to fill the `Login` field.

Say you have a JSON document with `name` member and you want to use it to populate the `Login` field,
for this you can use [field tags](FIXME).

**Listing 3: Using Field Tags**

```go
09     data := []byte(`
10     {
11         "name": "elliot",
12         "uid": 1000
13     }`)
14 
15     type User struct {
16         Login string `json:"name"`
17         UID   int
18     }
19 
20     var u User
21     if err := json.Unmarshal(data, &u); err != nil {
22         return err
23     }
24     fmt.Printf("%+v\n", u) // {Login:elliot UID:1000}
25     return nil
26 }
```

Listing 3 shows how to use field tags.
On lines 10-13 you have JSON document with a `name` member.
On lines 15-18 you define the `User` struct.
On line 16 you use the field tag `json:"name"` to tell `encoding/json` to use the `name` member to fill the `Login` field.
On lines 20-23 you unmarshal the JSON document to the `u` variable and on line 24 you print it.
You can see that the `Login` field is filled from the `name` JSON member.

The field tags used by `encoding/json` has a mini-langugae, [read the documentation](FIXME) to learn more.

## Missing vs Zero Values

### Conclusion

Contact me at [miki@ardanlabs.com](mailto:miki@ardanlabs.com).


