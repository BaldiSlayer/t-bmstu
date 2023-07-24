package timus

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Baldislayer/t-bmstu/pkg/repository"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
	"time"
)

type Task struct {
	ID   string
	Name string
}

type Timus struct {
	Name string
}

func getMiddle(start *goquery.Selection, end string) string {
	if start == nil {
		return ""
	}

	text := ""
	currentElement := start

	for currentElement.Length() != 0 && !(strings.HasPrefix(currentElement.Text(), end)) {
		if currentElement.Is("div.problem_par") {
			text += currentElement.Text() + " "
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

func (t *Timus) GetName() string {
	return t.Name
}

func (t *Timus) GetLanguages() []string {
	return []string{"FreePascal 2.6",
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
}

func (t *Timus) GetProblem(taskID string) (repository.Task, error) {
	url := fmt.Sprintf("https://acm.timus.ru/problem.aspx?space=1&locale=ru&num=%s", taskID)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	if doc.Find("div.problem_content").Length() == 0 {
		return repository.Task{}, err
	}

	task := repository.Task{}

	problemContent := doc.Find("div.problem_content")

	// Get problem information
	task.Name = strings.Split(problemContent.Find("h2.problem_title").Text(), ". ")[1]

	// Get constraints
	limitsText := problemContent.Find("div.problem_limits").Text()
	limitsTextSlice := strings.Split(limitsText, "Ограничение ")[1:]

	timeConstraint := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(limitsTextSlice[0], "времени: ", ""), "секунды", ""))
	memoryConstraint := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(limitsTextSlice[1], "памяти: ", ""), "МБ", ""))

	task.Constraints = map[string]string{
		"time":   timeConstraint,
		"memory": memoryConstraint,
	}

	problemTextDiv := problemContent.Find("div#problem_text")

	// Get condition
	condition := getMiddle(
		problemTextDiv.Find("div.problem_par").First(),
		"Исходные данные",
	)
	task.Condition = strings.TrimSpace(condition)

	// Get input data
	task.InputData = getMiddle(
		problemTextDiv.Find("h3.problem_subtitle:contains('Исходные данные')").Next(),
		"Результат",
	)

	// Get output data
	task.OutputData = getMiddle(
		problemContent.Find("h3.problem_subtitle:contains('Результат')").Next(),
		"Пример",
	)

	// Get tests
	testsTable := problemTextDiv.Find("table.sample")
	task.Tests = map[string]interface{}{
		"tests": parseTableToJSON(testsTable),
	}

	// Get source
	task.Source = problemContent.Find("div.problem_source").Text()

	task.AdditionalInfo = ""

	return task, nil
}

func (t *Timus) Submit(login string, id string, SourceCode string, Language string,
	contestId int, contestTaskId int) error {
	currentTime := time.Now()
	currentTimeString := currentTime.Format("2006-01-02 15:04:05")
	submission := repository.Submission{
		SenderLogin:    login,
		TaskID:         id,
		TestingSystem:  t.GetName(),
		Code:           SourceCode,
		Language:       Language,
		ContestTaskID:  contestTaskId,
		ContestID:      contestId,
		SubmissionTime: currentTimeString,
		SVerdictID:     "-",
	}

	// TODO get err from here
	repository.AddSubmission(submission)

	return nil
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
