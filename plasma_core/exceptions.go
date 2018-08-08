package plasma_core

import "errors"

type TxAlreadySpentException struct {
	s string
}

func (e *TxAlreadySpentException)Error() string {
	return e.s
}

type InvalidTxSignatureException struct {
	s string
}

func (e *InvalidTxSignatureException)Error() string {
	return e.s
}

type InvalidBlockSignatureException struct {
	s string
}

func (e *InvalidBlockSignatureException)Error() string {
	return e.s
}

type TxAmountMismatchException struct {
	s string
}

func (e *TxAmountMismatchException)Error() string {
	return e.s
}

type InvalidBlockMerkleException struct {
	s string
}

func (e *InvalidBlockMerkleException)Error() string {
	return e.s
}

func New(errType interface{}, text string) error {
	switch errType.(type) {
	case TxAlreadySpentException:
		return &TxAlreadySpentException{text}
	case InvalidTxSignatureException:
		return &InvalidTxSignatureException{text}
	case InvalidBlockSignatureException:
		return &InvalidBlockSignatureException{text}
	case TxAmountMismatchException:
		return &TxAmountMismatchException{text}
	case InvalidBlockMerkleException:
		return &InvalidBlockMerkleException{text}
	}
	return errors.New(text)
}