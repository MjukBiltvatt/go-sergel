go-sergel is a simple Go wrapper for the Sergel SMS REST API.
# Getting started
## Installation
```shell
go get github.com/MjukBiltvatt/go-sergel 
```
## Importing
```go
import "github.com/MjukBiltvatt/go-sergel"
```
## Usage
Before sending an SMS with `go-sergel`, you must first create a new Sergel client.
```go
client, err := sergel.NewClient(sergel.NewClientParams{
    Username:           "username",
    Password:           "password", 
    PlatformID:         "platform_id",
    PlatformPartnerID:  "platform_partner_id",
    URL:                "https://your-sergel-url.com",
    CountryCode:        "+46",
})
if err != nil {
    handleErr(err)
}
```
All field values of `sergel.NewClientParams` should be accessible to you via Sergel. The URL field can be including or excluding a trailing `/`. The field `CountryCode` is optional, but without it, you cannot send messages with a receiving phone number that starts with the character '0'.

You can start sending messages as soon as your client has been created. To send a regular MT (mobile-terminated) message, follow the example below.

```go
if err := client.Send("sender", "+46000000000", "message"); err != nil {
    handleErr(err)
}
fmt.Println("Message was sent!")
```

### Complete example
```go
client, err := sergel.NewClient(sergel.NewClientParams{
    Username: "username",
    Password: "password", 
    PlatformID: "platform_id",
    PlatformPartnerID: "platform_partner_id",
    URL: "https://your-sergel-url.com",
    CountryCode: "+46",
})
if err != nil {
    fmt.Println(err)
    return
}

if err := client.Send("sender", "0700000000", "Hello world!"); err != nil {
    fmt.Println(err)
    return
}
```