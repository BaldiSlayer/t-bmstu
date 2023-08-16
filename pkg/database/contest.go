package database

import (
	"context"
	"github.com/jackc/pgx/v4"
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
	)
	if err != nil {
		return Contest{}, err
	}

	return contest, nil
}
