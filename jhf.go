package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jackrendor/jhf/resources"
)

var bannerString string = `
     ██ ██   ██ ███████ 
     ██ ██   ██ ██      
     ██ ███████ █████ 
██   ██ ██   ██ ██    
 █████  ██   ██ ██ 
                   
            	Jack Hash Finder by ​@jackrendor
`

func main() {
	defer color.Unset()
	normalColor := color.New(color.FgCyan)
	greenColor := color.New(color.FgHiGreen, color.Bold)
	redColor := color.New(color.FgHiRed, color.Bold)
	onlyfound := flag.Bool("onlyfound", false, "Print only found hashes")
	nobanner := flag.Bool("nobanner", false, "Do not print banner")
	fileInput := flag.String("file", "", "Read hashes from file, line by line")
	norun := flag.Bool("norun", false, "Do not send any requests")
	nocolor := flag.Bool("nocolor", false, "Disable color output")
	flag.Parse()

	if *nocolor {
		color.NoColor = true
	}

	var data []string
	if len(*fileInput) > 0 {
		fileD, openErr := os.Open(*fileInput)
		if openErr != nil {
			log.Println("[main] [os.Open]:", openErr.Error())
			return
		}
		// Read from file
		scanner := bufio.NewScanner(fileD)
		scanner.Split(bufio.ScanLines)

		// Read line by line
		for scanner.Scan() {
			value := strings.TrimSpace(scanner.Text())
			data = append(data, value)
		}
		fileD.Close()
	} else {
		data = flag.Args()
	}

	if !*nobanner {
		normalColor.Println(bannerString)
	}
	if *norun {
		return
	}
	for _, element := range resources.Crack(data) {
		if element.Solved {
			greenColor.Printf("[Found]\t%s:%s\n", element.Hash, element.Value)
		} else if !*onlyfound {
			redColor.Printf("[N/A]\t%s\n", element.Hash)
		}
	}
}
