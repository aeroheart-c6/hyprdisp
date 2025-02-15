package profiles

import (
	"encoding/hex"
	"fmt"
	"strings"

	"crypto/sha3"
)

type Display struct {
	ID   string
	Name string
	Desc string
}

func (d Display) String() string {
	return fmt.Sprintf("[%s, %s, %s]",
		d.ID,
		d.Name,
		d.Desc,
	)
}

func CreateID(displays []Display) {
	var identifiers []string = make([]string, 0, len(displays))

	for _, display := range displays {
		identifiers = append(identifiers, display.String())
	}

	var body string = strings.Join(identifiers, "|")
	var hash [32]byte = sha3.Sum256([]byte(body))

	fmt.Printf("Hash of the display combination: %v is %v\n",
		body,
		hex.EncodeToString(hash[:])[:12],
	)
}
