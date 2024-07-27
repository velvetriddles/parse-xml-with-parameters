# Educational Task: HTTP Server Testing in Go

This project is an educational task designed to teach testing of an HTTP server in the Go programming language. The focus is on creating and validating test cases for parsing query parameters, handling various edge cases, and simulating HTTP requests and responses.

## Overview

The project involves implementing a simple HTTP server that handles search queries. The main components include:

1. **Search Server**: Handles incoming HTTP GET requests with query parameters for search operations.
2. **Search Client**: Sends search requests to the search server and processes the responses.
3. **Test Cases**: Validates the server's behavior under different conditions, such as invalid query parameters, incorrect dataset files, and various order fields.


## Files

- `main.go`: Contains the implementation of the search server and client.
- `main_test.go`: Contains the test cases for the search server and client.

## Usage

### Running the Server

To run the search server, execute the following command:

```sh
go run .
```

To run tests, execute the following command: 

```sh
go test -v 
```

You can check the coverage of the tests in cover.html file


