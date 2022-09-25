package main

import "nftservice/server"

func main() {
	server := server.NewServer()

	server.Start(":18888")
}
