package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/viper"
	"log"
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

	// пересмотреть `contest_task_id`, нужно ли его хранить?
	// если нужно, то нужно ли хранить в таблице contests JSONB tasks?
	// скорее всего нет

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
            verdict VARCHAR(255),
            execution_time VARCHAR(255),
            memory_used VARCHAR(255),
            test VARCHAR(255),
            language VARCHAR(255),
            submission_number VARCHAR(255),
            status INTEGER                                                    
        );

        CREATE TABLE IF NOT EXISTS contests (
            id SERIAL PRIMARY KEY,
            title VARCHAR(255),
            access JSONB,
            participants JSONB,
            results JSONB,
            tasks JSONB,
            group_owner INTEGER
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
		   groups        JSONB,
		   tasks 		 JSONB
		  );

		CREATE TABLE IF NOT EXISTS groups (
		    id SERIAL PRIMARY KEY,
		    title TEXT,
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
			group_owner  int
		)

		err := rows.Scan(&id, &title, &access, &participants, &results, &tasks, &group_owner)
		if err != nil {
			fmt.Println("Не удалось получить данные строки:", err)
			return []Contest{}
		}
		contests = append(contests, Contest{id, title, access, participants,
			results, tasks, group_owner})
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
		&contest.GroupOwner,
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

func GetVerditctsOfContestTask(login string, contestID, contestTaskID int) ([]Submission, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)

	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, sender_login, task_id, testing_system, code, submission_time, contest_id,
		contest_task_id, verdict, language, execution_time, memory_used, test, submission_number, status
		FROM submissions
		WHERE sender_login = $1 AND contest_id = $2 AND contest_task_id = $3
		ORDER BY id DESC
	`

	rows, err := conn.Query(context.Background(), query, login, contestID, contestTaskID)
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
	err = tx.QueryRow(context.Background(), "INSERT INTO groups (title, students, teachers, admins) VALUES ($1, $2, $3, $4) RETURNING id",
		group.Title, group.Students, group.Teachers, group.Admins).Scan(&group.ID)
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
		return nil, err
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), `
		SELECT groups
		FROM users
		WHERE username = $1
	`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type GroupReference struct {
		Role    string `json:"role"`
		GroupID int    `json:"group_id"`
	}

	var groupReferences []GroupReference
	if rows.Next() {
		if err := rows.Scan(&groupReferences); err != nil {
			return nil, err
		}
	}

	var groups []Group
	for _, ref := range groupReferences {
		group, err := GetGroupByID(ref.GroupID)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func GetGroupByID(groupID int) (Group, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return Group{}, err
	}
	defer conn.Close(context.Background())

	var group Group
	conn.QueryRow(context.Background(), `
		SELECT id, title, students, teachers, admins, invite_code
		FROM groups
		WHERE id = $1
	`, groupID).Scan(&group.ID, &group.Title, &group.Students, &group.Teachers, &group.Admins, &group.InviteCode)
	if err != nil {
		return Group{}, err
	}
	return group, nil
}

func GetGroupContests(groupOwner int) ([]Contest, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), `
		SELECT id, title, access, participants, results, tasks, group_owner
		FROM contests
		WHERE group_owner = $1
	`, groupOwner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contests []Contest
	for rows.Next() {
		var contest Contest
		if err := rows.Scan(&contest.ID, &contest.Title, &contest.Access, &contest.Participants, &contest.Results, &contest.Tasks, &contest.GroupOwner); err != nil {
			return nil, err
		}
		contests = append(contests, contest)
	}

	return contests, nil
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

func GetSubmitsWithStatus(testsystem string, status int) ([]Submission, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)

	if err != nil {
		return nil, err
	}

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
