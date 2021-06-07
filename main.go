package main

import (
	"github.com/Nirbodha/Astralyx/network/packet"
	"github.com/Nirbodha/Astralyx/network/types"
)

func main() {
	var p packet.Packet
	p.Create(types.Byte(0), types.ByteArray(types.Byte(1)))
	fmt.Println(p.Bytes())
}
