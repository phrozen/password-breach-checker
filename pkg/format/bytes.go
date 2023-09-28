package format

import "fmt"

func Bytes(bytes uint64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%dB", bytes)
	}
	var i int
	fbytes := float64(bytes)
	for i = -1; fbytes > 1024; i++ {
		fbytes /= 1024
	}
	return fmt.Sprintf("%.02f%cB", fbytes, "KMGTPE"[i])
}
