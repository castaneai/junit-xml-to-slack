package main

import (
	"io"
	"encoding/xml"
	"os"
	"log"
	"encoding/json"
	"fmt"
)

func main() {
	report, err := parseJUnitXML(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	payload, err := toSlackPayload(report)

	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(payload); err != nil {
		log.Fatalln(err)
	}
}

type TestSuite struct {
	Name string `xml:"name,attr"`
	TestCount int `xml:"tests,attr"`
	FailureCount int `xml:"failures,attr"`
	ErrorCount int `xml:"errors,attr"`
	Time float64 `xml:"time,attr"`
	TestCases []*TestCase `xml:"testcase"`
}

type TestCase struct {
	Name string `xml:"name,attr"`
	ClassName string `xml:"classname,attr"`
	Time float64 `xml:"time,attr"`
	Failures []*TestFailure `xml:"failure"`
	Skipped []*TestSkipped `xml:"skipped"`
	Errors []*TestError `xml:"error"`
}

func (tc *TestCase) IsError() bool {
	return len(tc.Errors) > 0
}

func (tc *TestCase) IsFailure() bool {
	return len(tc.Failures) > 0
}

func (tc *TestCase) IsSkipped() bool {
	return len(tc.Skipped) > 0
}

func (tc *TestCase) Emoji() string {
	if tc.IsError() || tc.IsFailure() {
		return "❌"
	}
	if tc.IsSkipped() {
		return "⏩"
	}
	return "✅"
}

func (tc *TestCase) Color() string {
	if tc.IsError() {
		return "#e74c3c"
	}
	if tc.IsFailure() {
		return "#f1c40f"
	}
	return "#2ecc71"
}

type TestFailure struct {
	Type string `xml:"type,attr"`
	Message string `xml:"message,attr"`
	Body string `xml:",chardata"`
}

type TestError struct {
	Type string `xml:"type,attr"`
	Message string `xml:"message,attr"`
	Body string `xml:",chardata"`
}

type TestSkipped struct {}

func parseJUnitXML(jr io.Reader) (*TestSuite, error) {
	dec := xml.NewDecoder(jr)

	report := &TestSuite{}
	if err := dec.Decode(&report); err != nil {
		return nil, err
	}
	return report, nil
}

type SlackPayload struct {
	Text string `json:"text,omitempty"`
	Attachments []*SlackAttachment `json:"attachments,omitempty"`
}

type SlackAttachment struct {
	Title string `json:"title,omitempty"`
	TitleLink string `json:"title_link,omitempty"`
	Color string `json:"color,omitempty"`
	Footer string `json:"footer,omitempty"`
	AuthorName string `json:"author_name,omitempty"`
}

func toSlackPayload(suite *TestSuite) (*SlackPayload, error) {
	p := &SlackPayload{}
	p.Text = "UnitTest suite"
	for _, tcase := range suite.TestCases {
		if tcase.IsError() || tcase.IsFailure() {
			at := newTestCaseAttachment(suite, tcase)
			p.Attachments = append(p.Attachments, at)
		}
	}
	return p, nil
}

func newTestCaseAttachment(suite *TestSuite, tc *TestCase) *SlackAttachment {
	return &SlackAttachment{
		Title: fmt.Sprintf("%s %s", tc.Emoji(), tc.Name),
		AuthorName: tc.ClassName,
		Color: tc.Color(),
		Footer: fmt.Sprintf("%.2f ms", tc.Time * 1000),
	}
}