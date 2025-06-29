package payslips

import (
	"fmt"
//	"payroll/src/3s"
	"strconv"
	"time"

	"github.com/go-pdf/fpdf"
)

// PayslipPDFData holds data for generating payslip PDF
type PayslipPDFData struct {
  //  User            *user.User
    Period          *PayrollPeriod
    Payslip         *Payslip
    Attendances     []Attendance
    Overtimes       []Overtime
    Reimbursements  []Reimbursement
    SalaryBaseOnAttended float64
}

func GeneratePDF(data PayslipPDFData, outputPath string) error {
    // Initialize PDF
    pdf := fpdf.New("P", "mm", "A4", "")
    pdf.SetMargins(10, 15, 10)
    pdf.AddPage()

    // Set font
    pdf.SetFont("Times", "", 12)

    // Header
    pdf.SetFont("Times", "B", 16)
    pdf.CellFormat(0, 10, "SLIP GAJI", "", 1, "C", false, 0, "")
    pdf.SetFont("Times", "", 12)
    pdf.CellFormat(0, 7, "PT. Dealls Indonesia", "", 1, "C", false, 0, "")
    pdf.Ln(5)
    pdf.SetFont("Times", "", 10)
   // pdf.CellFormat(0, 6, fmt.Sprintf("Employee: %s (ID: %d)", data.User.Name, data.User.ID), "", 1, "L", false, 0, "")
    pdf.CellFormat(0, 6, fmt.Sprintf("Period: %s - %s", data.Period.StartDate.Format("02/01/2006"), data.Period.EndDate.Format("02/01/2006")), "", 1, "L", false, 0, "")
    pdf.Ln(5)

    // Attendance Table
    pdf.SetFont("Times", "B", 10)
    pdf.CellFormat(0, 7, "Attendance Details", "", 1, "L", false, 0, "")
    colWidths := []float64{30, 30, 30, 30} // Date, Check-in, Check-out, Duration
    pdf.SetFillColor(220, 220, 220)
    pdf.CellFormat(colWidths[0], 7, "Date", "1", 0, "C", true, 0, "")
    pdf.CellFormat(colWidths[1], 7, "Check-in", "1", 0, "C", true, 0, "")
    pdf.CellFormat(colWidths[2], 7, "Check-out", "1", 0, "C", true, 0, "")
    pdf.CellFormat(colWidths[3], 7, "Duration (Hours)", "1", 0, "C", true, 0, "")
    pdf.Ln(7)

    pdf.SetFont("Times", "", 10)
    if len(data.Attendances) == 0 {
        pdf.CellFormat(120, 7, "No attendance records", "1", 1, "C", false, 0, "")
    } else {
        for _, att := range data.Attendances {
            checkOut := "N/A"
            duration := "N/A"
            if !att.CheckOut.IsZero() {
                checkOut = att.CheckOut.Format("02/01/2006 15:04")
                if att.CheckIn.Truncate(24*time.Hour).Equal(att.CheckOut.Truncate(24*time.Hour)) {
                    duration = fmt.Sprintf("%.2f", att.CheckOut.Sub(att.CheckIn).Hours())
                }
            }
            pdf.CellFormat(colWidths[0], 7, att.CheckIn.Format("02/01/2006"), "1", 0, "C", false, 0, "")
            pdf.CellFormat(colWidths[1], 7, att.CheckIn.Format("15:04"), "1", 0, "C", false, 0, "")
            pdf.CellFormat(colWidths[2], 7, checkOut, "1", 0, "C", false, 0, "")
            pdf.CellFormat(colWidths[3], 7, duration, "1", 0, "C", false, 0, "")
            pdf.Ln(7)
        }
    }
    pdf.Ln(5)

    // Overtime Table
    pdf.SetFont("Times", "B", 10)
    pdf.CellFormat(0, 7, "Overtime Details", "", 1, "L", false, 0, "")
    pdf.SetFillColor(220, 220, 220)
    pdf.CellFormat(colWidths[0], 7, "Date", "1", 0, "C", true, 0, "")
    pdf.CellFormat(colWidths[1], 7, "Hours", "1", 0, "C", true, 0, "")
    pdf.CellFormat(colWidths[2], 7, "Rate per Hour", "1", 0, "C", true, 0, "")
    pdf.CellFormat(colWidths[3], 7, "Total", "1", 0, "C", true, 0, "")
    pdf.Ln(7)

    pdf.SetFont("Times", "", 10)
    if len(data.Overtimes) == 0 {
        pdf.CellFormat(120, 7, "No overtime records", "1", 1, "C", false, 0, "")
    } else {
        overtimeRate := (13000000.0/ 160.0) * 2
        for _, ot := range data.Overtimes {
            pdf.CellFormat(colWidths[0], 7, ot.OvertimeDate.Format("02/01/2006"), "1", 0, "C", false, 0, "")
            pdf.CellFormat(colWidths[1], 7, strconv.Itoa(ot.Hours), "1", 0, "C", false, 0, "")
            pdf.CellFormat(colWidths[2], 7, FormatRupiah(overtimeRate), "1", 0, "C", false, 0, "")
            pdf.CellFormat(colWidths[3], 7, FormatRupiah(float64(ot.Hours)*overtimeRate), "1", 0, "C", false, 0, "")
            pdf.Ln(7)
        }
    }
    pdf.Ln(5)

    // Reimbursement Table
    pdf.SetFont("Times", "B", 10)
    pdf.CellFormat(0, 7, "Reimbursement Details", "", 1, "L", false, 0, "")
    colWidths = []float64{90, 30} // Description, Amount
    pdf.SetFillColor(220, 220, 220)
    pdf.CellFormat(colWidths[0], 7, "Description", "1", 0, "C", true, 0, "")
    pdf.CellFormat(colWidths[1], 7, "Amount", "1", 0, "C", true, 0, "")
    pdf.Ln(7)

    pdf.SetFont("Times", "", 10)
    if len(data.Reimbursements) == 0 {
        pdf.CellFormat(120, 7, "No reimbursement records", "1", 1, "C", false, 0, "")
    } else {
        for _, r := range data.Reimbursements {
            desc := "N/A"
            if r.Description != nil {
                desc = *r.Description
            }
            rowHeight := Max(7, CalculateMultiCellHeight(pdf, desc, colWidths[0], 4))
            y := pdf.GetY()
            pdf.MultiCell(colWidths[0], 4, desc, "1", "L", false)
            pdf.Rect(10, y, colWidths[0], rowHeight, "D")
            pdf.SetXY(10+colWidths[0], y)
            pdf.CellFormat(colWidths[1], rowHeight, FormatRupiah(r.Amount), "1", 0, "C", false, 0, "")
            pdf.SetXY(10, y+rowHeight)
        }
    }
    pdf.Ln(5)

    // Summary Table
    pdf.SetFont("Times", "B", 10)
    pdf.CellFormat(0, 7, "Summary", "", 1, "L", false, 0, "")
    colWidths = []float64{90, 30}
    pdf.SetFillColor(220, 220, 220)
    pdf.CellFormat(colWidths[0], 7, "Item", "1", 0, "L", true, 0, "")
    pdf.CellFormat(colWidths[1], 7, "Amount", "1", 0, "C", true, 0, "")
    pdf.Ln(7)

    pdf.SetFont("Times", "", 10)
    summaryItems := []struct {
        Item  string
        Amount float64
    }{
        {"Base Salary (Full Month)", data.Payslip.BaseSalary},
        {"Salary Based on Attendance", data.SalaryBaseOnAttended},
        {"Overtime Pay", data.Payslip.OvertimePay},
        {"Reimbursement Pay", data.Payslip.ReimbursementPay},
        {"Total Take-Home Pay", data.Payslip.TotalPay},
    }
    for _, item := range summaryItems {
        pdf.CellFormat(colWidths[0], 7, item.Item, "1", 0, "L", false, 0, "")
        pdf.CellFormat(colWidths[1], 7, FormatRupiah(item.Amount), "1", 0, "C", false, 0, "")
        pdf.Ln(7)
    }

    // Footer
    pdf.Ln(10)
    pdf.SetFont("Times", "", 10)
    pdf.CellFormat(0, 6, fmt.Sprintf("Generated on: %s", time.Now().Format("02/01/2006")), "", 1, "L", false, 0, "")
    pdf.SetFont("Times", "B", 10)
    pdf.CellFormat(0, 6, "HR Department", "", 1, "R", false, 0, "")

    // Save PDF
    err := pdf.OutputFileAndClose(outputPath)
    if err != nil {
        return fmt.Errorf("failed to save PDF: %w", err)
    }

    return nil
}