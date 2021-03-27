package cli

import (
	"bytes"
	"strconv"
	"text/template"
	"time"

	"github.com/google/go-github/v34/github"
)

type templateParams struct {
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

func ConstructPullRequest(pulls []*github.PullRequest) (string, string, error) {
	var title, body string
	var params templateParams
	params = setTimeParams(params)
	params.Pulls = pulls

	titleBuffer := new(bytes.Buffer)
	titleTemplate, err := template.New("title").Parse("Release {{ .Year }}-{{ .Month }}-{{ .Date }} {{ .Hour }}:{{ .Minute }}")
	if err != nil {
		return title, body, err
	}
	err = titleTemplate.Execute(titleBuffer, params)
	if err != nil {
		return title, body, err
	}
	title = titleBuffer.String()

	bodyTemplateString := `# Pull Requests
{{ range $i, $pull := .Pulls }}
* {{ $pull.Title }} (#{{ $pull.Number }}) @{{ $pull.User.Login }}
{{- end }}`

	bodyBuffer := new(bytes.Buffer)
	bodyTemplate, err := template.New("body").Parse(bodyTemplateString)
	if err != nil {
		return title, body, err
	}
	err = bodyTemplate.Execute(bodyBuffer, params)
	if err != nil {
		return title, body, err
	}
	body = bodyBuffer.String()

	return title, body, nil
}

func setTimeParams(params templateParams) templateParams {
	currentTime := time.Now()

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
