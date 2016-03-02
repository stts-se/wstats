# wstats

`wstats` is a module for computing word statistics on wikipedia data. It is NOT ready for proper use, so use at your own risk

Usage:
  
     $ go run wstats.go <path> <limit>*
       <path> wikimedia dump (file or url, xml or xml.bz2)
       <limit> limit number of pages to read (optional)
   	
Example usage:

     $ go run wstats.go https://dumps.wikimedia.org/svwiki/latest/svwiki-latest-pages-articles-multistream.xml.bz2 10000

The program will print running progress and basic statistics to standard error.<br/>
A complete word frequency list will be printed to standard out.

Wikipedia dumps: https://dumps.wikimedia.org/backup-index.html

<br/>
Xml parsing inspired by: http://blog.davidsingleton.org/parsing-huge-xml-files-with-go
