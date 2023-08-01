package repository

import (
	"context"
	"encoding/json"
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
		
		CREATE TABLE IF NOT EXISTS users (
		   username      TEXT PRIMARY KEY,
		   password_hash TEXT,
		   last_name     TEXT,
		   first_name    TEXT,
		   email         TEXT,
		   group_name    TEXT,
		   role          TEXT,
		   solved_tasks  TEXT[],
		   groups        JSONB	
		  );

		CREATE TABLE IF NOT EXISTS groups (
		    id SERIAL PRIMARY KEY,
		    title TEXT,
		    contests INTEGER[],
		    students TEXT[],
		    teachers TEXT[],
		    admins TEXT[],
		    invite_code TEXT
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

	rows, err := conn.Query(context.Background(), "SELECT * FROM contests ORDER BY id DESC")
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

func GetContestInfoById(id int) (Contest, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return Contest{}, err
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
		return Contest{}, err
	}

	return contest, nil
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

func CreateUser(user User) error {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	query := `
  INSERT INTO users (username, password_hash, last_name, first_name, email, group_name, role, solved_tasks, groups)
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
 `

	_, err = conn.Exec(context.Background(), query, user.Username, user.PasswordHash, user.LastName, user.FirstName, user.Email, user.Group, user.Role, user.SolvedTasks, user.Groups)
	if err != nil {
		return err
	}

	return nil
}

func CheckIfUserExists(login string) (bool, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return false, err
	}
	defer conn.Close(context.Background())

	var count int

	query := `
  SELECT COUNT(*) FROM users WHERE username = $1;
 `

	err = conn.QueryRow(context.Background(), query, login).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func AuthenticateUser(login, password string) (bool, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return false, err
	}
	defer conn.Close(context.Background())

	var count int

	query := `
  SELECT COUNT(*) FROM users WHERE username = $1 AND password_hash = $2;
 `

	err = conn.QueryRow(context.Background(), query, login, password).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func GetUserRole(username string) (string, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return "", err
	}
	defer conn.Close(context.Background())

	var role string

	query := `
  SELECT role FROM users WHERE username = $1;
 `

	err = conn.QueryRow(context.Background(), query, username).Scan(&role)
	if err != nil {
		return "", err
	}

	return role, nil
}

func AddUserToGroup(username string, groupId int, role string) error {
	// TODO проверка был би там уже пользователь?
	// если да, то можно его просто редиректнуть в эту группу

	type GroupInfo struct {
		GroupID int    `json:"group_id"`
		Role    string `json:"role"`
	}

	groupInfo := GroupInfo{
		GroupID: groupId,
		Role:    role,
	}

	groupJSON, err := json.Marshal(groupInfo)
	if err != nil {
		return err
	}

	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer conn.Close(context.Background())

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	_, err = conn.Exec(context.Background(), "UPDATE users SET groups = groups || $1 WHERE username = $2", groupJSON, username)

	if err != nil {
		tx.Rollback(context.Background())
		fmt.Println(err)
		return fmt.Errorf("failed to update user group fields: %w", err)
	}

	switch role {
	case "student":
		_, err = tx.Exec(context.Background(), "UPDATE groups SET students = array_append(students, $1) WHERE id = $2",
			username, groupId)
	case "teacher":
		_, err = tx.Exec(context.Background(), "UPDATE groups SET teachers = array_append(teachers, $1) WHERE id = $2",
			username, groupId)
	case "admin":
		_, err = tx.Exec(context.Background(), "UPDATE groups SET admins = array_append(admins, $1) WHERE id = $2",
			username, groupId)
	default:
		return fmt.Errorf("unknown role: %s", role)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func AddGroupWithMembers(group Group, memberUsernames []json.RawMessage) error {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer conn.Close(context.Background())

	// Начало транзакции
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Вставка новой группы в таблицу groups
	err = tx.QueryRow(context.Background(), "INSERT INTO groups (title, contests, students, teachers, admins) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		group.Title, group.Contests, group.Students, group.Teachers, group.Admins).Scan(&group.ID)
	if err != nil {
		tx.Rollback(context.Background())
		return fmt.Errorf("failed to insert new group: %w", err)
	}

	// Обновление полей group для каждого пользователя из списка никнеймов

	usernameRoleData := struct {
		Username string `json:"username"`
		Role     string `json:"role"`
	}{}

	for _, usernameData := range memberUsernames {
		err = json.Unmarshal(usernameData, &usernameRoleData)
		if err != nil {
			tx.Rollback(context.Background())
			return fmt.Errorf("failed to unmarshal username json: %w", err)
		}

		username := usernameRoleData.Username
		role := usernameRoleData.Role

		err = AddUserToGroup(username, group.ID, role)
		if err != nil {
			return err
		}

		if err != nil {
			tx.Rollback(context.Background())
			return fmt.Errorf("failed to add member %s to group: %w", username, err)
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func GetUserGroups(username string) ([]Group, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %w", err)
	}
	defer conn.Close(context.Background())

	// Запрос для получения всех групп пользователя
	rows, err := conn.Query(context.Background(), "SELECT id, title, contests, students, teachers, admins FROM groups WHERE $1 = ANY (students)", username)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer rows.Close()

	var groups []Group

	// Итерация по результатам запроса
	for rows.Next() {
		var group Group
		err := rows.Scan(&group.ID, &group.Title, &group.Contests, &group.Students, &group.Teachers, &group.Admins)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении результатов запроса: %w", err)
		}

		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов запроса: %w", err)
	}

	return groups, nil
}

func GetGroupContests(groupId int) ([]Contest, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT contests FROM groups WHERE id=$1", groupId)
	if err != nil {
		return nil, fmt.Errorf("error executing the query: %w", err)
	}
	defer rows.Close()

	var contests []int
	if rows.Next() {
		err := rows.Scan(&contests)
		if err != nil {
			return nil, fmt.Errorf("error scanning results: %w", err)
		}

		// Convert the int array to Contest struct array
		var contestList []Contest
		for _, contestID := range contests {
			contest, err := GetContestInfoById(contestID)
			if err != nil {
				return nil, err
			}
			contestList = append(contestList, contest)
		}

		return contestList, nil
	}

	return nil, fmt.Errorf("group with ID %d not found", groupId)
}

func CheckInviteCode(inviteCode string) (bool, int, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return false, 0, fmt.Errorf("failed to connect to the database: %w", err)
	}
	defer conn.Close(context.Background())

	var groupID int
	err = conn.QueryRow(context.Background(), "SELECT id FROM groups WHERE invite_code = $1", inviteCode).Scan(&groupID)
	if err == pgx.ErrNoRows {
		return false, 0, nil // No group found with the given invite_code
	} else if err != nil {
		return false, 0, fmt.Errorf("error while querying the database: %w", err)
	}

	return true, groupID, nil
}
