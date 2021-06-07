package packet

import (
	"errors"
	"github.com/Nirbodha/Astralyx/network/types"
)

func (p *Packet) Create(ID types.Byte, Data ...types.Encodable) {
	p.ID = ID
	for _, v := range Data {
		p.Data = append(p.Data, v.Encode()...)
	}
	p.Data = append(p.Data, types.Byte(10))
	return
}

func (p Packet) Bytes() (Converted types.ByteArray) {
	x := 0
	Converted = make(types.ByteArray, 2 + len(p.Data))
	Size := types.VariableInteger(2 + len(p.Data)).Encode()
	for _, v := range Size {
		Converted[x] = v
		x += 1
	}
	Converted[x] = types.Byte(p.ID)
	for i, v := range p.Data {
		Converted[x + 1 + i] = v
	}
	return
}

func (p *Packet) Convert(Read types.ByteArray) (Error error) {
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
		v.Decode(Read[0:x])
		if (int(v) == len(Read)) {
			break
		}
		x++
	}
	p.ID = types.Byte(Read[x])
	p.Data = Read[x+1:]
	return
}
