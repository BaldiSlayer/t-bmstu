package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

var (
	DbURL = ""
)

func CreateTables(dbUsername, dbPassword, dbHost, dbName string) error {
	// TODO миграции

	DbURL = fmt.Sprintf(
		"postgresql://%s:%s@%s:5432/%s",
		dbUsername,
		dbPassword,
		dbHost,
		dbName,
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
            group_owner INTEGER,
            start_time TIMESTAMP,
    		duration INTERVAL
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

		/*CREATE TABLE IF NOT EXISTS forum_posts (
		    id SERIAL PRIMARY KEY,
		    title TEXT,
		    text TEXT,
		    tags TEXT[]
		)*/
    `)
	if err != nil {
		return err
	}

	fmt.Println("Tables created successfully!")
	return nil
}
