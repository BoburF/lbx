package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/BoburF/lbx/common"
	"github.com/BoburF/lbx/config"
	"github.com/BoburF/lbx/lb"
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

	var cfg config.LbxConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Println("Failed to parse config:", err)
		return
	}

	lbx, err := lb.NewLoadBalancer(cfg.Servers, cfg.RetryTimeInMinutes)
	if err != nil {
		fmt.Println("Failed make NewLoadBalancer:", err)
		return
	}

	lbx.Server("localhost", cfg.Port)
}
