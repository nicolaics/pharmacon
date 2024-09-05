package user

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/nicolaics/pos_pharmacy/service/auth"
	"github.com/nicolaics/pos_pharmacy/types"
)

type Store struct {
	db          *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db,}
}

func (s *Store) GetUserByName(name string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM user WHERE name = ? ", name)

	if err != nil {
		return nil, err
	}

	user := new(types.User)

	for rows.Next() {
		user, err = scanRowIntoUser(rows)

		if err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("customer not found")
	}

	return user, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM user WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	user := new(types.User)

	for rows.Next() {
		user, err = scanRowIntoUser(rows)

		if err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO user (name, password, admin, phone_number) VALUES (?, ?, ?, ?)",
		user.Name, user.Password, user.Admin, user.PhoneNumber)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteUser(user *types.User) error {
	_, err := s.db.Exec("DELETE FROM user WHERE id = ?", user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetAllUsers() ([]types.User, error) {
	rows, err := s.db.Query("SELECT * FROM user")

	if err != nil {
		return nil, err
	}

	users := make([]types.User, 0)

	for rows.Next() {
		user, err := scanRowIntoUser(rows)

		if err != nil {
			return nil, err
		}

		users = append(users, *user)
	}

	return users, nil
}

func (s *Store) UpdateLastLoggedIn(id int) error {
	_, err := s.db.Exec("UPDATE user SET last_logged_in = ? WHERE id = ? ",
		time.Now(), id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) ModifyUser(id int, user types.User) error {
	columns := "name = ?, password = ?, admin = ?, phone_number = ?"

	_, err := s.db.Exec(fmt.Sprintf("UPDATE user SET %s WHERE id = ? ", columns),
		user.Name, user.Password, user.Admin, user.PhoneNumber, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) SaveToken(userId int, tokenDetails *types.TokenDetails) error {
	tokenExp := time.Unix(tokenDetails.TokenExp, 0) //converting Unix to UTC(to Time object)

	query := "INSERT INTO verify_token(user_id, uuid, expired_at) VALUES (?, ?, ?)"
	_, err := s.db.Exec(query, userId, tokenDetails.UUID, tokenExp)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteToken(givenUuid string) error {
	query := "DELETE FROM verify_token WHERE uuid = ?"
	_, err := s.db.Exec(query, givenUuid)
	if err != nil {
		return err
	}

	return nil
}

// TODO: change into from sql, not redis
// TODO: using SELECT COUNT(*) WHERE name = ? AND
// TODO: token = ? AND TIMESTAMPDIFF(HOUR, created_at, NOW()) <= ?
// TODO: check using the count, if count 0, means not valid
// TODO: delete the token
func (s *Store) ValidateUserToken(w http.ResponseWriter, r *http.Request, needAdmin bool) (*types.User, error) {
	accessDetails, err := auth.ExtractTokenFromClient(r)
	if err != nil {
		return nil, err
	}

	query := "SELECT user_id FROM verify_token WHERE uuid = ? AND user_id = ? AND expired_at <= ?"
	rows, err := s.db.Query(query, accessDetails.UUID, accessDetails.UserID, time.Now())
	if err != nil {
		return nil, err
	}

	var userId int

	for rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			delErr := s.DeleteToken(accessDetails.UUID)
			if delErr != nil {
				return nil, fmt.Errorf("delete error: %v", delErr)
			}

			return nil, fmt.Errorf("token expired, log in again")
		}
	}

	// check if user exist
	user, err := s.GetUserByID(userId)
	if err != nil {
		return nil, err
	}

	// if the account must be admin
	if needAdmin {
		if !user.Admin {
			return nil, fmt.Errorf("unauthorized! not admin")
		}
	}

	return user, nil
}


func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.Admin,
		&user.PhoneNumber,
		&user.LastLoggedIn,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	user.LastLoggedIn = user.LastLoggedIn.Local()
	user.CreatedAt = user.CreatedAt.Local()

	return user, nil
}