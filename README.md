payload
=======

`payload` is a Go package that allows you to attach a key/value store
payload to the end of any file. It can then read the payload into a
`map[string][]byte` regardless of what the data before it is.

cmd/genpayload
--------------

The genpayload command takes any number of file paths as arguments,
traversing directories recursively and attaching regular files as-is, to
build a payload. The output may be appended to a file to attach the
payload to it.

PACKAGE DOCUMENTATION
---------------------

### package payload

    import "github.com/boomlinde/payload"

Package payload provides functions to read and write an arbitrary set of
keys and values from and to a file. The payload may be appended to any
file and will be identifiable by reading the MAGIC string from the end
of it.

### CONSTANTS

    const MAGIC = "PAYLOADS"

MAGIC is the magic string appended at the end of the payload.

### VARIABLES

    var MagicError = errors.New("Magic string is does not match the expected string")

MagicError represents a failure to read the magic string from the
payload.

### TYPES

    type Payload map[string][]byte

A Payload is a representation of a payload to append or read.

    func IgnoreMissing(p Payload, err error) (Payload, error)

IgnoreMissing, wrapped around a function returning a Payload and an
Error, will ignore any MagicError thrown and instead return an empty
Payload.

    func Load(r io.ReadSeeker) (Payload, error)

Load will load a payload appended to the end of the ReadSeeker. It may
return a MagicError if the magic string is missing.

    func LoadFile(path string) (Payload, error)

Load, given a path, will resolve any symbolic links of that path and
return the payload of that path.

    func LoadSelf() (Payload, error)

LoadSelf will load a payload from Args\[0\]

    func (p Payload) Dump(w io.Writer) error

Dump will write a Payload to an io.Writer
