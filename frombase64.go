// These file convert from base64 to ...
package base64

import ()

// ToHex creates a new byte array (string) of ascii hex characters (0-9a-f)
// from base64-encoded ascii input text
func ToHex(data []byte) (result []byte) {
	bytesNeeded := actualLength(data)
	chan8 := make(chan byte)
	result = make([]byte, 2*bytesNeeded)

	go EightBits(data, chan8, bytesNeeded)
	i := 0
	for next8 := range chan8 {
		result[i], result[1+i] = Bin2HexPair(next8)
		i += 2
	}
	return result
}

// ToBinary creates a new byte array from base64-encoded ascii input text
// Each input byte contributes 6 bits, each output byte contains 8 bits.
func ToBinary(data []byte) (result []byte) {
	bytesNeeded := actualLength(data)
	chan8 := make(chan byte)
	result = make([]byte, bytesNeeded)
	go EightBits(data, chan8, bytesNeeded)
	i := 0
	for next8 := range chan8 {
		result[i] = next8
		i++
	}
	return result
}

// actualLength computes the number of output binary bytes that will be output
// when converting from base64 to 'binary'.  Where 'binary' is any format at all.
// That is, use the number of '=' chars at the end to calculate amt of actual data.
func actualLength(data []byte) (bytesNeeded int) {
	dataLen := len(data)
	if 0 == dataLen {
		Fatalf("No base64 data given to decode\n")
	}
	// Each input string is a series of 4-byte Groups, yielding 3 bytes out
	// followed by one of three cases:
	//  A)  2 ascii bytes and 2 '=' yielding 1 output byte   [1]
	//  B)  3 ascii bytes and 1 '=' yielding 2 output bytes  [2]
	//  C)  Nothing, yielding nothing.                       [0]

	realDataLen := dataLen
	if 0 < realDataLen {
		for '=' == data[realDataLen-1] { // Reverse index until we find non-padding.
			realDataLen--
		}
	}
	numEquals := dataLen - realDataLen
	numGroupBytes := (realDataLen / 4) * 3
	numAddBytes := []int{0, 2, 1} // bytes to add per # '='
	bytesNeeded = numGroupBytes + numAddBytes[numEquals]

	//fmt.Printf("8bits: %d + %d = %3d bytes needed\n", numGroupBytes, numAddBytes[numEquals], bytesNeeded)
	return bytesNeeded
}

// GetBinaryStream returns a readonly channel with binary data converted from the base64 input string.
// 1) Works okay.  2) OBSOLETE.
func GetBinaryStream(data []byte) (channel8 chan byte) {
	bytesNeeded := actualLength(data)
	channel8 = make(chan byte)
	go EightBits(data, channel8, bytesNeeded)
	return channel8
}

// EightBits reads ascii base64 string converting each 8 bits to a binary byte.
// data is the input.  The output is binary bytes out on channel8. The number of bytes
// output is computed in advance and is bytesNeeded.
func EightBits(data []byte, channel8 chan<- byte, bytesNeeded int) {
	lenData := len(data)
	numIn := 0     // Number of bytes input
	dataBytes := 0 // Number of bytes scanned
	outBytes := 0  // Num bytes sent out on channel

	var val byte
	var next8 byte

	// dataBytes is nearly always == numIn.
	// numIn is incremented for 0 padding at the end.
	// dataBytes is only incremented for data bytes.
	val = 0
	for outBytes < bytesNeeded {
		normal := numIn < lenData
		switch dataBytes % 4 {
		case 0:
			val = nextIn(data, &dataBytes, normal)
			next8 = val << 2 // use 6 of 6
			// so far, 6 out of 8 bits put in next8
		case 1:
			val = nextIn(data, &dataBytes, normal)
			next8 |= 3 & (val >> 4) // use 2 of 6
			channel8 <- next8
			outBytes += 1
		case 2:
			// 4 left over in 4 lo bits of val
			next8 = val << 4
			//fmt.Printf("2:8 %02x— ", next8)
			val = nextIn(data, &dataBytes, normal)
			next8 |= 0xF & (val >> 2) // use 4 of 6
			channel8 <- next8
			outBytes += 1
			// val has 2 lo bits, # 0 & 1
		case 3:
			next8 = val << 6
			val = nextIn(data, &dataBytes, normal)
			next8 |= val & 0x3F
			channel8 <- next8
			outBytes += 1
		}
		//fmt.Printf(" [%d] <%02x, %c>\n", (dataBytes%4)-1, next8, next8)
		numIn += 1
	}
	close(channel8)
	return
}

// nextIn gets 6 bits from the input, as ascii byte base64.
// returns padding, too, as zero.
func nextIn(data []byte, dataBytes *int, normal bool) byte {
	if normal {
		retval := Base64ToBinary(data[*dataBytes])
		//fmt.Printf("'%c' -> 0x%02x  ", data[*dataBytes], retval)
		(*dataBytes)++
		return retval
	} else {
		return 0
	}
}

// Base64ToBinary converts ascii byte to number.  'm' ==> 38
func Base64ToBinary(ascii byte) byte {
	switch {
	case ('A' <= ascii) && (ascii <= 'Z'):
		return ascii - 'A'
	case ('a' <= ascii) && (ascii <= 'z'):
		return 26 + ascii - 'a'
	case ('0' <= ascii) && (ascii <= '9'):
		return 52 + ascii - '0'
	case '+' == ascii:
		return 62
	case '-' == ascii: // The alternative safe alphabet "base64url"
		return 62
	case '/' == ascii:
		return 63
	case '_' == ascii: // The alternative safe alphabet "base64url"
		return 63
	case '=' == ascii:
		return 0 // Right??
	default:
		Bugf("The Byte is not base64: %02X  %c\n", ascii, ascii)
	}
	return 0 // Bogus for compiler
}
