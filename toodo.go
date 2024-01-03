package toodo

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/taylormonacelli/navyloss"
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
	periodSeconds, err := navyloss.PeriodToSeconds(opts.Args.Period)
	if err != nil {
		return err
	}

	mdfind := NewMDFind(opts.Root, time.Duration(periodSeconds)*time.Second, opts.FileExtensions)
	mdfind.ExpandHomeDir()

	cmd := mdfind.BuildCommand()

	// Conditionally add -onlyin switch based on --root flag
	if opts.Root != "" {
		cmd.AddArgument("-onlyin", opts.Root)
	}

	slog.Debug("debug command", "command", cmd.String())

	output, err := cmd.Run()
	if err != nil {
		log.Fatalf("Error running mdfind command: %v", err)

		if stderr := cmd.GetStderr(); stderr != "" {
			fmt.Printf("mdfind command stderr: %s\n", stderr)
		}
		return err
	}

	cmd.Stdout = cmd.GetStdout()
	cmd.Stderr = cmd.GetStderr()

	fmt.Println(output)

	return nil
}
