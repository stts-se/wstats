# wstats

`wstats` is used for parsing wikimedia dump files on the fly into word frequency lists.

The program will print running progress and basic statistics to standard error.
A complete word frequency list will be printed to standard out.

There is no guarantee that all words of the actual Wikipedia article texts are counted, partly because there is no expansion of templates, etc, of the dump file.


Usage:

    $ go run wstats.go <flags> <wikipedia dump path (file or url, xml or xml.bz2)>

Cmd line flags:

     -pl int    page limit: limit number of pages to read (optional, default = unset)
     -mf int    min freq: lower limit for word frequencies to be printed (optional, default = 0)
     -h(elp)    help: print help message

Example usage:

     $ go run wstats.go -pl 10000 https://dumps.wikimedia.org/svwiki/latest/svwiki-latest-pages-articles-multistream.xml.bz2 

The program will print running progress and basic statistics to standard error.<br/>
A complete word frequency list will be printed to standard out (limited by min freq, if set).

Wikipedia dumps: https://dumps.wikimedia.org/backup-index.html

<br/>
Xml parsing inspired by: http://blog.davidsingleton.org/parsing-huge-xml-files-with-go

<br/>

## List of (some) Wikipedia dump files

    https://dumps.wikimedia.org/elwiki/latest/elwiki-latest-pages-articles-multistream.xml.bz2
    https://dumps.wikimedia.org/ruwiki/latest/ruwiki-latest-pages-articles-multistream.xml.bz2


## wfreqs2gnuplot.scala
Scala script to generate a gnuplot eps file from the wstats output (word frequency list).
Requires: scala, gnuplot

---
_This work was supported by the Swedish Post and Telecom Authority (PTS) through the grant "Wikispeech – en användargenererad talsyntes på Wikipedia" (2016–2017)._
