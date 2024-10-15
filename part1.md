## JSON - The Fine Print: Part 1 - Serialization & JSON
+++
title = "JSON - The Fine Print: Part 1 - Serialization & JSON"
date = "FIXME"
tags = ["golang", "json"]
categories = ["golang"]
url = "FIXME"
author = "mikit"
+++


### Introduction

Everybody knows [JSON](https://www.json.org/), it's a simple serialization format and the default format for REST APIs.
But, there are several details you need to get right to work effectively with JSON and avoid common mistakes.

This article will focus on serialization and JSON, and will look at the `encoding/json` API.

### Serialization

Before diving into JSON, I want to take a look at serialization in general and discuss common mistakes I've seen in the wild.

Serialization is the process of converting a value in Go, say a slice of ints, into bytes on one side and converting this sequence of bytes back to a Go value on the other side.

Which begs the question: What problem does serialization solve?

The answer is that under the hood, computers store everything in bytes. When you need to transfer data between two pieces of code that don't share memory, you first need to serialize the data and then transfer it. You'll use serialization in network operations, saving data to disk, or a database, and more.
