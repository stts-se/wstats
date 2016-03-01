package main

import (
	"bufio"
	"compress/bzip2"
	"encoding/xml"
	"fmt"
	"github.com/stts-se/wstats/util"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// start: xml parsing
// http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/
type Redirect struct {
	Title string `xml:"title,attr"`
}
type Page struct {
	Title string   `xml:"title"`
	Redir Redirect `xml:"redirect"`
	Text  string   `xml:"revision>text"`
}

// end: xml parsing

func convert(l util.XString) util.XString {
	result := l.Value
	for _, repl := range tokenReplacements {
		result = repl.From.ReplaceAllString(result, repl.To)
	}
	return util.XString{result}.Trim().ToLower()
}

// start: pre-compiled regexps
type Replacement struct {
	From *regexp.Regexp
	To   string
}

var tokenReplacements = []Replacement{
	Replacement{regexp.MustCompile("http://[^\\s]+"), ""},
	Replacement{regexp.MustCompile("&lt;!--"), "<!--"},
	Replacement{regexp.MustCompile("--&gt;"), "-->"},
	Replacement{regexp.MustCompile("<!--[^>]+-->"), ""},
	Replacement{regexp.MustCompile("(&lt;|<)/?ref( |(&gt;|>)).*$"), ""},
	Replacement{regexp.MustCompile("&quot;"), "\""},
	Replacement{regexp.MustCompile("&amp;"), "&"},
	Replacement{regexp.MustCompile("^ *\\* *"), ""},
	Replacement{regexp.MustCompile("&[a-z]+;"), ""},
	Replacement{regexp.MustCompile("<[^>]+>"), ""},
	Replacement{regexp.MustCompile("\\{\\{[^}]+(\\}\\}|$)"), ""},
	Replacement{regexp.MustCompile("[{}]"), ""},
	Replacement{regexp.MustCompile("\\[\\[Kategori:"), "[["},
	Replacement{regexp.MustCompile("\\[\\[[A-Za-z]+:([^|\\]]+\\|)+"), "[["},
	Replacement{regexp.MustCompile("\\[\\[([^|\\]]+)\\|?\\]\\]"), "$1"},
	Replacement{regexp.MustCompile("\\[\\[(?:[^|\\]]+)\\|([^|\\]]+)\\]\\]"), "$1"},
	Replacement{regexp.MustCompile("\\[\\[(?:[^|\\]]+)(?:\\|(?:[^|\\]]+))*\\|([^|\\]]+)\\]\\]"), "$1"},
	Replacement{regexp.MustCompile("[\\[\\]]+"), ""},
	Replacement{regexp.MustCompile("==+"), ""},
	Replacement{regexp.MustCompile(" ' "), " "},
	Replacement{regexp.MustCompile("(: | :)"), " "},
	Replacement{regexp.MustCompile("[\\]\\[!\"”#$%&()*+,./;<=>?@\\^_`{|}~\\s\u00a0–]+"), " "},
	Replacement{regexp.MustCompile("(( |^)'+|'+( |$))"), " "},
	Replacement{regexp.MustCompile("( *- | - *)"), " "},
}
var lineReplacements = []Replacement{
	Replacement{regexp.MustCompile("&lt;"), "<"},
	Replacement{regexp.MustCompile("&gt;"), ">"},
	Replacement{regexp.MustCompile("&quot;"), "\""},
	Replacement{regexp.MustCompile("&amp;"), "&"},
	Replacement{regexp.MustCompile("^ *<text[^>]*>"), ""},
	Replacement{regexp.MustCompile("#REDIRECT "), ""},
	Replacement{regexp.MustCompile("^ *:;?"), ""},
}
var skipRe = regexp.MustCompile("^ *(!|\\||<|\\{\\||&|<redirect[^>]+>).*")

// end: pre-compiled regexps

func tokenizeLine(l0 util.XString) util.XArray {
	splittable := convert(l0)
	return splittable.SplitWhiteSpace()
}

func preFilterLine(l util.XString) util.XString {
	result := l.Value
	for _, repl := range lineReplacements {
		result = repl.From.ReplaceAllString(result, repl.To)
	}
	return util.XString{result}
}

func skip(l util.XString) bool {
	return (!strings.HasPrefix(l.Trim().Value, "<page") && !strings.HasPrefix(l.Trim().Value, "<text") && (l.MatchesRe(skipRe) || l.Contains("[[Användar") || l.Contains("<comment>")))
}

func lIntRoundToString(i int) string {
	if i > 100000 {
		return fmt.Sprintf("%7.2fM", (float64(i) / float64(1000000)))
	} else if i > 1000 {
		return fmt.Sprintf("%7.2fK", (float64(i) / float64(1000)))
	} else {
		return fmt.Sprintf("%7d.00", i)
	}
}

func lIntPrettyPrint(i int) string {
	return fmt.Sprintf("%12s", util.XString{fmt.Sprintf("%d", i)}.
		ReplaceAll("([0-9])([0-9]{3})([0-9]{3})([0-9]{3})$", "$1,$2,$3,$4").
		ReplaceAll("([0-9])([0-9]{3})([0-9]{3})$", "$1,$2,$3").
		ReplaceAll("([0-9])([0-9]{3})$", "$1,$2").Value)
}

func printProgress(nPages int, nLines int, nWords int) {
	pp := lIntRoundToString(nPages)
	ww := lIntRoundToString(nWords)
	ll := lIntRoundToString(nLines)
	numbers := fmt.Sprintf("%s pgs, %s lns, %s wds", pp, ll, ww)
	withPadding := fmt.Sprintf("\r%-52s\r", numbers)
	fmt.Fprint(os.Stderr, withPadding)
}

func clearProgress() {
	withPadding := fmt.Sprintf("\r%-45s\r", " ")
	fmt.Fprint(os.Stderr, withPadding)
}

func tokenizeText(text string) (int, int, map[util.XString]int) {
	var nLines = 0
	var nLinesSkipped = 0
	wordFreqs := make(map[util.XString]int)
	for _, l0 := range strings.Split(text, "\n") {
		nLines++
		line := preFilterLine(util.XString{l0})
		if skip(line) {
			nLinesSkipped++
		} else {
			words := tokenizeLine(line)
			if len(words.Value) > 0 {
				for _, word := range words.Value {
					wordFreqs[word]++
				}
			}
		}
	}
	return nLines, nLinesSkipped, wordFreqs
}

func loadXml(output io.Writer, path string, pageLimit int, logAt int) (int, int, int, int, int, map[util.XString]int) {
	var decoder *xml.Decoder
	if strings.HasPrefix(path, "http") {
		response, err := http.Get(path)
		if err != nil {
			log.Fatal(err)
		}
		if response.StatusCode != 200 {
			log.Fatal(response.Status + " " + path)
		}
		defer response.Body.Close()
		if strings.HasSuffix(path, "bz2") {
			bz := bzip2.NewReader(response.Body)
			decoder = xml.NewDecoder(bz)
		} else {
			decoder = xml.NewDecoder(response.Body)
		}
	} else {
		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		if strings.HasSuffix(path, "bz2") {
			bz := bzip2.NewReader(file)
			decoder = xml.NewDecoder(bz)
		} else {
			decoder = xml.NewDecoder(file)
		}
	}

	nLines := 0
	nLinesSkipped := 0
	nPages := 0
	nRedirects := 0
	nWords := 0
	wordFreqs := make(map[util.XString]int)

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		if pageLimit > 0 && nPages >= pageLimit {
			clearProgress()
			log.Println(fmt.Sprintf("Break called at %d pages (for debugging)", nPages))
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "page" {
				var p Page
				nPages++
				decoder.DecodeElement(&p, &se)
				var text = p.Text
				var title = p.Title
				if len(title) > 0 {
					text = title + "\n" + text
				}
				var redirect = p.Redir.Title
				if len(redirect) > 0 {
					nRedirects++
				} else {
					nL, nLS, wFs := tokenizeText(text)
					nLines += nL
					nLinesSkipped += nLS
					for w, f := range wFs {
						nWords += f
						wordFreqs[w] += f
					}
				}
				if nPages%logAt == 0 {
					printProgress(nPages, nLines, nWords)
				}
			}
		}
	}
	return nPages, nRedirects, nLines, nLinesSkipped, nWords, wordFreqs
}

