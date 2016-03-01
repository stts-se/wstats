# wstats

`wstats` is a sketch of/place holder for a module to compute word statistics on wikipedia data. It is NOT ready for proper use, so use at your own risk.

Usage:
  
     $ go run wikistats.go <path> <limit>*
       <path> wikimedia dump (file or url, xml or xml.bz)
       <limit> limit number of pages to read (optional)
   	
Example usage:

     $ go run wikistats.go https://dumps.wikimedia.org/svwiki/latest/svwiki-latest-pages-articles-multistream.xml.bz2 10000

The program will print running progress and basic statistics to standard error. A complete word frequency list will be printed to standard out.


<br/>
Xml parsing inspired by : http://blog.davidsingleton.org/parsing-huge-xml-files-with-go
