package vlc


import (
	"fmt"
	"strings"
	"unicode/utf8"
	"strconv"
)

const chunkSize = 8


type BinaryChunk string

type BinaryChunks []BinaryChunk


// splitByChunks splits binary string by chunks with given size,
// i.g. : '100101011001010110010101' -> '10010101 10010101 10010101'
func splitByChunks(bStr string, chunkSize int) BinaryChunks {
	strLen := utf8.RuneCountInString(bStr)
	chunksCount := strLen / chunkSize
	

	if strLen % chunkSize != 0{
		chunksCount++
	}
	res := make(BinaryChunks, 0, chunksCount)
	
	var buf strings.Builder

	for i, ch := range bStr{
		buf.WriteString(string(ch))
		
		if (i+1) % chunkSize == 0{
			res = append(res, BinaryChunk(buf.String()))
			buf.Reset()	
		}
	}
	
	if buf.Len() != 0{
		lastChunk := buf.String()
		lastChunk += strings.Repeat("0", chunkSize-len(lastChunk))
		
		res = append(res, BinaryChunk(lastChunk))

	}

	return res;
}

func (bcs BinaryChunks) ToString() string{
	const sep = " "
	bLen := len(bcs)
	

	switch bLen{
		case 0:
			return ""
		case 1:
			return string(bcs[0])

	}

	var buf strings.Builder

	for _, bChunk := range bcs{
		buf.WriteString(string(bChunk))
	}

	return buf.String()
}

func (bcs BinaryChunks) Bytes() []byte{
	res := make([]byte, 0, len(bcs))

	for _, bc := range bcs{
		res = append(res, bc.Byte())
	}

	return res
}



func (bc BinaryChunk) Byte() byte{
	num, err := strconv.ParseUint(string(bc), 2, chunkSize)
	if err != nil{
		panic("Can't parse binary chunk...")
	
	}

	return byte(num)
}



func NewBinChunks(codes []byte) BinaryChunks{
	res := make(BinaryChunks, 0, len(codes))	

	for _, code := range codes{
		res = append(res, NewBinChunk(code))
	}

	return res
}


func NewBinChunk(code byte) BinaryChunk{
	return BinaryChunk(fmt.Sprintf("%08b", code))
}
