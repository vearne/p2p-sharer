package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/vearne/p2p-sharer/consts"
	"github.com/vearne/p2p-sharer/models"
	"github.com/vearne/p2p-sharer/utils"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate BT seed files",
	Long:  "",
	Run:   genSeedFile,
}

var trackerAddr string
var filePath string
var seedDirPath string

const (
	PieceSize             int = 1 * 1024 * 1024
	DefaultReaderBuffSize     = 1024 * 4
	BuffSize                  = 1024
)

func init() {
	genCmd.Flags().StringVar(&trackerAddr, "tracker", "localhost:3456", "address of tracker")

	genCmd.Flags().StringVar(&filePath, "filePath", "/tmp/Motrix-1.4.1.dmg",
		"file path of file which will be shared")

	genCmd.Flags().StringVar(&seedDirPath, "seedPath", "./",
		"seed file dir path")
	rootCmd.AddCommand(genCmd)

}

func genSeedFile(cmd *cobra.Command, args []string) {
	log.Println("[start]gen Seed file")
	log.Println("trackerAddr:", trackerAddr)
	log.Println("filePath:", filePath)
	var err error
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	buffReader := bufio.NewReaderSize(f, DefaultReaderBuffSize)
	var info models.SeedInfo
	_, info.FileName = path.Split(filePath)
	info.TrackerAddr = trackerAddr
	info.CreatedAt = time.Now()
	info.Pieces = make([]*models.PieceInfo, 0)

	totalLength := 0
	buff := make([]byte, BuffSize)
	index := 0

	memBuffer := bytes.NewBuffer(make([]byte, 0))
	for {
		n, err := buffReader.Read(buff)
		if err == io.EOF {
			if memBuffer.Len() != 0 {
				info.Pieces = append(info.Pieces, createPiece(memBuffer, index))
				index++
				memBuffer.Reset()
			}
			break
		} else {
			totalLength += n
			memBuffer.Write(buff[0:n])
			if memBuffer.Len() == PieceSize {
				info.Pieces = append(info.Pieces, createPiece(memBuffer, index))
				index++
				memBuffer.Reset()

			}
		}

	}
	info.Length = totalLength
	log.Println("size of info.Pieces", len(info.Pieces))
	log.Println("totalLength", totalLength)
	// write to seed file
	seedFilePath := filepath.Join(seedDirPath, info.FileName+consts.SeedPostfix)
	sf, err := os.Create(seedFilePath)
	defer sf.Close()
	if err != nil {
		log.Fatal(err)
	}
	bt, _ := json.MarshalIndent(info, "", "\t")
	sf.Write(bt)
	log.Println(string(bt))
	log.Println("[end]gen Seed file")
}

func createPiece(memBuffer *bytes.Buffer, index int) *models.PieceInfo {
	hashValue := utils.BlockHash(memBuffer.Bytes())
	pieceInfo := models.PieceInfo{index, memBuffer.Len(), hashValue}
	return &pieceInfo
}
