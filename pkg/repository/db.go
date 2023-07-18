package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"time"
)

var (
	DbURL = ""
)

func CreateTables() error {
	// TODO миграции

	DbURL = fmt.Sprintf(
		"postgresql://%s:%s@localhost:5432/%s",
		viper.GetString("DBUsername"),
		viper.GetString("DBPassword"),
		viper.GetString("DBName"),
	)
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS submissions (
            id SERIAL PRIMARY KEY,
            sender_login VARCHAR(255),
            task_id VARCHAR(255),
            testing_system VARCHAR(255),
            code TEXT,
            submission_time TIMESTAMPTZ,
            contest_id INTEGER,
            contest_task_id INTEGER,
            language VARCHAR(255),
            sverdict_id VARCHAR(255)
        );

        CREATE TABLE IF NOT EXISTS submissions_verdicts (
            id SERIAL PRIMARY KEY,
            sender_login VARCHAR(255),
            task_id VARCHAR(255),
            testing_system VARCHAR(255),
            code TEXT,
            submission_time TIMESTAMPTZ,
            contest_id INTEGER,
            contest_task_id INTEGER,
            verdict VARCHAR(255),
            execution_time VARCHAR(255),
            memory_used VARCHAR(255),
            test VARCHAR(255),
            language VARCHAR(255),
            submission_number VARCHAR(255)
        );

        CREATE TABLE IF NOT EXISTS unverified_submissions (
            submission_id VARCHAR(50) PRIMARY KEY,
            external_submission_id VARCHAR(50),
            testing_system VARCHAR(255),
            judge_id VARCHAR(255)
        );

        CREATE TABLE IF NOT EXISTS contests (
            id SERIAL PRIMARY KEY,
            title VARCHAR(255),
            access JSONB,
            participants JSONB,
            results JSONB,
            tasks JSONB
        );

        CREATE TABLE IF NOT EXISTS tasks (
            id VARCHAR(100) PRIMARY KEY,
            name VARCHAR(255),
            constraints JSONB,
            condition TEXT,
            input_data TEXT,
            output_data TEXT,
            source TEXT,
            additional_info TEXT,
            tests JSONB
        );
    `)
	if err != nil {
		return err
	}

	fmt.Println("Tables created successfully!")
	return nil
}

func AddSubmission(submission Submission) {
	// Добавляем в базу данных
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	// логика - после того, как я принял запрос - я кидаю в таблицу посылок с вердиктами (SubmissionVerdict)
	// затем я кидаю его в таблицу еще неотправленных посылок (submissions)
	// на этом все заканчивается

	verdict := SubmissionVerdict{
		SenderLogin:      submission.SenderLogin,
		TaskID:           submission.TaskID,
		TestingSystem:    submission.TestingSystem,
		Code:             []byte(submission.Code),
		SubmissionTime:   submission.SubmissionTime,
		ContestID:        submission.ContestID,
		ContestTaskID:    submission.ContestTaskID,
		Language:         submission.Language,
		Verdict:          "Waiting",
		ExecutionTime:    "-",
		MemoryUsed:       "-",
		Test:             "-",
		SubmissionNumber: "-",
	}

	var verdictID int
	err = conn.QueryRow(context.Background(), `
    INSERT INTO submissions_verdicts (sender_login, task_id, testing_system, code, submission_time, contest_id, contest_task_id, language, verdict, execution_time, memory_used, test, submission_number)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    RETURNING id
`, verdict.SenderLogin, verdict.TaskID, verdict.TestingSystem, verdict.Code, verdict.SubmissionTime, verdict.ContestID, verdict.ContestTaskID, verdict.Language, verdict.Verdict, verdict.ExecutionTime, verdict.MemoryUsed, verdict.Test, verdict.SubmissionNumber).Scan(&verdictID)
	if err != nil {
		log.Printf("Failed to insert submission verdict: %v", err)
	}

	submission.SVerdictID = strconv.Itoa(verdictID)

	_, err = conn.Exec(context.Background(), `
    INSERT INTO submissions (sender_login, task_id, testing_system, code, submission_time, contest_id, contest_task_id, language, sverdict_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `,
		submission.SenderLogin, submission.TaskID, submission.TestingSystem, submission.Code, submission.SubmissionTime, submission.ContestID, submission.ContestTaskID, submission.Language, submission.SVerdictID)
	if err != nil {
		log.Printf("Failed to insert submission: %v", err)
	}
}

func GetVerdicts(username string, taskId string, testingSystem string) []SubmissionVerdict {
	// TODO возвращать еще и ошибку
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), `
        SELECT id, sender_login, task_id, testing_system, code, submission_time, contest_id, contest_task_id, verdict, language, execution_time, memory_used, test, submission_number
        FROM submissions_verdicts
        WHERE sender_login = $1 AND task_id = $2 AND testing_system = $3
		ORDER BY id DESC
    `, username, taskId, testingSystem)
	if err != nil {
		log.Printf("Failed to query submissions verdicts: %v", err)
	}
	defer rows.Close()

	var verdicts []SubmissionVerdict
	for rows.Next() {
		var verdict SubmissionVerdict
		var submissionTime time.Time
		err := rows.Scan(&verdict.ID, &verdict.SenderLogin, &verdict.TaskID, &verdict.TestingSystem, &verdict.Code, &submissionTime, &verdict.ContestID, &verdict.ContestTaskID, &verdict.Verdict, &verdict.Language, &verdict.ExecutionTime, &verdict.MemoryUsed, &verdict.Test, &verdict.SubmissionNumber)
		if err != nil {
			log.Printf("Failed to fetch submission verdict: %v", err)
		}
		verdict.SubmissionTime = submissionTime.Format("2006-01-02 15:04:05")
		verdicts = append(verdicts, verdict)
	}

	return verdicts
}

func GetContests() []Contest {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	var (
		contests []Contest
	)

	rows, err := conn.Query(context.Background(), "SELECT * FROM contests")
	if err != nil {
		fmt.Println("Не удалось выполнить запрос:", err)
		return []Contest{}
	}
	defer rows.Close()

	// Перебор результатов запроса
	for rows.Next() {
		var (
			id           int
			title        string
			access       map[string]interface{}
			participants map[string]interface{}
			results      map[string]interface{}
			tasks        map[string]interface{}
		)

		err := rows.Scan(&id, &title, &access, &participants, &results, &tasks)
		if err != nil {
			fmt.Println("Не удалось получить данные строки:", err)
			return []Contest{}
		}
		contests = append(contests, Contest{id, title, access, participants,
			results, tasks})
	}

	return contests
}

func GetContestInfo(id int) Contest {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	var contest Contest
	err = conn.QueryRow(context.Background(), "SELECT * FROM contests WHERE id = $1", id).Scan(
		&contest.ID,
		&contest.Title,
		&contest.Access,
		&contest.Participants,
		&contest.Results,
		&contest.Tasks,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Println("Запись с указанным ID не найдена")
			return Contest{}
		}
		fmt.Println("Ошибка при выполнении запроса:", err)
		return Contest{}
	}

	return contest
}

func TaskExist(taskID string) (bool, Task, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	// Выполнение запроса на проверку наличия задачи
	var task Task
	err = conn.QueryRow(context.Background(), "SELECT * FROM tasks WHERE id = $1", taskID).Scan(
		&task.ID,
		&task.Name,
		&task.Constraints,
		&task.Condition,
		&task.InputData,
		&task.OutputData,
		&task.Source,
		&task.AdditionalInfo,
		&task.Tests,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, Task{}, nil
		}
		return false, Task{}, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}

	return true, task, nil
}

func AddProblem(task Task) error {
	// Установка соединения с базой данных PostgreSQL
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}
	defer conn.Close(context.Background())

	// Выполнение SQL-запроса на добавление задания
	_, err = conn.Exec(context.Background(), `
		INSERT INTO tasks (id, name, constraints, condition, input_data, output_data, source, additional_info, tests)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		task.ID, task.Name, task.Constraints, task.Condition, task.InputData, task.OutputData, task.Source, task.AdditionalInfo, task.Tests)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении SQL-запроса: %w", err)
	}

	return nil
}

