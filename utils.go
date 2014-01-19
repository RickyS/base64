package base64

import (
	"fmt"
	"log"
	"os"
)

// Utility functions: HexPair2Num|Bin2HexPair|Fatalf|Bugf

// Xterm colored output.  Hope it works okay on your system.  Used by test code.
const Red = "\x1b[31m"
const Boldred = "\x1b[1m\x1b[31m"
const Boldgreen = "\x1b[1m\x1b[32m"
const Boldmagenta = "\x1b[1m\x1b[35m"
const Boldyellow = "\x1b[1m\x1b[33m"
const BoldyellowOnblack = "\x1b[40m\x1b[1m\x1b[33m"

const Boldblue = "\x1b[1m\x1b[34m"
const Normal = "\x1b(B\x1b[m\x1b[0m"

var std *log.Logger

func init() {
	std = log.New(os.Stderr, "", log.Lshortfile) // whatever.go:48 message...
}

// Bugf formats according to a format specifier and writes to standard output.
// Colored, stdout, no file name (NOT a logger)
func Bugf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, Boldmagenta+format+Normal, a...)
}

//**********  LOG TO STDERR:  *************//

// Fatalf is a logger equivalent to Printf(STDERR) followed by a call to os.Exit(1).
// This is a hack from Logger.  'std' is a logger I create in an init(), above
func Fatalf(format string, v ...interface{}) {
	std.Output(2, fmt.Sprintf(Boldred+format+Normal+"\n", v...))
	os.Exit(1)
}

// Warningf is a logger equivalent to Fprintf(STDERR), but in color...
func Warningf(format string, v ...interface{}) {
	std.Output(2, fmt.Sprintf(Boldred+format+Normal+"\n", v...))
}

////////////

// Convert a character in hex to a binary integer.  "F" ==> 15, which is also 0x0F
func HexChar2Num(c byte) uint8 {
	switch {
	case (c >= 'a') && (c <= 'f'):
		return 10 + (c - 'a')
	case (c >= 'A') && (c <= 'F'):
		return 10 + (c - 'A')
	case (c >= '0') && (c <= '9'):
		return c - '0'
	default:
		Fatalf("Invalid hex character %02x, %c\n", c, c)
	}
	return 0 // Bogus return to satisfy compiler.
}

// Convert the 2 hex chars to a number.  "a1" ==> 161, which is also 0xA1
// Ought to be called HexPair2Binary
func HexPair2Num(hi byte, lo byte) uint8 {
	return (16 * HexChar2Num(hi)) + HexChar2Num(lo)
}

// HexBin2Char converts a 4-bit integer to a hex character. 0xB => 'b'
func HexBin2Char(hex byte) byte {
	switch {
	case (0 <= hex) && (hex <= 9):
		return '0' + hex
	case (10 <= hex) && (hex <= 15):
		return 'a' + hex - 10
	default:
		Bugf("Hex byte bigger than 15 0x%02X\n", hex)
	}
	return 0 // bogus return for the compiler.
}

// Bin2HexPair converts any byte to 2 ascii hex bytes. 0xAB => ("A", "B")
func Bin2HexPair(next8 byte) (hi byte, lo byte) {
	hi = HexBin2Char(next8 >> 4)
	lo = HexBin2Char(next8 & 0xF)
	return hi, lo
}

///////////////////////  Used only by test program decode.go: ////////////////////
// Compare a byte array to a string.  Byte-by-byte. Not byte-by-rune.
func EqualBS(a []byte, bs string) (eq bool) {
	b := []byte(bs) // so for ... range gets byte-by-byte
	l := len(a)
	if l != len(b) {
		return false
	}

	for i := 0; i < l; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
