package customer

import (
	"database/sql"
	"fmt"

	"github.com/nicolaics/pos_pharmacy/types"
	"github.com/nicolaics/pos_pharmacy/utils"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetCustomerByName(name string) (*types.Customer, error) {
	rows, err := s.db.Query("SELECT * FROM customer WHERE name = ? ", name)

	if err != nil {
		return nil, err
	}

	customer := new(types.Customer)

	for rows.Next() {
		customer, err = scanRowIntoCustomer(rows)

		if err != nil {
			return nil, err
		}
	}

	if customer.ID == 0 {
		return nil, fmt.Errorf("customer not found")
	}

	return customer, nil
}

func (s *Store) GetCashierByID(id int) (*types.Customer, error) {
	rows, err := s.db.Query("SELECT * FROM customer WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	customer := new(types.Customer)

	for rows.Next() {
		customer, err = scanRowIntoCustomer(rows)

		if err != nil {
			return nil, err
		}
	}

	if customer.ID == 0 {
		return nil, fmt.Errorf("customer not found")
	}

	return customer, nil
}

func (s *Store) CreateCashier(customer types.Customer) error {
	_, err := s.db.Exec("INSERT INTO customer (name) VALUES (?)",
					customer.Name)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) FindCustomerID(customerName string) (int, error) {
	rows, err := s.db.Query("SELECT * FROM customer WHERE name = ? ", customerName)

	if err != nil {
		return -1, err
	}

	customer := new(types.Customer)

	for rows.Next() {
		customer, err = scanRowIntoCustomer(rows)

		if err != nil {
			return -1, err
		}
	}

	return customer.ID, nil
}

func scanRowIntoCustomer(rows *sql.Rows) (*types.Customer, error) {
	customer := new(types.Customer)

	err := rows.Scan(
		&customer.ID,
		&customer.Name,
		&customer.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return customer, nil
}