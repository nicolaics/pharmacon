import requests as rq
import random


def create_users(BACKEND_ROOT, TOKEN):
    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/user/register",
        json={
            "name": "Adam adam",
            "password": "adampassword",
            "admin": True,
            "phoneNumber": "0819288176",
            "adminPassword": "test",
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        },
    )
    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/user/register",
        json={
            "name": "Bob bob",
            "password": "bobpassword",
            "admin": False,
            "phoneNumber": "0819288326",
            "adminPassword": "test",
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        },
    )
    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/user/register",
        json={
            "name": "Charlie charlie",
            "password": "charliepassword",
            "admin": False,
            "phoneNumber": "0813318326",
            "adminPassword": "test",
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        },
    )
    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/user/register",
        json={
            "name": "Delta delta",
            "password": "deltapassword",
            "admin": True,
            "phoneNumber": "081921296",
            "adminPassword": "test",
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        },
    )


def create_company_profile(BACKEND_ROOT, TOKEN):
    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/company-profile",
        json={
            "name": "Zulu",
            "address": "zulu abc",
            "businessNumber": "12345",
            "pharmacsit": "AA",
            "pharmacistLicenseNumber": "1239901",
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        },
    )


def create_customers(BACKEND_ROOT, TOKEN):
    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/customer",
        json={"name": "Graph"},
        headers={
            "Authorization": "Bearer " + TOKEN,
        },
    )
    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/customer",
        json={"name": "Alpha"},
        headers={
            "Authorization": "Bearer " + TOKEN,
        },
    )
    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/customer",
        json={"name": "Beta"},
        headers={
            "Authorization": "Bearer " + TOKEN,
        },
    )
    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/customer",
        json={"name": "Gama"},
        headers={
            "Authorization": "Bearer " + TOKEN,
        },
    )


def create_doctors(BACKEND_ROOT, TOKEN):
    r = rq.post(f"http://{BACKEND_ROOT}/api/v1/doctor", json={"name": "Dr. Gray"}, headers={
            "Authorization": "Bearer " + TOKEN,
        })
    r = rq.post(f"http://{BACKEND_ROOT}/api/v1/doctor", json={"name": "Dr. Jay"}, headers={
            "Authorization": "Bearer " + TOKEN,
        })
    r = rq.post(f"http://{BACKEND_ROOT}/api/v1/doctor", json={"name": "Dr. Awesome"}, headers={
            "Authorization": "Bearer " + TOKEN,
        })
    r = rq.post(f"http://{BACKEND_ROOT}/api/v1/doctor", json={"name": "Dr. Pole"}, headers={
            "Authorization": "Bearer " + TOKEN,
        })
    r = rq.post(f"http://{BACKEND_ROOT}/api/v1/doctor", json={"name": "Dr. Ulala"}, headers={
            "Authorization": "Bearer " + TOKEN,
        })
    r = rq.post(f"http://{BACKEND_ROOT}/api/v1/doctor", json={"name": "Dr. Oscar"}, headers={
            "Authorization": "Bearer " + TOKEN,
        })


def create_patients(BACKEND_ROOT, TOKEN):
    r = rq.post(f"http://{BACKEND_ROOT}/api/v1/patient", json={"name": "Yankee"}, headers={
            "Authorization": "Bearer " + TOKEN,
        })
    r = rq.post(f"http://{BACKEND_ROOT}/api/v1/patient", json={"name": "Awesome"}, headers={
            "Authorization": "Bearer " + TOKEN,
        })
    r = rq.post(f"http://{BACKEND_ROOT}/api/v1/patient", json={"name": "Two"}, headers={
            "Authorization": "Bearer " + TOKEN,
        })


