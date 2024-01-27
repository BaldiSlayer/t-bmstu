package acmp

import (
	"fmt"
	"github.com/Baldislayer/t-bmstu/pkg/database"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"
)

type ACMP struct {
	Name string
}

func (t *ACMP) Init() {

}

func (t *ACMP) GetName() string {
	return t.Name
}

func (t *ACMP) CheckLanguage(language string) bool {
	languagesDict := map[string]struct{}{
		"MinGW GNU C++ 13.1.0":        struct{}{},
		"Python 3.11.0":               struct{}{},
		"PascalABC.NET 3.8.3":         struct{}{},
		"Java SE JDK 16.0.1":          struct{}{},
		"Free Pascal 3.2.2":           struct{}{},
		"Borland Delphi 7.0":          struct{}{},
		"Microsoft Visual C++ 2017":   struct{}{},
		"Microsoft Visual C# 2017":    struct{}{},
		"Microsoft Visual Basic 2017": struct{}{},
		"PyPy3.9 v7.3.9":              struct{}{},
		"Go 1.16.3":                   struct{}{},
		"Node.js 19.0.0":              struct{}{},
	}

	_, exist := languagesDict[language]

	if !exist {
		return false
	}

	return true
}

func (t *ACMP) GetLanguages() []string {
	return []string{
		"MinGW GNU C++ 13.1.0",
		"Python 3.11.0",
		"PascalABC.NET 3.8.3",
		"Java SE JDK 16.0.1",
		"Free Pascal 3.2.2",
		"Borland Delphi 7.0",
		"Microsoft Visual C++ 2017",
		"Microsoft Visual C# 2017",
		"Microsoft Visual Basic 2017",
		"PyPy3.9 v7.3.9",
		"Go 1.16.3",
		"Node.js 19.0.0",
	}
}

func (t *ACMP) Submitter(wg *sync.WaitGroup, ch chan<- database.Submission) {
	defer wg.Done()

	myToACMPDict := map[string]string{
		"MinGW GNU C++ 13.1.0":        "CPP",
		"Python 3.11.0":               "PY",
		"PascalABC.NET 3.8.3":         "PP",
		"Java SE JDK 16.0.1":          "JAVA",
		"Free Pascal 3.2.2":           "PAS",
		"Borland Delphi 7.0":          "DPR",
		"Microsoft Visual C++ 2017":   "CXX",
		"Microsoft Visual C# 2017":    "CS",
		"Microsoft Visual Basic 2017": "BAS",
		"PyPy3.9 v7.3.9":              "PYPY",
		"Go 1.16.3":                   "GO",
		"Node.js 19.0.0":              "JS",
	}
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	// Выполнение первого запроса для аутентификации
	loginURL := "https://acmp.ru/index.asp?main=enter&r=30147517425972369652497"
	loginData := url.Values{
		"lgn":      {"aukseu228"},
		"password": {"M6v-rzz-Hgm-Skg"},
	}
	client.PostForm(loginURL, loginData)

	for {
		submissions, err := database.GetSubmitsWithStatus(t.GetName(), 0)
		if err != nil {
			fmt.Println(err)
		}

		// перебираем все решения
		for _, submission := range submissions {
			fileData := url.Values{
				"lang":   {myToACMPDict[submission.Language]},
				"source": {string(submission.Code)},
			}
			id, err := Submit(client,
				fmt.Sprintf("https://acmp.ru/index.asp?main=update&mode=upload&id_task=%s", submission.TaskID),
				fileData,
				submission.TaskID)
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
		}

		// TODO acmp придется тестировать на rps защиту
		time.Sleep(time.Second * 2)
	}
}

