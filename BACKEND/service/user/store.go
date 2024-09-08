package user

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nicolaics/pos_pharmacy/logger"
	"github.com/nicolaics/pos_pharmacy/service/auth"
	"github.com/nicolaics/pos_pharmacy/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByName(name string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM user WHERE name = ? ", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
	defer rows.Close()

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

func (s *Store) DeleteUser(user *types.User, deletedById int) error {
	data, err := s.GetUserByID(user.ID)
	if err != nil {
		return err
	}

	err = logger.WriteLog("delete", "user", deletedById, data.ID, data)
	if err != nil {
		return fmt.Errorf("error write log file")
	}

	_, err = s.db.Exec("DELETE FROM user WHERE id = ?", user.ID)
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
	defer rows.Close()

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

func (s *Store) ModifyUser(id int, user types.User, deletedById int) error {
	data, err := s.GetUserByID(user.ID)
	if err != nil {
		return err
	}

	writeData := map[string]interface{}{
		"previous_data": data,
	}

	err = logger.WriteLog("delete", "user", deletedById, data.ID, writeData)
	if err != nil {
		return fmt.Errorf("error write log file")
	}

	query := `UPDATE user SET name = ?, password = ?, admin = ?, phone_number = ? 
				WHERE id = ?`
	_, err = s.db.Exec(query,
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

func (s *Store) DeleteToken(userId int) error {
	query := "DELETE FROM verify_token WHERE user_id = ?"
	_, err := s.db.Exec(query, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) ValidateUserToken(w http.ResponseWriter, r *http.Request, needAdmin bool) (*types.User, error) {
	accessDetails, err := auth.ExtractTokenFromClient(r)
	if err != nil {
		return nil, err
	}

	log.Println(time.Now().UTC().Format("2006-01-02 15:04:05"))

	query := "SELECT user_id FROM verify_token WHERE uuid = ? AND user_id = ? AND expired_at >= ?"
	rows, err := s.db.Query(query, accessDetails.UUID, accessDetails.UserID, time.Now().UTC().Format("2006-01-02 15:04:05"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userId int

	for rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			return nil, err
		}
	}

	log.Println(userId)

	// check if user exist
	user, err := s.GetUserByID(userId)
	if err != nil {
		delErr := s.DeleteToken(accessDetails.UserID)
		if delErr != nil {
			return nil, fmt.Errorf("delete error: %v", delErr)
		}

		return nil, fmt.Errorf("token expired, log in again")
		// return nil, err
	}

	// if the account must be admin
	// if needAdmin {
	// 	if !user.Admin {
	// 		return nil, fmt.Errorf("unauthorized! not admin")
	// 	}
	// }

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
