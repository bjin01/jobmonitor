package pkg_updates

import "runtime"

func getGoroutineID() uint64 {
	var buffer [31]byte
	written := runtime.Stack(buffer[:], false)
	index := 10
	negative := buffer[index] == '-'
	if negative {
		index = 11
	}
	id := uint64(0)
	for index < written {
		byte := buffer[index]
		if byte == ' ' {
			break
		}
		if byte < '0' || byte > '9' {
			panic("could not get goroutine ID")
		}
		id *= 10
		id += uint64(byte - '0')
		index++
	}
	if negative {
		id = -id
	}
	return id
}
