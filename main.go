package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BoburF/lbx/common"
)

func main() {
	pathToConfig := flag.String("file", "config.json", "config file for lbx")
	flag.Parse()

	filePath, err := common.ParseFilePath(*pathToConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	data, err := os.ReadFile(filePath.AbsolutePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(data)
}
