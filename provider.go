package sergel

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

// provider is the struct containing methods that directly
// communicates with the Sergel API. The members of this
// struct are set in NewClient.
type provider struct {
	username          string
	password          string
	platformID        string
	platformPartnerID string
	baseUrl           string
}

// mtPayload contains the data required when sending mt
// messages using the mt method on the provider. These
// are not set by the caller but rather, is deduced from
// various other data such as the one present in the
// provider struct.
type mtPayload struct {
	Source            string `json:"source"`
	Destination       string `json:"destination"`
	UserData          string `json:"userData"`
	PlatformID        string `json:"platformId"`
	PlatformPartnerID string `json:"platformPartnerId"`
	UseDeliveryReport bool   `json:"useDeliveryReport"`
}

// errorResponse maps a potential error that could be
// returned from the Sergel API. Even though the documentation
// states that errors are returned in a different format.
type errorResponse struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
}

// providerCallback is a function that will be run with
// a http.Reponse value after one of the HTTP abstrations
// has received a response.
type providerCallback func(*http.Response) error

// mt makes a POST request to the Sergel API and tells it to
// send out a SMS text message. If the returned status code
// begins with the digit 2, then no error is returned. Otherwise,
// the method will attempt to parse the response from Sergel and
// deduces the correct error to return from errors.go.
func (p provider) mt(sender, to, message string) error {
	payload := mtPayload{
		Source:            sender,
		Destination:       to,
		UserData:          message,
		UseDeliveryReport: false,
		PlatformID:        p.platformID,
		PlatformPartnerID: p.platformPartnerID,
	}
	if err := p.post("/sms/send", payload, func(r *http.Response) error {
		if r.StatusCode/100 == 2 {
			return nil
		}
		var errResponse errorResponse
		if err := json.NewDecoder(r.Body).Decode(&errResponse); err != nil {
			return err
		}
		if isSergelError(errResponse.Status) {
			return sergelErr(errResponse.Status)
		}
		return ErrUnknown
	}); err != nil {
		return err
	}
	return nil
}

// post is an abstraction of the HTTP POST request that
// is implemented in order to coher to the DRY principle
// if more features are implemented on the provider. The
// callback parameter is ran after the closing of the
// response body is deferred, so that the caller does not
// need to concern itself with that.
func (p provider) post(uri string, data interface{}, callback providerCallback) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"POST",
		p.url(uri),
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", p.auth())
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return callback(resp)
}

// url constructs a full url based on the providers
// baseURL member and the provided URI.
func (p provider) url(uri string) string {
	return fmt.Sprintf("%s%s", p.baseUrl, uri)
}

// auth constructs an authentication string which
// cohers to Basic Authentication. The returned value
// is to be set as the Authorization header.
func (p provider) auth() string {
	credentials := fmt.Sprintf("%s:%s", p.username, p.password)
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	return fmt.Sprintf("Basic %s", encoded)
}