def create_suppliers(BACKEND_ROOT, TOKEN):
    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/supplier",
        json={
            "name": "Yes",
            "address": "no",
            "companyPhoneNumber": str(random.randint(1000000, 5000000)),
            "contactPersonName": "Al",
            "contactPersonNumber": str(random.randint(1000000, 5000000)),
            "terms": "ok",
            "vendorIsTaxable": True,
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        }
    )

    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/supplier",
        json={
            "name": "No",
            "address": "nodsafadfadgadsfgafgadsfdas",
            "companyPhoneNumber": str(random.randint(1000000, 5000000)),
            "contactPersonName": "Ald",
            "contactPersonNumber": str(random.randint(1000000, 5000000)),
            "terms": "no",
            "vendorIsTaxable": True,
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        }
    )

    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/supplier",
        json={
            "name": "Okay",
            "address": "dasfdasfasdfasdfasdfadsfwefrf",
            "companyPhoneNumber": str(random.randint(1000000, 5000000)),
            "contactPersonName": "Aldde",
            "contactPersonNumber": str(random.randint(1000000, 5000000)),
            "terms": "no",
            "vendorIsTaxable": False,
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        }
    )


def create_medicines(BACKEND_ROOT, TOKEN):
    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/medicine",
        json={
            "barcode": "1111",
            "name": "Ativan",
            "qty": 1000,
            "firstUnit": "BTL",
            "firstSubtotal": 11000,
            "firstDiscount": 1000,
            "firstPrice": 10000,
            "secondUnit": "BOX",
            "secondSubtotal": 112000,
            "secondPrice": 112000,
            "description": "test",
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        }
    )

    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/medicine",
        json={
            "barcode": "2222",
            "name": "Panadol",
            "qty": 500,
            "firstUnit": "TAB",
            "firstSubtotal": 3320,
            "firstPrice": 3320,
            "description": "test",
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        }
    )

    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/medicine",
        json={
            "barcode": "3333",
            "name": "Rhinos",
            "qty": 320,
            "firstUnit": "STP",
            "firstSubtotal": 19200,
            "firstDiscount": 200,
            "firstPrice": 19000,
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        }
    )


def create_invoices(BACKEND_ROOT, TOKEN):
    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/invoice",
        json={
            "number": 1,
            "customerId": 1,
            "subtotal": 10000,
            "discount": 100,
            "totalPrice": 9000,
            "paidAmount": 10000,
            "changeAmount": 1000,
            "paymentMethodName": "Cash",
            "invoiceDate": "2024-09-23",
            "medicineLists": [
                {
                    "medicineBarcode": "2222",
                    "medicineName": "Panadol",
                    "qty": 9,
                    "unit": "TAB",
                    "price": 3320,
                    "subtotal": 90000,
                }
            ],
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        }
    )

    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/invoice",
        json={
            "number": 2,
            "customerId": 5,
            "subtotal": 37000,
            "totalPrice": 37000,
            "paidAmount": 37000,
            "changeAmount": 0,
            "paymentMethodName": "Trasnfer",
            "invoiceDate": "2024-09-21",
            "medicineLists": [
                {
                    "medicineBarcode": "1111",
                    "medicineName": "Ativan",
                    "qty": 1,
                    "unit": "BTL",
                    "price": 11000,
                    "discount": 1000,
                    "subtotal": 30000,
                },
                {
                    "medicineBarcode": "1237738438",
                    "medicineName": "Rhinos",
                    "qty": 1,
                    "unit": "STP",
                    "price": 19000,
                    "discount": 3000,
                    "subtotal": 7000,
                },
            ],
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        }
    )


