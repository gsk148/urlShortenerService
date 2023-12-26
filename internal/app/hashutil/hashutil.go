package hashutil

import (
	"crypto/md5"
	"encoding/base64"
)

func Encode(data []byte) string {
	hash := md5.Sum(data)
	base64Hash := base64.RawURLEncoding.EncodeToString(hash[:])
	return base64Hash[:7]
}
