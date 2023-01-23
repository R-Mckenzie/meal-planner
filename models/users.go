package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/R-Mckenzie/meal-planner/tokens"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const hmacSecret = "hmac-secret-key"

type User struct {
	ID           int
	Email        string
	passHash     string
	Remember     string
	rememberHash string
	createdAt    time.Time
	updatedAt    time.Time
	deletedAt    sql.NullTime
}

type UserService struct {
	db     *sql.DB
	hmac   tokens.HMAC
	errLog *log.Logger
	infLog *log.Logger
}

func NewUserService(db *sql.DB, errLog, infLog *log.Logger) *UserService {
	return &UserService{
		db:     db,
		hmac:   tokens.NewHMAC(hmacSecret),
		errLog: errLog,
		infLog: infLog,
	}
}

func (us *UserService) Create(email, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:        email,
		passHash:     string(hash),
		Remember:     "",
		rememberHash: "",
		createdAt:    time.Now().UTC(),
		updatedAt:    time.Now().UTC(),
	}

	err = us.db.QueryRow("INSERT INTO users (email, hash, remember_hash, created_at, updated_at) VALUES($1, $2, $3, $4, $5) RETURNING id",
		user.Email, user.passHash, user.rememberHash, user.createdAt.Format(time.RFC3339), user.updatedAt.Format(time.RFC3339)).Scan(&user.ID)
	var pgErr *pq.Error
	if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Error code for constraint violation. email must be unique
		return nil, errors.New("A user already exists with this email")
	}
	if err != nil {
		us.errLog.Println("Problem creating new user", email, err)
		return nil, errors.New("There was a problem signing you u")
	}

	us.infLog.Print("Successfully added user: ", email)
	return user, nil
}

func (us *UserService) Update(user *User) error {
	_, err := us.db.Exec("UPDATE users SET email=$2, hash=$3, remember_hash=$4, updated_at=$5 WHERE id=$1",
		user.ID, user.Email, user.passHash, user.rememberHash, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	user, err := us.GetByEmail(email)
	if err != nil {
		us.errLog.Printf("Problem finding user with email %q. %v\n", email, err)
		return nil, errors.New("No user found with email " + email)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.passHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, errors.New("The password you have entered is incorrect")
		}
		return nil, errors.New("There was a problem logging you in")
	}

	return user, nil
}

func (us *UserService) GenerateRemember(user *User) error {
	token, err := tokens.RememberToken()
	if err != nil {
		return err
	}
	rememberHash := us.hmac.Hash(token)
	fmt.Println(rememberHash)
	user.Remember = token
	user.rememberHash = rememberHash
	err = us.Update(user)
	if err != nil {
		us.errLog.Printf("Problem updating user %q. %v", user.Email, err)
	}

	return nil
}

func (us *UserService) GetById(id int) (*User, error) {
	row := us.db.QueryRow("SELECT * FROM users WHERE id=$1", id)
	u := &User{}
	err := row.Scan(&u.ID, &u.Email, &u.passHash, &u.createdAt, &u.updatedAt, &u.deletedAt, &u.rememberHash)
	if err != nil {
		us.errLog.Printf("Error getting user by ID (%d): %v", id, err)
		return nil, err
	}
	return u, nil
}

func (us *UserService) GetByEmail(email string) (*User, error) {
	row := us.db.QueryRow("SELECT * FROM users WHERE email=$1", email)
	u := &User{}
	err := row.Scan(&u.ID, &u.Email, &u.passHash, &u.createdAt, &u.updatedAt, &u.deletedAt, &u.rememberHash)
	if err != nil {
		us.errLog.Printf("Error getting user by email (%q): %v", email, err)
		return nil, err
	}
	return u, nil
}

func (us *UserService) GetByRemember(remember string) (*User, error) {
	rHash := us.hmac.Hash(remember)
	row := us.db.QueryRow("SELECT * FROM users WHERE remember_hash=$1", rHash)
	u := &User{}
	err := row.Scan(&u.ID, &u.Email, &u.passHash, &u.createdAt, &u.updatedAt, &u.deletedAt, &u.rememberHash)
	if err != nil {
		us.errLog.Printf("Error getting user by remember token (%q): %v", remember, err)
		return nil, err
	}
	return u, nil
}