def create_po_invoices(BACKEND_ROOT, TOKEN):
    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/invoice/purchase-order",
        json={
            "number": 2,
            "companyId": 1,
            "supplierId": 3,
            "totalItems": 10,
            "invoiceDate": "2024-03-12",
            "purchaseOrderMedicineList": [
                {
                    "medicineBarcode": "1111",
                    "medicineName": "Ativan",
                    "orderQty": 7,
                    "receivedQty": 3,
                    "unit": "BTL",
                },
                {
                    "medicineBarcode": "3333",
                    "medicineName": "Rhinos",
                    "orderQty": 10,
                    "receivedQty": 10,
                    "unit": "TAB",
                    "remarks": "dfalks",
                },
            ],
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        }
    )

    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/invoice/purchase-order",
        json={
            "number": 2,
            "companyId": 1,
            "supplierId": 3,
            "totalItems": 10,
            "invoiceDate": "2024-107-12",
            "purchaseOrderMedicineList": [
                {
                    "medicineBarcode": "2222",
                    "medicineName": "Panadol",
                    "orderQty": 10,
                    "receivedQty": 10,
                    "unit": "TAB",
                    "remarks": "oke",
                }
            ],
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        }
    )


def create_prescription(BACKEND_ROOT, TOKEN):
    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/prescription",
        json={
            "invoice": {
                "number": 1,
                "userName": "admin",
                "customerName": "Alpha",
                "totalPrice": 30000,
                "invoiceDate": "2024-07-10",
            },
            "number": 100,
            "prescriptionDate": "2024-06-30",
            "patientName": "Yankee",
            "doctorName": "Dr. Play",
            "qty": 1,
            "price": 100000,
            "totalPrice": 100000,
            "prescriptionMedicineList": [
                {
                    "medicineBarcode": "1111",
                    "medicineName": "Ativan",
                    "qty": 7,
                    "unit": "BTL",
                    "price": 46200,
                    "subtotal": 46200,
                },
                {
                    "medicineBarcode": "3333",
                    "medicineName": "Rhinos",
                    "qty": 10,
                    "unit": "TAB",
                    "price": 2200,
                    "discount": 200,
                    "subtotal": 2000,
                },
            ],
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        }
    )

    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/prescription",
        json={
            "invoice": {
                "number": 1,
                "userName": "admin",
                "customerName": "Alpha",
                "totalPrice": 30000,
                "invoiceDate": "2024-07-10",
            },
            "number": 100,
            "prescriptionDate": "2024-06-30",
            "patientName": "Bpay",
            "doctorName": "Dr. Gray",
            "qty": 1,
            "price": 100000,
            "totalPrice": 100000,
            "prescriptionMedicineList": [
                {
                    "medicineBarcode": "2222",
                    "medicineName": "Panadol",
                    "qty": 7,
                    "unit": "TAB",
                    "price": 12900,
                    "subtotal": 12900,
                }
            ],
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        }
    )


def create_production(BACKEND_ROOT, TOKEN):
    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/production",
        json={
            "batchNumber": 2,
            "producedMedicineBarcode": random.randint(1, 10000),
            "producedMedicineName": "S2",
            "producedQty": random.randint(1, 200),
            "productionDate": "2024-01-10",
            "updatedToStock": False,
            "updatedToAccount": False,
            "totalCost": random.randint(1, 100000),
            "productionMedicineList": [
                {
                    "medicineBarcode": "1111",
                    "medicineName": "Ativan",
                    "qty": 7,
                    "unit": "BTL",
                    "cost": 46200,
                },
                {
                    "medicineBarcode": "2222",
                    "medicineName": "Panadol",
                    "qty": 10,
                    "unit": "TAB",
                    "cost": 2200,
                },
            ],
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        }
    )

    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/production",
        json={
            "batchNumber": 210,
            "producedMedicineBarcode": random.randint(1, 10000),
            "producedMedicineName": "S1",
            "producedQty": random.randint(1, 200),
            "productionDate": "2024-04-20",
            "updatedToStock": False,
            "updatedToAccount": False,
            "totalCost": random.randint(1, 100000),
            "productionMedicineList": [
                {
                    "medicineBarcode": "3333",
                    "medicineName": "Rhinos",
                    "qty": 10,
                    "unit": "TAB",
                    "cost": 462000,
                },
                {
                    "medicineBarcode": "2222",
                    "medicineName": "Panadol",
                    "qty": 10,
                    "unit": "TAB",
                    "cost": 2200,
                },
            ],
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        }
    )


