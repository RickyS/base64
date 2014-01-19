package base64

import (
	"fmt"
)

// printBase64Chart just prints a chart of the base64 values.
func PrintBase64Chart() {
	var val, line, num uint8
	fmt.Printf("-val- b64   -val- b64   -val- b64   -val- b64\n")
	for line = 0; line < 16; line++ {
		for num = 0; num < 49; num += 16 {
			val = line + num
			fmt.Printf("%2d %02x  %c   |", val, val, Num2Base64(val))
		}
		fmt.Println()
	}
	fmt.Println()
}