func (t *ACMP) Checker(wg *sync.WaitGroup, ch chan<- database.Submission) {
	// жесткий парсинг таблицы результатов
	defer wg.Done()

	for {
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

		pageNum := 0
		for len(submissions) != 0 {
			currentUrl := fmt.Sprintf("https://acmp.ru/index.asp?main=status&id_mem=%d&id_res=0&id_t=0&page=%d", 333835, pageNum)

			result, err := http.Get(currentUrl)
			if err != nil {
				fmt.Println("Error:", err)
			}
			defer result.Body.Close()

			utf8Reader, err := decodeWindows1251(result.Body)
			if err != nil {
				log.Fatal(err)
			}

			doc, err := goquery.NewDocumentFromReader(utf8Reader)
			if err != nil {
				log.Fatal(err)
			}

			table := doc.Find("table.main.refresh[align='center']")
			table.Find("tr").Each(func(index int, rowHtml *goquery.Selection) {
				columns := rowHtml.Find("td")
				idStr := columns.Eq(0).Text()

				for _, submissionID := range submissionsIDs {
					if idStr == submissionID {
						// Это строка с нужным id, вы можете выполнить здесь нужные действия 5 6 7 8
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

						verdict := strings.TrimSpace(columns.Eq(5).Text())
						test := strings.TrimSpace(columns.Eq(6).Text())
						executionTime := strings.TrimSpace(columns.Eq(7).Text())
						memoryUsed := strings.TrimSpace(columns.Eq(8).Text())

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

			pageNum++
		}

		time.Sleep(time.Second * 2)
	}
}

func (t *ACMP) GetProblem(taskID string) (database.Task, error) {
	taskURL := fmt.Sprintf("https://acmp.ru/index.asp?main=task&id_task=%s", taskID)

	resp, err := http.Get(taskURL)
	if err != nil {
		fmt.Println("Error:", err)
		return database.Task{}, err
	}
	defer resp.Body.Close()

	utf8Reader, err := decodeWindows1251(resp.Body)
	if err != nil {
		log.Fatal(err)
		return database.Task{}, err
	}

	doc, err := goquery.NewDocumentFromReader(utf8Reader)
	if err != nil {
		log.Fatal(err)
		return database.Task{}, err
	}

	var content string
	startTag := "h1"
	endTag := "<h4><i>Для отправки решения задачи необходимо <a href=\"/inc/register.asp\">зарегистрироваться</a> и авторизоваться!</i></h4>"

	taskName := ""

	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
		taskName = s.Text()
	})

	Constraints := map[string]string{}
	ended := false
	doc.Find(startTag).NextUntil(endTag).Each(func(i int, s *goquery.Selection) {
		if s.Text() != "" {
			if !ended {
				ended = strings.Contains(s.Text(), "Для отправки решения задачи")

				if !ended {
					if strings.Contains(s.Text(), "Время") {
						parts := strings.Split(s.Text(), "Память:")
						if len(parts) >= 2 {
							timePart := strings.TrimSpace(parts[0])
							memoryPart := strings.TrimSpace(parts[1])

							time := extractValuePref(timePart, "Время:")
							memory := extractValueSuf(memoryPart, "Мб")

							Constraints = map[string]string{
								"time":   time,
								"memory": memory,
							}
						} else {
							fmt.Println("Не удалось извлечь время и память.")
						}

					} else {
						elem, err := s.Html()
						if err != nil {
							return
						}
						content += elem
					}
				}
			}
		}
	})

	var Condition string
	centerFound := false
	h2InputFound := false
	doc.Find(startTag).NextUntil("<h2>Входные данные</h2>").Each(func(i int, s *goquery.Selection) {
		if !h2InputFound {
			if s.Is("h2") && s.Text() == "Входные данные" {
				h2InputFound = true
				return
			}

			if centerFound {
				elem, _ := s.Html()
				Condition += elem
			}

			if s.Is("center") {
				centerFound = true
			}
		}
	})

	// newDoc, err := goquery.NewDocumentFromReader(utf8Reader)
	if err != nil {
		log.Fatal(err)
		return database.Task{}, err
	}
	Input := getMiddle(doc.Find("h2:contains('Входные данные')").First(),
		"Выходные данные")
	Output := getMiddle(doc.Find("h2:contains('Выходные данные')").First(),
		"Пример")

	tests := parseTableToJSON(doc.Find("h2:contains('Пример')").First().Next())

	return database.Task{
		Name:        taskName,
		Condition:   Condition,
		Constraints: Constraints,
		InputData:   Input,
		OutputData:  Output,
		Tests: map[string]interface{}{
			"tests": tests,
		},
	}, nil
}
