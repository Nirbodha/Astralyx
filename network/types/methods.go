package types

import (
	"math"
	"errors"
)

type Encodable interface {
	Encode() ByteArray // Yes, that means we'll have to declare our own types. We cannot make methods of standard data types :(
}

func (s String) Encode() (Result ByteArray) {
	Result = ByteArray(s)
	return
}

func (s *String) Decode(Encoded ByteArray) (Error error) { // No errors should actually come out from this function. Just keep it like this.
	result := make([]byte, len(Encoded))
	for i, v := range result {
		result[i] = v
	}
	*s = String(string(result))
	return nil
}

func (v VariableInteger) Encode() (Result ByteArray) { //Please be careful if you edit this. You can hang your computer.
	Value := int32(v)
	for {
		Byte := Byte(Value & 0b01111111)
		Value = int32(uint32(Value) >> 7)
		if (Value != 0) {
			Byte |= 0b10000000
		}
		Result = append(Result, Byte)
		if (Value == 0) {
			break		
		}
	}
	return Result
}

func (v *VariableInteger) Decode(Encoded ByteArray) (Error error) { //Now THIS is a dangerous function. This is more probable to hang your computer
	BytesRead := int32(0)
	Result := int32(0)
	var Read Byte
	for {
		Read = Encoded[BytesRead]
		Value := int32(Read & 0b01111111)
		Result |= (Value << (7 * BytesRead))
		BytesRead++
		if BytesRead >= 5 {
			Error = errors.New("VariableInteger shouldn't be THIS large.")
		}
		if (Read & 0b10000000) == 0 {
			break
		}
	}
	*v = VariableInteger(Result)
	return	
}

func (b Byte) Encode() (Result ByteArray) {
	Result = ByteArray{Byte(b)}
	return
}

func (b *Byte) Decode(Encoded ByteArray) (Error error) {
	*b = Byte(Encoded[0]) // That's all folks
	return nil
}

func (f Float) Encode() (Result ByteArray) {
	Result = Integer(math.Float32bits(float32(f))).Encode()
	return
}

func (f *Float) Decode(Encoded ByteArray) (Error error) {
	var i Integer
	*f = Float(math.Float32frombits(uint32(i)))
	return nil
}

func (i Integer) Encode() (Result ByteArray) {
	Unsigned := uint32(i)
	Result = ByteArray{Byte(Unsigned>>24), Byte(Unsigned>>16), Byte(Unsigned>>8), Byte(Unsigned)}
	return
}

func (i *Integer) Decode(Encoded ByteArray) (Error error) {
	Bytes := Encoded[0:4]
	*i = Integer(int32(Bytes[0]) << 24 | int32(Bytes[1]) << 16 | int32(Bytes[2]) << 8 | int32(Bytes[3]))
	return nil
}
