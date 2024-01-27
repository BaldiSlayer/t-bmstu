package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"time"
)

func AddSubmission(submission Submission) (int, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	verdict := Submission{
		SenderLogin:      submission.SenderLogin,
		TaskID:           submission.TaskID,
		TestingSystem:    submission.TestingSystem,
		Code:             submission.Code,
		SubmissionTime:   submission.SubmissionTime,
		ContestID:        submission.ContestID,
		ContestTaskID:    submission.ContestTaskID,
		Language:         submission.Language,
		Verdict:          "Waiting",
		ExecutionTime:    "-",
		MemoryUsed:       "-",
		Test:             "-",
		SubmissionNumber: "-",
		Status:           0,
	}

	var id int
	err = conn.QueryRow(context.Background(), `
		INSERT INTO submissions (sender_login, task_id, testing_system, code, submission_time, contest_id, contest_task_id, language, verdict, execution_time, memory_used, test, submission_number, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id;
	`,
		verdict.SenderLogin, verdict.TaskID, verdict.TestingSystem, verdict.Code, verdict.SubmissionTime, verdict.ContestID, verdict.ContestTaskID, verdict.Language, verdict.Verdict, verdict.ExecutionTime, verdict.MemoryUsed, verdict.Test, verdict.SubmissionNumber, verdict.Status).Scan(&id)

	if err != nil {
		log.Printf("Failed to insert submission verdict: %v", err)
		return -1, err
	}

	return id, nil
}

func GetVerdicts(username string, taskId string, testingSystem string, contestId int, contestTaskId int) []Submission {
	// TODO возвращать еще и ошибку
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), `
        SELECT id, sender_login, task_id, testing_system, code, submission_time, contest_id, contest_task_id, verdict, language, execution_time, memory_used, test, submission_number, status
        FROM submissions
        WHERE sender_login = $1 AND task_id = $2 AND testing_system = $3 AND contest_id = $4 AND contest_task_id = $5
		ORDER BY id DESC
    `, username, taskId, testingSystem, contestId, contestTaskId)
	if err != nil {
		log.Printf("Failed to query submissions verdicts: %v", err)
	}
	defer rows.Close()

	var verdicts []Submission
	for rows.Next() {
		var verdict Submission
		var submissionTime time.Time
		err := rows.Scan(&verdict.ID, &verdict.SenderLogin, &verdict.TaskID, &verdict.TestingSystem, &verdict.Code, &submissionTime, &verdict.ContestID, &verdict.ContestTaskID, &verdict.Verdict, &verdict.Language, &verdict.ExecutionTime, &verdict.MemoryUsed, &verdict.Test, &verdict.SubmissionNumber, &verdict.Status)
		if err != nil {
			log.Printf("Failed to fetch submission verdict: %v", err)
		}
		verdict.SubmissionTime = submissionTime.Format("2006-01-02 15:04:05")
		verdicts = append(verdicts, verdict)
	}

	return verdicts
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

func GetSubmitsWithStatus(testsystem string, status int) ([]Submission, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)

	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	query := `
		SELECT id, sender_login, task_id, testing_system, code, submission_time, contest_id,
		contest_task_id, verdict, language, execution_time, memory_used, test, submission_number, status
		FROM submissions
		WHERE status = $1 AND testing_system = $2
		ORDER BY id ASC
	`

	rows, err := conn.Query(context.Background(), query, status, testsystem)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verdicts []Submission
	for rows.Next() {
		var verdict Submission
		var submissionTime time.Time
		err := rows.Scan(&verdict.ID, &verdict.SenderLogin, &verdict.TaskID, &verdict.TestingSystem, &verdict.Code, &submissionTime, &verdict.ContestID, &verdict.ContestTaskID, &verdict.Verdict, &verdict.Language, &verdict.ExecutionTime, &verdict.MemoryUsed, &verdict.Test, &verdict.SubmissionNumber, &verdict.Status)
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

func UpdateSubmissionData(submission Submission) error {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	query := `
        UPDATE submissions
        SET
            sender_login = $1,
            task_id = $2,
            testing_system = $3,
            code = $4,
            submission_time = $5,
            contest_id = $6,
            contest_task_id = $7,
            verdict = $8,
            language = $9,
            execution_time = $10,
            memory_used = $11,
            test = $12,
            submission_number = $13,
            status = $14
        WHERE id = $15
    `

	_, err = conn.Exec(
		context.Background(),
		query,
		submission.SenderLogin,
		submission.TaskID,
		submission.TestingSystem,
		submission.Code,
		submission.SubmissionTime,
		submission.ContestID,
		submission.ContestTaskID,
		submission.Verdict,
		submission.Language,
		submission.ExecutionTime,
		submission.MemoryUsed,
		submission.Test,
		submission.SubmissionNumber,
		submission.Status,
		submission.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func GetSubmissionCode(id int) (string, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return "", err
	}
	defer conn.Close(context.Background())

	var submission Submission
	err = conn.QueryRow(context.Background(), "SELECT Code FROM submissions WHERE id = $1", id).Scan(&submission.Code)
	if err != nil {
		return "", fmt.Errorf("query execution failed: %v", err)
	}

	return string(submission.Code), nil
}
