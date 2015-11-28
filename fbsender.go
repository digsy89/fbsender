package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type Config struct {
	Name string `json:"name"`
	Tps  int    `json:"transaction_per_sec"`
	Src  string `json:"source"`
	Dst  string `json:"destination"`
}

func ConfigParse(filename string) []Config {
	config_str, err := ioutil.ReadFile(filename)
	check(err)

	var data map[string][]Config
	json.Unmarshal([]byte(config_str), &data)
	check(err)

	return data["configurations"]
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func send(id int, c Config, tup chan int) {
	// open data source
	fmt.Printf("[sender %d] reading from %s\n", id, c.Src)
	fin, err := os.Open(c.Src)
	check(err)
	reader := bufio.NewReader(fin)

	// read data
	lines := make([]string, 0)
	line, isPrefix, err := reader.ReadLine()
	for err == nil && !isPrefix {
		s := string(line)
		lines = append(lines, s)
		line, isPrefix, err = reader.ReadLine()
	}

	if err != io.EOF {
		fmt.Println(err)
		return
	}

	// open data destination
	fout, err := os.Create(c.Dst)
	check(err)
	defer fout.Close()
	writer := bufio.NewWriter(fout)

	// send data
	stime := time.Second / time.Duration(c.Tps)
	fout.Sync()
	for i := 0; i < len(lines); i++ {
		_, err := writer.WriteString(lines[i] + "\n")
		check(err)
		time.Sleep(stime)
		writer.Flush()
	}
}

func main() {
	args := os.Args
	config_filename := args[1]
	configs := ConfigParse(config_filename)
	l := make(chan int)
	for i := range configs {
		go send(i, configs[i], l)
	}
	x := <-l
	fmt.Println(x)
}
