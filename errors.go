package cfb

import "errors"

var (
	ErrInsufficientData   = errors.New("insufficient data")
	ErrWrongSector        = errors.New("wrong sector")
	ErrInvalidSectorChain = errors.New("invalid sector chain")
	ErrInvalidSeek        = errors.New("invalid seek operation")
	ErrInvalidOffset      = errors.New("invalid offset")
	ErrWrongObjectType    = errors.New("wrong object type")
	ErrInvalidObject      = errors.New("invalid object")
	ErrObjectNotFound     = errors.New("object not found")
	ErrValidation         = errors.New("validation error")
)
