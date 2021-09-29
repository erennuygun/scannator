package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	scan "github.com/erennuygun/scannator/models"
	"github.com/fatih/color"
)

var wg sync.WaitGroup

const (
	BANNER_HEADER = `
 ███████╗ ██████╗ █████╗ ███╗   ██╗███╗   ██╗ █████╗ ████████╗ ██████╗ ██████╗ 
 ██╔════╝██╔════╝██╔══██╗████╗  ██║████╗  ██║██╔══██╗╚══██╔══╝██╔═══██╗██╔══██╗
 ███████╗██║     ███████║██╔██╗ ██║██╔██╗ ██║███████║   ██║   ██║   ██║██████╔╝
 ╚════██║██║     ██╔══██║██║╚██╗██║██║╚██╗██║██╔══██║   ██║   ██║   ██║██╔══██╗
 ███████║╚██████╗██║  ██║██║ ╚████║██║ ╚████║██║  ██║   ██║   ╚██████╔╝██║  ██║
 ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═══╝╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝
																				
`
	BANNER_AUTHOR = "▪ github.com/erennuygun\n"
	BANNER_SEP    = "________________________________________________________________________________\n"
)

func main() {
	startTime := time.Now()
	fmt.Println(BANNER_HEADER, BANNER_AUTHOR, BANNER_SEP)

	// Declarate Directory Scanner Flags
	dirScan := flag.NewFlagSet("dir", flag.ExitOnError)
	url := dirScan.String("u", "", "Add Target Domain")
	wordlist := dirScan.String("w", "dir-list", "Add Wordlist")
	excludes := dirScan.String("x", "", "Exclude Status Codes")
	threads := dirScan.Int("t", 1, "Add Threads ")

	// Declarate Subdomain Scanner Flags
	subScan := flag.NewFlagSet("dir", flag.ExitOnError)
	urlSub := subScan.String("u", "", "Add Target Domain")
	wlSub := subScan.String("w", "sub-list", "Add Wordlist")
	excludeSub := subScan.String("x", "", "Exclude Status Codes")
	threadSub := subScan.Int("t", 1, "Add Threads ")

	// CLI Args
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "dir":
			dirScan.Parse(os.Args[2:])
			scan.DirScan(*url, *wordlist, *excludes, *threads)
		case "sub":
			subScan.Parse(os.Args[2:])
			scan.SubScan(*urlSub, *wlSub, *excludeSub, *threadSub)
		case "-h", "--help":
			color.Yellow(`Usage of ./scannator:`)
			fmt.Println(`
	dir
		Directory Scan 
	sub
		Subdomain Scan`)
			color.Cyan(`
	For Example:
		.\scannator dir -u https://victim.co
		.\scannator sub -u victim.co`)
		default:
			fmt.Println("Please Read Help !")
		}
	} else {
		color.Yellow("Please read help with -h or --help ")
	}

	elapsedTime := time.Now().Sub(startTime)
	color.Magenta("\n [*] - Scan terminated in %v\n", elapsedTime)
}