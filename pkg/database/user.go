package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
)

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

type Profile struct {
	Username  string
	LastName  string
	FirstName string
	Email     string
}

func GetInfoForProfilePage(username string) (User, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return User{}, fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	var user User
	err = conn.QueryRow(context.Background(), "SELECT username, last_name, first_name, email FROM users WHERE username = $1", username).Scan(&user.Username, &user.LastName, &user.FirstName, &user.Email)
	if err != nil {
		return User{}, fmt.Errorf("query execution failed: %v", err)
	}

	return user, nil
}
