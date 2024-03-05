/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"lambdacrate-cli/cmd"
	"lambdacrate-cli/lib"
	"log"
)

func main() {
	cmd.Execute()

}

func init() {
	log.Println("generating config file")
	err := lib.CreateConfigFile()
	if err != nil {
		log.Fatal("Failed to generate config file for application", err)
	}

}
