package timus

import (
	"errors"
	"fmt"
	"github.com/Baldislayer/t-bmstu/pkg/database"
	"github.com/PuerkitoBio/goquery"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

type Task struct {
	ID   string
	Name string
}

type Timus struct {
	Name           string
	m              sync.RWMutex
	LanguagesMap   map[string]string
	LanguagesSlice []string
}

func (t *Timus) Init() {
	err := t.SetLanguages()
	if err != nil {
		log.Fatal("TImus init failed")
	}
}

func (t *Timus) GetName() string {
	return t.Name
}

func (t *Timus) CheckLanguage(language string) bool {
	t.m.RLock()
	defer t.m.RUnlock()

	_, exist := t.LanguagesMap[language]
	if !exist {
		return false
	}

	return true
}

func (t *Timus) GetLanguages() []string {
	t.m.RLock()
	defer t.m.RUnlock()

	return t.LanguagesSlice
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
					id, err := t.Submit(account.JudgeID, account.Name, submission)
					if err != nil {
						fmt.Println(err)
						// скипаем эту посылку
						// подумать, может ее и удалить еще?
						// было бы круто наверное кидать ее в конец очереди, но мне чет лень это писать
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

// Submit - функция, которая отправляет посылку
func (t *Timus) Submit(judgeId string, accountName string, submission database.Submission) (string, error) {
	url_ := "https://acm.timus.ru/submit.aspx"

	t.m.RLock()
	val, exist := t.LanguagesMap[submission.Language]
	t.m.RUnlock()

	if !exist {
		return "-1", errors.New("No such language")
	}

	resp, err := http.PostForm(url_, url.Values{
		"action":     {"submit"},
		"SpaceID":    {"1"},
		"JudgeID":    {judgeId},
		"Language":   {val},
		"ProblemNum": {submission.TaskID},
		"Source":     {string(submission.Code)},
	})

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

	// ошибка не знаем такой язык
	if strings.Contains(htmlContent, "Unknown language") {
		// значит надо идти парсить
		t.SetLanguages()

		// теперь надо расстоянием Левенштейна найти наиболее похожее
		t.m.RLock()
		nearestLanguage := getNearest(submission.Language, t.LanguagesSlice)
		t.m.RUnlock()

		submission.Language = nearestLanguage
		// надо пойти в базу данных и поменять язык

		return t.Submit(judgeId, accountName, submission)
	}

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

// SetLanguages - устанавливает поля LanguagesMap и LanguagesSlice
// структуры Timus
func (t *Timus) SetLanguages() error {
	languages, err := parseLanguages()
	if err != nil {
		return err
	}

	t.m.Lock()
	defer t.m.Unlock()
	t.LanguagesMap = languages

	t.LanguagesSlice = nil
	for elem := range languages {
		t.LanguagesSlice = append(t.LanguagesSlice, elem)
	}
	sort.Strings(t.LanguagesSlice)

	return nil
}

func getNearest(language string, languages []string) string {
	minDistance := 10000
	nearestLanguage := ""

	for _, lang := range languages {
		distance := levenshtein.DistanceForStrings([]rune(language), []rune(lang), levenshtein.DefaultOptions)
		if distance < minDistance {
			minDistance = distance
			nearestLanguage = lang
		}
	}

	return nearestLanguage
}
