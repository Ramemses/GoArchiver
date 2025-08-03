package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"strings"
	"io"
	"path/filepath"
	"archiver/lib/compression/vlc"
	"archiver/lib/compression"
	"archiver/lib/compression/vlc/table/shanon_fano"
	"archiver/lib/compression/vlc/table/haffman"
)


var unpackCmd = &cobra.Command{
	Use: "unpack",
	Short: "Unpack file",
	Run: unpack,
}


//var vlcUnpackCmd = &cobra.Command{
//	Use: "vlc",
//	Short: "Unpack file using variable-length code",
//	Run: unpack,
	
//}

// TODO: Take extension from file
const unpackedExtension = "txt"
//var ErrEmptyPath = errors.New("path to file is not specified")

func unpack(cmd *cobra.Command, args []string){
	var decoder compression.Decoder

	if len(args) == 0 || args[0] == ""{
		handleError(ErrEmptyPath)
	}

	method := cmd.Flag("method").Value.String()

	switch method{
		case "shanon_fano":
			decoder = vlc.New(shanon_fano.NewGenerator())
		case "haffman":
			decoder = vlc.New(haffman.NewGenerator())
	} 


	filePath := args[0]
	r, err:= os.Open(filePath)
	if err != nil{
		handleError(err)
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil{
		handleError(err)
	}

	unpacked := decoder.Decode(data)

	err = os.WriteFile(unpackedFileName(filePath), []byte(unpacked), 0644)
	if err != nil{
		handleError(err)
	}
}


// TODO: refactoring
func unpackedFileName(path string) string{
	// /path/to/file/myFile.vlc -> myFile.txt
	fileName := filepath.Base(path) // myFile.vlc
	ext := filepath.Ext(fileName) // .vlc
	baseName := strings.TrimSuffix(fileName, ext) // myFile.vlc -> myFile


	return baseName + "." + unpackedExtension // return myFile + . + txt
}




func init(){
	rootCmd.AddCommand(unpackCmd)

	unpackCmd.Flags().StringP("method", "m", "", "decompression method: vlc")

	if err := unpackCmd.MarkFlagRequired("method"); err != nil{
		handleError(err)
	}

}

