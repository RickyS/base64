package base64

// Data2Base64 creates a new byte array with base64 content copied from an
// array of (possibly binary) input data.  Each input byte results
// in 8 bits of information, spread over several output bytes.  This is because
// the output has 6 bits per byte, in base64 ascii format; plus padding.
func Data2Base64(data []byte) []byte {
	dataLen := len(data)
	numBitsIn := 8 * dataLen
	sextLen := (5 + (numBitsIn)) / 6 // number of sextets from data.
	addLen := equals[dataLen%3]      // number of padding chars '='
	result := make([]byte, sextLen+addLen)
	//fmt.Printf("result len %d\n", len(result))
	chan6 := make(chan byte) // byte is synonym for uint8
	go SixBitsBinary(data, chan6)
	i := 0
	for next6 := range chan6 {
		result[i] = Num2Base64(next6) // Num2Base64 cnvts small 6-bit integer
		// to an ascii byte in base64 coding.
		i++
	}
	for k := 0; k < addLen; k++ {
		result[i] = '='
		i++
	}
	return result
}

//sixBits eats bits from the input bytes, returning, by chan, 6 bits at a time, for base64.
//The data is not copied to any intermediate array.
func SixBitsBinary(data []byte, channel6 chan<- uint8) {
	var val, next6 uint8
	var num6, maxNumBits int
	var lenData = len(data)

	maxNumBits = 8 * lenData // Each byte is 8 bits of binary data.
	num6 = 0

	// 3 bytes is 24 bits, which is 4 nibbles of 6 bits each.
	i := 0
	for (6 * num6) < maxNumBits {
		//fmt.Printf("num6 %2d, maxNumBits %3d\n", num6, maxNumBits)
		normal := i < lenData
		switch num6 % 4 {
		case 0:
			val = nextDataByte(data, &i, normal)
			next6 = val >> 2
			channel6 <- next6
		case 1:
			next6 = ((3 & val) << 4) // upper 2 bits of the 6-bit nibble
			val = nextDataByte(data, &i, normal)
			next6 |= ((0xF0 & val) >> 4)
			channel6 <- next6
		case 2:
			next6 = ((0xF & val) << 2)
			val = nextDataByte(data, &i, normal)
			next6 |= (0xC0 & val) >> 6
			channel6 <- next6
		case 3:
			channel6 <- 0x3F & val
		}
		num6 += 1 // number of 6-bit units sent.
		//fmt.Printf("num6 %d\n", num6)
	}
	close(channel6)
	return

}

func nextDataByte(data []byte, i *int, normal bool) byte {
	if normal {
		retval := data[*i]
		(*i)++
		return retval
	} else {
		return 0
	}
}
