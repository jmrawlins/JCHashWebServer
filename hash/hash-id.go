package hash

import (
	"strconv"
)

type HashId uint64

func HashIdFromString(str string) (HashId, error) {
	id, err := strconv.ParseUint(str, 10, 64)
	return HashId(id), err
}
