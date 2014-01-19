package main

// TEST CONVERSION of base64 to hex strings and to binary.

import (
	"base64" // Ricky's package
	"bytes"
	sysbase64 "encoding/base64" // part of golang distribution, for comparison testing.
	"fmt"                       // part of golang
)

const red = base64.Red
const boldred = base64.Boldred
const normal = base64.Normal

func main() {
	colors := []string{red, boldred, normal}
	for _, c := range colors {
		fmt.Printf("%q %T\n", c, c)
	}
	testvec := []string{
		"Zg==", "Zm8=", "Zm9v", "Zm9vYg==", "Zm9vYmE=", "Zm9vYmFy",
		"Zm9vYmFyeQ==", "Zm9vYmFyaXQ=", "Zm9vYmFyaXRz",
		"AAAAAA==", "/////w==",
		"SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t",
	}
	ansvec := []string{
		"f", "fo", "foo", "foob", "fooba", "foobar",
		"foobary", "foobarit", "foobarits",
		"\x00\x00\x00\x00", "\xFF\xFF\xFF\xFF",
		"I'm killing your brain like a poisonous mushroom",
		//"49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d",
	}

	fmt.Printf("\n %s**** Enhanced Test Vector from rfc4648: Decode ****%s\n", boldred, normal)
	for i := 0; i < len(testvec); i++ {
		teststr := []byte(testvec[i])
		ansstr := ansvec[i]
		fmt.Printf("—————>%q<—————\n", teststr)

		mine := base64.ToBinary(teststr) // My code.  Converts to binary.  Which is often text in ascii or utf8. Whatever.
		fmt.Printf("%q  ", mine)
		for _, c := range mine {
			fmt.Printf("%02x ", c)
		}
		fmt.Println()
		theirs, err := sysbase64.StdEncoding.DecodeString(testvec[i]) // golang package.
		if nil != err {
			fmt.Printf("Error %v\n", err)
		}

		if !bytes.Equal(mine, theirs) {
			fmt.Printf("\t\t %s****  %q  Not Equal!  %q  ****%s\n", boldred, mine, theirs, normal)
		}

		if !base64.EqualBS(theirs, ansstr) {
			fmt.Printf("\t\t %s**** Program Error ****%s\n", boldred, normal)
		}
	}

	// Input and expected answers for conversion base64 => hex strings:
	testhex := []string{
		"Zm9vYmFyeQ==", "Zm9vYmFyaXQ=", "Zm9vYmFyaXRz",
		"AAAAAA==", "/////w==",
		"SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t",
	}
	anshex := []string{
		"666f6f62617279", "666f6f6261726974", "666f6f626172697473",
		"00000000", "ffffffff",
		// "I'm killing your brain like a poisonous mushroom",
		"49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d",
	}
	fmt.Printf("\n %s**** Miscellaneous base64 => hex strings. Decode ****%s\n", boldred, normal)

	for i := 0; i < len(testhex); i++ {
		teststr := []byte(testhex[i])
		fmt.Printf("—————>%q<—————\n", teststr)
		ansstr := anshex[i]
		mine := base64.ToHex(teststr)
		fmt.Printf("%q\n", mine)

		if !base64.EqualBS(mine, ansstr) {
			fmt.Printf("%q\n", ansstr)
			fmt.Printf("\t\t %s**** Not Equal ****%s\n", boldred, normal)
		}

	}

}
