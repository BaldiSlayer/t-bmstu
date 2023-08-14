package timus

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Baldislayer/t-bmstu/pkg/database"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func getMiddle(start *goquery.Selection, end string) string {
	if start == nil {
		return ""
	}

	text := ""
	currentElement := start

	for currentElement.Length() != 0 && !(strings.HasPrefix(currentElement.Text(), end)) {
		if currentElement.Is("div.problem_par") || currentElement.Is("div.problem_centered_picture") {
			htmlContent, err := currentElement.Html()
			if err != nil {
				log.Fatal(err)
			}

			htmlContent = strings.Replace(htmlContent, "/image", "https://acm.timus.ru/image", -1)
			text += htmlContent + " "
		}
		currentElement = currentElement.Next()
	}

	return strings.TrimSpace(text)
}

func parseTableToJSON(table *goquery.Selection) string {
	tests := []map[string]string{}

	rows := table.Find("tr")
	rows.Each(func(i int, row *goquery.Selection) {
		if i != 0 { // Skip the header row
			cells := row.Find("td")
			if cells.Length() == 2 {
				inputData := cells.Eq(0).Find("pre").Text()
				outputData := cells.Eq(1).Find("pre").Text()

				test := map[string]string{
					"input":  strings.TrimSpace(inputData),
					"output": strings.TrimSpace(outputData),
				}
				tests = append(tests, test)
			}
		}
	})

	jsonTests, err := json.MarshalIndent(tests, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return string(jsonTests)
}

// Submit - функция, которая отправляет посылку
func Submit(judgeId string, accountName string, submission database.Submission) (string, error) {
	d := map[string]string{
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

	url_ := "https://acm.timus.ru/submit.aspx"

	val, exist := d[submission.Language]

	if !exist {
		return "-1", errors.New("No such language")
	}

	r := url.Values{
		"action":     {"submit"},
		"SpaceID":    {"1"},
		"JudgeID":    {judgeId},
		"Language":   {val},
		"ProblemNum": {submission.TaskID},
		"Source":     {string(submission.Code)},
	}

	resp, err := http.PostForm(url_, r)

	if err != nil {
		return "-1", err
	}

	// TODO научиться проверять есть ли ошибка на самом тимусе

	defer resp.Body.Close()

	// Чтение тела ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "-1", err
	}

	// Вывод тела ответа в консоль
	htmlContent := string(body)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return "-1", err
	}

	taskID := submission.TaskID

	foundedId := "-1"

	doc.Find("table.status.status_nofilter tr").Each(func(i int, row *goquery.Selection) {
		if foundedId == "-1" {

			idValue := row.Find("td.id").Text()
			coderValue := row.Find("td.coder a").Text()
			problemValue := row.Find("td.problem a").Text()

			problemValue = strings.Split(problemValue, ".")[0]

			if coderValue == accountName && problemValue == taskID {
				foundedId = idValue
			}
		}
	})

	return foundedId, nil
}

func constructURL(id string, count int) string {
	return fmt.Sprintf("https://acm.timus.ru/status.aspx?space=1&count=%d&from=%s", count, id)
}

func GetTaskList(count int) ([]Task, error) {
	doc, err := goquery.NewDocument("https://acm.timus.ru/problemset.aspx?space=1&page=all&locale=ru")
	if err != nil {
		return nil, err
	}

	var tasks []Task

	// Найти таблицу по классу и выполнить парсинг строк
	doc.Find("table.problemset tr.content").Each(func(i int, s *goquery.Selection) {
		if i < count {
			id := s.Find("td").Eq(1).Text()
			name := s.Find("td.name a").Text()

			// Если name содержит перенос строки, уберите его с помощью strings.TrimSpace
			name = strings.TrimSpace(name)

			if id != "" && name != "" {
				idBytes := []byte("timus" + id)
				tasks = append(tasks, Task{
					ID:   base64.StdEncoding.EncodeToString(idBytes),
					Name: name,
				})
			}
		}
	})

	return tasks, nil
}

func endChecking(verdict string) bool {
	if verdict == "Compilation error" || verdict == "Wrong answer" || verdict == "Accepted" ||
		verdict == "Time limit exceeded" || verdict == "Memory limit exceeded" || verdict == "Runtime error (non-zero exit code)" ||
		verdict == "Runtime error" {
		return true
	}
	return false
}
