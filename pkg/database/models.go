package database

import (
	"encoding/json"
	"time"
)

type Submission struct {
	ID               int    `json:"id"`
	SenderLogin      string `json:"sender_login"`
	TaskID           string `json:"task_id"`
	TestingSystem    string `json:"testing_system"`
	Code             []byte `json:"code"`
	SubmissionTime   string `json:"submission_time"`
	ContestID        int    `json:"contest_id"`
	ContestTaskID    int    `json:"contest_task_id"`
	Verdict          string `json:"verdict"`
	Language         string `json:"language"`
	ExecutionTime    string `json:"execution_time"`
	MemoryUsed       string `json:"memory_used"`
	Test             string `json:"test"`
	SubmissionNumber string `json:"submission_number"`
	Status           int    `json:"status"` // Статус проверки 0 - не отправлено, 1 - проверяется, 2 - проверка окончена
}

type Contest struct {
	ID           int                    `json:"id"`
	Title        string                 `json:"title"`
	Access       map[string]interface{} `json:"access"`
	Participants map[string]interface{} `json:"participants"`
	Results      map[string]interface{} `json:"results"`
	Tasks        map[string]interface{} `json:"tasks"`
	GroupOwner   int                    `json:"group_owner"`
	StartTime    time.Time              `json:"start_time"`
	Duration     time.Duration          `json:"duration"`
}

type Task struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Constraints    map[string]string      `json:"constraints"`
	Condition      string                 `json:"condition"`
	InputData      string                 `json:"input_data"`
	OutputData     string                 `json:"output_data"`
	Source         string                 `json:"source"`
	AdditionalInfo string                 `json:"additional_info"`
	Tests          map[string]interface{} `json:"tests"`
}

type User struct {
	Username     string            `json:"username"`
	PasswordHash string            `json:"password_hash"`
	LastName     string            `json:"last_name"`
	FirstName    string            `json:"first_name"`
	Email        string            `json:"email"`
	Group        string            `json:"group_name"`
	Role         string            `json:"role"`
	SolvedTasks  []string          `json:"solved_tasks"`
	Groups       []json.RawMessage `json:"groups"`
	Tasks        []json.RawMessage `json:"tasks"`
}

type Group struct {
	ID         int      `json:"id"`
	Title      string   `json:"title"`
	Students   []string `json:"students"`
	Teachers   []string `json:"teachers"`
	Admins     []string `json:"admins"`
	InviteCode []byte   `json:"invite_code"`
}
