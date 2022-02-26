package main

import (
	"fmt"
	//"github.com/spf13/cobra"
	//"github.com/spf13/viper"
	"os"
	"webcrawler/rest"
)

// In case I have time to stored server parameters in a config file
// and rely on cobra/viper to handle it (and command line)
var cfgFile = ""

func initConfig() {
	if cfgFile != "" {
		// enable ability to specify config file via flag
	}
}

func init() {
	//cobra.OnInitialize(initConfig)
}

func main() {
	err := rest.RunServer()
	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(1)
	}
	fmt.Printf("Exiting main...\n")
}
