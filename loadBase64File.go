package base64

import (
	"bufio"
	"fmt"
	"os"
)

// Read in the base64 file, line by line, converting base64 to raw data,
// and return the byte array of raw data.
func LoadConvertFile(filename string) (data []byte) {
	fileB64, err := os.Open(filename) // For read access.
	if err != nil {
		Fatalf("Error: %v\n", err)
	}
	defer fileB64.Close()

	scanner := bufio.NewScanner(fileB64)
	numBytes := 0
	numLines := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		d := ToBinary(line) // Convert bas64 to original data
		data = append(data, d...)
		numLines++            // On *nix, number of line-ending chars.
		numBytes += len(line) // Number of base64 bytes in the file.  Not output size.
	}
	errs := scanner.Err()
	if errs != nil {
		Fatalf("Reading %s: %v ", filename, errs)
	}
	fmt.Printf("Read in %s. Bytes+lines %6d  ——>  %6d bytes in memory.\n",
		filename, numLines+numBytes, len(data))
	// On *nix, numLines+numBytes is the file size.
	// return a slice (reference) to the data.
	return data[:]
}
