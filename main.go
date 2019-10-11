package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var host = flag.String("host", "localhost", "The hostname or IP to connect to; defaults to \"localhost\".")
var port = flag.Int("port", 8080, "The port to connect to; default to 8080")

func main() {
	flag.Parse()
	dest := *host + ":" + strconv.Itoa(*port)
	fmt.Printf("Connectiong to %s...\n", dest)
	conn, err := net.Dial("tcp", dest)
	if err != nil {
		if _, t := err.(*net.OpError); t {
			fmt.Println("Error Connecting.")
		} else {
			fmt.Println("Unknown Error" + err.Error())
		}
		os.Exit(1)
	}
	defer conn.Close()
	for {
		log.Println("go")
		buf := make([]byte, 1024)
		cont, err := bufio.NewReader(conn).Read(buf)
		log.Println(cont)
		if err != nil {
			log.Println("Error reader" + err.Error())
			return
		}
		fmt.Println(string(buf))
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reader" + err.Error())
			return
		}
		handleInput(conn, text)
	}

}
func handleInput(conn net.Conn, text string) error {
	text = strings.TrimSuffix(text, "\n")
	fmt.Printf("%v\n", []byte("\n"))
	fmt.Printf("%+v\n", []byte(text))
	if text == "exit" {
		fmt.Println("Connection Closed")
		_, err := fmt.Fprintf(conn, text)
		if err != nil {
			return err
		}
		conn.Close()
		os.Exit(2)
		return nil
	} else if strings.HasPrefix(text, "download ") {
		fmt.Println("downloading...")
		_, err := fmt.Fprintf(conn, text)
		if err != nil {
			return err
		}
		reader := bufio.NewReader(conn)
		//conttemp , _ := reader.Read(make([]byte, 10))
		//fmt.Println(conttemp)
		file := strings.Replace(text, " ", "", -1)[len("download") : len(text)-1]
		fmt.Println(file)
		buf, err := reader.ReadBytes('\n')
		fmt.Println(buf[:len(buf)-1])
		if err != nil {
			return err
		}
		data := binary.BigEndian.Uint64(buf[:len(buf)-1])
		log.Printf("data %v\n", data)
		filebuf := make([]byte, data)
		log.Println(len(filebuf))
		cont, err := reader.Read(filebuf)
		log.Printf("log %v\n", cont)
		fmt.Printf("filebuf %v\n", filebuf)
		fmt.Println(file)
		err = ioutil.WriteFile(file, filebuf, 0644)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		fmt.Println("Kill me")
		return nil
	} else {
		_, err := fmt.Fprintf(conn, text)
		if err != nil {
			return err
		}
	}
	return nil
}
