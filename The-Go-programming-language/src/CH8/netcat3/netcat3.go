package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	//익명함수
	go func() {
		io.Copy(os.Stdout, conn) // Note: ignoring errors
		log.Println("done")
		done <- struct{}{} //a send statement
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done //result is discared
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
