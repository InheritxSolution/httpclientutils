
# HTTP Utility Package

This package provides a flexible and idiomatic way to make HTTP requests in Go. It is designed to be easy to use, extensible, and aligned with Go's standard library conventions. The package supports functional options for configuring requests, dynamic body handling, and response resolution.

---

## Features

- **Functional Options**: Configure HTTP requests using a clean and flexible API.
- **Dynamic Body Handling**: Supports JSON, XML, strings, and raw bytes for request bodies.
- **Response Resolution**: Automatically unmarshal JSON or XML responses into Go structs.
- **Basic Authentication**: Easily add basic authentication to requests.
- **TLS Configuration**: Customize TLS settings for secure requests.
- **Timeout Support**: Set timeouts for requests to avoid hanging.

---

## Usage

### Basic GET Request

```go
package main

import (
	"fmt"
	"github.com/InheritxSolution/httpclientutils"
)

func main() {
	statusCode, headers, body, err := httpclientutils.MakeHTTPRequest(
		httpclientutils.WithMethod("GET"),
		httpclientutils.WithURL("https://example.com"),
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Status Code: %d\n", statusCode)
	fmt.Printf("Headers: %v\n", headers)
	fmt.Printf("Body: %s\n", body)
}
```

### POST Request with JSON Body

```go
package main

import (
	"fmt"
	"github.com/InheritxSolution/httpclientutils"
)

func main() {
	type RequestBody struct {
		Key string `json:"key"`
	}
	type ResponseBody struct {
		Message string `json:"message"`
	}

	var respBody ResponseBody
	statusCode, headers, body, err := httpclientutils.MakeHTTPRequest(
		httpclientutils.WithMethod("POST"),
		httpclientutils.WithURL("https://example.com/api"),
		httpclientutils.WithBody(RequestBody{Key: "value"}),
		httpclientutils.WithResolveResponse(&respBody),
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Status Code: %d\n", statusCode)
	fmt.Printf("Headers: %v\n", headers)
	fmt.Printf("Response: %+v\n", respBody)
	fmt.Printf("Body: %s\n", body)
}
```

### Request with Basic Authentication

```go
package main

import (
	"fmt"
	"github.com/InheritxSolution/httpclientutils"
)

func main() {
	statusCode, headers, body, err := httpclientutils.MakeHTTPRequest(
		httpclientutils.WithMethod("GET"),
		httpclientutils.WithURL("https://example.com/protected"),
		httpclientutils.WithBasicAuth("user", "pass"),
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Status Code: %d\n", statusCode)
	fmt.Printf("Headers: %v\n", headers)
	fmt.Printf("Body: %s\n", body)
}
```

### Request with Timeout

```go
package main

import (
	"fmt"
	"github.com/InheritxSolution/httpclientutils"
	"time"
)

func main() {
	statusCode, headers, body, err := httpclientutils.MakeHTTPRequest(
		httpclientutils.WithMethod("GET"),
		httpclientutils.WithURL("https://example.com/slow"),
		httpclientutils.WithTimeout(5*time.Second),
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Status Code: %d\n", statusCode)
	fmt.Printf("Headers: %v\n", headers)
	fmt.Printf("Body: %s\n", body)
}
```

---

## Available Options

| Option                        | Description                                                                 |
|-------------------------------|-----------------------------------------------------------------------------|
| `WithMethod(method string)`   | Sets the HTTP method (e.g., `GET`, `POST`).                                 |
| `WithURL(url string)`         | Sets the request URL.                                                       |
| `WithBody(body interface{})`  | Sets the request body (supports JSON, XML, strings, and raw bytes).         |
| `WithHeaders(headers map[string]string)` | Adds custom headers to the request.                                |
| `WithTLSConfig(config *tls.Config)` | Sets the TLS configuration for the request.                          |
| `WithTimeout(timeout time.Duration)` | Sets a timeout for the request.                                     |
| `WithBasicAuth(username, password string)` | Adds basic authentication to the request.                     |
| `WithResolveResponse(resp interface{})` | Automatically unmarshals the response into the provided struct.    |
| `WithResolveXMLToJSON(resp interface{})` | Converts XML responses to JSON and unmarshals into the provided struct. |
| `WithDisableEscapeHTML(disable bool)` | Disables HTML escaping for JSON marshaling.                      |

---

## Response Handling

The `MakeHTTPRequest` function returns the following:

- **Status Code**: The HTTP status code of the response.
- **Headers**: The response headers.
- **Body**: The raw response body as a byte slice.
- **Error**: Any error that occurred during the request.

If `WithResolveResponse` is used, the response body is automatically unmarshaled into the provided struct. For XML responses, `WithResolveXMLToJSON` can be used to convert the XML to JSON before unmarshaling.

---

## Error Handling

Errors are wrapped with context to make debugging easier. Common errors include:

- `failed to prepare request body`: Indicates an issue with marshaling the request body.
- `failed to create request`: Indicates an issue with creating the HTTP request.
- `request timed out`: Indicates that the request exceeded the specified timeout.
- `failed to send request`: Indicates an issue with sending the request.
- `failed to read response body`: Indicates an issue with reading the response body.
- `failed to resolve response`: Indicates an issue with unmarshaling the response.

---