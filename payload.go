// Package payload provides functions to read and write an arbitrary
// set of keys and values from and to a file. The payload may be appended
// to any file and will be identifiable by reading the MAGIC string from
// the end of it.
package payload

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// MagicError represents a failure to read the magic string from the payload.
var MagicError = errors.New("Magic string is does not match the expected string")

// MAGIC is the magic string appended at the end of the payload.
const MAGIC = "PAYLOADS"

// A Payload is a representation of a payload to append or read.
type Payload map[string][]byte

func writechunk(w io.Writer, data []byte) error {
	err := binary.Write(w, binary.LittleEndian, int64(len(data)))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func readchunk(r io.Reader) ([]byte, error) {
	var l int64
	if err := binary.Read(r, binary.LittleEndian, &l); err != nil {
		return nil, err
	}
	data := make([]byte, l)
	_, err := io.ReadFull(r, data)
	return data, err
}

// Load will load a payload appended to the end of the ReadSeeker.
// It may return a MagicError if the magic string is missing.
func Load(r io.ReadSeeker) (Payload, error) {
	p := make(Payload)

	if _, err := r.Seek(-int64(len(MAGIC)), 2); err != nil {
		return nil, err
	}

	magic := make([]byte, len(MAGIC))
	if _, err := io.ReadFull(r, magic); err != nil {
		return p, MagicError
	}

	if string(magic) != MAGIC {
		return p, MagicError
	}

	if _, err := r.Seek(-int64(len(MAGIC)+8), 2); err != nil {
		return nil, err
	}

	var offset int64
	if err := binary.Read(r, binary.LittleEndian, &offset); err != nil {
		return nil, err
	}
	if _, err := r.Seek(-(offset + int64(8+len(MAGIC))), 2); err != nil {
		return nil, err
	}

	var n int64
	if err := binary.Read(r, binary.LittleEndian, &n); err != nil {
		return nil, err
	}

	for i := int64(0); i < n; i++ {
		k, err := readchunk(r)
		if err != nil {
			return nil, err
		}

		v, err := readchunk(r)
		if err != nil {
			return nil, err
		}

		p[string(k)] = v
	}
	return p, nil
}

// Load, given a path, will resolve any symbolic links of that path
// and return the payload of that path.
func LoadFile(path string) (Payload, error) {
	path, err := filepath.EvalSymlinks(os.Args[0])
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Load(f)
}

// LoadSelf will load a payload from Args[0]
func LoadSelf() (Payload, error) {
	return LoadFile(os.Args[0])
}

// IgnoreMissing, wrapped around a function returning a Payload and an Error,
// will ignore any MagicError thrown and instead return an empty Payload.
func IgnoreMissing(p Payload, err error) (Payload, error) {
	if err == MagicError {
		return p, nil
	} else {
		return p, err
	}
}

// Dump will write a Payload to an io.Writer
func (p Payload) Dump(w io.Writer) error {
	keys := make([]string, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	total := int64(8)
	if err := binary.Write(w, binary.LittleEndian, int64(len(keys))); err != nil {
		return err
	}

	for _, key := range keys {
		total += int64(8 + len(key) + 8 + len(p[key]))

		if err := writechunk(w, []byte(key)); err != nil {
			return err
		}
		if err := writechunk(w, p[key]); err != nil {
			return err
		}
	}

	if err := binary.Write(w, binary.LittleEndian, total); err != nil {
		return err
	}
	if _, err := io.WriteString(w, MAGIC); err != nil {
		return err
	}

	return nil
}
