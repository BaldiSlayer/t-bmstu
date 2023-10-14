package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"time"
)

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
		&contest.StartTime,
		&contest.Duration,
	)
	if err != nil {
		return Contest{}, err
	}

	return contest, nil
}

func CreateContest(title string, access map[string]interface{}, groupOwner int, startTime time.Time, duration time.Duration) (int, error) {
	conn, err := pgx.Connect(context.Background(), DbURL)
	if err != nil {
		return 0, err
	}
	defer conn.Close(context.Background())

	var contestID int
	err = conn.QueryRow(
		context.Background(),
		"INSERT INTO contests (title, access, group_owner, start_time, duration) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		title,
		access,
		groupOwner,
		startTime,
		duration,
	).Scan(&contestID)
	if err != nil {
		return 0, err
	}

	return contestID, nil
}
