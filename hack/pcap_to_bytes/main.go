package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Output the DHCP packets of a PCAP file in bytes to be used in tests")
		fmt.Println("This is a workaround to not require pcap for building")
		fmt.Println("Specify pcap files as parameters to this")
		return
	}
	p, err := pcap.OpenOffline(os.Args[1])
	if err != nil {
		panic(err)
	}
	packetSource := gopacket.NewPacketSource(p, p.LinkType())
	for packet := range packetSource.Packets() {
		dhcpLayer := packet.Layers()[3]
		stringRep := fmt.Sprintf("%+v", dhcpLayer.LayerContents())
		fmt.Println(strings.ReplaceAll(stringRep, " ", ", "))
	}
}
