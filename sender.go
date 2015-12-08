package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

type sender interface {
	send(data []string, tup chan int)
}

type file struct {
	Dst  string
	Tps  int
	Name string
}

type socket struct {
	Dst  string
	Tps  int
	Name string
}

func (f file) send(data []string, tup chan int) {
	// open data destination
	fout, err := os.Create(f.Dst)
	check(err)
	writer := bufio.NewWriter(fout)

	fmt.Printf("[%s] sending to pipe[%s]\n", f.Name, f.Dst)
	// send data
	stime := time.Second / time.Duration(f.Tps)
	fout.Sync()
	for i := 0; i < len(data); i++ {
		_, err := writer.WriteString(data[i] + "\n")
		check(err)
		time.Sleep(stime)
		writer.Flush()
	}
	close(tup)
}

func (s socket) send(data []string, tup chan int) {
	conn, err := net.Dial("tcp", s.Dst)
	check(err)
	writer := bufio.NewWriter(conn)

	fmt.Printf("[%s] sending to socket[%s]\n", s.Name, s.Dst)
	// send data
	stime := time.Second / time.Duration(s.Tps)
	for i := 0; i < len(data); i++ {
		_, err := writer.WriteString(data[i] + "\n")
		check(err)
		time.Sleep(stime)
		writer.Flush()
	}
	close(tup)
}
