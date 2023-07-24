package main

import (
	flag "github.com/spf13/pflag"
	ft "github.com/steffenfritz/FileTrove"
	"time"
)

// TSStartedFormated is the formated timestamp when FileTrove was started
var TSStartedFormated string

func init() {
	TSStarted := time.Now()
	TSStartedFormated = TSStarted.Format("2006-01-02_15:04:05")
}

func main() {
	inDir := flag.StringP("indir", "i", "", "Input directory to work on.")
	getSiegfried := flag.BoolP("download-siegfried", "s", false, "Download siegfried's pronom database.")
	getNSRL := flag.BoolP("download-nsrl", "n", false, "Download NSRL database.")
	outputFile := flag.StringP("output", "o", *inDir+"_"+TSStartedFormated, "Name of the file where the results are written to.")

	flag.Parse()

	if *getSiegfried {
		ft.GetSiegfriedDB()
	}
	if *getNSRL {
		ft.GetSiegfriedDB()
	}


}
