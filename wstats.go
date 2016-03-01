package main

import (
	"bufio"
	"compress/bzip2"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// start: util
func SplitWhiteSpace(s string) []string {
	splitted := strings.Split(s, " ")
	result := make([]string, 0)
	for _, v0 := range splitted {
		v := strings.TrimSpace(v0)
		if len(v) > 0 {
			result = append(result, v)
		}
	}
	return result
}

func SortByWordCount(wordFrequencies map[string]int) FreqList {
	pl := make(FreqList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Freq{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Freq struct {
	Key   string
	Value int
}
type FreqList []Freq

func (p FreqList) Len() int           { return len(p) }
func (p FreqList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p FreqList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// end: util

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

func convert(s string) string {
	result := s
	for _, repl := range tokenReplacements {
		result = repl.From.ReplaceAllString(result, repl.To)
	}
	return strings.ToLower(strings.TrimSpace(result))
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

func tokenizeLine(l string) []string {
	l = convert(l)
	return SplitWhiteSpace(l)
}

func preFilterLine(l string) string {
	result := l
	for _, repl := range lineReplacements {
		result = repl.From.ReplaceAllString(result, repl.To)
	}
	return result
}

func skip(l string) bool {
	l = strings.TrimSpace(l)
	return (!strings.HasPrefix(l, "<page") && !strings.HasPrefix(l, "<text") && (skipRe.MatchString(l) || strings.Contains(l, "[[Användar") || strings.Contains(l, "<comment>")))
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
	result := fmt.Sprintf("%d", i)
	replacements := []Replacement{
		Replacement{regexp.MustCompile("([0-9])([0-9]{3})([0-9]{3})([0-9]{3})$"), "$1,$2,$3,$4"},
		Replacement{regexp.MustCompile("([0-9])([0-9]{3})([0-9]{3})$"), "$1,$2,$3"},
		Replacement{regexp.MustCompile("([0-9])([0-9]{3})$"), "$1,$2"},
	}
	for _, repl := range replacements {
		result = repl.From.ReplaceAllString(result, repl.To)
	}
	return result
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

func tokenizeText(text string) (int, int, map[string]int) {
	var nLines = 0
	var nLinesSkipped = 0
	wordFreqs := make(map[string]int)
	for _, l0 := range strings.Split(text, "\n") {
		nLines++
		line := preFilterLine(l0)
		if skip(line) {
			nLinesSkipped++
		} else {
			words := tokenizeLine(line)
			if len(words) > 0 {
				for _, word := range words {
					wordFreqs[word]++
				}
			}
		}
	}
	return nLines, nLinesSkipped, wordFreqs
}

func loadXml(output io.Writer, path string, pageLimit int, logAt int) (int, int, int, int, int, map[string]int) {
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
	wordFreqs := make(map[string]int)

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		if pageLimit > 0 && nPages >= pageLimit {
			clearProgress()
			log.Println(fmt.Sprintf("!! Break called at %d pages (limit set by user)", nPages))
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
		fmt.Fprintln(os.Stderr, "   \t<limit> limit number of pages to read (optional)")
		fmt.Fprintln(os.Stderr, "EXAMPLE\tgo run wikistats.go https://dumps.wikimedia.org/svwiki/latest/svwiki-latest-pages-articles-multistream.xml.bz2 10000")
		os.Exit(1)
	}

	output := bufio.NewWriter(os.Stdout)

	log.Print("*** RUNNING wikistats.main() ***")

	start := time.Now()
	defer output.Flush()

	path := os.Args[1]
	log.Print("Path  : ", path)

	pageLimit := -1
	if len(os.Args) == 3 {
		p, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		pageLimit = p
	}
	if pageLimit > 0 {
		log.Print("Limit : ", pageLimit)
	} else {
		log.Print("Limit : ", "None")
	}

	logAt := 100
	nPages, nRedirects, nLines, nLinesSkipped, nWords, wordFreqs := loadXml(output, path, pageLimit, logAt)

	loaded := time.Now()

	for _, pair := range SortByWordCount(wordFreqs) {
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

	log.Print("Load took             : ", fmt.Sprintf("%12v\n", loadDur))
	log.Print("Print took            : ", fmt.Sprintf("%12v\n", printDur))
	log.Print("Total dur             : ", fmt.Sprintf("%12v\n", totalDur))

	log.Print("No. of pages          : ", lIntPrettyPrint(nPages))
	log.Print("No. of redirects      : ", lIntPrettyPrint(nRedirects))
	log.Print("No. of lines          : ", lIntPrettyPrint(nLines))
	log.Print("No. of skipped lines  : ", lIntPrettyPrint(nLinesSkipped))
	log.Print("No. of words          : ", lIntPrettyPrint(nWords))
	log.Print("No. of unique words   : ", lIntPrettyPrint(len(wordFreqs)))

}
