package invoice

import (
	"database/sql"
	"fmt"
	"strconv"
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

func (s *Store) GetInvoiceByID(id int) (*types.Invoice, error) {
	query := "SELECT * FROM invoice WHERE id = ? AND deleted_at IS NULL"
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoice := new(types.Invoice)

	for rows.Next() {
		invoice, err = scanRowIntoInvoice(rows)

		if err != nil {
			return nil, err
		}
	}

	if invoice.ID == 0 {
		return nil, fmt.Errorf("invoice not found")
	}

	return invoice, nil
}

func (s *Store) GetInvoiceID(number int, userId int, customerId int, totalPrice float64, invoiceDate time.Time) (int, error) {
	query := `SELECT id FROM invoice 
				WHERE number = ? AND user_id = ? AND customer_id = ? AND 
				total_price = ? AND invoice_date = ? AND deleted_at IS NULL`

	rows, err := s.db.Query(query, number, userId, customerId, totalPrice, invoiceDate)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var invoiceId int

	for rows.Next() {
		err = rows.Scan(&invoiceId)
		if err != nil {
			return 0, err
		}
	}

	if invoiceId == 0 {
		return 0, fmt.Errorf("invoice not found")
	}

	return invoiceId, nil
}

func (s *Store) GetInvoicesByNumber(number int) ([]types.Invoice, error) {
	query := "SELECT * FROM invoice WHERE number LIKE ? AND deleted_at IS NULL"

	searchVal := "%"
	for _, val := range strconv.Itoa(number) {
		if string(val) != " " {
			searchVal += (string(val) + "%")
		}
	}

	rows, err := s.db.Query(query, searchVal)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := make([]types.Invoice, 0)

	for rows.Next() {
		invoice, err := scanRowIntoInvoice(rows)
		if err != nil {
			return nil, err
		}

		invoices = append(invoices, *invoice)
	}

	return invoices, nil
}

func (s *Store) GetInvoicesByDate(startDate time.Time, endDate time.Time) ([]types.InvoiceListsReturnPayload, error) {
	query := `SELECT invoice.id, invoice.number, 
					user.name, customer.name, 
					invoice.subtotal, invoice.discount, 
					invoice.tax, invoice.total_price, 
					payment_method.name, 
					invoice.description, invoice.invoice_date 
					FROM invoice 
					JOIN user ON user.id = invoice.user_id 
					JOIN customer ON customer.id = invoice.customer_id 
					JOIN payment_method ON payment_method.id = invoice.payment_method_id 
				WHERE (invoice.invoice_date BETWEEN DATE(?) AND DATE(?)) 
					AND invoice.deleted_at IS NULL 
				ORDER BY invoice.invoice_date DESC`
	rows, err := s.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := make([]types.InvoiceListsReturnPayload, 0)

	for rows.Next() {
		invoice, err := scanRowIntoInvoiceLists(rows)
		if err != nil {
			return nil, err
		}

		invoices = append(invoices, *invoice)
	}

	return invoices, nil
}

func (s *Store) GetInvoicesByDateAndNumber(startDate time.Time, endDate time.Time, number int) ([]types.InvoiceListsReturnPayload, error) {
	query := `SELECT COUNT(*) 
				FROM invoice 
				WHERE (invoice_date BETWEEN DATE(?) AND DATE(?)) 
				AND number = ? 
				AND deleted_at IS NULL 
				ORDER BY invoice_date DESC`

	row := s.db.QueryRow(query, startDate, endDate, number)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var count int

	err := row.Scan(&count)
	if err != nil {
		return nil, err
	}

	invoices := make([]types.InvoiceListsReturnPayload, 0)

	if count == 0 {
		query = `SELECT invoice.id, invoice.number, 
					user.name, customer.name, 
					invoice.subtotal, invoice.discount, 
					invoice.tax, invoice.total_price, 
					payment_method.name, 
					invoice.description, invoice.invoice_date 
					FROM invoice 
					JOIN user ON user.id = invoice.user_id 
					JOIN customer ON customer.id = invoice.customer_id 
					JOIN payment_method ON payment_method.id = invoice.payment_method_id 
					WHERE (invoice.invoice_date BETWEEN DATE(?) AND DATE(?)) 
					AND invoice.number LIKE ? 
					AND invoice.deleted_at IS NULL 
					ORDER BY invoice.invoice_date DESC`

		searchVal := "%"
		for _, val := range strconv.Itoa(number) {
			if string(val) != " " {
				searchVal += (string(val) + "%")
			}
		}

		rows, err := s.db.Query(query, startDate, endDate, searchVal)
		if err != nil {
			return nil, err
		}
		defer rows.Close()


		for rows.Next() {
			invoice, err := scanRowIntoInvoiceLists(rows)
			if err != nil {
				return nil, err
			}

			invoices = append(invoices, *invoice)
		}

		return invoices, nil
	}

	query = `SELECT invoice.id, invoice.number, 
					user.name, customer.name, 
					invoice.subtotal, invoice.discount, 
					invoice.tax, invoice.total_price, 
					payment_method.name, 
					invoice.description, invoice.invoice_date 
					FROM invoice 
					JOIN user ON user.id = invoice.user_id 
					JOIN customer ON customer.id = invoice.customer_id 
					JOIN payment_method ON payment_method.id = invoice.payment_method_id 
					WHERE (invoice.invoice_date BETWEEN DATE(?) AND DATE(?)) 
					AND invoice.number = ? 
					AND invoice.deleted_at IS NULL 
					ORDER BY invoice.invoice_date DESC`
	
	rows, err := s.db.Query(query, startDate, endDate, number)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		invoice, err := scanRowIntoInvoiceLists(rows)

		if err != nil {
			return nil, err
		}

		invoices = append(invoices, *invoice)
	}

	return invoices, nil
}

