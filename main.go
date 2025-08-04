package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/alecthomas/kong"
)

var CLI struct {
	In           string `arg:"" optional:"" help:"input file. must be in slides format"`
	Out          string `short:"o" help:"output file"`
	Reference    string `short:"R" help:"show reference information for a specific directive"`
	ReferenceAll bool   `short:"r" help:"show reference information for all directives"`
}

var help = map[string]string{
	"!global": "(at beginning of a block) establish a block for global options",
	"!pfoot":  "<string> (global) set a string to be the primary footer value. NOTE: for the footer to show, both !pfoot and !sfoot must be present",
	"!sfoot":  "<string> (global) text to be set for the secondary footer",
	"!exec":   "<file> (global) JS file to be executed every time the loop finishes. It's `main` method will be called.",
	"!styles": "<file> (global) CSS file to be included in the slideshow",
	"!title":  "(slide) sets the current slide blocks type to title-slide",
	"!time":   "<float> (global OR slide) sets the duration for the entire slideshow or for the current slide, depending on the scope in which it is invoked",
	"!image":  "<file> (slide) sets the image to show on the slide",
	"!id":     "<string> (slide) sets the current id for the slide. By default each slide has no id. Useful for scripting.",
	"!noautoplay": "(global) disable auto playing. each slide can be toggled by enabling the `active` attribute. only for debugging",
}

func rightPad(s string, leng int, ch string) string {
	if len(s) >= leng {
		return s
	}
	return s + strings.Repeat(ch, leng-len(s))
}

func main() {
	kong.Parse(&CLI)

	if CLI.Reference != "" {
		if val, ok := help[CLI.Reference]; ok {
			fmt.Printf("%s: %s\n", CLI.Reference, val)
		}
		os.Exit(0)
	}

	if CLI.ReferenceAll {
		var list sort.StringSlice
		for k, v := range help {
			list = append(list, fmt.Sprintf("%s%s\n", rightPad(k+":", 13, " "), v))
		}
		list.Sort()
		for _, v := range list {
			fmt.Print(v)
		}
		os.Exit(0)
	}

	b, err := os.ReadFile(CLI.In)

	log.SetFlags(log.LstdFlags | log.Lmsgprefix)
	log.SetPrefix("infodisplay: ")
	if err != nil {
		log.Fatalf("couldn't read slides file %q: %s", CLI.In, err)
	}

	ss, err := slidesFromString(string(b))

	if err != nil {
		log.Fatalf("parse: %s", err)
	}

	html, err := ss.toHTML()

	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(func() string {
		if CLI.Out != "" {
			return CLI.Out
		} else {
			return "index.html"
		}
	}(), html, 0777)
}
