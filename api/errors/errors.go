package errors

import (
	"errors"
	"fmt"
)

// ErrInvalidAppContext signals an invalid context passed to the routing system
var ErrInvalidAppContext = errors.New("invalid app context")

// ErrInvalidJSONRequest signals an error in json request formatting
var ErrInvalidJSONRequest = errors.New("invalid json request")

// ErrValidation signals an error in validation
var ErrValidation = errors.New("validation error")

// ErrInvalidSignatureHex signals a wrong hex value was provided for the signature
var ErrInvalidSignatureHex = errors.New("invalid signature, could not decode hex value")

// ErrTxGenerationFailed signals an error generating a transaction
var ErrTxGenerationFailed = errors.New("transaction generation failed")

// ErrInvalidSenderAddress signals a wrong format for sender address was provided
var ErrInvalidSenderAddress = errors.New("invalid hex sender address provided")

// ErrInvalidReceiverAddress signals a wrong format for receiver address was provided
var ErrInvalidReceiverAddress = errors.New("invalid hex receiver address provided")

// ErrInvalidTxFields signals that one or more field of a transaction are invalid
type ErrInvalidTxFields struct {
	Message string
	Reason  string
}

// Error returns the string message of the ErrInvalidTxFields custom error struct
func (eitx *ErrInvalidTxFields) Error() string {
	return fmt.Sprintf("%s : %s", eitx.Message, eitx.Reason)
}
