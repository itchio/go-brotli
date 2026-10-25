package enc

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/itchio/dskompress/brotli"
)

func TestSimpleString(T *testing.T) {
	testCompress([]byte("Hello Hello Hello, Hello Hello Hello"), T)
}

func TestShortString(T *testing.T) {
	s := []byte("The quick brown fox")
	l := len(s)

	// Brotli will not compress arrays shorter than 3 characters
	for ; l > 3; l-- {
		testCompress(s[:l], T)
	}
}

func testCompress(s []byte, T *testing.T) {
	T.Logf("Compressing: %s\n", s)

	encoded, cerr := CompressBuffer(s, nil)
	if cerr != nil {
		T.Error(cerr)
	}

	check("Buffer decompress", s, decompress(bytes.NewReader(encoded), T), T)
}

// Run roundtrip tests from Brotli repository
func TestRoundtrip(T *testing.T) {
	inputs := []string{
		"../testdata/alice29.txt",
		"../testdata/asyoulik.txt",
		"../testdata/lcet10.txt",
		"../testdata/plrabn12.txt",
		"../c/enc/encode.c",
		"../c/common/dictionary.h",
	}

	for _, file := range inputs {
		input, err := ioutil.ReadFile(file)
		if err != nil {
			T.Fatal(err)
		}

		for _, quality := range []int{1, 6, 9, 11} {
			T.Logf("Roundtrip testing %s at quality %d", file, quality)

			options := &BrotliWriterOptions{
				Quality: quality,
			}

			bro, err := CompressBuffer(input, options)
			if err != nil {
				T.Error(err)
			}
			T.Logf("  Compressed from %d to %d bytes, %.1f%%", len(input), len(bro), (float32(len(bro))/float32(len(input)))*100)

			check("Buffer decompress", input, decompress(bytes.NewReader(bro), T), T)

			// Stream compress
			buffer := new(bytes.Buffer)
			testCompressStream(options, input, buffer, T)
			check("Stream compress", input, decompress(buffer, T), T)

			// Stream roundtrip
			reader, writer := io.Pipe()
			go testCompressStream(options, input, writer, T)
			check("Stream roundtrip", input, decompress(reader, T), T)
		}
	}
}

func decompress(r io.Reader, T *testing.T) []byte {
	br, err := brotli.NewReader(r, nil)
	if err != nil {
		T.Fatal(err)
	}

	out, err := ioutil.ReadAll(br)
	if err != nil {
		T.Error(err)
	}

	return out
}

func testCompressStream(options *BrotliWriterOptions, input []byte, writer io.Writer, T *testing.T) {
	bwriter := NewBrotliWriter(writer, options)
	n, err := bwriter.Write(input)
	if err != nil {
		T.Error(err)
	}

	err = bwriter.Close()
	if err != nil {
		T.Error(err)
	}

	if n != len(input) {
		T.Error("Not all input was consumed")
	}
}

func check(test string, input, output []byte, T *testing.T) {
	if len(input) != len(output) {
		T.Errorf("  %s: Length of decompressed output (%d) doesn't match input (%d)", test, len(output), len(input))
	}

	if !bytes.Equal(input, output) {
		T.Errorf("  %s: Input does not match decompressed output", test)
	}
}
