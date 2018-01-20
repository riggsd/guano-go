package guano

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

// Chunk is a RIFF/WAVE subchunk descriptor
type Chunk struct {
	Id   [4]byte
	Size uint32
}

// IdStr returns the subchunk ID as an ASCII string
func (c *Chunk) IdStr() (id string) {
	return string(c.Id[:])
}

func (c *Chunk) String() string {
	return fmt.Sprintf("%T{Id:%s Size:%X}", *c, string(c.Id[:]), c.Size)
}

// RiffHeader is the initial header of a RIFF file
type RiffHeader struct {
	Chunk
	Format [4]byte
}

// FormatStr returns the RIFF format as an ASCII string
func (h *RiffHeader) FormatStr() (format string) {
	return string(h.Format[:])
}

func (h *RiffHeader) String() string {
	return fmt.Sprintf("%T{Id:%s Size:%X Format:%s}", *h, string(h.Id[:]), h.Size, h.Format)
}

// Guano represents a .WAV file with GUANO metadata
type Guano struct {
	Filename string
	Fields   map[string]string
}

// New constructs a new fully-initialized Guano struct
func New() *Guano {
	return &Guano{Fields: make(map[string]string)}
}

// Read reads a named file, parsing its GUANO metadata
func Read(filename string) (*Guano, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}
	defer f.Close()
	return ReadFile(f)
}

// ReadFile reads a File, parsing its GUANO metadata
func ReadFile(f *os.File) (*Guano, error) {

	// parse the RIFF file header
	riffHeader := new(RiffHeader)
	err := binary.Read(f, binary.LittleEndian, riffHeader)
	if err != nil {
		return nil, err
	}
	if id := riffHeader.IdStr(); id != "RIFF" {
		return nil, fmt.Errorf("expected \"RIFF\" header ID, found %q", id)
	}
	if format := riffHeader.FormatStr(); format != "WAVE" {
		return nil, fmt.Errorf("expected RIFF format \"WAVE\", found %q", format)
	}
	//log.Printf("%v\n", riffHeader) // DEBUG

	// parse the WAVE fmt_ chunk
	fmtChunk := new(Chunk)
	err = binary.Read(f, binary.LittleEndian, fmtChunk)
	if err != nil {
		return nil, err
	}
	if id := fmtChunk.IdStr(); id != "fmt " {
		return nil, fmt.Errorf("expected \"fmt \" chunk, found %q", id)
	}
	//log.Printf("%v\n", fmtChunk) // DEBUG
	fmtData := make([]byte, fmtChunk.Size)
	n, err := f.Read(fmtData)
	if err != nil {
		return nil, err
	}
	if n != int(fmtChunk.Size) {
		return nil, fmt.Errorf("expected %d bytes, read %d", fmtChunk.Size, n)
	}

	var guanoData []byte

	// parse all the remaining chunks...
	for {
		chunk := new(Chunk)
		err = binary.Read(f, binary.LittleEndian, chunk)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		data := make([]byte, chunk.Size)
		n, err := f.Read(data) // FIXME: word align
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n != int(chunk.Size) {
			return nil, fmt.Errorf("expected %d bytes, read %d", chunk.Size, n) // FIXME: word align
		}
		//log.Printf("%v\n", chunk) // DEBUG

		if chunk.IdStr() == "guan" {
			guanoData = data
		}

	}

	if guanoData == nil {
		return nil, fmt.Errorf("no \"guan\" chunk found")
	}
	return ParseGuanoString(string(guanoData))
}

// ParseGuanoString parses a UTF-8 string as GUANO metadata
func ParseGuanoString(s string) (*Guano, error) {
	g := New()

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t\r\n\x00")
		if line == "" {
			continue
		}
		i := strings.Index(line, ":")
		if i < 1 {
			return nil, fmt.Errorf("failed parsing GUANO field %q", line)
		}
		k, v := strings.TrimSpace(line[:i]), strings.TrimSpace(line[i+1:])
		g.Fields[k] = v
	}

	return g, nil
}
