package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
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
	first := true
	buf := make([]byte, 1024)
	_, err = bufio.NewReader(conn).Read(buf)
	if err != nil {
		log.Println("Error Reading from Server Try Again Later " + err.Error())
		return
	}
	fmt.Println(string(buf))
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reader" + err.Error())
			return
		}
		if first {
			_, err := fmt.Fprintf(conn, text)
			if err != nil {
				return
			}
			first = false
		} else {
			handleInput(conn, text)
		}
		buf := make([]byte, 1024)
		_, err = bufio.NewReader(conn).Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Goodbye")
				return
			}
			log.Println("Error Reading from Server Try Again Later" + err.Error())
			return
		}
		fmt.Print(string(buf))

	}

}
func handleInput(conn net.Conn, text string) error {
	text = strings.TrimSuffix(text, "\n")
	if text == "exit" {
		fmt.Println("Connection Closed")
		_, err := fmt.Fprintf(conn, text)
		if err != nil {
			return err
		}
		return nil
	} else if strings.HasPrefix(text, "download ") {
		fmt.Println("downloading...")
		_, err := fmt.Fprintf(conn, text)
		if err != nil {
			return err
		}
		reader := bufio.NewReader(conn)
		file := strings.Replace(text, " ", "", -1)[len("download") : len(text)-1]
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}
		data := read_int32(buf[:len(buf)-1])
		//add a send a got to make sure I recived the whole file
		if data != -1 {
			fmt.Fprintf(conn, "yup\n")
			filebuf := make([]byte, data)
			_, err = reader.Read(filebuf)
			fmt.Println(file)
			err = ioutil.WriteFile(file, filebuf, 0644)
			if err != nil {
				log.Println(err.Error())
				return err
			}
			return nil
		} else {
			fmt.Fprintf(conn, "yup\n")
			return nil
		}
	} else {
		_, err := fmt.Fprintf(conn, text)
		if err != nil {
			return err
		}
	}
	return nil
}

func read_int32(data []byte) (ret int64) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &ret)
	return
}
