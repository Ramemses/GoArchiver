package shanon_fano

import (
	"sort"
	"math"
	"fmt"
	"strings"

	"archiver/lib/compression/vlc/table"
)

type Generator struct{}

type charStat map[rune]int

type code struct{
	Char rune
	Quantity int
	Bits uint32
	Size int
}	

type encodingTable map[rune]code

func NewGenerator() Generator{
	return Generator{}
}


func (g Generator) NewTable(text string) table.EncodingTable{
			
		encTable := build(text)	
		return table.EncodingTable(encTable.Export())
}

func (et encodingTable) Export() map[rune]string{
	res := make(map[rune]string)

	for k, v := range et{
		byteStr := fmt.Sprintf("%b", v.Bits)

		lenDiff := v.Size - len(byteStr)

		if lenDiff > 0{
			byteStr = strings.Repeat("0", lenDiff) + byteStr
		}
		res[k] = byteStr
	}

	return res
}


func build(str string) encodingTable{
	
	stat := newCharStat(str)
	codes := make([]code, 0, len(stat))

	for ch, qty := range stat{
		codes = append(codes, code{
				Char: ch,
				Quantity: qty,
			})
	}

	sort.Slice(codes, func(i, j int) bool {
		if codes[i].Quantity != codes[j].Quantity{
			return codes[i].Quantity > codes[j].Quantity
		}
		return codes[i].Char < codes[j].Char	
	})


	assignCodes(codes)

	res := make(encodingTable)
	
	for _, code := range codes{
		res[code.Char] = code
	}
	
	return res
}


func assignCodes(codes []code){
	if len(codes) == 0{
		return 
	}
	if len(codes) == 1{
		if codes[0].Size == 0{
			codes[0].Bits <<= 1
			codes[0].Size++

		}
		return 
	}

	//divide codes
	divider := bestDividePosition(codes)


	//add 0 or 1
	for i := 0; i < len(codes); i++{
		codes[i].Bits <<= 1
		codes[i].Size++

		if i>= divider{
			codes[i].Bits |= 1
			
		}

	}

	assignCodes(codes[:divider])
	assignCodes(codes[divider:])

}


func bestDividePosition(codes []code)int{
	// a b | c f d => diff < newDiff

	left := 0
	total := 0
	prevDiff := math.MaxInt	
	bestPosition := 0


	for _, code := range codes{
		total += code.Quantity
	}	


	for i := 0; i < len(codes) - 1; i++{
		left += codes[i].Quantity
		right := total - left
		
		diff := abs(right - left)
		if diff >= prevDiff{
			break
		}

		prevDiff = diff
		bestPosition = i + 1
	}
	
	return bestPosition
}

func abs(x int)int{
	if x < 0{
		x = -x
	}
	return x
}


func newCharStat(text string) charStat{

	res := make(charStat)

	for _, ch := range text{
		res[ch]++
	}
	
	return res
}
