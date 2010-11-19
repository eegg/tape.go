package main

import (
	"./tape"
	"fmt"
	)

func main() {
	t, _ := tape.NewTapeFromFilename("sample.txt")

	i := 0

	for {
		byte, err := t.ReadByte()
		if err != nil {
			return
		} else {
			fmt.Printf("%c", byte)
			i++

			if i == 150 {
				ok := t.Rewind(75)
				if !ok {
					fmt.Println("ERROR!")
				}
			}
		}
	}
}