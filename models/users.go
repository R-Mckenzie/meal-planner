package models

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id        int
	email     string
	password  string // DOES NOT GET SAVED TO DB
	passHash  string
	createdAt time.Time
	updatedAt time.Time
	deletedAt sql.NullTime
}

type UserService struct {
	db *sql.DB
	errLog     *log.Logger
	infLog     *log.Logger
}

func NewUserService(db *sql.DB, errLog, infLog *log.Logger) *UserService {
	return &UserService{
		db: db,
		errLog: errLog,
		infLog: infLog,
	}
}

func (us *UserService) Create(email, password string) error {
	hash , err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = us.db.Exec("INSERT INTO users (email, hash, created_at, updated_at) VALUES($1, $2, $3, $3)",
		email, hash, time.Now().UTC().Format(time.RFC3339))
	var pgErr *pq.Error
	if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Error code for constraint violation. email must be unique
		return errors.New("A user already exists with this email")
	} 
	if err != nil {
		us.errLog.Println("Problem creating new user", email, err)
		return errors.New("There was a problem signing you u")
	}

	us.infLog.Print("Successfully added user: ", email)
	return nil
}

func (us *UserService) GetById(id int) (*User, error) {
	row := us.db.QueryRow("SELECT * FROM users WHERE id=$1", id)
	u := &User{password: ""}
	err := row.Scan(&u.id, &u.email, &u.passHash, &u.createdAt, &u.updatedAt, &u.deletedAt)
	if err != nil {
		us.errLog.Println("Error getting user by ID: ", err)
		return nil, err
	}
	return u, nil
}
