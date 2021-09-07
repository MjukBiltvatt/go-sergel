package sergel

import (
	"fmt"
	"regexp"
)

type Client interface {
	Send(sender string, to string, message string) error
	SetCountryCode(countryCode string) error
}

type NewClientParams struct {
	Username          string
	Password          string
	PlatformID        string
	PlatformPartnerID string
	URL               string
}

func NewClient(params NewClientParams) Client {
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
	}
}

type client struct {
	countryCode string
	provider
}

func (c *client) SetCountryCode(countryCode string) error {
	if countryCode[0] != '+' {
		return ErrBadCountryCode
	}
	c.countryCode = countryCode
	return nil
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

func (c client) formatReceiver(to string) (string, error) {
	if c.countryCode != "" {
		if countryCodedTo, err := c.prependReceiverCountryCode(to); err != nil {
			return "", err
		} else {
			to = countryCodedTo
		}
	}
	reg, err := regexp.Compile("[^,+0-9]+")
	if err != nil {
		return "", fmt.Errorf("failed to compile regexp: %v", err.Error())
	}
	to = reg.ReplaceAllString(to, "")
	return to, nil
}

func (c client) prependReceiverCountryCode(to string) (string, error) {
	if to[0] == '+' {
		return to, nil
	}
	return fmt.Sprintf("%s%s", c.countryCode, to[1:]), nil
}
