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
client := sergel.NewClient(sergel.NewClientParams{
    Username: "username",
    Password: "password", 
    PlatformID: "platform_id",
    PlatformPartnerID: "platform_partner_id",
    URL: "https://your-sergel-url.com",
})
```
All field values of `sergel.NewClientParams` should be accessible to you via Sergel. The URL field can be including or excluding a trailing `/`. 

You can start sending messages as soon as your client has been created. To send a regular MT (mobile-terminated) message, follow the example below. Note that this example requires you to specify the country code for each message, such as `+46`. The receiving phone number therefor cannot start with a zero.

```go
if err := client.Send("sender", "+46000000000", "message"); err != nil {
    handleErr(err)
}
fmt.Println("Message was sent!")
```

To set a standard country code for the client, which will allow you to use a zero as the first character in the receiver argument above, simple call the following method.

```go 
if err := client.SetCountryCode("+46"); err != nil {
    handleErr(err)
}
```

If the argument to `SetCountryCode` does not begin with a `+` character, then an error will be returned.

### Complete example
```go
client := sergel.NewClient(sergel.NewClientParams{
    Username: "username",
    Password: "password", 
    PlatformID: "platform_id",
    PlatformPartnerID: "platform_partner_id",
    URL: "https://your-sergel-url.com",
})

if err := client.SetCountryCode("+46"); err != nil {
    fmt.Println(err)
    return
}

if err := client.Send("sender", "0700000000", "Hello world!"); err != nil {
    fmt.Println(err)
    return
}
```