package util

import (
	"regexp"
	"sort"
	"strings"
)

// start: libraries
type XString struct {
	Value string
}

func (s XString) String() string {
	return s.Value
}

func (s XString) ReplaceAll(fromRe string, to string) XString {
	re := regexp.MustCompile(fromRe)
	return XString{re.ReplaceAllString(s.Value, to)}
}

func (s XString) MatchesRe(re *regexp.Regexp) bool {
	return re.MatchString(s.Value)
}

func (s XString) Contains(substring string) bool {
	return strings.Contains(s.Value, substring)
}

func (s XString) ToLower() XString {
	return XString{strings.ToLower(s.Value)}
}

func (s XString) Trim() XString {
	return XString{strings.TrimSpace(s.Value)}
}

func (s XString) SplitWhiteSpace() XArray {
	splitted := strings.Split(s.Value, " ")
	result := make([]XString, 0)
	for _, v0 := range splitted {
		v := strings.TrimSpace(v0)
		if len(v) > 0 {
			result = append(result, XString{v})
		}
	}
	return XArray{result}
}

type XArray struct {
	Value []XString
}

func (x XArray) Filter(f func(XString) bool) XArray {
	vsf := make([]XString, 0)
	for _, v := range x.Value {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return XArray{vsf}
}

func (x XArray) Map(f func(XString) string) XArray {
	vsm := make([]XString, len(x.Value))
	for i, v := range x.Value {
		vsm[i] = XString{f(v)}
	}
	return XArray{vsm}
}
func SortByWordCount(wordFrequencies map[XString]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key   XString
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// end: libraries
