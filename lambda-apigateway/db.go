package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	Db *sql.DB
}

type User struct {
	ID                  int    `json:"id"`
	Username            string `json:"username"`
	Email               string `json:"email"`
	PasswordFromPayload string `json:"password"`
	PasswordHash        string `json:"-"`
}

func NewMysql(host, user, password, port, database string) (*Mysql, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.New("error connecting to the database")
	}

	return &Mysql{
		Db: db,
	}, nil
}

func (mysql *Mysql) GetUsers() ([]User, error) {
	rows, err := mysql.Db.Query("SELECT id, username, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (mysql *Mysql) CreateUser(user User) error {
	_, err := mysql.Db.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", user.Username, user.Email, user.PasswordFromPayload)
	if err != nil {
		return err
	}

	return nil
}

func (mysql *Mysql) GetUserByID(id int) (User, error) {
	var user User
	err := mysql.Db.QueryRow("SELECT id, username, email FROM users WHERE id = ?", id).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (mysql *Mysql) UpdateUser(user User) error {
	res, err := mysql.Db.Exec("UPDATE users SET username = ?, email = ?, password_hash = ? WHERE id = ?", user.Username, user.Email, user.PasswordFromPayload, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows were affected")
	}

	return nil
}