func main() {

	// Download data here: https://dumps.wikimedia.org
	// Valid input:
	//   xml file : XXwiki-YYYYMMDD-pages-articles-multistream.xml
	//   bz2 file : XXwiki-YYYYMMDD-pages-articles-multistream.xml.bz2
	//   xml url  : implemented by not likely to be used...
	//   bz2 url  : https://dumps.wikimedia.org/svwiki/latest/svwiki-latest-pages-articles-multistream.xml.bz2

	if len(os.Args) != 2 && len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "USAGE\tgo run wikistats.go <path> <limit>*")
		fmt.Fprintln(os.Stderr, "   \t<path> wikimedia dump (file or url, xml or xml.bz)")
		fmt.Fprintln(os.Stderr, "   \t<limit> limit number of pages (optional)")
		fmt.Fprintln(os.Stderr, "EXAMPLE\tgo run wikistats.go https://dumps.wikimedia.org/svwiki/latest/svwiki-latest-pages-articles-multistream.xml.bz2 10000")
		os.Exit(1)
	}

	output := bufio.NewWriter(os.Stdout)

	log.Print("*** RUNNING wikistats.main() ***")

	start := time.Now()
	defer output.Flush()

	path := os.Args[1]
	pageLimit := -1
	if len(os.Args) == 3 {
		p, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		pageLimit = p
	}

	log.Print("Path : ", path)

	logAt := 100
	nPages, nRedirects, nLines, nLinesSkipped, nWords, wordFreqs := loadXml(output, path, pageLimit, logAt)

	loaded := time.Now()

	for _, pair := range util.SortByWordCount(wordFreqs) {
		if pair.Value > 1 {
			fmt.Fprintf(output, "%d\t%s\n", pair.Value, pair.Key)
		}
	}

	output.Flush()
	end := time.Now()

	clearProgress()

	loadDur := loaded.Sub(start) - loaded.Sub(start)%time.Millisecond
	printDur := end.Sub(loaded) - end.Sub(loaded)%time.Millisecond
	totalDur := end.Sub(start) - end.Sub(start)%time.Millisecond

	log.Print("Inläsning tog   : ", fmt.Sprintf("%12v\n", loadDur))
	log.Print("Utskrift tog    : ", fmt.Sprintf("%12v\n", printDur))
	log.Print("Total tid       : ", fmt.Sprintf("%12v\n", totalDur))

	log.Print("Antal sidor     : ", lIntPrettyPrint(nPages))
	log.Print("Varav redirects : ", lIntPrettyPrint(nRedirects))
	log.Print("Antal rader     : ", lIntPrettyPrint(nLines))
	log.Print("Skippade rader  : ", lIntPrettyPrint(nLinesSkipped))
	log.Print("Antal löpord    : ", lIntPrettyPrint(nWords))
	log.Print("Antal unika ord : ", lIntPrettyPrint(len(wordFreqs)))

}
