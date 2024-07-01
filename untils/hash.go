package untils

import (
	"crypto/sha256"
	"fmt"
)

func Sha256(stringStr string) string {
	hash := sha256.New()
	hash.Write([]byte(stringStr))

	return fmt.Sprintf("%x", hash.Sum(nil))
}
