package repository

type Submission struct {
	ID             int    `json:"id"`
	SenderLogin    string `json:"sender_login"`
	TaskID         string `json:"task_id"`
	TestingSystem  string `json:"testing_system"`
	Code           string `json:"code"`
	SubmissionTime string `json:"submission_time"`
	ContestID      int    `json:"contest_id"`
	ContestTaskID  int    `json:"contest_task_id"`
	Language       string `json:"language"`
	SVerdictID     string `json:"sverdict_id"`
}

type SubmissionVerdict struct {
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
}

type UnverifiedSubmission struct {
	SubmissionID         string `json:"submission_id"`
	ExternalSubmissionID string `json:"external_submission_id"`
	TestingSystem        string `json:"testing_system"`
	JudgeId              string `json:"judge_id"`
}

type Contest struct {
	ID           int                    `json:"id"`
	Title        string                 `json:"title"`
	Access       map[string]interface{} `json:"access"`
	Participants map[string]interface{} `json:"participants"`
	Results      map[string]interface{} `json:"results"`
	Tasks        map[string]interface{} `json:"tasks"`
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
