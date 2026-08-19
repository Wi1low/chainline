package types

import (
	"encoding/hex"
	"fmt"
)

type Address [20]uint8

func (a Address) ToSlice() []byte {
	return a[:]
}
func (a Address) String() string {
	return hex.EncodeToString(a.ToSlice())
}
func AddressFromBytes(data []byte) Address {
	if len(data) != 20 {
		panic(fmt.Sprintf("given bytes with length %d should be 20", len(data)))
	}
	return Address(data)
}
