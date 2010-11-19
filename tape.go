package tape

import (
	"bufio"
	"os"
	)

const nodeLen = 100
// Max index of node.contents is 99

type node struct {
	contents [nodeLen]byte
	prev *node
	next *node
}

type Tape struct {
	frontNode *node    // The node in which we're writing to the front of the queue
	frontIndex int     // The index of the byte in `contents` to be written next

	readNode *node     // The node from which we're reading
	readIndex int      // The index into read.contents of the byte to be read next

	src *bufio.Reader  // The source from which we're reading
}

func NewTape(src *bufio.Reader) *Tape {
	frontNode := &node{}
	return &Tape{ frontNode: frontNode, readNode: frontNode, src: src }
}

func NewTapeFromFile(src *os.File) *Tape {
	return NewTape(bufio.NewReader(src))
}

func NewTapeFromFilename(filename string) (*Tape, os.Error) {
	file, err := os.Open(filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	return NewTapeFromFile(file), nil
}

func (this *Tape) ReadByte() (byte, os.Error) {
	// Get the byte under the read head.
	// If the source Buffer returns an error, return it.

	if this.frontNode == this.readNode && this.frontIndex == this.readIndex {
		// We reading from the front.
		// get the next from the source (and record it)
		outByte, ok := this.src.ReadByte()

		if ok != nil {
			return 0, ok
		}


		// If we're out of space (i.e. frontIndex == nodeLen)
		// then allocate a new node.
		if this.frontIndex == nodeLen {
			newNode := &node{ prev: this.frontNode }
			this.frontNode.next = newNode
			this.frontNode = newNode
			this.frontIndex = 0
		}

		// Now record it
		this.frontNode.contents[this.frontIndex] = outByte
		this.frontIndex++

		// Keep the read head up to date
		this.readNode = this.frontNode
		this.readIndex = this.frontIndex

		return outByte, nil
	}

	// Else just read from our record

	if this.readIndex == nodeLen {
		// Swap out the read node, get the next one
		this.readIndex = 0
		this.readNode = this.readNode.next
	}

	outByte := this.readNode.contents[this.readIndex]
	this.readIndex++

	return outByte, nil
}


func (this *Tape) Rewind(howMany int) (ok bool) {
	// Rewind the read head by `howMany` bytes.
	// Return `false` if rewinding this far is not possible.

	// Rewind as many nodes as we can
	for ; howMany >= nodeLen; howMany -= nodeLen {
		this.readNode = this.readNode.prev

		if this.readNode == nil {
			return false
		}
	}

	this.readIndex -= howMany

	if this.readIndex < 0 {

		this.readNode = this.readNode.prev

		if this.readNode == nil {
			return false
		}

		this.readIndex += nodeLen
	}

	return true
}