package purchaseinvoice

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/nicolaics/pos_pharmacy/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetPurchaseInvoiceByNumber(number int) (*types.PurchaseInvoice, error) {
	rows, err := s.db.Query("SELECT * FROM purchase_invoice WHERE number = ?", number)
	if err != nil {
		return nil, err
	}

	purchaseInvoice := new(types.PurchaseInvoice)

	for rows.Next() {
		purchaseInvoice, err = scanRowIntoPurchaseInvoice(rows)

		if err != nil {
			return nil, err
		}
	}

	if purchaseInvoice.ID == 0 {
		return nil, fmt.Errorf("purchase invoice not found")
	}

	return purchaseInvoice, nil
}

func (s *Store) GetPurchaseInvoiceByID(id int) (*types.PurchaseInvoice, error) {
	rows, err := s.db.Query("SELECT * FROM purchase_invoice WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	purchaseInvoice := new(types.PurchaseInvoice)

	for rows.Next() {
		purchaseInvoice, err = scanRowIntoPurchaseInvoice(rows)

		if err != nil {
			return nil, err
		}
	}

	if purchaseInvoice.ID == 0 {
		return nil, fmt.Errorf("purchase invoice not found")
	}

	return purchaseInvoice, nil
}

func (s *Store) CreatePurchaseInvoice(purchaseInvoice types.PurchaseInvoice) error {
	fields := "number, company_id, supplier_id, subtotal, discount, tax, total_price, description, "
	fields += "cashier_id, invoice_date"
	values := "?"

	for i := 0; i < 9; i++ {
		values += ", ?"
	}

	_, err := s.db.Exec(fmt.Sprintf("INSERT INTO purchase_invoice (%s) VALUES (%s)", fields, values),
						purchaseInvoice.Number, purchaseInvoice.CompanyID, purchaseInvoice.SupplierID,
						purchaseInvoice.Subtotal, purchaseInvoice.Discount, purchaseInvoice.Tax,
						purchaseInvoice.TotalPrice, purchaseInvoice.Description, purchaseInvoice.CashierID,
						purchaseInvoice.InvoiceDate)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) CreatePurchaseMedicineItems(purchaseMedItem types.PurchaseMedicineItem) error {
	fields := "purchase_invoice_id, medicine_id, qty, unit_id, "
	fields += "purchase_price, purchase_discount, purchase_tax, "
	fields += "subtotal, batch_number, expired_date"
	values := "?"

	for i := 0; i < 9; i++ {
		values += ", ?"
	}

	_, err := s.db.Exec(fmt.Sprintf("INSERT INTO purchase_medicine_items (%s) VALUES (%s)", fields, values),
						purchaseMedItem.PurchaseInvoiceID, purchaseMedItem.MedicineID, purchaseMedItem.Qty,
						purchaseMedItem.UnitID, purchaseMedItem.PurchasePrice, purchaseMedItem.PurchaseDiscount,
						purchaseMedItem.PurchaseTax, purchaseMedItem.Subtotal, purchaseMedItem.BatchNumber,
						purchaseMedItem.ExpDate)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetPurhcaseInvoices(startDate time.Time, endDate time.Time) ([]types.PurchaseInvoice, error) {
	query := fmt.Sprintf("SELECT * FROM purchase_invoice WHERE invoice_date BETWEEN DATE('%s') AND DATE('%s') ORDER BY invoice_date DESC",
				startDate, endDate)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	purchaseInvoices := make([]types.PurchaseInvoice, 0)

	for rows.Next() {
		purchaseInvoice, err := scanRowIntoPurchaseInvoice(rows)

		if err != nil {
			return nil, err
		}

		purchaseInvoices = append(purchaseInvoices, *purchaseInvoice)
	}

	return purchaseInvoices, nil
}

func (s *Store) GetPurhcaseMedicineItems(purchaseInvoiceId int) ([]types.PurchaseMedicineItemsReturn, error) {
	query := "SELECT "

	query += "pmi.id, medicine.name, pmi.qty, unit.unit, pmi.purchase_price, pmi.purchase_discount, "
	query += "pmi.purchase_tax, pmi.subtotal, pmi.batch_numbebr, pmi.expired_date "

	query += "FROM purchase_medicine_items as pmi "
	query += "JOIN purchase_invoice as pi ON pmi.purchase_invoice_id = pi.id "
	query += "JOIN medicine ON pmi.medicine_id = medicine.id "
	query += "JOIN unit ON pmi.unit_id = unit.id "
	query += "WHERE pi.id = ? "

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	purchaseMedicineItems := make([]types.PurchaseMedicineItemsReturn, 0)

	for rows.Next() {
		purchaseMedicineItem, err := scanRowIntoPurchaseMedicineItems(rows)

		if err != nil {
			return nil, err
		}

		purchaseMedicineItems = append(purchaseMedicineItems, *purchaseMedicineItem)
	}

	return purchaseMedicineItems, nil
}

func (s *Store) DeletePurchaseInvoice(purchaseInvoice *types.PurchaseInvoice) error {
	_, err := s.db.Exec("DELETE FROM purchase_medicine_items WHERE purchase_invoice_id = ? ", purchaseInvoice.ID)
	if err != nil {
		return err
	}
	
	_, err = s.db.Exec("DELETE FROM purchase_invoice WHERE id = ?", purchaseInvoice.ID)
	if err != nil {
		return err
	}

	return nil
}

/*
func (s *Store) ModifyCustomer(id int, newName string) error {
	_, err := s.db.Exec("UPDATE customer SET name = ? WHERE id = ? ", newName, id)

	if err != nil {
		return err
	}

	return nil
}
*/

func scanRowIntoPurchaseInvoice(rows *sql.Rows) (*types.PurchaseInvoice, error) {
	purchaseInvoice := new(types.PurchaseInvoice)

	err := rows.Scan(
		&purchaseInvoice.ID,
		&purchaseInvoice.Number,
		&purchaseInvoice.CompanyID,
		&purchaseInvoice.SupplierID,
		&purchaseInvoice.Subtotal,
		&purchaseInvoice.Discount,
		&purchaseInvoice.Tax,
		&purchaseInvoice.TotalPrice,
		&purchaseInvoice.Description,
		&purchaseInvoice.CashierID,
		&purchaseInvoice.InvoiceDate,
		&purchaseInvoice.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	purchaseInvoice.InvoiceDate = purchaseInvoice.InvoiceDate.Local()
	purchaseInvoice.CreatedAt = purchaseInvoice.CreatedAt.Local()

	return purchaseInvoice, nil
}

func scanRowIntoPurchaseMedicineItems(rows *sql.Rows) (*types.PurchaseMedicineItemsReturn, error) {
	purchaseMedicineItem := new(types.PurchaseMedicineItemsReturn)

	err := rows.Scan(
		&purchaseMedicineItem.ID,
		&purchaseMedicineItem.MedicineName,
		&purchaseMedicineItem.Qty,
		&purchaseMedicineItem.Unit,
		&purchaseMedicineItem.Price,
		&purchaseMedicineItem.Discount,
		&purchaseMedicineItem.Tax,
		&purchaseMedicineItem.Subtotal,
		&purchaseMedicineItem.BatchNumber,
		&purchaseMedicineItem.ExpDate,
	)

	if err != nil {
		return nil, err
	}

	purchaseMedicineItem.ExpDate = purchaseMedicineItem.ExpDate.Local()

	return purchaseMedicineItem, nil
}