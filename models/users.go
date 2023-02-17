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

type UserService interface {
	Authenticate(email, password string) (*User, error)
	GenerateRemember(user *User) error
	userDB
}

type userDB interface {
	// Methods fot getting a user
	ByID(id int) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for editing a user
	Create(email, password string) error
	Update(user *User) error
	//Delete(id int) error
}

type userService struct {
	userDB
	hmac tokens.HMAC
	iLog *log.Logger
	eLog *log.Logger
}

type userPG struct {
	db   *sql.DB
	hmac tokens.HMAC
	iLog *log.Logger
	eLog *log.Logger
}

func newUserPG(db *sql.DB, hmac tokens.HMAC, iLog, eLog *log.Logger) (*userPG, error) {
	return &userPG{
		db:   db,
		hmac: hmac,
		iLog: iLog,
		eLog: eLog,
	}, nil
}

func NewUserService(db *sql.DB, iLog, eLog *log.Logger) (UserService, error) {
	hmac := tokens.NewHMAC(hmacSecret)
	pgdb, err := newUserPG(db, hmac, iLog, eLog)
	if err != nil {
		return nil, fmt.Errorf("in NewUserService: %w", err)
	}
	return &userService{
		userDB: pgdb,
		hmac:   hmac,
		iLog:   iLog,
		eLog:   eLog,
	}, nil
}

// === USER SERVICE
func (us *userService) Authenticate(email, password string) (*User, error) {
	user, err := us.userDB.ByEmail(email)
	if err != nil {
		us.eLog.Printf("tried to authenticate %q, but there is no stored user with this email\n", email)
		return nil, errors.New("No user found with email " + email)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.passHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			us.eLog.Printf("tried to authenticate %q, but the given password hash does not match\n", email)
			return nil, errors.New("The password you have entered is incorrect")
		}
		return nil, errors.New("There was a problem logging you in")
	}
	us.iLog.Printf("User %q authenticated\n", email)
	return user, nil
}

func (us *userService) GenerateRemember(user *User) error {
	token, err := tokens.RememberToken()
	if err != nil {
		return fmt.Errorf("in UserService.GenerateRemember: %w", err)
	}
	rememberHash := us.hmac.Hash(token)
	user.Remember = token
	user.rememberHash = rememberHash
	err = us.Update(user)
	if err != nil {
		return fmt.Errorf("in UserService.GenerateRemember: %w", err)
	}
	return nil
}

// === DB SERVICE
func (pg *userPG) Create(email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return fmt.Errorf("in UserPG.Create: %w", err)
	}

	token, err := tokens.RememberToken()
	if err != nil {
		return fmt.Errorf("in UserPG.Create: %w", err)
	}
	rememberHash := pg.hmac.Hash(token)

	user := &User{
		Email:        email,
		passHash:     string(hash),
		Remember:     token,
		rememberHash: rememberHash,
		createdAt:    time.Now().UTC(),
		updatedAt:    time.Now().UTC(),
	}

	err = pg.db.QueryRow("INSERT INTO users (email, hash, remember_hash, created_at, updated_at) VALUES($1, $2, $3, $4, $5) RETURNING id",
		user.Email, user.passHash, user.rememberHash, user.createdAt.Format(time.RFC3339), user.updatedAt.Format(time.RFC3339)).Scan(&user.ID)
	var pgErr *pq.Error
	if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Error code for constraint violation. email must be unique
		log.Println(err)
		return errors.New("A user already exists with this email")
	}
	if err != nil {
		log.Println(err)
		return errors.New("There was a problem signing you up")
	}

	pg.iLog.Printf("User %q created\n", email)
	return nil
}

func (pg *userPG) Update(user *User) error {
	_, err := pg.db.Exec("UPDATE users SET email=$2, hash=$3, remember_hash=$4, updated_at=$5 WHERE id=$1",
		user.ID, user.Email, user.passHash, user.rememberHash, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("in UserPG.Update: %w", err)
	}
	pg.iLog.Printf("User %q updated\n", user.Email)
	return nil
}

func (pg *userPG) ByID(id int) (*User, error) {
	row := pg.db.QueryRow("SELECT * FROM users WHERE id=$1", id)
	u := &User{}
	err := row.Scan(&u.ID, &u.Email, &u.passHash, &u.createdAt, &u.updatedAt, &u.deletedAt, &u.rememberHash)
	if err != nil {
		return nil, fmt.Errorf("in UserPG.ByID: %w", err)
	}
	return u, nil
}

func (pg *userPG) ByEmail(email string) (*User, error) {
	row := pg.db.QueryRow("SELECT * FROM users WHERE email=$1", email)
	u := &User{}
	err := row.Scan(&u.ID, &u.Email, &u.passHash, &u.createdAt, &u.updatedAt, &u.deletedAt, &u.rememberHash)
	if err != nil {
		return nil, fmt.Errorf("in UserPG.ByEmail: %w", err)
	}
	return u, nil
}

func (us *userService) ByRemember(remember string) (*User, error) {
	rHash := us.hmac.Hash(remember)
	return us.userDB.ByRemember(rHash)
}

func (pg *userPG) ByRemember(remember string) (*User, error) {
	row := pg.db.QueryRow("SELECT * FROM users WHERE remember_hash=$1", remember)
	u := &User{}
	err := row.Scan(&u.ID, &u.Email, &u.passHash, &u.createdAt, &u.updatedAt, &u.deletedAt, &u.rememberHash)
	if err != nil {
		return nil, fmt.Errorf("in UserPG.ByRemember: %w", err)
	}
	return u, nil
}
