package main

import (
	"fmt"
	"net"
	"strings"
	"math"
	"errors"
	//"io"
)



// Declarations/Logger

//These aren't variables within the logger itself, but the logger will use these still.
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

// Declarations/Client

var Clients []Client

type Client struct {
	Name string // This might be changed later to a Google Authentication token instead.
	IP string
	Connection net.Conn
	Coordinates [2]float32 //At most, this should take up 8 bytes.
	Incoming chan []byte
}

func (c *Client) Disconnect(Reason string) {
	c.Connection.Close()
	for i, v := range Clients {
		if v.IP == c.IP { //Replace this with c.Name when Google Authentication is possible
			Clients[i] = Clients[len(Clients)-1]
			Clients[len(Clients) - 1] = Client{}
			Clients = Clients[:len(Clients) - 1]
		}
	}
	console.General("Client Disconnected", c.IP + " has disconnected. Reason: " + Reason)
	return
}

//Read and Write are low-level functions.

func (c *Client) Read() (int, error) { 
	Data := make([]byte, 4096)
	BytesRead, Err := c.Connection.Read(Data)
	Data[BytesRead] = byte(10)
	c.Incoming <- Data
	return BytesRead, Err

}

func (c Client) Write(Bytes []byte) (int, error) {
	if c.Connection == nil {
		return 0, nil
	}
	BytesSent, Err := c.Connection.Write(Bytes)
	if (Err != nil) {
		return BytesSent, Err
	}
	return BytesSent, nil
}

func (c *Client) Broadcast(Packet []byte) {
	for _, v := range Clients {
		v.Write([]byte(string(Packet)))
	}
	fmt.Println(strings.ReplaceAll(string(Packet), "\n", "")) //This is mainly just chatting functionality. If they try to move, some jargon will be written in the terminal instead. We'll have to handle that later. We ddon't have packet decoding capabilities (yet)

}
// Declaration/Types (These are needed so methods are possible.)

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
	fmt.Println(p.ID, Converted[x])
	for i, v := range p.Data {
		Converted[x + 1 + i] = v
	}
	return
}


func (p *Packet) Convert(Read []byte) (Error error) {
	x := int(1)
	for x <= 7 {
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




// Example chatroom. Uncomment if you want to test this out. This is not the best chatroom out there. Use at your own risk. This doesn't use packets unfortunately. My makeshift client is netcat, and unless I type my messages, convert them into bytes, and count up the total packet size, it'll not work.
/*

func ConnectionTidbits(c Client) {
	c.Write([]byte("Welcome! Please enter a name.\n"))
	GetName(&c)
	for {
		if c.Connection == nil {
			return
		}
		_, Err := c.Read()
		if (Err != nil) {
			c.Disconnect(Err.Error())
			return
		} 
		readData := <- c.Incoming
		readData = []byte(strings.ReplaceAll(string(readData), "\n", ""))
		//c.Broadcast(<-c.Incoming)
		c.Broadcast([]byte(Green + "\n[" + c.Name + "] " + Default + string(readData)))
	}

}

func init() {
	Clients = make([]Client, 10)
}


func main() {
	console.Notify("Astralyx v0.x", "Server up!")
	Line, Err := net.Listen("tcp", ":37669")
	if (Err != nil) {
		console.Critical("Server Startup Failed", "Error: " + Err.Error())
	}
	for {
		Conn, Err := Line.Accept()
		if (Err != nil) {
			console.Critical("Connection with a client failed", "Error: " + Err.Error() )
		}
		client := Client{}
		client.IP = Conn.RemoteAddr().String()
		client.Connection = Conn
		client.Incoming = make(chan []byte, 4096)
		Clients = append(Clients, client)
		console.Notify("New connection!", "Connection from " + client.IP)
		go ConnectionTidbits(client)

	}
}

//Miscellaneous Functions

func GetName(c *Client) {
	c.Write([]byte("Enter your name: "))
	c.Read()
	Name := string(<-c.Incoming)
	Name = strings.ReplaceAll(Name, "\n", "")
	Name = strings.ReplaceAll(Name, " ", "")
	Name = strings.ReplaceAll(Name, "	", "")
	Checking := []byte(Name)
	for _, v := range Checking {
		if v != 0 {
			console.General(c.IP, "Their name is now " + Name + ".")
			c.Name = Name
			return
		}	
	}
	c.Write([]byte("This isn't a valid name. Do note that whitespace of any form will be deleted."))
	GetName(c)
}
*/