func GetTaskNameByID(taskID string) (string, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return "", err
	}
	defer conn.Close(context.Background())

	query := "SELECT name FROM tasks WHERE id = $1"
	var name string
	err = conn.QueryRow(context.Background(), query, taskID).Scan(&name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("задача с ID '%s' не найдена", taskID)
		}
		return "", fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}

	return name, nil
}

func GetVerditctsOfContestTask(login string, contestID, contestTaskID int) ([]SubmissionVerdict, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)

	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, sender_login, task_id, testing_system, code, submission_time, contest_id,
		contest_task_id, verdict, language, execution_time, memory_used, test, submission_number
		FROM submissions_verdicts
		WHERE sender_login = $1 AND contest_id = $2 AND contest_task_id = $3
		ORDER BY id DESC
	`

	rows, err := conn.Query(context.Background(), query, login, contestID, contestTaskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verdicts []SubmissionVerdict
	for rows.Next() {
		var verdict SubmissionVerdict
		var submissionTime time.Time
		err := rows.Scan(&verdict.ID, &verdict.SenderLogin, &verdict.TaskID, &verdict.TestingSystem, &verdict.Code, &submissionTime, &verdict.ContestID, &verdict.ContestTaskID, &verdict.Verdict, &verdict.Language, &verdict.ExecutionTime, &verdict.MemoryUsed, &verdict.Test, &verdict.SubmissionNumber)
		if err != nil {
			log.Printf("Failed to fetch submission verdict: %v", err)
		}
		verdict.SubmissionTime = submissionTime.Format("2006-01-02 15:04:05")
		verdicts = append(verdicts, verdict)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return verdicts, nil
}
