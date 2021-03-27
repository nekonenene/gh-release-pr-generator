package cli

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/google/go-github/v34/github"
)

type templateParameters struct {
	Time         time.Time
	Year         string
	YearShort    string
	Month        string
	MonthShort   string
	Date         string
	DateShort    string
	Weekday      string
	WeekdayShort string
	Hour         string
	HourShort    string
	Minute       string
	MinuteShort  string
	Second       string
	SecondShort  string
	Pulls        []*github.PullRequest
}

const (
	titleTemplateStringDefault = "Release {{ .Year }}-{{ .Month }}-{{ .Date }} {{ .Hour }}:{{ .Minute }}"
	bodyTemplateStringDefault  = `# Pull Requests
{{ range $i, $pull := .Pulls }}
* {{ $pull.Title }} (#{{ $pull.Number }}) @{{ $pull.User.Login }}
{{- end }}`
)

func ConstructTitleAndBody(pulls []*github.PullRequest) (string, string, error) {
	var title, body string
	var titleTemplateString, bodyTemplateString string
	var templateParams templateParameters

	templateParams = setTimeParams(templateParams)
	templateParams.Pulls = pulls

	titleTemplateString, bodyTemplateString, err := readTemplateFile()
	if err != nil {
		return title, body, err
	}
	if titleTemplateString == "" {
		titleTemplateString = titleTemplateStringDefault
	}
	if bodyTemplateString == "" {
		bodyTemplateString = bodyTemplateStringDefault
	}

	titleBuffer := new(bytes.Buffer)
	titleTemplate, err := template.New("title").Parse(titleTemplateString)
	if err != nil {
		return title, body, err
	}
	err = titleTemplate.Execute(titleBuffer, templateParams)
	if err != nil {
		return title, body, err
	}
	title = titleBuffer.String()

	bodyBuffer := new(bytes.Buffer)
	bodyTemplate, err := template.New("body").Parse(bodyTemplateString)
	if err != nil {
		return title, body, err
	}
	err = bodyTemplate.Execute(bodyBuffer, templateParams)
	if err != nil {
		return title, body, err
	}
	body = bodyBuffer.String()

	return title, body, nil
}

func readTemplateFile() (string, string, error) {
	var title, body string

	if params.TemplatePath == "" {
		return title, body, nil
	}

	file, err := os.Open(params.TemplatePath)
	if err != nil {
		return title, body, err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return title, body, err
	}

	fileString := string(buf)
	splitedString := strings.Split(fileString, "\n")
	title = splitedString[0]
	body = strings.Join(splitedString[1:], "\n")

	return title, body, nil
}

func setTimeParams(params templateParameters) templateParameters {
	currentTime := time.Now()

	params.Time = currentTime
	params.Year = currentTime.Format("2006")
	params.YearShort = currentTime.Format("06")
	params.Month = currentTime.Format("01")
	params.MonthShort = currentTime.Format("1")
	params.Date = currentTime.Format("02")
	params.DateShort = currentTime.Format("2")
	params.Weekday = currentTime.Format("Monday")
	params.WeekdayShort = currentTime.Format("Mon")
	params.Hour = currentTime.Format("15")
	params.HourShort = strconv.Itoa(currentTime.Hour())
	params.Minute = currentTime.Format("04")
	params.MinuteShort = currentTime.Format("4")
	params.Second = currentTime.Format("05")
	params.SecondShort = currentTime.Format("5")

	return params
}
