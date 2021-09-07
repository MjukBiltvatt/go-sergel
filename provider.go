package sergel

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type provider struct {
	username          string
	password          string
	platformID        string
	platformPartnerID string
	baseUrl           string
}

type mtPayload struct {
	Source            string `json:"source"`
	Destination       string `json:"destination"`
	UserData          string `json:"userData"`
	PlatformID        string `json:"platformId"`
	PlatformPartnerID string `json:"platformPartnerId"`
	UseDeliveryReport bool   `json:"useDeliveryReport"`
}

type errorResponse struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
}

type providerCallback func(*http.Response) error

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

func (p provider) url(uri string) string {
	return fmt.Sprintf("%s%s", p.baseUrl, uri)
}

func (p provider) auth() string {
	credentials := fmt.Sprintf("%s:%s", p.username, p.password)
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	return fmt.Sprintf("Basic %s", encoded)
}
