package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
)

func main() {
	IOCopy()
}

func IOCopy() {
	f, err := os.OpenFile("planets.txt", os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	buf.WriteString("\nsun")

	// err = ioutil.WriteFile("planets.txt", buf.Bytes(), os.ModeAppend)
	// if err != nil {
	// 	panic(err)
	// }

	_, err = io.Copy(f, buf)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(os.Stdout, f)
	if err != nil {
		panic(err)
	}
}

func OpenFile2TestReader() {
	f, err := os.Open("planets.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		log.Println(sc.Text())
	}
	if err = sc.Err(); err != nil {
		panic(err)
	}
}

func TestChanWriter() {
	cw := NewWriter()

	go func() {
		cw.Write([]byte("this line"))
		cw.Write([]byte(" and this line too"))
		defer cw.Close()
	}()
	for b := range cw.Chan() {
		log.Println(string(b))
	}

}

type chanWriter struct {
	ch chan byte
}

func NewWriter() *chanWriter {
	return &chanWriter{make(chan byte, 1024)}
}

func (cw chanWriter) Chan() <-chan byte {
	return cw.ch
}
func (cw chanWriter) Write(p []byte) (int, error) {
	n := 0
	for _, b := range p {
		cw.ch <- b
		n++
	}
	return n, nil
}

func (cw chanWriter) Close() error {
	close(cw.ch)
	log.Println("closing ch channel")
	return nil
}
