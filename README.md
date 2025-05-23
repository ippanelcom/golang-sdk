# ippanel Go SDK

The ippanel Go SDK provides a simple client to interact with the ippanel API for sending messages. With this SDK you can send messages via a general web service, via a predefined pattern, or with a verification OTP (VOTP). The client is designed with simplicity in mind and uses Go's standard `net/http` client.

## Installation

Use Go modules to install the SDK:

```bash
go get github.com/ippanelcom/golang-sdk@v1.0.0
```

## Usage

First, import the package in your Go code:

```go
package main

import (
    "fmt"
    "log"
    "github.com/ippanelcom/golang-sdk/ippanel" // update the import path accordingly
)

func main() {
    // Create a new client instance with your API key.
    client := ippanel.NewClient("YOUR_API_KEY")

    // --- Sending a web service message ---
    // Example: Sending a basic web service message.
    webResp, err := client.SendWebservice("Hello from Webservice!", "+983000505", []string{"+989123456789", "+989356789012"})
    if err != nil {
        log.Fatalf("Error sending web service message: %v", err)
    }
    fmt.Printf("Webservice Response: %+v\n", webResp)

    // --- Sending a pattern message ---
    // Example: Sending a message using a predefined pattern.
    patternResp, err := client.SendPattern("PATTERN_CODE", "+983000505", "+989123456789", map[string]interface{}{
        "param1": "value1",
        "param2": "value2",
    })
    if err != nil {
        log.Fatalf("Error sending pattern message: %v", err)
    }
    fmt.Printf("Pattern Response: %+v\n", patternResp)

    // --- Sending a VOTP message ---
    // Example: Sending a verification OTP message.
    votpResp, err := client.SendVOTP(123456, "+989356789012")
    if err != nil {
        log.Fatalf("Error sending VOTP message: %v", err)
    }
    fmt.Printf("VOTP Response: %+v\n", votpResp)
}
```

Be sure to replace the placeholder values (`YOUR_API_KEY`, `SENDER_NUMBER`, `RECIPIENT_NUMBER`, `PATTERN_CODE`, etc.) with your actual values.

## Configuration & Customization

The SDK creates an HTTP client with a default timeout of 10 seconds. If needed, you can modify the client’s settings by using your own `http.Client` and setting it to the `HTTPClient` field of the `Client` struct. For example:

```go
client := ippanel.NewClient("YOUR_API_KEY")
client.HTTPClient = &http.Client{
    Timeout: 30 * time.Second, // customize timeout as needed
}
```

Refer to the source file [`ippanel/sdk.go`](./ippanel/sdk.go) for additional details.

## Source

For a detailed view of the implementation, see [ippanel/sdk.go](./ippanel/sdk.go).
