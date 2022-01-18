package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/alexflint/go-arg"
	scriptish "github.com/ganbarodigital/go_scriptish"
)

var args struct {
	File     string `arg:"-f,--file" help:"file containing URLs to scan"`
	WordList string `arg:"-w,--wordlist,required" help:"file path to wordlist"`
	Ports    string `arg:"-p,--ports" help:"ports to pass to httprobe"`
	Threads  string `default:"15" arg:"-t,--threads" help:"threads to use with gobuster"`
	BadCodes string `default:"404" arg:"-b,--blacklist" help:"blacklist HTTP status codes for gobuster"`
}

func main() {
	arg.MustParse(&args)

	now := time.Now()
	urls, err := scriptish.ExecPipeline(
		scriptish.CatFile(args.File),
		scriptish.Exec("httprobe", "-c", "30"),
	).Strings()

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	for _, url := range urls {
		fmt.Printf("===> %s\n", url)
		cmd := exec.Command(
			"gobuster", "dir", "--quiet", "--useragent", "noleak", "--expanded", "--follow-redirect",
			"--url", url, "--wordlist", args.WordList, "--threads", args.Threads,
			"--status-codes-blacklist", args.BadCodes,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

	end := time.Since(now)
	fmt.Printf("\nTime elapsed: %s\n", end)
}
