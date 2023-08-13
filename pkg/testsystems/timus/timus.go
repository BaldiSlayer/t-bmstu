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
	"sync"
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

func (t *Timus) GetName() string {
	return t.Name
}

func (t *Timus) CheckLanguage(language string) bool {
	languagesDict := map[string]struct{}{
		"FreePascal 2.6":      struct{}{},
		"Visual C 2019":       struct{}{},
		"Visual C++ 2019":     struct{}{},
		"Visual C 2019 x64":   struct{}{},
		"Visual C++ 2019 x64": struct{}{},
		"GCC 9.2 x64":         struct{}{},
		"G++ 9.2 x64":         struct{}{},
		"Clang++ 10 x64":      struct{}{},
		"Java 1.8":            struct{}{},
		"Visual C# 2019":      struct{}{},
		"Python 3.8 x64":      struct{}{},
		"PyPy 3.8 x64":        struct{}{},
		"Go 1.14 x64":         struct{}{},
		"Ruby 1.9":            struct{}{},
		"Haskell 7.6":         struct{}{},
		"Scala 2.11":          struct{}{},
		"Rust 1.58 x64":       struct{}{},
		"Kotlin 1.4.0":        struct{}{},
	}

	_, exist := languagesDict[language]

	if !exist {
		return false
	}

	return true
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

func (t *Timus) GetProblem(taskID string) (database.Task, error) {
	taskUrl := fmt.Sprintf("https://acm.timus.ru/problem.aspx?space=1&locale=ru&num=%s", taskID)

	doc, err := goquery.NewDocument(taskUrl)
	if err != nil {
		log.Fatal(err)
	}

	if doc.Find("div.problem_content").Length() == 0 {
		return database.Task{}, err
	}

	task := database.Task{}

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

func (t *Timus) Submitter(wg *sync.WaitGroup, ch chan<- database.Submission) {
	defer wg.Done()

	type Account struct {
		Name      string    `json:"name"`
		JudgeID   string    `json:"judge_id"`
		UsageTime time.Time `json:"usage_time"`
	}

	// TODO добавить считывание с файла
	accounts := []Account{
		{
			Name:    "$tup1d2281337",
			JudgeID: "342187EL",
		},
	}

	timeDifference := 11 * time.Second

	for i := range accounts {
		accounts[i].UsageTime = time.Now().Add(-timeDifference)
	}

	for {
		submissions, err := database.GetSubmitsWithStatus(t.GetName(), 0)
		if err != nil {
			fmt.Println(err)
		}

		// перебираем все решения
		for _, submission := range submissions {
			// перебираем аккаунты
			for i, account := range accounts {
				if elapsedTime := time.Now().Sub(account.UsageTime); elapsedTime > timeDifference {
					id, err := Submit(account.JudgeID, account.Name, submission)
					if err != nil {
						fmt.Println(err)
						// скипаем эту посылку
						// подумать, может ее и удалить еще?
						continue
					}

					// теперь надо передать по каналу, что был изменен статус этой задачи
					submission.Status = 1
					submission.Verdict = "Compiling"
					submission.SubmissionNumber = id
					ch <- submission

					// устанавливаем время
					accounts[i].UsageTime = time.Now()
				}
			}
		}

		time.Sleep(time.Second)
	}
}

func constructURL(id string, count int) string {
	return fmt.Sprintf("https://acm.timus.ru/status.aspx?space=1&count=%d&from=%s", count, id)
}

func (t *Timus) Checker(wg *sync.WaitGroup, ch chan<- database.Submission) {
	defer wg.Done()

	for {
		// получение отправленных, но еще не прошедших проверку посылок
		submissions, err := database.GetSubmitsWithStatus(t.GetName(), 1)

		if err != nil {
			fmt.Println(err)
		}

		submissionsDict := make(map[string]database.Submission)
		submissionsIDs := make([]string, 0)

		for _, submission := range submissions {
			submissionsDict[submission.SubmissionNumber] = submission
			submissionsIDs = append(submissionsIDs, submission.SubmissionNumber)
		}

		for len(submissions) != 0 {
			count := 50
			url := constructURL(submissions[0].SubmissionNumber, count)

			doc, err := goquery.NewDocument(url)
			if err != nil {
				log.Fatal(err)
			}

			doc.Find("tr").Each(func(index int, rowHtml *goquery.Selection) {
				class, exists := rowHtml.Attr("class")
				if exists && (class == "even" || class == "odd") {
					// Получение значения id посылки
					idStr := strings.TrimSpace(rowHtml.Children().First().Text())
					if err != nil {
						log.Println("Error converting id:", err)
						return
					}

					if _, exists := submissionsDict[idStr]; exists {
						// удаление из словаря и списка
						submission, exists := submissionsDict[idStr]
						if !exists {
							log.Println("Submission with ID not found:", idStr)
							return
						}
						delete(submissionsDict, idStr)

						for i, id := range submissionsIDs {
							if id == idStr {
								submissionsIDs = append(submissionsIDs[:i], submissionsIDs[i+1:]...)
								break
							}
						}

						submissions = submissions[1:]

						verdict := strings.TrimSpace(rowHtml.Children().Eq(5).Text())
						test := strings.TrimSpace(rowHtml.Children().Eq(6).Text())
						executionTime := strings.TrimSpace(rowHtml.Children().Eq(7).Text())
						memoryUsed := strings.TrimSpace(rowHtml.Children().Eq(8).Text())

						submission.Verdict = verdict
						submission.Test = test
						submission.ExecutionTime = executionTime
						submission.MemoryUsed = memoryUsed

						if endChecking(verdict) {
							submission.Status = 2
						}

						ch <- submission
					}
				}
			})

		}

		time.Sleep(time.Second)
	}
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
