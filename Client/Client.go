package main

import (
	"fmt"
	"net"
	"bufio"
	"strconv"
	"math"
	"errors"
	"bytes"
)


var Default  = "\033[0m"
var Red    = "\033[31m"
var Green  = "\033[32m"
var Yellow = "\033[33m"


type Logger struct {
	//Empty for now. If logfiles are really needed, this will have a string variable within it containing the location of the logfile with respect to the Golang executable.
}

func (Logger) Critical(Title string, Info string) (int, error) {
	return fmt.Println(Red + "[" + Title + "] " + Default + Info)
}

func (Logger) General(Title string, Info string) (int, error) {
	return fmt.Println(Yellow + "[" + Title + "] " + Default + Info)
}

func (Logger) Notify(Title string, Info string) (int, error) {
	return fmt.Println(Green + "[" + Title + "] " + Default + Info)
}

var console Logger

type (
	String string
	VariableInteger int32 //Why is this needed? Verifying the integrity of a packet, of course.
	Byte uint8
	Float float32
	Integer int32
)

type Encodable interface {
	Encode() []byte // Yes, that means we'll have to declare our own types. We cannot make methods of standard data types :(
}

func (s String) Encode() (Result []byte) {
	Result = []byte(s)
	return
}

func (s *String) Decode(Encoded []byte) (Error error) { // No errors should actually come out from this function. Just keep it like this.
	*s = String(string(Encoded))
	return nil
}

func (v VariableInteger) Encode() (Result []byte) { //Please be careful if you edit this. You can hang your computer.
	Value := int32(v)
	for {
		Byte := byte(Value & 0b01111111)
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

func (v *VariableInteger) Decode(Encoded []byte) (Error error) { //Now THIS is a dangerous function. This is more probable to hang your computer
	BytesRead := int32(0)
	Result := int32(0)
	var Read byte
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

func (b Byte) Encode() (Result []byte) {
	Result = []byte{byte(b)}
	return
}

func (b *Byte) Decode(Encoded []byte) (Error error) {
	*b = Byte(Encoded[0]) // That's all folks
	return nil
}

func (f Float) Encode() (Result []byte) {
	Result = Integer(math.Float32bits(float32(f))).Encode()
	return
}

func (f *Float) Decode(Encoded []byte) (Error error) {
	var i Integer
	*f = Float(math.Float32frombits(uint32(i)))
	return nil
}

func (i Integer) Encode() (Result []byte) {
	Unsigned := uint32(i)
	Result = []byte{byte(Unsigned>>24), byte(Unsigned>>16), byte(Unsigned>>8), byte(Unsigned)}
	return
}

func (i *Integer) Decode(Encoded []byte) (Error error) {
	Bytes := Encoded[0:4]
	*i = Integer(int32(Bytes[0]) << 24 | int32(Bytes[1]) << 16 | int32(Bytes[2]) << 8 | int32(Bytes[3]))
	return nil
}


// Declarations/Packets

type Packet struct {
	ID Byte
	Data []byte
}


func (p *Packet) Create(ID Byte, Data ...Encodable) {
	p.ID = ID
	for _, v := range Data {
		p.Data = append(p.Data, v.Encode()...)
	}

	return
}

func (p *Packet) ToBytes() (Converted []byte) {
	x := 0
	Converted = make([]byte, 2 + len(p.Data))
	Size := VariableInteger(2 + len(p.Data)).Encode()
	for _, v := range Size {
		Converted[x] = v
		x += 1
	}
	Converted[x] = byte(p.ID)
	for i, v := range p.Data {
		Converted[x + 1 + i] = v
	}
	return
}


func (p *Packet) Convert(Read []byte) (Error error) {
	if (len(Read) == 0 || len(Read) == 1 ) {
		Error = errors.New("There is nothing to convert")
		return
	}
	x := int(1)
	for  x < len(Read) - 1 && x <= 7 {
		var v VariableInteger
		if (x == 7) {
			Error = errors.New("Size byte(s) do(es) not add up with the size of the packet")
			return
		}
		v.Decode(Read[0:x])
		if (int(v) == len(Read)) {
			break
		}
		x++
	}
	p.ID = Byte(Read[x])
	p.Data = Read[x+1:]
	return
}


func MakePacketByUser() (p Packet) {
	var Final bytes.Buffer
	var getID func() Byte
	var getData func() []byte
	getID = func() Byte {
		getID := func() (ID Byte) {
			fmt.Println("Enter a packet ID")
			var input []byte
			fmt.Scanln(&input)
			result, Err :=  strconv.Atoi(string(input))
			if (Err != nil || result > 255 || result < 0) {
				fmt.Println("This is not a valid packet ID. Try again.")
				ID = getID()
				return
			}
			ID = Byte(byte(result))
			return
		}
		
		return getID()
	}
	getData = func() []byte {
		getData := func() (Data []byte) {
			fmt.Println("Enter a byte for data.")
			var input []byte
			fmt.Scanln(&input)
			result, Err := strconv.Atoi(string(input))
			if (Err != nil || result > 255 || result < 0) {
				fmt.Println("This is not a valid byte. Try again.")
				Data = getData()
				return
			}
			Final.WriteByte(byte(result))
			fmt.Println("Would you like to add another byte? Type in 'y' if you do.")
			var inputa string
			fmt.Scanln(&inputa)
			if (inputa=="y") {
				Data = getData()
				return
			}
			Data = Final.Bytes()
			return
		}
		return getData()
	}
	p.Create(getID(), String(string(getData()) + "\n"))
	return
}

func main() {
	p := MakePacketByUser()
	fmt.Println(p.ToBytes())
	Connection, Err := net.Dial("tcp", ":37669")
	if (Err != nil) {
		fmt.Println(Err.Error())
		return
	}
	Connection.Write(p.ToBytes())
	Data, Err := bufio.NewReader(Connection).ReadBytes(byte(10))
	if (Err != nil) {
		fmt.Println(Err.Error())
		return
	}
	fmt.Println(Data)
}
