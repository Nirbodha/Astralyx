package packets

import (
	//"fmt"
	"errors"
	"github.com/Nirbodha/Astralyx/network/types"
)

type Packet struct {
	ID types.Byte
	Data types.ByteArray
}

func (p *Packet) Create(ID types.Byte, Data ...types.Encodable) {
	p.ID = ID
	for _, v := range Data {
		p.Data = append(p.Data, v.Encode()...)
	}
	p.Data = append(p.Data, Byte(10))
	return
}

func (p Packet) Bytes() (Converted []Byte) {
	x := 0
	Converted = make([]Byte, 2 + len(p.Data))
	Size := types.VariableInteger(2 + len(p.Data)).Encode()
	for _, v := range Size {
		Converted[x] = v
		x += 1
	}
	Converted[x] = Byte(p.ID)
	for i, v := range p.Data {
		Converted[x + 1 + i] = v
	}
	return
}

func (p *Packet) Convert(Read ByteArray) (Error error) {
	if (len(Read) < 2) {
		Error = errors.New("There is nothing to convert")
		return
	}
	x := int(1)
	for x < len(Read) - 1 {
		var v types.VariableInteger
		if (x >= 7) {
			Error = errors.New("Size byte(s) do(es) not add up with the size of the packet")
			return
		}
		v.Decode(Read[x])
		if (int(v) == len(Read)) {
			break
		}
		x++
	}
	p.ID = Byte(Read[x])
	p.Data = Read[x+1:]
	return
}
