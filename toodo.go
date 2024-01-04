package toodo

import (
	"fmt"
	"log/slog"

	"github.com/jessevdk/go-flags"
	mdf "github.com/taylormonacelli/toodo/mdfind"
)

var opts struct {
	Args struct {
		Period string `description:"Time period parameter in the format 1y, 1M10w1m, 10M, 10m, 200s, 34d, 1y23d, 2d20s, etc."`
	} `positional-args:"yes" required:"yes"`
	LogFormat      string `long:"log-format" choice:"text" choice:"json" default:"text" description:"Log format"`
	Verbose        []bool `short:"v" long:"verbose" description:"Show verbose debug information, each -v bumps log level"`
	logLevel       slog.Level
	FileExtensions []string `short:"e" long:"file-extension" description:"File extension to search for" required:"false"`
	Root           string   `short:"r" long:"root" description:"Root directory for search" required:"false"`
	DryRun         bool     `short:"d" long:"dry-run" description:"Don't run mdfind, useful with -vv to see what mdfind command was generated" required:"false"`
}

func Execute() int {
	if err := parseFlags(); err != nil {
		return 1
	}

	if err := setLogLevel(); err != nil {
		return 1
	}

	if err := setupLogger(); err != nil {
		return 1
	}

	if err := run(); err != nil {
		slog.Error("run failed", "error", err)
		return 1
	}

	return 0
}

func parseFlags() error {
	_, err := flags.Parse(&opts)
	return err
}

func run() error {
	mdfind := mdf.NewMDFind(opts.Root, opts.Args.Period, opts.FileExtensions)
	cmdString := mdfind.String()
	fmt.Println(cmdString)

	if !opts.DryRun {
		_, err := mdfind.Run()
		if err != nil {
			fmt.Printf("Command failed: %s\n", cmdString)
			return err
		}
	}

	return nil
}
