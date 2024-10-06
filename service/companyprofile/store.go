package companyprofile

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/nicolaics/pos_pharmacy/logger"
	"github.com/nicolaics/pos_pharmacy/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetCompanyProfileByName(name string) (*types.CompanyProfile, error) {
	query := `SELECT * FROM self_company_profile WHERE name = ?`
	rows, err := s.db.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	companyProfile := new(types.CompanyProfile)

	for rows.Next() {
		companyProfile, err = scanRowIntoCompanyProfile(rows)

		if err != nil {
			return nil, err
		}
	}

	if companyProfile.ID == 0 {
		return nil, fmt.Errorf("company profile not found")
	}

	return companyProfile, nil
}

func (s *Store) GetCompanyProfileByID(id int) (*types.CompanyProfile, error) {
	query := `SELECT * FROM self_company_profile WHERE id = ?`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	companyProfile := new(types.CompanyProfile)

	for rows.Next() {
		companyProfile, err = scanRowIntoCompanyProfile(rows)
		if err != nil {
			return nil, err
		}
	}

	if companyProfile.ID == 0 {
		return nil, fmt.Errorf("company profile not found")
	}

	return companyProfile, nil
}

func (s *Store) CreateCompanyProfile(companyProfile types.CompanyProfile) error {
	query := `INSERT INTO self_company_profile 
				(name, address, last_modified_by_user_id) 
				VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query,
		companyProfile.Name, companyProfile.Address, companyProfile.LastModifiedByUserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetCompanyProfile() (*types.CompanyProfileReturn, error) {
	query := `SELECT c.id, c.name, c.last_modified, 
					user.name 
					FROM self_company_profile AS c 
					JOIN user ON user.id = c.last_modified_by_user_id`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	companyProfile := new(types.CompanyProfileReturn)

	for rows.Next() {
		companyProfile, err = scanRowIntoCompanyProfileReturn(rows)
		if err != nil {
			return nil, err
		}
	}

	return companyProfile, nil
}

// func (s *Store) DeleteCompanyProfile(cpid int, userId int) error {
// 	query := "UPDATE self_company_profile SET deleted_at = ?, deleted_by_user_id = ? WHERE id = ?"
// 	_, err := s.db.Exec(query, time.Now(), userId, cpid)
// 	if err != nil {
// 		return err
// 	}

// 	data, err := s.GetCompanyProfileByID(cpid)
// 	if err != nil {
// 		return err
// 	}

// 	err = logger.WriteLog("delete", "company-profile", userId, data.ID, data)
// 	if err != nil {
// 		return fmt.Errorf("error write log file")
// 	}

// 	return nil
// }

func (s *Store) ModifyCompanyProfile(id int, user *types.User, newCompanyProfile types.CompanyProfile) error {
	data, err := s.GetCompanyProfileByID(id)
	if err != nil {
		return err
	}

	err = logger.WriteLog("modify", "company-profile", user.Name, data.ID, map[string]interface{}{"previous_data": data})
	if err != nil {
		return fmt.Errorf("error write log file")
	}
	
	query := `UPDATE self_company_profile SET 
			name = ?, address = ?, last_modified = ?, last_modified_by_user_id = ? WHERE id = ?`

	_, err = s.db.Exec(query,
		newCompanyProfile.Name, newCompanyProfile.Address, time.Now(), user.ID, id)
	if err != nil {
		return err
	}



	return nil
}

func scanRowIntoCompanyProfile(rows *sql.Rows) (*types.CompanyProfile, error) {
	companyProfile := new(types.CompanyProfile)

	err := rows.Scan(
		&companyProfile.ID,
		&companyProfile.Name,
		&companyProfile.Address,
		&companyProfile.CreatedAt,
		&companyProfile.LastModified,
		&companyProfile.LastModifiedByUserID,
	)

	if err != nil {
		return nil, err
	}

	companyProfile.CreatedAt = companyProfile.CreatedAt.Local()
	companyProfile.LastModified = companyProfile.LastModified.Local()

	return companyProfile, nil
}

func scanRowIntoCompanyProfileReturn(rows *sql.Rows) (*types.CompanyProfileReturn, error) {
	companyProfile := new(types.CompanyProfileReturn)

	err := rows.Scan(
		&companyProfile.ID,
		&companyProfile.Name,
		&companyProfile.Address,
		&companyProfile.LastModified,
		&companyProfile.LastModifiedByUserName,
	)

	if err != nil {
		return nil, err
	}

	companyProfile.LastModified = companyProfile.LastModified.Local()

	return companyProfile, nil
}
