package helpers

// Everything minecraft related

import (
	"crypto/md5"
	"fmt"
	"strings"
)

// Username to Offline UUID (with hyphen)
func OfflineUUID(username string) string {
	sum := md5.Sum([]byte("OfflinePlayer:" + username))
	sum[6] = (sum[6] & 0x0f) | 0x30
	sum[8] = (sum[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", sum[0:4], sum[4:6], sum[6:8], sum[8:10], sum[10:16])
}

// Removes hyphens from UUID
func RemoveHyphensUUID(uuid string) string {
	return strings.ToLower(strings.ReplaceAll(uuid, "-", ""))
}
