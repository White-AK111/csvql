package main

import (
	"csvql/config"

	"fmt"
	"log"

	"github.com/go-git/go-git/v5"
	"go.uber.org/zap"
)

func main() {
	log.Printf("Start load configuration.\n")
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Eror on load config %v\n", err)
	}

	cfg.App.Logger.Info("Configuration file successfully load.",
		zap.String("configuration file", cfg.App.ConfigPath),
		zap.String("source CSV file", cfg.App.FilePath),
		zap.String("delimiter of CSV file", cfg.App.Delimiter),
	)

	fmt.Printf("Source CSV file: %s\n", cfg.App.FilePath)

	// Getting the latest commit on the current branch
	r, err := git.PlainOpen("../")
	ref, err := r.Head()
	commit, err := r.CommitObject(ref.Hash())
	fmt.Printf("\nInformation about last commit\n%v\n:", commit)

	// flushes buffer, if any
	defer cfg.App.Logger.Sync()
}
