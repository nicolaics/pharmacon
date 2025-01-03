package pdf

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/nicolaics/pharmacon/config"
	"github.com/nicolaics/pharmacon/constants"
	"github.com/nicolaics/pharmacon/types"
	"github.com/nicolaics/pharmacon/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/go-pdf/fpdf"
)

func CreatePrescriptionPDF(presc types.PrescriptionPDFReturn, prescStore types.PrescriptionStore, prevFileName string) (string, error) {
	directory, err := filepath.Abs("static/pdf/prescription/")
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(directory, 0744); err != nil {
		return "", err
	}

	pdf, err := initPrescriptionPdf()
	if err != nil {
		return "", err
	}

	err = createPrescriptionHeader(pdf, presc.Doctor.Name)
	if err != nil {
		return "", err
	}

	err = createPrescriptionInfo(pdf, presc)
	if err != nil {
		return "", err
	}

	usage := map[string]string{
		"sakit":           "[sS][0-9]+.*",
		"Batuk dan Pilek": "[bB][pP][0-9]+.*",
		"maag":            "[mM][0-9]+.*",
		"sk":              "[sS][kK][0-9]+.*",
		"pusing":          "[pP][0-9]+.*",
	}

	for _, setItem := range presc.MedicineSets {
		if setItem.Usage == "" {
			for idx, med := range setItem.MedicineItems {
				for key, val := range usage {
					r := regexp.MustCompile(val)
					use := r.FindAllString(med.MedicineName, -1)

					if len(use) != 0 {
						setItem.Usage = key
						setItem.MedicineItems = removeMedicineFromList(setItem.MedicineItems, idx)
						break
					}
				}
			}
		}
	}

	err = createPrescriptionData(pdf, presc.MedicineSets)
	if err != nil {
		return "", err
	}

	for i := 0; i < pdf.PageCount(); i++ {
		pdf.SetPage(i + 1)

		err = createPrescriptionFooter(pdf)
		if err != nil {
			return "", err
		}
	}

	fileName := prevFileName

	if prevFileName == "" {
		fileName = "p-" + utils.GenerateRandomCodeAlphanumeric(8) + "-" + utils.GenerateRandomCodeAlphanumeric(8) + ".pdf"
		isExist, err := prescStore.IsPDFUrlExist("prescription", fileName)
		if err != nil {
			return "", err
		}

		for isExist {
			fileName = "p-" + utils.GenerateRandomCodeAlphanumeric(8) + "-" + utils.GenerateRandomCodeAlphanumeric(8) + ".pdf"
			isExist, err = prescStore.IsPDFUrlExist("prescription", fileName)
			if err != nil {
				return "", err
			}
		}
	}

	err = pdf.OutputFileAndClose(directory + "\\" + fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func initPrescriptionPdf() (*fpdf.Fpdf, error) {
	s, _ := filepath.Abs("static/assets/font/")

	pdf := fpdf.NewCustom(&fpdf.InitType{
		OrientationStr: "P",
		UnitStr:        "cm",
		SizeStr:        "10.8x21.3",
		Size: fpdf.SizeType{
			Wd: constants.PRESC_WIDTH,
			Ht: constants.PRESC_HEIGHT,
		},
		FontDirStr: s,
	})

	pdf.SetMargins(0.3, 0.3, 0.3)
	pdf.SetAutoPageBreak(true, constants.PRESC_MARGIN)

	pdf.AddUTF8Font("Jolly", constants.REGULAR, "Jolly.ttf")
	pdf.AddUTF8Font("Mockine", constants.REGULAR, "Mockine.ttf")
	pdf.AddUTF8Font("Arial", constants.REGULAR, "Arial.TTF")
	pdf.AddUTF8Font("Arial", constants.BOLD, "ArialBD.TTF")
	pdf.AddUTF8Font("Arial", constants.ITALIC, "ArialI.TTF")
	pdf.AddUTF8Font("Abuget", constants.REGULAR, "Abuget.ttf")
	pdf.AddUTF8Font("Pristina", constants.REGULAR, "Pristina.ttf")
	pdf.AddUTF8Font("Bree", constants.REGULAR, "bree-serif-regular.ttf")
	pdf.AddUTF8Font("Bree", constants.BOLD, "Bree Serif Bold.ttf")
	pdf.AddUTF8Font("Calibri", constants.REGULAR, "Calibri.TTF")
	pdf.AddUTF8Font("Calibri", constants.BOLD, "CalibriBold.TTF")
	pdf.AddUTF8Font("Aller", constants.BOLD, "aller.bold.ttf")
	pdf.AddUTF8Font("Deco", constants.REGULAR, "A780-Deco Regular.ttf")
	pdf.AddUTF8Font("Ameretto", constants.ITALIC, "Ameretto Extended Italic.ttf")

	pdf.AddPage()

	if pdf.Error() != nil {
		return nil, fmt.Errorf("error init presc pdf: %v", pdf.Error())
	}

	return pdf, nil
}

func createPrescriptionHeader(pdf *fpdf.Fpdf, doctor string) error {
	pdf.SetXY((constants.PRESC_MARGIN + 0.1), 0.3)

	pdf.Image(config.Envs.CompanyLogoURL, pdf.GetX(), pdf.GetY(), constants.PRESC_LOGO_WIDTH, constants.PRESC_LOGO_HEIGHT, false, "", 0, "")

	pdf.SetTextColor(constants.GREEN_R, constants.GREEN_G, constants.GREEN_B)
	pdf.SetDrawColor(constants.GREEN_R, constants.GREEN_G, constants.GREEN_B)

	pdf.SetXY(1.7, 0.4)
	pdf.SetFont("Bree", constants.BOLD, 16)
	pdf.CellFormat(4, constants.PRESC_STD_CELL_HEIGHT, config.Envs.CompanyName, "", 0, "L", false, 0, "")

	pdf.SetFont("Deco", constants.REGULAR, 10)
	pdf.CellFormat(0, constants.PRESC_STD_CELL_HEIGHT, "(Berdiri sejak 1999)", "", 1, "L", false, 0, "")

	pdf.SetX(1.7)
	pdf.SetFont("Calibri", constants.REGULAR, 10)
	pdf.CellFormat(0, constants.PRESC_STD_CELL_HEIGHT, config.Envs.CompanyAddress, "", 1, "L", false, 0, "")

	pdf.SetXY(1.7, (pdf.GetY() - 0.05))
	pdf.SetFont("Calibri", constants.REGULAR, 10)
	phoneNumber := fmt.Sprintf("No. Telp : %s | WhatsApp : %s", config.Envs.CompanyPhoneNumber, config.Envs.CompanyWhatsAppNumber)
	pdf.CellFormat(0, constants.PRESC_STD_CELL_HEIGHT, phoneNumber, "", 1, "L", false, 0, "")

	pdf.SetY(pdf.GetY() + 0.1)
	pdf.SetFont("Calibri", constants.REGULAR, 6)
	pharmacist := fmt.Sprintf("Apoteker : %s No. SIPA : %s", config.Envs.Pharmacist, config.Envs.PharmacistLicenseNumber)
	pdf.CellFormat(0, 0.25, pharmacist, "", 1, "C", false, 0, "")

	if doctor == constants.MAIN_DOCTOR {
		pdf.SetY(pdf.GetY() + 0.05)
		pdf.SetFont("Aller", constants.BOLD, 7)
		doctorText := fmt.Sprintf("%s   No. STR : %s", config.Envs.MainDoctor, config.Envs.MainDoctorLicenseNumber)
		pdf.CellFormat(0, 0.3, doctorText, "", 1, "C", false, 0, "")
	}

	pdf.SetXY(8.2, (pdf.GetY() - 0.05))
	pdf.SetFont("Bree", constants.BOLD, 16)
	pdf.CellFormat(1.4, 0.6, "ITER", "", 1, "R", false, 0, "")

	pdf.SetLineWidth(0.05)
	pdf.Line(0, pdf.GetY(), constants.PRESC_WIDTH, pdf.GetY())

	pdf.SetY(pdf.GetY() + 0.08)
	pdf.Line(0, pdf.GetY(), constants.PRESC_WIDTH, pdf.GetY())

	pdf.SetY(pdf.GetY() + 0.08)

	if pdf.Error() != nil {
		return fmt.Errorf("error create presc pdf header: %v", pdf.Error())
	}

	return nil
}

func createPrescriptionInfo(pdf *fpdf.Fpdf, presc types.PrescriptionPDFReturn) error {
	var caser = cases.Title(language.Indonesian)

	pdf.SetTextColor(constants.GREEN_R, constants.GREEN_G, constants.GREEN_B)

	pdf.SetXY(constants.PRESC_MARGIN, pdf.GetY())
	pdf.SetFont("Ameretto", constants.ITALIC, 13)
	pdf.CellFormat(0, constants.PRESC_STD_CELL_HEIGHT, "SALINAN RESEP", "", 1, "C", false, 0, "")

	pdf.SetTextColor(constants.BLACK_R, constants.BLACK_G, constants.BLACK_B)

	// Number
	{
		pdf.SetFont("Calibri", constants.REGULAR, constants.PRESC_DATA_FONT_SZ)
		pdf.CellFormat(0.7, constants.PRESC_STD_CELL_HEIGHT, "No.", "", 0, "L", false, 0, "")

		pdf.SetX(1.6)
		pdf.SetFont("Calibri", constants.REGULAR, constants.PRESC_DATA_FONT_SZ)
		pdf.CellFormat(0.3, constants.PRESC_STD_CELL_HEIGHT, ":", "", 0, "L", false, 0, "")

		pdf.SetFont("Arial", constants.REGULAR, constants.PRESC_DATA_FONT_SZ)
		pdf.CellFormat(4.5, 0.45, strconv.Itoa(presc.Number), "", 0, "LM", false, 0, "")
	}

	// Date
	{
		pdf.SetFont("Calibri", constants.REGULAR, constants.PRESC_DATA_FONT_SZ)
		pdf.CellFormat(1, constants.PRESC_STD_CELL_HEIGHT, " Tgl", "", 0, "L", false, 0, "")

		pdf.SetX(pdf.GetX() + 0.3)
		pdf.SetFont("Calibri", constants.REGULAR, constants.PRESC_DATA_FONT_SZ)
		pdf.CellFormat(0.3, constants.PRESC_STD_CELL_HEIGHT, ":", "", 0, "L", false, 0, "")

		pdf.SetFont("Arial", constants.REGULAR, constants.PRESC_DATA_FONT_SZ)
		pdf.CellFormat(0, 0.45, presc.Date.Format("02-01-2006"), "", 1, "LM", false, 0, "")
	}

	// Patient
	{
		presc.Patient.Name = caser.String(presc.Patient.Name)

		pdf.SetY(pdf.GetY() - 0.05)
		pdf.SetFont("Calibri", constants.REGULAR, constants.PRESC_DATA_FONT_SZ)
		pdf.CellFormat(1.2, constants.PRESC_STD_CELL_HEIGHT, "Untuk", "", 0, "L", false, 0, "")

		pdf.SetX(1.6)
		pdf.SetFont("Calibri", constants.REGULAR, constants.PRESC_DATA_FONT_SZ)
		pdf.CellFormat(0.3, constants.PRESC_STD_CELL_HEIGHT, ":", "", 0, "L", false, 0, "")

		pdf.SetFont("Arial", constants.REGULAR, constants.PRESC_DATA_FONT_SZ)
		pdf.CellFormat(4.5, 0.45, presc.Patient.Name, "", 0, "LM", false, 0, "")

		pdf.SetFont("Calibri", constants.REGULAR, constants.PRESC_DATA_FONT_SZ)
		pdf.CellFormat(1.2, constants.PRESC_STD_CELL_HEIGHT, " Umur", "", 0, "L", false, 0, "")

		pdf.SetX(pdf.GetX() + 0.1)
		pdf.SetFont("Calibri", constants.REGULAR, constants.PRESC_DATA_FONT_SZ)
		pdf.CellFormat(0.3, constants.PRESC_STD_CELL_HEIGHT, ":", "", 0, "L", false, 0, "")

		if presc.Patient.Age > 0 {
			pdf.SetFont("Arial", constants.REGULAR, constants.PRESC_DATA_FONT_SZ)
			pdf.CellFormat(0, 0.45, strconv.Itoa(presc.Patient.Age), "", 1, "LM", false, 0, "")
		} else {
			pdf.CellFormat(0, 0.45, ".....................", "", 1, "LM", false, 0, "")
		}
	}

	// Doctor
	{
		presc.Doctor.Name = caser.String(presc.Doctor.Name)

		pdf.SetFont("Calibri", constants.REGULAR, constants.PRESC_DATA_FONT_SZ)
		pdf.CellFormat(1.3, constants.PRESC_STD_CELL_HEIGHT, "Dari dr.", "", 0, "L", false, 0, "")

		pdf.SetX(1.6)
		pdf.SetFont("Calibri", constants.REGULAR, constants.PRESC_DATA_FONT_SZ)
		pdf.CellFormat(0.3, constants.PRESC_STD_CELL_HEIGHT, ":", "", 0, "L", false, 0, "")

		pdf.SetFont("Arial", constants.REGULAR, constants.PRESC_DATA_FONT_SZ)
		pdf.CellFormat(0, 0.45, presc.Doctor.Name, "", 1, "LM", false, 0, "")
	}

	if pdf.Error() != nil {
		return fmt.Errorf("error create presc info: %v", pdf.Error())
	}

	return nil
}

func createPrescriptionData(pdf *fpdf.Fpdf, medicineSets []types.PrescriptionSetItemReturn) error {
	var caser = cases.Title(language.Indonesian)

	pdf.SetLineWidth(0.02)
	pdf.SetDrawColor(constants.BLACK_R, constants.BLACK_G, constants.BLACK_B)
	pdf.SetTextColor(constants.BLACK_R, constants.BLACK_G, constants.BLACK_B)

	firstSet := true
	nameRegex := regexp.MustCompile(`[0-9]+`)

	pdf.SetY(pdf.GetY())

	startMedicineX := 1.0

	_, pageBreakTrigger := pdf.GetAutoPageBreak()
	pageBottomMargin := constants.PRESC_HEIGHT - pageBreakTrigger

	for _, medicineSet := range medicineSets {
		pdf.SetX(0.7)

		// add (2 + margin) for safety margin
		if (pdf.GetY() + (constants.PRESC_STD_CELL_HEIGHT*float64(len(medicineSet.MedicineItems)) + 2.0) + constants.PRESC_MARGIN) > pageBottomMargin {
			pdf.AddPage()

			// change top margin into 0.5
			pdf.SetY(0.5)
		} else {
			if firstSet {
				pdf.SetY(pdf.GetY() + 0.3)
				firstSet = false
			} else {
				pdf.SetY(pdf.GetY() - constants.PRESC_MARGIN)
				firstSet = false
			}
		}

		pdf.SetFont(constants.PRESC_R_SLASH_FONT, constants.REGULAR, constants.PRESC_R_SLASH_FONT_SZ)
		cellWidth := pdf.GetStringWidth("R|") + constants.PRESC_MARGIN
		pdf.CellFormat(cellWidth, constants.PRESC_STD_CELL_HEIGHT, "R|", "", 0, "L", false, 0, "")

		for _, medicine := range medicineSet.MedicineItems {
			pdf.SetXY(startMedicineX, pdf.GetY()+0.1)

			medicine.MedicineName = caser.String(medicine.MedicineName)
			nameSplit := nameRegex.FindAllStringIndex(medicine.MedicineName, -1)

			if len(nameSplit) == 0 {
				pdf.SetFont(constants.PRESC_MED_FONT, constants.REGULAR, constants.PRESC_MED_FONT_SZ)
				cellWidth = pdf.GetStringWidth(medicine.MedicineName) + constants.PRESC_MARGIN
				pdf.CellFormat(cellWidth, constants.PRESC_STD_CELL_HEIGHT, medicine.MedicineName, "", 0, "L", false, 0, "")

				pdf.SetX(pdf.GetX() + 0.1)
			} else {
				startIdx := 0
				for _, idx := range nameSplit {
					if idx[0] == 0 {
						pdf.SetFont(constants.PRESC_MED_QTY_UNIT_FONT, constants.REGULAR, constants.PRESC_MED_QTY_UNIT_FONT_SZ)
						text := medicine.MedicineName[startIdx:idx[1]]
						cellWidth = pdf.GetStringWidth(text)
						pdf.CellFormat(cellWidth, constants.PRESC_STD_CELL_HEIGHT, text, "", 0, "L", false, 0, "")

						startIdx = idx[1]

						continue
					}

					pdf.SetFont(constants.PRESC_MED_FONT, constants.REGULAR, constants.PRESC_MED_FONT_SZ)
					text := medicine.MedicineName[startIdx:idx[0]]
					cellWidth = pdf.GetStringWidth(text)
					pdf.CellFormat(cellWidth, constants.PRESC_STD_CELL_HEIGHT, text, "", 0, "L", false, 0, "")

					pdf.SetFont(constants.PRESC_MED_QTY_UNIT_FONT, constants.REGULAR, constants.PRESC_MED_QTY_UNIT_FONT_SZ)
					text = medicine.MedicineName[idx[0]:idx[1]]
					cellWidth = pdf.GetStringWidth(text)
					pdf.CellFormat(cellWidth, constants.PRESC_STD_CELL_HEIGHT, text, "", 0, "L", false, 0, "")

					startIdx = idx[1]
				}
			}

			if medicine.QtyString == "" && medicine.Unit == "" {
				pdf.Ln(-1)
			} else {
				fractionIdx := strings.Index(medicine.QtyString, "/")

				pdf.SetX(pdf.GetX() + 0.1)

				if fractionIdx != -1 {
					qtySplit := strings.Split(medicine.QtyString, "/")

					pdf.SetX(pdf.GetX() + 0.3)
					pdf.SetFont(constants.PRESC_MED_QTY_UNIT_FONT, constants.REGULAR, constants.PRESC_MED_QTY_UNIT_FONT_SZ)
					pdf.SubWrite(constants.PRESC_STD_CELL_HEIGHT, qtySplit[0], 10, 4, 0, "")

					pdf.SetX(pdf.GetX() - 0.05)
					pdf.SetFont(constants.PRESC_R_SLASH_FONT, constants.REGULAR, constants.PRESC_MED_QTY_UNIT_FONT_SZ)
					pdf.CellFormat(pdf.GetStringWidth("/"), constants.PRESC_STD_CELL_HEIGHT, "/", "", 0, "L", false, 0, "")

					pdf.SetX(pdf.GetX() - 0.05)
					pdf.SetFont(constants.PRESC_MED_QTY_UNIT_FONT, constants.REGULAR, constants.PRESC_MED_QTY_UNIT_FONT_SZ)
					pdf.SubWrite(constants.PRESC_STD_CELL_HEIGHT, qtySplit[1], 10, -2.5, 0, "")

					pdf.SetXY((pdf.GetX() + 0.1), (pdf.GetY() + 0.05))
				} else {
					pdf.SetFont(constants.PRESC_MED_QTY_UNIT_FONT, constants.REGULAR, constants.PRESC_MED_QTY_UNIT_FONT_SZ)
					cellWidth = pdf.GetStringWidth(medicine.QtyString)
					pdf.CellFormat(cellWidth, constants.PRESC_STD_CELL_HEIGHT, medicine.QtyString, "", 0, "L", false, 0, "")
				}

				pdf.SetX(pdf.GetX() + 0.1)
				pdf.SetFont(constants.PRESC_MED_QTY_UNIT_FONT, constants.REGULAR, constants.PRESC_MED_QTY_UNIT_FONT_SZ)
				cellWidth = pdf.GetStringWidth(medicine.Unit)
				pdf.CellFormat(cellWidth, constants.PRESC_STD_CELL_HEIGHT, strings.ToLower(medicine.Unit), "", 1, "L", false, 0, "")
			}

			pdf.SetX(startMedicineX)
		}

		pdf.SetFont(constants.PRESC_MF_DOSE_FONT, constants.REGULAR, constants.PRESC_MF_DOSE_FONT_SZ)
		pdf.CellFormat(0, 0.6, strings.ToLower(medicineSet.Mf), "", 1, "L", false, 0, "")

		pdf.SetFont(constants.PRESC_MF_DOSE_FONT, constants.REGULAR, (constants.PRESC_MF_DOSE_FONT_SZ + 1))
		pdf.SetX((constants.PRESC_WIDTH / 2) - 2)
		fractionIdx := strings.Index(medicineSet.Dose, "/")
		if fractionIdx != -1 {
			doseSplit := strings.Split(medicineSet.Dose, "/")

			aDay := strings.Split(doseSplit[0], "x")
			aDay[0] = strings.TrimSpace(aDay[0])
			aDay[1] = strings.TrimSpace(aDay[1])

			doseOne := aDay[0] + " x "

			cellWidth = pdf.GetStringWidth(doseOne) + constants.PRESC_MARGIN
			pdf.CellFormat(cellWidth, constants.PRESC_STD_CELL_HEIGHT, doseOne, "", 0, "L", false, 0, "")

			pdf.SetXY((pdf.GetX() - 0.15), (pdf.GetY() + 0.05))
			pdf.SubWrite(constants.PRESC_STD_CELL_HEIGHT, aDay[1], 7, 5.5, 0, "")

			pdf.SetX(pdf.GetX() - 0.08)
			pdf.CellFormat(pdf.GetStringWidth("/"), constants.PRESC_STD_CELL_HEIGHT, "/", "", 0, "L", false, 0, "")

			pdf.SetX(pdf.GetX() - 0.07)
			pdf.SubWrite(constants.PRESC_STD_CELL_HEIGHT, doseSplit[1], 7, -1, 0, "")

			pdf.SetX(pdf.GetX() + 0.15)
		} else {
			cellWidth = pdf.GetStringWidth(medicineSet.Dose) + constants.PRESC_MARGIN
			pdf.CellFormat(cellWidth, constants.PRESC_STD_CELL_HEIGHT, medicineSet.Dose, "", 0, "L", false, 0, "")
		}

		pdf.SetFont(constants.PRESC_MF_DOSE_FONT, constants.REGULAR, constants.PRESC_MF_DOSE_FONT_SZ)
		cellWidth = pdf.GetStringWidth(medicineSet.SetUnit)
		pdf.CellFormat(cellWidth, constants.PRESC_STD_CELL_HEIGHT, strings.ToLower(medicineSet.SetUnit), "", 0, "L", false, 0, "")

		if medicineSet.ConsumeTime != "" {
			pdf.SetFont(constants.PRESC_MF_DOSE_FONT, constants.REGULAR, constants.PRESC_MF_DOSE_FONT_SZ)
			cellWidth = pdf.GetStringWidth(medicineSet.ConsumeTime) + constants.PRESC_MARGIN
			pdf.CellFormat(cellWidth, constants.PRESC_STD_CELL_HEIGHT, strings.ToLower(medicineSet.ConsumeTime), "", 0, "L", false, 0, "")
		}

		pdf.SetFont(constants.PRESC_MF_DOSE_FONT, constants.REGULAR, constants.PRESC_MF_DOSE_FONT_SZ)
		cellWidth = pdf.GetStringWidth(strings.ToLower(medicineSet.Usage)) + constants.PRESC_MARGIN
		pdf.CellFormat(cellWidth, constants.PRESC_STD_CELL_HEIGHT, strings.ToLower(medicineSet.Usage), "", 1, "L", false, 0, "")

		pdf.Line(0.5, (pdf.GetY() + constants.PRESC_MARGIN), (constants.PRESC_WIDTH - 2), (pdf.GetY() + constants.PRESC_MARGIN))

		pdf.SetX(constants.PRESC_WIDTH - 2.05)
		pdf.SetFont(constants.PRESC_DET_FONT, constants.REGULAR, constants.PRESC_DET_FONT_SZ)
		cellWidth = pdf.GetStringWidth(strings.ToLower(medicineSet.Det)) + constants.PRESC_MARGIN
		if strings.ToLower(medicineSet.Det) == "nedet" {
			pdf.CellFormat(cellWidth, constants.PRESC_STD_CELL_HEIGHT, strings.ToLower(medicineSet.Det), "", 1, "L", false, 0, "")
		} else {
			det := fmt.Sprintf("det: %s", strings.ToLower(medicineSet.Det))
			cellWidth = pdf.GetStringWidth(det) + constants.PRESC_MARGIN
			pdf.CellFormat(cellWidth, constants.PRESC_STD_CELL_HEIGHT, det, "", 1, "L", false, 0, "")
		}

		pdf.SetY(pdf.GetY() + constants.PRESC_MARGIN)
	}

	if pdf.Error() != nil {
		return fmt.Errorf("error create presc info: %v", pdf.Error())
	}

	return nil
}

func createPrescriptionFooter(pdf *fpdf.Fpdf) error {
	pdf.SetLineWidth(0.05)
	pdf.SetDrawColor(constants.GREEN_R, constants.GREEN_G, constants.GREEN_B)

	pdf.SetXY(0.3, 17.5)
	pdf.SetFont("Calibri", constants.BOLD, constants.PRESC_DATA_FONT_SZ)
	pdf.CellFormat(1, 0.8, "H", "1", 1, "LM", false, 0, "")

	pdf.SetFont("Calibri", constants.BOLD, constants.PRESC_DATA_FONT_SZ)
	pdf.CellFormat(1, 0.8, "T", "1", 1, "LM", false, 0, "")

	pdf.SetFont("Calibri", constants.BOLD, constants.PRESC_DATA_FONT_SZ)
	pdf.CellFormat(1, 0.8, "K", "1", 1, "LM", false, 0, "")

	pdf.SetFont("Calibri", constants.BOLD, constants.PRESC_DATA_FONT_SZ)
	pdf.CellFormat(1, 0.8, "P", "1", 1, "LM", false, 0, "")

	pdf.SetXY(9, 20)
	pdf.SetFont("Calibri", constants.REGULAR, constants.PRESC_DATA_FONT_SZ)
	pdf.CellFormat(0, constants.PRESC_STD_CELL_HEIGHT, "PCC", "", 0, "CM", false, 0, "")

	if pdf.Error() != nil {
		return fmt.Errorf("error create presc footer: %v", pdf.Error())
	}

	return nil
}

func removeMedicineFromList(slice []types.PrescriptionMedicineItemReturn, s int) []types.PrescriptionMedicineItemReturn {
	return append(slice[:s], slice[s+1:]...)
}
