package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

// Errors
var ErrUserExists = errors.New("A user already exists with this email")

type User struct {
	ID       int
	Email    string
	PassHash []byte
	Role     int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type postgresUserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) UserStore {
	return &postgresUserStore{db}
}

type UserStore interface {
	Create(email, password, provider string) error
	GetByEmail(email string) (User, error) // Before authentication, we can find the user with the email provided at login
	GetByID(id int) (User, error)          // Once the user is authenticated, we can use the ID for future requests
}

func (us *postgresUserStore) Create(email, password, provider string) error {
	// Generate hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Println(err)
	}

	// Create User struct
	user := &User{
		Email:     email,
		PassHash:  hash,
		Role:      0,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = us.db.QueryRow("INSERT INTO users (email, hashed_pass, created, updated) VALUES($1, $2, $3, $4) RETURNING id",
		user.Email, user.PassHash, user.CreatedAt.Format(time.RFC3339), user.UpdatedAt.Format(time.RFC3339)).Scan(&user.ID)

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Error code for constraint violation. email must be unique
		log.Println(err)
		return ErrUserExists
	}
	if err != nil {
		log.Println(err)
		return errors.New("There was a problem signing you up")
	}

	fmt.Println(user)
	return nil
}

func (us *postgresUserStore) GetByEmail(email string) (User, error) {
	row := us.db.QueryRow("SELECT id, email, hashed_pass, role, created, updated FROM users where email=$1", email)
	var u = &User{}
	err := row.Scan(&u.ID, &u.Email, &u.PassHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return *u, fmt.Errorf("in GetByEmail: %w", err)
	}
	return *u, nil
}

func (us *postgresUserStore) GetByID(id int) (User, error) {
	row := us.db.QueryRow("SELECT id, email, hashed_pass, role, created, updated FROM users where id=$1", id)
	var u = &User{}
	err := row.Scan(&u.ID, &u.Email, &u.PassHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return *u, fmt.Errorf("in GetByID: %w", err)
	}
	return *u, nil
}
