package main

import (
	"flag"
	"fmt"

	"github.com/akakream/sailorsailor/ipfsPubsub"
	"github.com/akakream/sailorsailor/libp2pPubsub"
)

func main() {

	uselibp2p := flag.Bool("uselibp2p", false, "Set true if libp2p without ipfs agent should be used")
	flag.Parse()

	if *uselibp2p {
		fmt.Println("Using libp2p...")
		libp2pPubsub.Run()
	} else {
		fmt.Println("Using ipfs...")
		ipfsPubsub.Run()
	}

}
