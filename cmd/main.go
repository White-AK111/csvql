package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"csvql/config"
	"csvql/pkg/parser"
	"csvql/pkg/scanner"

	"github.com/go-git/go-git/v5"
	"go.uber.org/zap"
)

func main() {
	log.Printf("Start load configuration.\n")
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Eror on load config %v\n", err)
	}

	// flushes buffer, if any
	defer cfg.App.Logger.Sync()

	cfg.App.Logger.Info("Configuration file successfully load.",
		zap.String("configuration file", cfg.App.ConfigPath),
		zap.String("source CSV file", cfg.App.FilePath),
		zap.String("delimiter of CSV file", cfg.App.Delimiter),
	)

	fmt.Printf("\nSource CSV file: %s\n", cfg.App.FilePath)

	// Getting the latest commit on the current branch
	r, err := git.PlainOpen("../")
	if err != nil {
		cfg.App.Logger.Error("Error on open git repository.",
			zap.String("repository", "../"),
			zap.Error(err),
		)
	}
	ref, err := r.Head()
	if err != nil {
		cfg.App.Logger.Error("Error on get HEAD in git.",
			zap.String("repository", "../"),
			zap.Error(err),
		)
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		cfg.App.Logger.Error("Error on get last commit in git.",
			zap.String("repository", "../"),
			zap.Error(err),
		)
	}
	fmt.Printf("Last commit hash:%s\n", commit.Hash)

	ctxCurr, cancelCurr := context.WithCancel(context.Background())

	scannerFile, err := scanner.NewScanner(cfg.App.FilePath)
	if err != nil {
		cfg.App.Logger.Error("Error on init file scanner.",
			zap.Error(err),
		)
		os.Exit(1)
	}
	defer scannerFile.File.Close()

	err = scannerFile.GetHeaders(cfg.App.Delimiter, cfg.App.Comment)
	if err != nil {
		cfg.App.Logger.Error("Error on get headers from source CSV file.",
			zap.Error(err),
		)
		os.Exit(1)
	}
	fmt.Printf("\nHeaders in CSV file: %v\n", scannerFile.Headers)

	fmt.Printf("Input your query %s:", "(example: 'age > 20 AND status = \"sick\"')")
	scannerStd := bufio.NewScanner(os.Stdin)
	scannerStd.Scan() //
	query := scannerStd.Text()
	fmt.Printf("Your query is: %s\n", query)

	pser := parser.Parser{}
	err = pser.GetConditions(query)
	if err != nil {
		cfg.App.Logger.Error("Error on parse input query.",
			zap.String("query", query),
			zap.Error(err),
		)
	}
	fmt.Println(pser)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	for {
		select {
		case <-ctxCurr.Done():
			return
		case sig := <-sigCh:
			switch sig {
			case syscall.SIGINT:
				// Close while catch SIGINT
				log.Println("Catch SIGINT")
				cancelCurr()
				return
			default:
				log.Println("Catch other signal")
			}
		default:
			time.Sleep(60 * time.Second)
			log.Println("Next query ...")
		}
	}
}
