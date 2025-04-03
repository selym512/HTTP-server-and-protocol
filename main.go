package main

import "fmt"
import "os"
import "log"
import "io"

func main() {

	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	ch := getLinesChannel(file)
	for {
		line, ok := <-ch
		if !ok {
			break
		}
		fmt.Printf("read: %s\n", line)

	}
}
func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		var b []byte = make([]byte, 8)
		var line []byte
		for {
			blen, err := f.Read(b)
			b = b[:blen]
			if err == io.EOF {
				f.Close()
				close(ch)
				break
			}
			if err != nil {
				log.Println("File error")
			}
			for _, v := range b {
				if v == 10 {
					//fmt.Printf("read: %s\n", string(line))
					ch <- string(line)
					line = nil
					continue
				}
				line = append(line, v)
			}
		}
	}()
	return ch
}
