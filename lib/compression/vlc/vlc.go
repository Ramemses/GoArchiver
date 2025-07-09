package vlc

import (
	"strings"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"archiver/lib/compression/vlc/table"
)


type EncoderDecoder struct{
	tblGenerator table.Generator	
}
	
func New(tblGenerator table.Generator) EncoderDecoder{
	return EncoderDecoder{tblGenerator: tblGenerator}
}


func (ed EncoderDecoder) Encode(str string) []byte {

	
//haffman or shanon-fano table
	table := ed.tblGenerator.NewTable(str)
		
	encoded := encodeBin(str, table)

	return buildEncodeFile(table, encoded)
}


// prepareText prepares text to be fit for encode
// changes upper case letter to: ! + lower case letter
// i.g.: M -> !m


func buildEncodeFile(tbl table.EncodingTable, data string) []byte{
	encodedTable := encodeTable(tbl)

	var buf bytes.Buffer


	buf.Write(encodeInt(len(encodedTable)))
	buf.Write(encodeInt(len(data)))
	buf.Write(encodedTable)
	buf.Write(splitByChunks(data, chunkSize).Bytes())
	

	return buf.Bytes()

}


func (ed EncoderDecoder) Decode(encData []byte) string{
	//
	table, data  := parseFile(encData)

	return table.Decode(data)
	
}



func parseFile(data []byte) (table table.EncodingTable, date string){
	const (
		tableSizeBytesCount = 4
		dataSizeBytesCount = 4
	)

	tableSizeBinary, data := data[:tableSizeBytesCount], data[tableSizeBytesCount:]
	dataSizeBinary, data := data[:dataSizeBytesCount], data[dataSizeBytesCount:]

	
	tableSize := binary.BigEndian.Uint32(tableSizeBinary)
	dataSize := binary.BigEndian.Uint32(dataSizeBinary)

	tblBinary, data := data[:tableSize], data[tableSize:]
	
	
	return decodeTable(tblBinary), NewBinChunks(data).ToString()[:dataSize]
} 


func encodeInt(num int) []byte{
	
	res := make([]byte, 4)

	binary.BigEndian.PutUint32(res, uint32(num))

	return res
}


func encodeTable(tbl table.EncodingTable) []byte{
	var tableBuf bytes.Buffer


	if err := gob.NewEncoder(&tableBuf).Encode(tbl); err != nil{
		log.Fatal("can't serialize table...", err)
	}


	return tableBuf.Bytes()
}


func decodeTable(tblBinary []byte) (table.EncodingTable) {
	var tbl table.EncodingTable


	r := bytes.NewReader(tblBinary)
	if err := gob.NewDecoder(r).Decode(&tbl); err != nil{
		log.Fatal("can't serialize table...", err)
	}


	return tbl
}



//encodeBin encodes string into binary codes string withou spaces
func encodeBin(str string, table table.EncodingTable) string{
	var buf strings.Builder

	
	for _, ch := range str{
		buf.WriteString(bin(ch, table))
	}

	
	return buf.String()
}	

// bin uses character as a key for encodeTable and returns its binary code
func bin(ch rune, table table.EncodingTable) string{

	res, ok := table[ch]
	if !ok{
		panic("unknown character: " + "\\" + string(ch) + "\\")
	}

	return res
}




