package utils

import (
	"strconv"

	"github.com/spaolacci/murmur3"
)

func DomainHash(domain string) string {
	res := murmur3.Sum64([]byte(domain))
	sum := strconv.FormatUint(res, 16)
	return sum
}
