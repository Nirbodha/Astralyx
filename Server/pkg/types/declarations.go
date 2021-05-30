package types

import (
	"io"
	"encoding/binary"
	"errors"
)

type (
	ByteArray []byte
	VariableInt ByteArray
	Message string
	Position struct {
		X, Y uint16
	}
	Integer int32
)

