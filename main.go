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

//Setup default port and host
var host = flag.String("host", "localhost", "The hostname or IP to connect to; defaults to \"localhost\".")
var port = flag.Int("port", 8080, "The port to connect to; default to 8080")

func main() {
	flag.Parse()
	dest := *host + ":" + strconv.Itoa(*port)
	fmt.Printf("Connectiong to %s...\n", dest)
	//connect to server (localhost:8080) by default
	conn, err := net.Dial("tcp", dest)
	if err != nil {
		// if error connecting
		if _, t := err.(*net.OpError); t {
			fmt.Println("Error Connecting.")
		} else {
			fmt.Println("Unknown Error" + err.Error())
		}
		os.Exit(1)
	}
	//close connection when finish running
	defer conn.Close()
	//value for first read
	first := true
	buf := make([]byte, 1024)
	//read all input from server
	_, err = bufio.NewReader(conn).Read(buf)
	if err != nil {
		log.Println("Error Reading from Server Try Again Later " + err.Error())
		return
	}
	// print output from server
	fmt.Println(string(buf))
	for {
		// create reader for standard in
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		//read from user input
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
//function for processing input
func handleInput(conn net.Conn, text string) error {
	text = strings.TrimSuffix(text, "\n")
	if text == "exit" {
		fmt.Println("Connection Closed")
		_, err := fmt.Fprintf(conn, text)
		if err != nil {
			return err
		}
		return nil
	//process how to download files
	} else if strings.HasPrefix(text, "download ") {
		fmt.Println("downloading...")
		_, err := fmt.Fprintf(conn, text)
		if err != nil {
			return err
		}
		//create new reader on TCP socket
		reader := bufio.NewReader(conn)
		// find file name to save file to
		file := strings.Replace(text, " ", "", -1)[len("download") : len(text)-1]
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}
		// read size of file (in big endian)
		data := read_int32(buf[:len(buf)-1])
		//if file exist the size will not be -1
		if data != -1 {
			//send confirmation that you got file size
			fmt.Fprintf(conn, "yup\n")
			//create a buffer of file size
			filebuf := make([]byte, data)
			_, err = reader.Read(filebuf)
			fmt.Println(file)
			//write file to disk(clobber any file that exists with same name)
			err = ioutil.WriteFile(file, filebuf, 0644)
			if err != nil {
				log.Println(err.Error())
				return err
			}
			return nil
		} else {
			//repsond you got it and exit processing step
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
