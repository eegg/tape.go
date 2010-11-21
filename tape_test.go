package tape

import (
	"testing"
	"strings"
	"bufio"
	"bytes"
	"fmt"
	)

const sample_string = "Once upon a midnight dreary, while I pondered, weak and weary,\nOver many a quaint and curious volume of forgotten lore--\nWhile I nodded, nearly napping, suddenly there came a tapping,\nAs of some one gently rapping--rapping at my chamber door.\n\"'Tis some visitor,\" I muttered, \"tapping at my chamber door--\n                                  Only this and nothing more.\"\n\nAh, distinctly I remember, it was in the bleak December,\nAnd each separate dying ember wrought its ghost upon the floor.\nEagerly I wished the morrow;--vainly I had sought to borrow\nFrom my books surcease of sorrow--sorrow for the lost Lenore--\nFor the rare and radiant maiden whom the angels name Lenore--\n                                  Nameless here for evermore.\n\nAnd the silken sad uncertain rustling of each purple curtain\nThrilled me--filled me with fantastic terrors never felt before;\nSo that now, to still the beating of my heart, I stood repeating\n'Tis some visitor entreating entrance at my chamber door--\nSome late visitor entreating entrance at my chamber door;--\n                                  This it is and nothing more."


func TestBasicRead(test *testing.T) {
	t := NewTape(bufio.NewReader(strings.NewReader(sample_string)))

	bs := bytes.NewBufferString("")

	for {
		byte, err := t.ReadByte()
		if err != nil {
			return
		} else {
			fmt.Fprintf(bs, "%c", byte)
		}
	}

	out_string := bs.String()

	if out_string != sample_string {
		test.Error("The tape did not output the input!")
	}
}


func TestSmallRewind(test *testing.T) {
	t := NewTape(bufio.NewReader(strings.NewReader("hello, world!")))	

	test_read := func(expected_char byte) {
		out_char, err := t.ReadByte()
		if err != nil {
			test.Errorf("Unexpected end of characters; expecting '%c'", expected_char)
		} else if out_char != expected_char	{
			test.Errorf("Unexpected character '%c'; expecting '%c'", out_char, expected_char)
		}
	}

	test_str_read := func(expected_str string) {
		rdr := strings.Reader(expected_str)
		for {
			expected_byte, err := rdr.ReadByte()
			if err != nil { return }
			test_read(expected_byte)
		}
	}

	test_str_read("hel")
	t.Rewind(1)
	test_str_read("llo")
	test_str_read(", world")
	test_str_read("!")
	t.Rewind(6)
	test_str_read("world!")
}