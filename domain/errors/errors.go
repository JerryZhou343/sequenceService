package errors

import "github.com/pkg/errors"

var (
	ErrSeqNotEnough = errors.New("sequence not enough")
)