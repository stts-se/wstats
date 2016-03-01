package main

import (
	"fmt"
	"github.com/stts-se/wstats/util"
	"os"
	"strings"
	"testing"
)

var fsExp = "Xpctd: '%v' got: '%v'"
var fsDidntExp = "Didn't expect: '%v'"

func testConvert(input0 string) string {
	input := preFilterLine(util.XString{input0})
	var result string
	if skip(input) {
		result = ""
	} else {
		result = convert(input).Value
	}
	return result
}

func TestAll(t *testing.T) {

	tests := map[string]string{
		"hej. och e[]n apa":                                                            "hej och en apa",
		"Vid [[Teherankonferensen]] med Churchill [[och Roosevelt]] sades":             "Vid Teherankonferensen med Churchill och Roosevelt sades",
		"Vid [[Teherankonferensen|konferensen]] med Churchill [[och Roosevelt]] sades": "Vid konferensen med Churchill och Roosevelt sades",
		"<text xml:space=\"preserve\">{{Taxobox":                                       "",
		"I [[upplysningen]]s Europa under det sena 1700-talet började ordet återigen användas för att beskriva den egna trosuppfattningen. Välkända ateister såsom [[Baron d'Holbach]] (1770), Richard Carlile (1826), Charles Southwell (1842), Charles Bradlaugh (1876) och Anne Besant (1877) använde ordet ateism i betydelsen avsaknad av tro på Gud. Sedan dess har ateistiska tänkare och religionsfilosofer använt ordet i den betydelsen.&lt;ref&gt;Martin M ''Atheism. A Philosophical Justification'', Philadelphia 1990, sid 463ff&lt;/ref&gt;": "I upplysningens Europa under det sena 1700-talet började ordet återigen användas för att beskriva den egna trosuppfattningen Välkända ateister såsom Baron d'Holbach 1770 Richard Carlile 1826 Charles Southwell 1842 Charles Bradlaugh 1876 och Anne Besant 1877 använde ordet ateism i betydelsen avsaknad av tro på Gud Sedan dess har ateistiska tänkare och religionsfilosofer använt ordet i den betydelsen",
		"&lt;ref name=&quot;esa.un.org&quot;&gt;[http://esa.un.org/unpd/wpp/Excel-Data/population.htm/ &quot;World Population Prospects: The 2010 Revision&quot;] [[Förenta nationerna|United Nations]] (Department of Economic and Social Affairs, population division)&lt;/ref&gt;":                                                                                                                                                                                                                                                                       "",
		"{{Webbref | titel = How Space is Explored| url = http://adc.gsfc.nasa.gov/adc/education/space_ex/exploration.html| utgivare = NASA}}&lt;/ref&gt; Fysisk utforskning av rymden genomförs både med [[bemannade rymdfärder]] och av obemannade [[rymdsond]]er.":                                                                                                                                                                                                                                                                                       "",
		"* [http://www.fishbase.org/search.php?lang=Swedish Fishbase], en databas över 29 300 olika fiskarter, deras förekomst och vetenskapliga namn.":                                                                                                                                                                                                                                                                                                                                                                                                     "Fishbase en databas över 29 300 olika fiskarter deras förekomst och vetenskapliga namn",
		"[[Fil:House sparrow04.jpg|miniatyr|vänster|[[Gråsparv]]ens utbredningsområde har expanderat dramatiskt på grund av mänsklig aktivitet.&lt;ref&gt;{{Bokref |efternamn = Newton |förnamn = Ian |år = 2003 |titel = The Speciation and Biogeography of Birds |utgivningsort = Amsterdam |utgivare = Academic Press |isbn = 0-12-517375-X |sid = s. 463}}&lt;/ref&gt; ]]":                                                                                                                                                                              "Gråsparvens utbredningsområde har expanderat dramatiskt på grund av mänsklig aktivitet",
		"[[Kategori:Personer inom Sveriges näringsliv under 1700-talet]]</text>":              "Personer inom Sveriges näringsliv under 1700-talet",
		"      <comment>- externa död länkar + mall fotnoter</comment>":                       "",
		"* {{flaggbild|Norge}} Kommendör med kraschan av [[S:t Olavsorden|Sankt Olavsorden]]": "Kommendör med kraschan av Sankt Olavsorden",
		"I slutet av 1700-talet började Fredrik Blom sin bana som [[lärling]] hos en [[bildhuggare|amiralitetsbildhuggare]] i [[Karlskrona]], vilket så småningom förde honom vidare till [[Kungliga Akademien för de fria konsterna|Konstakademien]] i Stockholm. Bloms [[Mentorskap|mentor]] var [[amiral]] [[Carl August Ehrensvärd (1745–1800)|Carl August Ehrensvärd]], som under en period var utbildningschef för [[svenska marinen]] i [[Karlskrona]]. Under kriget mot [[Ryssland]] kom Blom 1808–1809 trots sin ställning som officer inte i direkt kontakt med krigshändelserna. Däremot kom han att ingå i [[Curt von Stedingk]]s [[stab]] vid förhandlingarna med [[Ryssland]] efter det svenska nederlaget i [[Finland]]. Denna position förde Blom till [[S:t Petersburg]] och [[tsar]] [[Alexander I av Ryssland|Alexanders]] [[hov (uppvaktning)|hov]], vilket måste imponerat på den unge Karlskronabon.":                  "I slutet av 1700-talet började Fredrik Blom sin bana som lärling hos en amiralitetsbildhuggare i Karlskrona vilket så småningom förde honom vidare till Konstakademien i Stockholm Bloms mentor var amiral Carl August Ehrensvärd som under en period var utbildningschef för svenska marinen i Karlskrona Under kriget mot Ryssland kom Blom 1808 1809 trots sin ställning som officer inte i direkt kontakt med krigshändelserna Däremot kom han att ingå i Curt von Stedingks stab vid förhandlingarna med Ryssland efter det svenska nederlaget i Finland Denna position förde Blom till S:t Petersburg och tsar Alexanders hov vilket måste imponerat på den unge Karlskronabon",
		"I slutet av 1700-talet började Fredrik Blom sin bana som [[lärling]] hos en [[bildhuggare|amiralitetsbildhuggare]] i [[Karlskrona]], vilket så småningom förde honom vidare till [[Kungliga Akademien för de fria konsterna|Konstakademien]] i Stockholm. Bloms [[Mentorskap|mentor]] var [[amiral]] [[Carl August Ehrensvärd (1745–1800)|Carl August Ehrensvärd]], som under en period var utbildningschef för [[svenska marinen]] i [[Karlskrona]]. Under kriget mot [[Ryssland]] kom Blom 1808–1809 trots sin ställning som officer inte i direkt kontakt med krigshändelserna. Däremot kom han att ingå i [[Curt von Stedingk]]s [[stab]] vid förhandlingarna med [[Ryssland]] efter det svenska nederlaget i [[Finland]]. Denna position förde Blom till [[S:t Petersburg|Sankt Petersburg]] och [[tsar]] [[Alexander I av Ryssland|Alexanders]] [[hov (uppvaktning)|hov]], vilket måste imponerat på den unge Karlskronabon.": "I slutet av 1700-talet började Fredrik Blom sin bana som lärling hos en amiralitetsbildhuggare i Karlskrona vilket så småningom förde honom vidare till Konstakademien i Stockholm Bloms mentor var amiral Carl August Ehrensvärd som under en period var utbildningschef för svenska marinen i Karlskrona Under kriget mot Ryssland kom Blom 1808 1809 trots sin ställning som officer inte i direkt kontakt med krigshändelserna Däremot kom han att ingå i Curt von Stedingks stab vid förhandlingarna med Ryssland efter det svenska nederlaget i Finland Denna position förde Blom till Sankt Petersburg och tsar Alexanders hov vilket måste imponerat på den unge Karlskronabon",
		"'''Jakarta''' (även '''Djakarta''', distriktnamn ''Jakarta Raya'' eller ''DKI Jakarta'', före [[1949]] ''Batavia'') är [[huvudstad]]en i [[Indonesien]] och är belägen på ön [[Java]]. Staden hade 8&amp;nbsp;839&amp;nbsp;247 invånare [[2008]]&lt;ref&gt;[http://www.kependudukancapil.go.id/index.php?option=com_content&amp;view=article&amp;id=4&amp;Itemid=63 Penduduk Provinsi DKI Jakarta: Penduduk Provinsi DKI Jakarta Januari 2008 (Demographics and Civil Records Service: Population of the Province of Jakarta January 2008]&lt;/ref&gt;. Storstadsregionen, som benämns ''[[Jabodetabekjur]]''":                                                                                                                                                                                                                                                                                                                      "Jakarta även Djakarta distriktnamn Jakarta Raya eller DKI Jakarta före 1949 Batavia är huvudstaden i Indonesien och är belägen på ön Java Staden hade 8839247 invånare 2008",
		"På öns östra del, vid Öresundskusten, finns [[Amager Strandpark]]&lt;!--stort S på danska--&gt; med en populär sandstrand. Området har omgestaltats, med en [[konstgjord ö]] och en [[lagun]] innanför. Nyinvigningen av parken, som funnits sedan 1934, ägde rum 2005.":                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            "På öns östra del vid Öresundskusten finns Amager Strandpark med en populär sandstrand Området har omgestaltats med en konstgjord ö och en lagun innanför Nyinvigningen av parken som funnits sedan 1934 ägde rum 2005",
		"<text xml:space=\"preserve\">[[Fil:Paul Heinrich Dietrich Baron d'Holbach Roslin.jpg|miniatyr|[[Baron d'Holbach]], [[Frankrike|fransk]] [[1700-talet|1700-tals]][[författare]], som var en av de första att beskriva sig själv som ateist, och som betytt mycket för ateismens utveckling.]]":                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       "fransk 1700-talsförfattare som var en av de första att beskriva sig själv som ateist och som betytt mycket för ateismens utveckling", // impossible parsing...!
		"&quot;De söner som de får räknas som äktfödda. Med [[Dödsstraff|döden straffas]] däremot den som har [[samlag]] med nästas hustru eller våldtar en [[jungfru]] eller plundrar grannens egendom eller gör honom orätt. Även om nordbor utmärker sig för gästfrihet, ligger svenskarna ett steg före. De räknar det som den värsta skam att neka resande gästvänskap, ja, det härskar en ivrig kapplöpning om vem som anses värdig att mottaga gästen. Där visas denna all möjlig vänlighet och så länge han önskar stanna förs han hem till den ena efter den andra av värdens vänner. Sådana vackra drag finns det bland deras sedvänjor&quot;.":                                                                                                                                                                                                                                                                                    "De söner som de får räknas som äktfödda Med döden straffas däremot den som har samlag med nästas hustru eller våldtar en jungfru eller plundrar grannens egendom eller gör honom orätt Även om nordbor utmärker sig för gästfrihet ligger svenskarna ett steg före De räknar det som den värsta skam att neka resande gästvänskap ja det härskar en ivrig kapplöpning om vem som anses värdig att mottaga gästen Där visas denna all möjlig vänlighet och så länge han önskar stanna förs han hem till den ena efter den andra av värdens vänner Sådana vackra drag finns det bland deras sedvänjor",
		"\"De söner som de får räknas som äktfödda. Med [[Dödsstraff|döden straffas]] däremot den som har [[samlag]] med nästas hustru eller våldtar en [[jungfru]] eller plundrar grannens egendom eller gör honom orätt. Även om nordbor utmärker sig för gästfrihet, ligger svenskarna ett steg före. De räknar det som den värsta skam att neka resande gästvänskap, ja, det härskar en ivrig kapplöpning om vem som anses värdig att mottaga gästen. Där visas denna all möjlig vänlighet och så länge han önskar stanna förs han hem till den ena efter den andra av värdens vänner. Sådana vackra drag finns det bland deras sedvänjor\".":                                                                                                                                                                                                                                                                                            "De söner som de får räknas som äktfödda Med döden straffas däremot den som har samlag med nästas hustru eller våldtar en jungfru eller plundrar grannens egendom eller gör honom orätt Även om nordbor utmärker sig för gästfrihet ligger svenskarna ett steg före De räknar det som den värsta skam att neka resande gästvänskap ja det härskar en ivrig kapplöpning om vem som anses värdig att mottaga gästen Där visas denna all möjlig vänlighet och så länge han önskar stanna förs han hem till den ena efter den andra av värdens vänner Sådana vackra drag finns det bland deras sedvänjor",
		"<text xml:space=\"preserve\">{| class=&quot;infobox&quot; style=&quot;font-size:90%;&quot; width=&quot;300&quot;": "",
		"<redirect title=\"Användbarhet\" />":                                                                              "",
		"Trots sitt namn är inte [[anarki]] och anarkism samma sak som [[kaos]]. Istället är anarkister ofta inriktade på lokal [[direktdemokrati]] som även skall gälla över ekonomin.&lt;ref&gt;{{webbref |url=http://www.mutualist.org/id107.html |titel=Carson, Kevin ‘’Studies in Mutualist Political Economy (2004) |hämtdatum= |format= |verk= }} Such a project requires self-organization at the grassroots level to build &quot;alternative social infrastructure.&quot; It entails things like producers' and consumers' co-ops, LETS systems and mutual banks, syndicalist industrial unions, tenant associations and rent strikes, neighborhood associations, (non-police affiliated) crime-watch and cop-watch programs, voluntary courts for civil arbitration, community-supported agriculture, etc.&lt;/ref&gt;": "trots sitt namn är inte anarki och anarkism samma sak som kaos istället är anarkister ofta inriktade på lokal direktdemokrati som även skall gälla över ekonomin",
		"Trots sitt namn är inte [[anarki]] och anarkism samma sak som [[kaos]]. Istället är anarkister ofta inriktade på lokal [[direktdemokrati]] som även skall gälla över ekonomin.<ref>{{webbref |url=http://www.mutualist.org/id107.html |titel=Carson, Kevin ‘’Studies in Mutualist Political Economy (2004) |hämtdatum= |format= |verk= }} Such a project requires self-organization at the grassroots level to build &quot;alternative social infrastructure.&quot; It entails things like producers' and consumers' co-ops, LETS systems and mutual banks, syndicalist industrial unions, tenant associations and rent strikes, neighborhood associations, (non-police affiliated) crime-watch and cop-watch programs, voluntary courts for civil arbitration, community-supported agriculture, etc.</ref>":             "trots sitt namn är inte anarki och anarkism samma sak som kaos istället är anarkister ofta inriktade på lokal direktdemokrati som även skall gälla över ekonomin",
		"[[1836]] gifte sig Charles Dickens med [[Catherine Hogarth]] och de fick tio barn varav ett dog då det var 8 månader. Samma år utsågs han till [[redaktör]] för ''[[Bentley's Miscellany]]''. Han behöll denna post till [[1839]] då han blev osams med ägaren. Hans framgång som romanförfattare fortsatte samtidigt. Han skrev ''[[Oliver Twist]]'' (1837–1839), ''[[Nicholas Nickleby]]'' (1838–1839), sedan ''[[Den gamla antikvitetshandeln]]'' och ''[[Barnaby Rudge]]'' som del av serien ''[[Mäster Humphreys klocka]]'' (1840–1841). Alla dessa publicerades i månatliga avsnitt innan de gavs ut som böcker.":                                                                                                                                                                                                  "1836 gifte sig Charles Dickens med Catherine Hogarth och de fick tio barn varav ett dog då det var 8 månader Samma år utsågs han till redaktör för Bentley's Miscellany Han behöll denna post till 1839 då han blev osams med ägaren Hans framgång som romanförfattare fortsatte samtidigt Han skrev Oliver Twist 1837 1839 Nicholas Nickleby 1838 1839 sedan Den gamla antikvitetshandeln och Barnaby Rudge som del av serien Mäster Humphreys klocka 1840 1841 Alla dessa publicerades i månatliga avsnitt innan de gavs ut som böcker",
	}

	for input, expect0 := range tests {
		expect := strings.ToLower(expect0)
		result := testConvert(input)
		if result != expect {
			t.Errorf(fsExp, expect, result)
		}
	}
	fmt.Fprint(os.Stderr, "[wstats_test] ", len(tests), " tests passed\n")
}