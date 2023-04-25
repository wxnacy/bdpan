package bdpan

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/mattn/go-runewidth"
)

type PrettyData struct {
	Name       string
	Value      string
	IsFillLeft bool
}

type Pretty interface {
	BuildPrintData() []PrettyData
}

type PrettyEnumerate interface {
	GetWriter() io.Writer
	List() []Pretty
}

type PrettyList []Pretty

func (p PrettyList) GetWriter() io.Writer {
	return os.Stdout
}

func (p PrettyList) List() []Pretty {
	return p
}

func PrettyPrintList(p PrettyEnumerate) {
	items := p.List()
	if len(items) == 0 {
		return
	}
	maxLengthMap := getMaxLengthMap(items)
	prettyMap := getPrettyMap(items[0])
	firstLine := ""
	for _, pd := range items[0].BuildPrintData() {
		name := pd.Name
		length, _ := maxLengthMap[pd.Name]
		p, _ := prettyMap[name]
		prettyName := fmtString(name, length, p.IsFillLeft)
		firstLine = firstLine + prettyName + " "

	}
	fmt.Fprintln(p.GetWriter(), firstLine)
	for _, pretty := range items {
		line := ""
		for _, pd := range pretty.BuildPrintData() {
			length, _ := maxLengthMap[pd.Name]
			line = line + fmtString(pd.Value, length, pd.IsFillLeft) + " "
		}
		fmt.Fprintln(p.GetWriter(), line)
	}
	fmt.Fprintf(p.GetWriter(), "Total: %d\n", len(items))
}

func getPrettyMap(pretty Pretty) map[string]PrettyData {
	m := make(map[string]PrettyData, 0)
	for _, p := range pretty.BuildPrintData() {
		m[p.Name] = p
	}
	return m
}

func getMaxLengthMap(items []Pretty) map[string]int {
	m := make(map[string]int, 0)
	for _, item := range items {
		data := item.BuildPrintData()
		for _, d := range data {
			name := d.Name
			leng, exit := m[name]
			leng = int(math.Max(float64(len(name)), float64(leng)))
			if exit {
				leng = int(math.Max(float64(runewidth.StringWidth(d.Value)), float64(leng)))
			}
			m[name] = leng
		}
	}
	return m
}

func fmtString(s string, length int, isFillLeft bool) string {
	var value string
	if isFillLeft {
		value = runewidth.FillLeft(s, length)
	} else {
		value = runewidth.FillRight(s, length)
	}
	return value
}
