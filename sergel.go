package sergel

import (
	"fmt"
	"regexp"
)

// Client is an abstraction of the interactions with go-sergel. This
// interface can be mocked for testing if required. To get started
// with mocking please visit mocks.go or inspect the README.
type Client interface {
	// Send sends an SMS message to the provided receiver, here named
	// 'to'. The message will have a source of 'sender' and will
	// contain the text specified in 'message'.
	Send(sender string, to string, message string) error
}

// NewClientParams is a struct containing all the fields to set
// up a new concrete client that communicates with Sergel.
type NewClientParams struct {
	Username          string // required
	Password          string // required
	PlatformID        string // required
	PlatformPartnerID string // required
	URL               string // required
	CountryCode       string
}

// NewClient returns a concrete implementation of the Client interface
func NewClient(params NewClientParams) (Client, error) {
	if err := validateNewClientParams(params); err != nil {
		return nil, err
	}
	if params.URL[len(params.URL)-1] == '/' {
		params.URL = params.URL[:len(params.URL)-1]
	}
	return &client{
		provider: provider{
			username:          params.Username,
			password:          params.Password,
			platformID:        params.PlatformID,
			platformPartnerID: params.PlatformPartnerID,
			baseUrl:           params.URL,
		},
		countryCode: params.CountryCode,
	}, nil
}

// validateNewClientParams validates the parameters sent into
// NewClient and returns an error if any required parameter is
// malformed.
func validateNewClientParams(params NewClientParams) error {
	if params.Username == "" {
		return ErrInvalidUsername
	}
	if params.Password == "" {
		return ErrInvalidPassword
	}
	if params.PlatformID == "" {
		return ErrInvalidPlatformID
	}
	if params.PlatformPartnerID == "" {
		return ErrInvalidPlatformPartnerID
	}
	// TODO: perform better validation for URL
	if params.URL == "" {
		return ErrInvalidBaseURL
	}
	return nil
}

// client is a concrete implementation of the Client interface
// which uses an underlying provided that communicates with
// Sergel.
type client struct {
	countryCode string
	provider
}

func (c client) Send(sender, to, message string) error {
	if err := c.validateSendParams(sender, to, message); err != nil {
		return err
	}
	// Transform the receiver number into one with a pre-pended country code
	to, err := c.formatReceiver(to)
	if err != nil {
		return err
	}
	return c.provider.mt(sender, to, message)
}

// validateSendParams makes sure that the parameters passed into
// the send method are not empty and returns a suitable error if
// any of them are.
func (c client) validateSendParams(sender, to, message string) error {
	if sender == "" {
		return ErrInvalidSender
	}
	if to == "" {
		return ErrInvalidReceiver
	}
	if message == "" {
		return ErrInvalidMessage
	}
	return nil
}

// formatReceiver formats the receiving phone number in such a
// way that will allow Sergel to understand it.
func (c client) formatReceiver(to string) (string, error) {
	if c.countryCode != "" {
		to = c.prependReceiverCountryCode(to)
	}
	reg, err := regexp.Compile("[^,+0-9]+")
	if err != nil {
		return "", fmt.Errorf("failed to compile regexp: %v", err.Error())
	}
	to = reg.ReplaceAllString(to, "")
	return to, nil
}

// prependReceiverCountryCode adds the standard country code to
// the phone number if a country code is not in place already
func (c client) prependReceiverCountryCode(to string) string {
	if to[0] == '+' {
		return to
	}
	return fmt.Sprintf("%s%s", c.countryCode, to[1:])
}
