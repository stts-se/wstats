# wstats

`wstats` is a module for computing word statistics on wikipedia data. It is NOT ready for proper use, so use at your own risk

The program will print running progress and basic statistics to standard error.\nA complete word frequency list will be printed to standard out.


Cmd line arguments:

     path to the wikimedia dump file (file or url, xml or xml.bz2) (required)
     -pl=int    page limit: limit number of pages to read (optional, default = unset)
     -fl=int    freq limit: lower limit for word frequencies to be printed (optional, default = 2)
     -h         help: print help message

Example usage:

    $ go run wstats.go https://dumps.wikimedia.org/svwiki/latest/svwiki-latest-pages-articles-multistream.xml.bz2 -pl=10000


The program will print running progress and basic statistics to standard error.<br/>
A complete word frequency list will be printed to standard out.

Wikipedia dumps: https://dumps.wikimedia.org/backup-index.html

<br/>
Xml parsing inspired by: http://blog.davidsingleton.org/parsing-huge-xml-files-with-go
