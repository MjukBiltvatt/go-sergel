package sergel

import "errors"

var (
	ErrBadCountryCode  = errors.New("bad country code")
	ErrInvalidSender   = errors.New("the specified sender is invalid")
	ErrInvalidReceiver = errors.New("the specified receiver is invalid")
	ErrInvalidMessage  = errors.New("the specified message is invalid")
)

var (
	ErrUnknown                                   = errors.New("unknown error, please contact support and include your whole transaction")
	ErrInvalidAuthentication                     = errors.New("invalid authentication, please check your username and password")
	ErrAccessDenied                              = errors.New("access denied, please check your username and password")
	ErrInvalidOrMissingPlatformID                = errors.New("invalid or missing platform id")
	ErrInvalidOrMissingPlatformPartnerID         = errors.New("invalid or missing platform partner id")
	ErrInvalidOrMissingCurrencyForPremiumMessage = errors.New("invalid or missing currency for premium message")
	ErrNoGatesAvailable                          = errors.New("no gates available, contact support and include your whole transaction")
	ErrSpecifiedGateUnavailable                  = errors.New("specififed gate unavailable")
	ErrUnableToAccessCredentials                 = errors.New("unable to access SMSC credentials")
)

var sergelErrMap = map[int]error{
	106000: ErrUnknown,
	106100: ErrInvalidAuthentication,
	106101: ErrAccessDenied,
	106102: ErrUnableToAccessCredentials,
	106200: ErrInvalidOrMissingPlatformID,
	106201: ErrInvalidOrMissingPlatformPartnerID,
	106202: ErrInvalidOrMissingCurrencyForPremiumMessage,
	106300: ErrNoGatesAvailable,
	106301: ErrSpecifiedGateUnavailable,
}

func isSergelError(resultCode int) bool {
	_, exists := sergelErrMap[resultCode]
	return exists
}

func sergelErr(resultCode int) error {
	return sergelErrMap[resultCode]
}
