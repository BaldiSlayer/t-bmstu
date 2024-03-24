package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
)

func CheckGroupExist(title string, inviteCode string) (bool, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return false, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer conn.Close(context.Background())

	var count int
	err = conn.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM groups WHERE title = $1 OR invite_code = $2", title, inviteCode).
		Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
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
	err = tx.QueryRow(context.Background(), "INSERT INTO groups (title, students, teachers, admins, invite_code) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		group.Title, group.Students, group.Teachers, group.Admins, group.InviteCode).Scan(&group.ID)
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
		SELECT id, title, access, participants, results, tasks, group_owner, start_time, duration
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
		if err := rows.Scan(&contest.ID, &contest.Title, &contest.Access, &contest.Participants, &contest.Results, &contest.Tasks, &contest.GroupOwner,
			&contest.StartTime, &contest.Duration); err != nil {
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