func (s *Store) GetInvoicesByDateAndUserID(startDate time.Time, endDate time.Time, uid int) ([]types.InvoiceListsReturnPayload, error) {
	query := `SELECT invoice.id, invoice.number, 
					user.name, customer.name, 
					invoice.subtotal, invoice.discount, 
					invoice.tax, invoice.total_price, 
					payment_method.name, 
					invoice.description, invoice.invoice_date 
					FROM invoice 
					JOIN user ON user.id = invoice.user_id 
					JOIN customer ON customer.id = invoice.customer_id 
					JOIN payment_method ON payment_method.id = invoice.payment_method_id 
				WHERE (invoice.invoice_date BETWEEN DATE(?) AND DATE(?)) AND user_id = ? 
					AND invoice.deleted_at IS NULL 
				ORDER BY invoice.invoice_date DESC`

	rows, err := s.db.Query(query, startDate, endDate, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := make([]types.InvoiceListsReturnPayload, 0)

	for rows.Next() {
		invoice, err := scanRowIntoInvoiceLists(rows)
		if err != nil {
			return nil, err
		}

		invoices = append(invoices, *invoice)
	}

	return invoices, nil
}

func (s *Store) GetInvoicesByDateAndCustomerID(startDate time.Time, endDate time.Time, cid int) ([]types.InvoiceListsReturnPayload, error) {
	query := `SELECT invoice.id, invoice.number, 
					user.name, customer.name, 
					invoice.subtotal, invoice.discount, 
					invoice.tax, invoice.total_price, 
					payment_method.name, 
					invoice.description, invoice.invoice_date 
					FROM invoice 
					JOIN user ON user.id = invoice.user_id 
					JOIN customer ON customer.id = invoice.customer_id 
					JOIN payment_method ON payment_method.id = invoice.payment_method_id 
				WHERE (invoice.invoice_date BETWEEN DATE(?) AND DATE(?)) AND customer_id = ? 
					AND invoice.deleted_at IS NULL 
				ORDER BY invoice.invoice_date DESC`

	rows, err := s.db.Query(query, startDate, endDate, cid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := make([]types.InvoiceListsReturnPayload, 0)

	for rows.Next() {
		invoice, err := scanRowIntoInvoiceLists(rows)
		if err != nil {
			return nil, err
		}

		invoices = append(invoices, *invoice)
	}

	return invoices, nil
}

func (s *Store) GetInvoicesByDateAndPaymentMethodID(startDate time.Time, endDate time.Time, pmid int) ([]types.InvoiceListsReturnPayload, error) {
	query := `SELECT invoice.id, invoice.number, 
					user.name, customer.name, 
					invoice.subtotal, invoice.discount, 
					invoice.tax, invoice.total_price, 
					payment_method.name, 
					invoice.description, invoice.invoice_date 
					FROM invoice 
					JOIN user ON user.id = invoice.user_id 
					JOIN customer ON customer.id = invoice.customer_id 
					JOIN payment_method ON payment_method.id = invoice.payment_method_id 
				WHERE (invoice.invoice_date BETWEEN DATE(?) AND DATE(?)) AND payment_method_id = ? 
					AND invoice.deleted_at IS NULL 
				ORDER BY invoice.invoice_date DESC`

	rows, err := s.db.Query(query, startDate, endDate, pmid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := make([]types.InvoiceListsReturnPayload, 0)

	for rows.Next() {
		invoice, err := scanRowIntoInvoiceLists(rows)
		if err != nil {
			return nil, err
		}

		invoices = append(invoices, *invoice)
	}

	return invoices, nil
}

func (s *Store) GetNumberOfInvoices() (int, error) {
	query := `SELECT COUNT(*) FROM invoice`
	row := s.db.QueryRow(query)
	if row.Err() != nil {
		return -1, row.Err()
	}

	var numberOfInvoices int

	err := row.Scan(&numberOfInvoices)
	if err != nil {
		return -1, err
	}

	return numberOfInvoices, nil
}

func (s *Store) CreateInvoice(invoice types.Invoice) error {
	values := "?"
	for i := 0; i < 12; i++ {
		values += ", ?"
	}

	query := `INSERT INTO invoice (
			number, user_id, customer_id, subtotal, discount, tax, 
			total_price, paid_amount, change_amount, payment_method_id, description, 
			invoice_date, last_modified_by_user_id
	) VALUES (` + values + `)`

	_, err := s.db.Exec(query,
		invoice.Number, invoice.UserID, invoice.CustomerID,
		invoice.Subtotal, invoice.Discount, invoice.Tax,
		invoice.TotalPrice, invoice.PaidAmount, invoice.ChangeAmount,
		invoice.PaymentMethodID, invoice.Description, invoice.InvoiceDate,
		invoice.LastModifiedByUserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateMedicineItems(medicineItem types.MedicineItems) error {
	values := "?"
	for i := 0; i < 6; i++ {
		values += ", ?"
	}

	query := `INSERT INTO medicine_items (
		invoice_id, medicine_id, qty, unit_id, price, discount, subtotal
	) VALUES (` + values + `)`
	_, err := s.db.Exec(query,
		medicineItem.InvoiceID, medicineItem.MedicineID, medicineItem.Qty,
		medicineItem.UnitID, medicineItem.Price, medicineItem.Discount,
		medicineItem.Subtotal)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetMedicineItems(invoiceId int) ([]types.MedicineItemReturnPayload, error) {
	query := `SELECT 
				mi.id, 
				medicine.barcode, medicine.name, 
				mi.qty, 
				unit.name, 
				mi.price, mi.discount, mi.subtotal 
				FROM medicine_items as mi 
				JOIN invoice ON mi.invoice_id = invoice.id 
				JOIN medicine ON mi.medicine_id = medicine.id 
				JOIN unit ON mi.unit_id = unit.id 
				WHERE invoice.id = ? AND invoice.deleted_at IS NULL`

	rows, err := s.db.Query(query, invoiceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	medicineItems := make([]types.MedicineItemReturnPayload, 0)

	for rows.Next() {
		medicineItem, err := scanRowIntoMedicineItems(rows)

		if err != nil {
			return nil, err
		}

		medicineItems = append(medicineItems, *medicineItem)
	}

	return medicineItems, nil
}

func (s *Store) DeleteInvoice(invoice *types.Invoice, user *types.User) error {
	query := "UPDATE invoice SET deleted_at = ?, deleted_by_user_id = ? WHERE id = ?"
	_, err := s.db.Exec(query, time.Now(), user.ID, invoice.ID)
	if err != nil {
		return err
	}

	data, err := s.GetInvoiceByID(invoice.ID)
	if err != nil {
		return err
	}

	err = logger.WriteLog("delete", "invoice", user.Name, data.ID, data)
	if err != nil {
		return fmt.Errorf("error write log file")
	}

	return nil
}

func (s *Store) DeleteMedicineItems(invoice *types.Invoice, user *types.User) error {
	data, err := s.GetMedicineItems(invoice.ID)
	if err != nil {
		return err
	}

	writeData := map[string]interface{}{
		"invoice":                invoice,
		"deleted_medicine_items": data,
	}

	err = logger.WriteLog("delete", "invoice", user.Name, invoice.ID, writeData)
	if err != nil {
		return fmt.Errorf("error write log file")
	}

	_, err = s.db.Exec("DELETE FROM medicine_items WHERE invoice_id = ? ", invoice.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) ModifyInvoice(invoiceId int, invoice types.Invoice, user *types.User) error {
	data, err := s.GetInvoiceByID(invoiceId)
	if err != nil {
		return err
	}

	writeData := map[string]interface{}{
		"previous_data": data,
	}

	err = logger.WriteLog("modify", "invoice", user.Name, data.ID, writeData)
	if err != nil {
		return fmt.Errorf("error write log file")
	}

	query := `UPDATE invoice SET 
			number = ?, user_id = ?, customer_id = ?, subtotal = ?, discount = ?, 
			tax = ?, total_price = ?, paid_amount = ?, change_amount = ?, 
			payment_method_id = ?, description = ?, invoice_date = ?, last_modified = ?,
			last_modified_by_user_id = ? 
			WHERE id = ? AND deleted_at IS NULL`

	_, err = s.db.Exec(query,
		invoice.Number, invoice.UserID, invoice.CustomerID,
		invoice.Subtotal, invoice.Discount, invoice.Tax,
		invoice.TotalPrice, invoice.PaidAmount, invoice.ChangeAmount,
		invoice.PaymentMethodID, invoice.Description, invoice.InvoiceDate,
		time.Now(), invoice.LastModifiedByUserID, invoiceId)
	if err != nil {
		return err
	}

	return nil
}

func scanRowIntoInvoice(rows *sql.Rows) (*types.Invoice, error) {
	invoice := new(types.Invoice)

	err := rows.Scan(
		&invoice.ID,
		&invoice.Number,
		&invoice.UserID,
		&invoice.CustomerID,
		&invoice.Subtotal,
		&invoice.Discount,
		&invoice.Tax,
		&invoice.TotalPrice,
		&invoice.PaidAmount,
		&invoice.ChangeAmount,
		&invoice.PaymentMethodID,
		&invoice.Description,
		&invoice.InvoiceDate,
		&invoice.CreatedAt,
		&invoice.LastModified,
		&invoice.LastModifiedByUserID,
		&invoice.DeletedAt,
		&invoice.DeletedByUserID,
	)

	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func scanRowIntoInvoiceLists(rows *sql.Rows) (*types.InvoiceListsReturnPayload, error) {
	invoice := new(types.InvoiceListsReturnPayload)

	err := rows.Scan(
		&invoice.ID,
		&invoice.Number,
		&invoice.UserName,
		&invoice.CustomerName,
		&invoice.Subtotal,
		&invoice.Discount,
		&invoice.Tax,
		&invoice.TotalPrice,
		&invoice.PaymentMethodName,
		&invoice.Description,
		&invoice.InvoiceDate,
	)

	if err != nil {
		return nil, err
	}

	invoice.InvoiceDate = invoice.InvoiceDate.Local()

	return invoice, nil
}

func scanRowIntoMedicineItems(rows *sql.Rows) (*types.MedicineItemReturnPayload, error) {
	medicineItem := new(types.MedicineItemReturnPayload)

	err := rows.Scan(
		&medicineItem.ID,
		&medicineItem.MedicineBarcode,
		&medicineItem.MedicineName,
		&medicineItem.Qty,
		&medicineItem.Unit,
		&medicineItem.Price,
		&medicineItem.Discount,
		&medicineItem.Subtotal,
	)

	if err != nil {
		return nil, err
	}

	return medicineItem, nil
}
