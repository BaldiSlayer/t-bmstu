package timus

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

var tasks []Task

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

func constructURL(id string, count int) string {
	return fmt.Sprintf("https://acm.timus.ru/status.aspx?space=1&count=%d&from=%s", count, id)
}

func GetTaskList(from int, count int) ([]Task, error) {
	doc, err := goquery.NewDocument("https://acm.timus.ru/problemset.aspx?space=1&page=all&locale=ru")
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		// Найти таблицу по классу и выполнить парсинг строк
		doc.Find("table.problemset tr.content").Each(func(i int, s *goquery.Selection) {
			id := s.Find("td").Eq(1).Text()
			name := s.Find("td.name a").Text()

			name = strings.TrimSpace(name)

			if id != "" && name != "" {
				idBytes := []byte("timus" + id)
				tasks = append(tasks, Task{
					ID:   base64.StdEncoding.EncodeToString(idBytes),
					Name: name,
				})
			}
		})
	}

	return tasks[from-1 : from+count-1], nil
}

func endChecking(verdict string) bool {
	if verdict == "Compilation error" || verdict == "Wrong answer" || verdict == "Accepted" ||
		verdict == "Time limit exceeded" || verdict == "Memory limit exceeded" || verdict == "Runtime error (non-zero exit code)" ||
		verdict == "Runtime error" {
		return true
	}
	return false
}

func parseLanguages() (map[string]string, error) {
	resp, err := http.Get("https://acm.timus.ru/submit.aspx")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	languages := make(map[string]string)
	doc.Find("select[name='Language']").Find("option").Each(func(i int, s *goquery.Selection) {
		value, _ := s.Attr("value")
		text := s.Text()
		if value != "" && text != "" {
			languages[text] = value
		}
	})

	return languages, nil
}
