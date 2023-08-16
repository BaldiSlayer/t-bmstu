package acmp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func decodeWindows1251(reader io.Reader) (io.Reader, error) {
	decoder := charmap.Windows1251.NewDecoder()
	return transform.NewReader(reader, decoder), nil
}

func saveToFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func extractValueSuf(part, prefix string) string {
	index := strings.Index(part, prefix)
	if index >= 0 {
		value := strings.TrimSpace(part[:index])
		return value
	}
	return ""
}

func extractValuePref(part, prefix string) string {
	index := strings.Index(part, prefix)
	if index >= 0 {
		value := strings.TrimSpace(part[index+len(prefix):])
		return strings.Replace(value, " сек.", "", -1)
	}
	return ""
}

func getMiddle(start *goquery.Selection, end string) string {
	if start == nil {
		return ""
	}

	text := ""
	currentElement := start.Next()

	for currentElement.Length() != 0 && !(strings.HasPrefix(currentElement.Text(), end)) {
		htmlContent, err := currentElement.Html()
		if err != nil {
			log.Fatal(err)
		}

		text += htmlContent + " "
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
			if cells.Length() == 3 {
				inputData, _ := cells.Eq(1).Html()
				outputData, _ := cells.Eq(2).Html()

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

// Submit - Функция для выполнения второго запроса
func Submit(client *http.Client, submitURL string, fileData url.Values, taskId string) (string, error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Добавление пустого файла в multipart форму
	formFile, _ := writer.CreateFormFile("fname", "")
	formFile.Write([]byte(""))

	for key, val := range fileData {
		_ = writer.WriteField(key, val[0])
	}
	writer.Close()

	req, err := http.NewRequest("POST", submitURL, &buffer)
	if err != nil {
		return "0", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return "0", err
	}
	defer resp.Body.Close()

	// https://acmp.ru/index.asp?main=status&id_mem=333835&id_res=0&id_t=1&page=0

	forIdUrl := fmt.Sprintf("https://acmp.ru/index.asp?main=status&id_mem=%d&id_res=0&id_t=%s&page=0", 333835, taskId)
	result, err := http.Get(forIdUrl)
	if err != nil {
		fmt.Println("Error:", err)
		return "1", err
	}
	defer result.Body.Close()

	utf8Reader, err := decodeWindows1251(result.Body)
	if err != nil {
		log.Fatal(err)
		return "1", err
	}

	doc, err := goquery.NewDocumentFromReader(utf8Reader)
	if err != nil {
		log.Fatal(err)
		return "1", err
	}

	table := doc.Find("table.main.refresh[align='center']")
	if table.Length() > 0 {
		// Найти первую строку таблицы, которая не является заголовком
		rows := table.Find("tr")
		for i := 1; i < rows.Length(); i++ { // Начинаем с 1, чтобы пропустить заголовок
			row := rows.Eq(i)
			columns := row.Find("td")
			if columns.Length() > 0 {
				id := columns.Eq(0).Text()
				fmt.Sprintf(id)
				return id, nil
			}
		}
	} else {
		fmt.Println("Table not found")
	}

	return "1", nil
}

func endChecking(verdict string) bool {
	if verdict == "Compilation error" || verdict == "Wrong answer" || verdict == "Accepted" ||
		verdict == "Time limit exceeded" || verdict == "Memory limit exceeded" || verdict == "Runtime error (non-zero exit code)" ||
		verdict == "Runtime error" {
		return true
	}
	return false
}

type Task struct {
	ID   string
	Name string
}

func removeLeadingZeros(s string) string {
	trimmed := strings.TrimLeft(s, "0")
	if trimmed == "" {
		return "0"
	}
	return trimmed
}

func GetTaskList(count int) ([]Task, error) {
	result, err := http.Get("https://acmp.ru/index.asp?main=tasks")
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

	var tasks []Task

	doc.Find("table.main tr.white").Each(func(i int, s *goquery.Selection) {
		if i < count {
			id := s.Find("td").Eq(0).Text()
			name := s.Find("td").Eq(1).Text()

			id = removeLeadingZeros(strings.TrimSpace(id))
			name = strings.TrimSpace(name)

			if id != "" && name != "" {
				idBytes := []byte("acmp" + id)
				tasks = append(tasks, Task{
					ID:   base64.StdEncoding.EncodeToString(idBytes),
					Name: name,
				})
			}
		}
	})

	return tasks, nil
}
