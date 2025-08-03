package cmd

import (
	"github.com/spf13/cobra"
	"errors"
	"os"
	"strings"
	"io"
	"path/filepath"
	"archiver/lib/compression/vlc"
	"archiver/lib/compression"
	"archiver/lib/compression/vlc/table/shanon_fano"
	"archiver/lib/compression/vlc/table/haffman"
)

var packCmd = &cobra.Command{
	Use: "pack",
	Short: "Pack file",
	Run: pack,
}



//var vlcPackCmd = &cobra.Command{
//	Use: "vlc",
//	Short: "Pack file using variable-length code",
//	Run: pack,
	
//}

const packedExtension = "vlc"

var ErrEmptyPath = errors.New("path to file is not specified")
func pack(cmd *cobra.Command, args []string){

	var encoder compression.Encoder

	if len(args) == 0 || args[0] == ""{
		handleError(ErrEmptyPath)
	}

	method := cmd.Flag("method").Value.String()
	
	switch method{
		case "shanon_fano":
			encoder = vlc.New(shanon_fano.NewGenerator())
		case "haffman":
			encoder = vlc.New(haffman.NewGenerator())
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

	packed := encoder.Encode(string(data))

	err = os.WriteFile(packedFileName(filePath), packed, 0644)
	if err != nil{
		handleError(err)
	}
}

func packedFileName(path string) string{
	// /path/to/file/myFile.txt -> myFile.vlc
	fileName := filepath.Base(path) // myFile.txt
	ext := filepath.Ext(fileName) // .txt
	baseName := strings.TrimSuffix(fileName, ext) // myFile.txt -> myFile


	return baseName + "." + packedExtension // return myFile + . + vlc
}



func init(){
	rootCmd.AddCommand(packCmd)


	packCmd.Flags().StringP("method", "m", "", "copmression method: vlc")
	if err := packCmd.MarkFlagRequired("method"); err != nil{
		handleError(err)
	}

}
	
