var accFreq = 0
var nWords = 0

def printWriter(f: String) = new java.io.PrintWriter(new java.io.File(f), "UTF-8")

def createPFile(pFile: String, datFile: String, epsFile: String): Unit = {
  val pw = printWriter(pFile)
  pw.println("""# Gnuplot script file for plotting data in file """" + datFile + """"
# This file is called gnuplot_test.p
set autoscale                       # scale axes automatically
unset log                              # remove any log-scaling
unset label                            # remove any previous labels
set style data linespoints
set xtic auto                          # set xtics automatically
set ytic auto                          # set ytics automatically
set title "Wikipedia lexicon coverage"
set xlabel "No. of words in lexicon"
set ylabel "Coverage (%)"
set xtics nomirror
set ytics nomirror
set border 3
set terminal postscript eps
set output """" + epsFile + """"
plot    """" + datFile + """" using 1:2 title "coverage" pt 7 ps 0.5 lw 1 lc rgb "blue" with linespoints
""")
  pw.flush()
  pw.close()
}

if (args.length!=2) {
  Console.err.println("Script to generate a gnuplot eps file from wstats output (word frequency list)")
  Console.err.println("- requires gnuplot")
  Console.err.println("ARGS: <InputWordFreqs> <OutputEpsFile>")
  System.exit(1)
}

val wFreqsFile=args(0)
val outputFile=args(1).replaceFirst("\\.[^.]+$","")
val outputDatFile=outputFile + ".dat"
val outputPFile=outputFile + ".p"
val outputEpsFile=outputFile + ".eps"

var totFreq = io.Source.fromFile(wFreqsFile).getLines.map(l => l.split("\t").head.toLong).sum

createPFile(outputPFile, outputDatFile, outputEpsFile)

Console.err.println("INPUT=" + wFreqsFile)
Console.err.println("TOTFREQ=" + totFreq)
//Console.err.println("DATFILE=" + outputDatFile)
//Console.err.println("GNUPLOT=" + outputPFile)

val pw = printWriter(outputDatFile)

pw.println("""# Text coverage Wikipedia
# No. wds	% Coverage	Word""")

for (l <- io.Source.fromFile(wFreqsFile).getLines) {
  val fs = l.split("\t", -1).toList
  val f = fs.head.toInt
  val w = fs(1)
  accFreq = accFreq + f
  nWords = nWords + 1
  if (nWords % 100 == 0 || accFreq == totFreq) {
    val coverage = (accFreq*100d)/totFreq
    val covS = f"$coverage%2.2f %%"
    pw.println(nWords + "\t" + covS + "\t" + w)
  }
}
pw.flush()
pw.close()


Console.err.println("Generating output file ... ")

import sys.process._
import scala.language.postfixOps
val gnuplotCmd = "gnuplot " + outputPFile
gnuplotCmd !

Console.err.println("OUTPUT=" + outputEpsFile)
