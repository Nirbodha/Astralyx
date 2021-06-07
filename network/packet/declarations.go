package packet

import (

	"github.com/Nirbodha/Astralyx/network/types"
)

type Packet struct {
	ID types.Byte
	Data types.ByteArray
}
