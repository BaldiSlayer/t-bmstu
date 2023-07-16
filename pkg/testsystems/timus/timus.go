package timus

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
)

var (
	Languages = []string{"FreePascal 2.6",
		"Visual C 2019",
		"Visual C++ 2019",
		"Visual C 2019 x64",
		"Visual C++ 2019 x64",
		"GCC 9.2 x64",
		"G++ 9.2 x64",
		"Clang++ 10 x64",
		"Java 1.8",
		"Visual C# 2019",
		"Python 3.8 x64",
		"PyPy 3.8 x64",
		"Go 1.14 x64",
		"Ruby 1.9",
		"Haskell 7.6",
		"Scala 2.11",
		"Rust 1.58 x64",
		"Kotlin 1.4.0",
	}

	Codes = map[string]string{
		"FreePascal 2.6":      "31",
		"Visual C 2019":       "63",
		"Visual C++ 2019":     "64",
		"Visual C 2019 x64":   "65",
		"Visual C++ 2019 x64": "66",
		"GCC 9.2 x64":         "67",
		"G++ 9.2 x64":         "68",
		"Clang++ 10 x64":      "69",
		"Java 1.8":            "32",
		"Visual C# 2019":      "61",
		"Python 3.8 x64":      "57",
		"PyPy 3.8 x64":        "71",
		"Go 1.14 x64":         "58",
		"Ruby 1.9":            "18",
		"Haskell 7.6":         "19",
		"Scala 2.11":          "33",
		"Rust 1.58 x64":       "72",
		"Kotlin 1.4.0":        "60",
	}
)

func GetProblem(task_id string) (string, error) {
	doc, err := goquery.NewDocument("https://acm.timus.ru/problem.aspx?space=1&locale=ru&num=" + task_id)
	if err != nil {
		return "", err
	}

	// Find the div with class "problem_content"
	problemContent := doc.Find("div.problem_content")

	// Get the HTML content of the div
	htmlContent, err := problemContent.Html()

	return htmlContent, err
}

func SendSubmission(judge_id string, language string, task_id string, code string) error {
	// TODO давай по новой, Леша, все хня
	// This is the function, that send solution to the timus
	url_ := "https://acm.timus.ru/submit.aspx"

	r := url.Values{
		"action":     {"submit"},
		"SpaceID":    {"1"},
		"JudgeID":    {judge_id},
		"Language":   {language},
		"ProblemNum": {task_id},
		"Source":     {code},
	}

	resp, err := http.PostForm(url_, r)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return err
}
