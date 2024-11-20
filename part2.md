## JSON - The Fine Print: Part 2 - Emitting JSON
+++
title = "JSON - The Fine Print: Part 2 - Emitting JSON"
date = "FIXME"
tags = ["golang", "json"]
categories = ["golang"]
url = "FIXME"
author = "mikit"
+++


### Introduction

In [part 1](https://www.ardanlabs.com/blog/2024/10/json-the-fine-print-part-1.html) we took a high level view on serialization in general and JSON in specific. 
In this part, we'll roll our sleeves and start working with JSON.
This part focuses on emitting JSON, you might think it's basic, but there much more to it than calling `json.Marshal`.

### `json.Marshal` vs `json.Encoder`

`encoding/json` has two main APIs: `json.Marshal` and `json.Encoder`.
`json.Marshal` returns a `[]byte` while `json.Encoder` will write to `io.Writer`.
The question is: When to use each one?

Typically, you'll use use both in an HTTP handler to return JSON to the caller.
Using `json.Encoder` on the `http.ResponseWriter` will allow you to write JSON without constructing a `[]byte` in memory.
But, if there are encoding errors, you cannot notify the client by setting the HTTP status code.
Once you start writing to `http.ResponseWriter` it will set the status code to `OK`.

Working with `json.Marshaler` will allows you to check for encoding error *before* you send the data to the client.
But, you create a `[]byte` in memory of all the data you want to send, and it might be big.

### Simple Marshaling

There are two easy ways to send JSON to the client, using a `map[string]any` and using a struct.

The simplest way is to use a map:

**Listing 1: Marshaling a map**

```go
25     vm := map[string]any{
26         "id":     "b70229443f8d489bbc733f13a9268f63",
27         "cpus":   4,
28         "memory": 32,
29     }
30 
31     json.NewEncoder(os.Stdout).Encode(vm)
```

Listing 1 shows how to marshal a `map[string]any`.
In lines 25-29 we create the VM map and on line 31 we marshal it to the standard output.
This code will print:

```
{"cpus":4,"id":"b70229443f8d489bbc733f13a9268f63","memory":32}
```

Using maps is easy and quick, but using structs allows you to define a schema for your outgoing messages.
Here's the same code using a struct:

**Listing 2: Marshaling a Struct**

```go
08 type VM struct {
09     ID     string
10     CPUs   int
11     Memory int
12 }
...
13 
15     vm := VM{
16         ID:     "b70229443f8d489bbc733f13a9268f63",
17         CPUs:   4,
18         Memory: 32,
19     }
20 
21     json.NewEncoder(os.Stdout).Encode(vm)
```

Listing 2 shows struct marshaling.
On lines 08-12 we define the `VM` struct and on lines 15-19 we create the `vm` variable.
Then, on line 21 we marshal the struct to the standard output.
This code will print:

```
{"ID":"b70229443f8d489bbc733f13a9268f63","CPUs":4,"Memory":32}
```

The properties are not in lower case (e.g. `ID` and not `id`) as custom in JSON.
If you want to change the name of the output properties, you need to use field tags.

### Field Tags

You can use field tags to tell `encoding/json` how to map the struct fields with the emitted JSON properties.

**Listing 3: Using Fields Tags**

```go
08 type VM struct {
09     ID     string `json:"id"`
10     CPUs   int    `json:"cpus"`
11     Memory int    `json:"memory"`
12 }
13 
...
15     vm := VM{
16         ID:     "b70229443f8d489bbc733f13a9268f63",
17         CPUs:   4,
18         Memory: 32,
19     }
20     json.NewEncoder(os.Stdout).Encode(vm)
```

Listing 3 shows how to use field tags to change the name of the output JSON properties.
On lines 08-12 we define the `VM` struct and add a field tag to each field in the struct.
On lines 15-19 we create a `vm` variable and on line 20 we marshal it to the standard output.

This code will print:

```
{"id":"b70229443f8d489bbc733f13a9268f63","cpus":4,"memory":32}
```

You can do more with field tags, see the [json.Marshal](https://pkg.go.dev/encoding/json#Marshal) function documentation for the full specification.
I'm going to focus on two things: `omitempty` and `-`.

If you add `omitempty` to a field tag, `encoding/json` won't write the field to the output if it has the zero value.
This allows you to save bandwidth, which eventually will save you money.

**Listing 4: omitempty**

```go
08 type VM struct {
09     ID     string `json:"id,omitempty"`
10     CPUs   int    `json:"cpus,omitempty"`
11     Memory int    `json:"memory,omitempty"`
12 }
13 
...
15     vm := VM{
16         ID:     "",
17         CPUs:   4,
18         Memory: 32,
19     }
20     json.NewEncoder(os.Stdout).Encode(vm)
```

Listing 4 shows how to use `omitempty`.
On lines 08-12 we define `VM` with `omitempty` in the fields tags.
On lines 15-19 we create a `vm` variable and on line 20 we marshal it to the standard output.

The code will print:
```
{"cpus":4,"memory":32}
```

Without `omitempty` this code will print:

```
{"id":"","cpus":4,"memory":32}
```

`omitempty` save 8 bytes, which might not seem a lot, but if you're sending many messages can amount to a lot of bandwidth.

Another struct tag that `encoding/json` recognizes is `-`, which tells it not to emit a specific field.
`-` is used to prevent leaking sensitive information, say a `Token` field in your `User` struct.
I see `-` as a sign that you don't have a good separation between your API layer and your business or data layers.

### Custom Serialization

JSON has a limited set of types, not all Go types can be mapped to JSON types.
For example, JSON does not have a "timestamp" type, while Go has `time.Time`.

Timestamps are serialized into JSON in two ways: A string such as `2025-01-02T15:35:47.990Z` or a number which is usually the number of seconds since January 1, 1970 UTC (known as epoch).

Let's try and see what Go does:

**Listing 5: Marshaling time.Time**

```go
09 type Log struct {
10     Time    time.Time
11     Level   string
12     Message string
13 }
...
16     l := Log{
17         Time:    time.Now().UTC(),
18         Level:   "ERROR",
19         Message: "divide by cucumber error",
20     }
21 
22     if err := json.NewEncoder(os.Stdout).Encode(l); err != nil {
23         fmt.Println("ERROR:", err)
24     }
```

Listing 5 shows marshaling a struct with a `time.Time` field.
On lines 09-13 we define `Log`, on lines 16-20 we create a variable and on line 22 we marshal `l` to stdout.

This code produces the following output without any error:

```
{"Time":"2024-11-19T05:33:22.774425457Z","Level":"ERROR","Message":"divide by cucumber error"}
```

`encoding/json` marshals `time.Time` into an [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339).
If you look at the documentation of [time.Time](https://pkg.go.dev/time#Time) you'll see it has a method called `MarshalJSON`.
This means that `time.Time` implements the [json.Marshaler](https://pkg.go.dev/encoding/json#Marshaler) interface.


_Note: Another type that is missing from JSON is a binary type. `encoding/json` will encode a `[]byte` into a [base64](https://en.wikipedia.org/wiki/Base64) encoded string._
`fmt.Printf("%08b\n", 13)` will print `00001101`._

#### Using json.Marshaler

You can implement `json.Marshaler` on your types to get custom JSON encoding.
Assume you have the following type:

**Listing 6: Value Type**

```go
09 type Unit string
10 
11 const (
12     Meter = "meter"
13     Inch  = "inch"
14 )
15 
16 type Value struct {
17     Unit   Unit
18     Amount float64
19 }
```

By default, `encoding/json` will marhshal a `Value` to a JSON object with a `Unit` and `Amount` properties.
But say that you want a `Value` to be encoded as `14.2inch` instead.

I use two steps when implementing `json.Marshaler`.
First step is converting the type to a type `encoding/json` know how to handle, the second step is to use `json.Marshal` to return the result.
Do not try to construct the output JSON by hand unless you have a really good reason, there are many edge cases you might miss that way.

**Listing 7: Implementing json.Marshaler**

```go
21 func (v Value) MarshalJSON() ([]byte, error) {
22     // Step 1: Convert to type known to encoding/json
23     s := fmt.Sprintf("%f%s", v.Amount, v.Unit)
24 
25     // Step 2: Use json.Marshal
26     return json.Marshal(s)
27 }
...
49     v := Value{
50         Unit:   Meter,
51         Amount: 2.1,
52     }
53 
54     data, err := json.Marshal(v)
55     if err != nil {
56         return err
57     }
58     fmt.Println(string(data)) // "2.1meter"
```

Listing 7 show how to implement `json.Marshaler` for `Value`.
On line 23 we convert v to a string and on line 25 we use `json.Marshal` to convert the string to JSON.
On lines 49-58 we create a Value and then encode it to JSON.

Note that `MarshalJSON` is defined with value pointer semantics.
But it'll work for pointer value semantics as well.
If you're not sure whey, checkout [this awesome video](https://www.youtube.com/watch?v=Z5cvLOrWlLM) by Bill.

### Streaming JSON

The JSON specification does not support streaming - sending one JSON object after another.
If you want to stream JSON, the common way is to send one JSON object per line.
This is known as [jsonlines](https://jsonlines.org/) or [ndjson](https://docs.mulesoft.com/dataweave/latest/dataweave-formats-ndjson).

Lucky for you, `json.Encoder` already does that for you. Here's an example:

**Listing 8: Streaming JSON**

```go
09 type Event struct {
10     Type string  `json:"type"`
11     X    float64 `json:"x"`
12     Y    float64 `json:"y"`
13 }
14 
15 func work() error {
16     events := []Event{
17         {"click", 100, 200},
18         {"move", 101, 202},
19     }
20 
21     enc := json.NewEncoder(os.Stdout)
22 
23     for _, e := range events {
24         if err := enc.Encode(e); err != nil {
25             return err
26         }
27     }
28     return nil
29 }
```

Listing 8 shows how to stream JSON.
On lines 09-13 we define `Event`.
On lines 16-19 we create a slice of two events.
On line 21 we create a JSON encoder and on lines 23-27 we use the same encoder to encode all the events.

The output of this code is:

```
{"type":"click","x":100,"y":200}
{"type":"move","x":101,"y":202}
```

The encoder encoded each JSON object in a single line and added a newline between each JSON object.
Of course, the receiving side should know to parse each line as a JSON object, which the JSON decoder does as well.

#### Streaming JSON with HTTP Chunked Transfer Encoding

HTTP version 1.1 added [chunked transfer encoding](https://en.wikipedia.org/wiki/Chunked_transfer_encoding).
This allows an HTTP server to send the response in chunks.
The server sets the HTTP header `Transfer-Encoding` to `chunked` and then write size followed by data.
In Go, you can send chunked data using an [http.ResponseController](https://pkg.go.dev/net/http#ResponseController).

**Listing 9: Streaming JSON Over HTTP** 

```go
41 func eventsHandler(w http.ResponseWriter, r *http.Request) {
42     ctrl := http.NewResponseController(w)
43 
44     enc := json.NewEncoder(w)
45     for evt := range queryEvents() {
46         if err := enc.Encode(evt); err != nil {
47             // Can't set error
48             slog.Error("JSON encode", "error", err)
49             return
50         }
51 
52         if err := ctrl.Flush(); err != nil {
53             slog.Error("flush", "error", err)
54             return
55         }
56     }
57 }
```

Listing 9 shows how to stream JSON in an HTTP handler.
On line 42 we create an `http.ResponseController` and on line 44 we create a `json.Encoder`.
On line 45 we iterate over the events, on line 46 we use `enc` to encode the event and on line 52 we call `Flush` that will send the current chunk of data.

You can use `curl` to view the raw HTTP response:

**Listing 10: Using curl to Call the Server**

```
01 $ curl --raw -i http://localhost:8080/events
02 HTTP/1.1 200 OK
03 Date: Tue, 19 Nov 2024 17:15:21 GMT
04 Content-Type: text/plain; charset=utf-8
05 Transfer-Encoding: chunked
06 
07 21
08 {"type":"click","x":100,"y":200}
09 
10 20
11 {"type":"move","x":101,"y":202}
12 
13 20
14 {"type":"move","x":102,"y":203}
15 
16 20
17 {"type":"move","x":103,"y":204}
18 
19 20
20 {"type":"move","x":104,"y":204}
21 
22 21
23 {"type":"click","x":104,"y":204}
24 
25 0
```

Listing 10 shows how to call the server and view the underlying chunked response.
On line 01 we use `curl` to call the server. The `--raw` flag tells `curl` to show the raw response and the `-i` switch tells `curl` to show the response HTTP headers.
On line 05 we see that we get chunked response and on lines 06-24 we see the chunks.
One line 25 we see the sentinel value of `0` tell the HTTP client there are no more chunks.

### Conclusion

Emitting JSON can be simple as `json.Marshal`, but if you need more sophisticated methods, `encoding/json` is ther for you with field tags, the `json.Marshaler` interface and streaming support.
Get to know the API and understand the pros and cons of using `io.Writer` vs a `[]byte` in your code.
