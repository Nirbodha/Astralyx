// IGNORE EVERYTHING HERE! It's just for testing's sake. It'll be deleted in a bit.
package main

import (
	"fmt"
	"github.com/Nirbodha/Astralyx/network/packet"
	"github.com/Nirbodha/Astralyx/network/types"
)

func main() {
	var p packet.Packet
	p.Create(types.Byte(0), types.Byte(1))
	fmt.Println(p.Bytes())
}
