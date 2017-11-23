package reader

import (
	"bufio"
	"encoding/binary"
	"log"
	"os"
	"strings"
)

func TextReader(i string) []string {
	fr, err := os.Open(i)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer fr.Close()

	/**
	lines, err := ioutil.ReadFile(i)
	for _, line := range lines {
		log.Print(line)
	}
	return lines
	*/
	var lines []string
	scanner := bufio.NewScanner(fr)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	return lines
}

func BinReader(i string) interface{} {
	fr, err := os.Open(i)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer fr.Close()

	info := new(interface{})
	err = binary.Read(fr, binary.LittleEndian, info)
	if err != nil {
		log.Fatal(err)
	}
	return info
}
