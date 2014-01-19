// Package base64 converts data to ascii strings in base64 format.
// see http://www.ietf.org/rfc/rfc4648.txt, and http://en.wikipedia.org/wiki/Base64
// Inputs are accepted as hex characters in memory, either by a pointer to the data returning a stream
// of bytes over a channel, or an array of ascii hex text returning an array of ascii base64 bytes.
// Also, another function accepts an array of raw memory bytes and returns a base64 array.
package base64

// 1. Convert hex to base64 and back.
// The string:
//   49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d
//
// should produce:
//   SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t
//
// Now use this code everywhere for the rest of the exercises. Here's a
// simple rule of thumb:
//   Always operate on raw bytes, never on encoded strings. Only use hex
//   and base64 for pretty-printing.

import (
	"log"
)

var equals []int = []int{0, 2, 1} // For counting '='

// NewHex2Base64 creates a new byte array with base64 content copied from the input string,
// which is ASSUMED to contain hex values in char form: "AB00FF1D".  Each input
// PAIR of bytes generates 8 bits of output, plus padding
func NewHex2Base64(data []byte) []byte {
	strLen := len(data)
	if 0 != (1 & strLen) {
		log.Fatalf("NewHex2Base64: data length must be even, is %d\n", strLen)
	}
	halfLen := strLen / 2
	sextLen := (5 + (halfLen * 8)) / 6 // number sextets
	addLen := equals[halfLen%3]        // number '='
	result := make([]byte, sextLen+addLen)
	chan6 := make(chan uint8)
	go SixBits(data, chan6)
	i := 0
	for next6 := range chan6 {
		result[i] = Num2Base64(next6)
		i++
	}

	// Add Padding '='
	for k := 0; k < addLen; k++ {
		result[i] = '='
		i++
	}
	return result
}

//sixBits eats bits from the string, returning, by chan, 6 bits at a time, for base64.
//The string is not copied.  The string is ASSUMED to contain ascii bytes representing
// hexadecimal values, like "001DFFea".
// This DOES NOT ADD PADDING '='. See above.
func SixBits(data []byte, channel6 chan<- uint8) {
	var val, next6 uint8
	var num6, maxNumBits int
	var lenData = len(data)

	if 0 != lenData%2 {
		log.Fatalf("sixBits: data length must be even, is %d\n", lenData)
	}

	maxNumBits = 4 * lenData // Each text char is a hex digit is 4 bits.
	num6 = 0

	// 3 bytes is 24 bits, which is 4 nibbles of 6 bits each.
	i := 0
	for (6 * num6) < maxNumBits {
		normal := i < lenData
		switch num6 % 4 {
		case 0:
			val = NextByte(data, i, normal)
			i += 2
			next6 = val >> 2
			channel6 <- next6
		case 1:
			next6 = ((3 & val) << 4) // upper 2 bits of the 6-bit nibble
			val = NextByte(data, i, normal)
			i += 2
			next6 |= ((0xF0 & val) >> 4)
			channel6 <- next6
		case 2:
			next6 = ((0xF & val) << 2)
			val = NextByte(data, i, normal)
			i += 2
			next6 |= (0xC0 & val) >> 6
			channel6 <- next6
		case 3:
			channel6 <- 0x3F & val
		}
		num6 += 1 // number of 6-bit units sent.
	}
	close(channel6)

	return

}

// NextByte combines 2 hex chars into a binary byte. "AB" => 0xAB
func NextByte(data []byte, i int, normal bool) uint8 {
	if normal {
		return HexPair2Num(data[i], data[i+1])
	} else {
		return 0 //  Extra zero bits to tack onto the end
	}
}

// Convert a small (6-bit) integer to a printable base64 character
// see: http://en.wikipedia.org/wiki/Base64
func Num2Base64(num uint8) byte {
	if 0 != (0xC0 & num) {
		Fatalf("Illegal base64 value %02x\n", num)
	}
	switch {
	case (0 <= num) && (num <= 25):
		return 'A' + num
	case (26 <= num) && (num <= 51):
		return 'a' + num - 26
	case (52 <= num) && (num <= 61):
		return '0' + num - 52
	case 62 == num:
		return '+'
	case 63 == num:
		return '/' // Not good for urls or file names.
	}
	return 0
}
