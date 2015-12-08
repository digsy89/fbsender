package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
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

func isSocketAddress(dst string) bool {
	socketAddressRegex := regexp.MustCompile("(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}):(\\d+)")
	return socketAddressRegex.MatchString(dst)
}

func loadData(configs []Config) map[string][]string {
	sources := make(map[string]int)
	for _, c := range configs {
		sources[c.Src] = sources[c.Src] + 1
	}

	data := make(map[string][]string)
	for s := range sources {
		fmt.Println("\033[1;32mloading data from\033[0m " + s)
		fin, err := os.Open(s)
		check(err)
		reader := bufio.NewReader(fin)

		// read data
		line, isPrefix, err := reader.ReadLine()
		for err == nil && !isPrefix {
			data[s] = append(data[s], string(line))
			line, isPrefix, err = reader.ReadLine()
		}

		if err != io.EOF {
			fmt.Println(err)
			return nil
		}
	}

	return data
}

func main() {
	args := os.Args
	config_filename := args[1]
	configs := ConfigParse(config_filename)

	data := loadData(configs)

	l := make(chan int)
	for _, c := range configs {
		var sdr sender
		if isSocketAddress(c.Dst) {
			sdr = socket{Dst: c.Dst, Tps: c.Tps, Name: c.Name}
		} else {
			sdr = file{Dst: c.Dst, Tps: c.Tps, Name: c.Name}
		}
		go sdr.send(data[c.Src], l)
	}
	x := <-l
	fmt.Println(x)
}