def create_purchase_invoice(BACKEND_ROOT, TOKEN):
    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/invoice/purchase",
        json={
            "number": random.randint(1, 100),
            "companyId": 1,
            "supplierId": random.randint(1, 3),
            "subtotal": random.randint(1, 100000),
            "discount": random.randint(1, 10000),
            "tax": random.randint(1, 200),
            "totalPrice": random.randint(1, 300000),
            "description": "test",
            "invoiceDate": "2024-05-12",
            "purchaseMedicineList": [
                {
                    "medicineBarcode": "1111",
                    "medicineName": "Ativan",
                    "qty": random.randint(1, 20),
                    "unit": "BTL",
                    "price": random.randint(1, 50000),
                    "discount": random.randint(1, 500),
                    "tax": random.randint(1, 500),
                    "subtotal": random.randint(1, 200000),
                    "batchNumber": "B1234",
                    "expDate": "2026-01-10",
                },
                {
                    "medicineBarcode": "3333",
                    "medicineName": "Rhinos",
                    "qty": random.randint(1, 20),
                    "unit": "TAB",
                    "price": random.randint(1, 50000),
                    "subtotal": random.randint(1, 200000),
                    "batchNumber": "B1235",
                    "expDate": "2026-02-10",
                },
            ],
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        }
    )

    r = rq.post(
        f"http://{BACKEND_ROOT}/api/v1/invoice/purchase",
        json={
            "number": random.randint(1, 100),
            "companyId": 1,
            "supplierId": random.randint(1, 3),
            "subtotal": random.randint(1, 100000),
            "discount": random.randint(1, 10000),
            "tax": random.randint(1, 200),
            "totalPrice": random.randint(1, 300000),
            "invoiceDate": "2024-06-12",
            "purchaseMedicineList": [
                {
                    "medicineBarcode": "2222",
                    "medicineName": "Panadol",
                    "qty": random.randint(1, 20),
                    "unit": "BTL",
                    "price": random.randint(1, 50000),
                    "discount": random.randint(1, 500),
                    "subtotal": random.randint(1, 200000),
                    "batchNumber": "B124",
                    "expDate": "2026-12-06",
                },
                {
                    "medicineBarcode": "3333",
                    "medicineName": "Rhinos",
                    "qty": random.randint(1, 20),
                    "unit": "TAB",
                    "price": random.randint(1, 50000),
                    "subtotal": random.randint(1, 200000),
                    "batchNumber": "B125",
                    "expDate": "2026-02-10",
                },
            ],
        },
        headers={
            "Authorization": "Bearer " + TOKEN,
        }
    )


def main():
    BACKEND_ROOT = input("enter backend host:port: ")
    TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXV0aG9yaXplZCI6dHJ1ZSwiZXhwaXJlZEF0IjoxNzI3NTQzNzkxLCJ0b2tlblV1aWQiOiIwMTkyMzcwZi0yMDNjLTcxYjMtYTU4Ny04YzIxNWExOTAwMGIiLCJ1c2VySWQiOjR9.XOFxx4RykwJ03oKFt1CMe7qKXQA8AS2jU7b9VIZ7tTs"
    create_users(BACKEND_ROOT, TOKEN)
    create_company_profile(BACKEND_ROOT, TOKEN)
    create_customers(BACKEND_ROOT, TOKEN)
    create_doctors(BACKEND_ROOT, TOKEN)
    create_patients(BACKEND_ROOT, TOKEN)
    create_suppliers(BACKEND_ROOT, TOKEN)
    create_medicines(BACKEND_ROOT, TOKEN)
    create_invoices(BACKEND_ROOT, TOKEN)
    create_po_invoices(BACKEND_ROOT, TOKEN)
    create_prescription(BACKEND_ROOT, TOKEN)
    create_production(BACKEND_ROOT, TOKEN)
    create_purchase_invoice(BACKEND_ROOT, TOKEN)

    print("DONE")


if __name__ == "__main__":
    main()
