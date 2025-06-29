package payslips

import (
	"bytes"
	"errors"
	"html/template"
	"time"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func GeneratePayslipHTML(user *User, period *PayrollPeriod, payslip *Payslip, attendances []Attendance, overtimes []Overtime, reimbursements []Reimbursement) (string, error) {
    // Format Rupiah
    p := message.NewPrinter(language.Indonesian)
    formatRupiah := func(value float64) string {
        return p.Sprintf("Rp%.2f", value)
    }

    // Template HTML
    const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1, h2 { text-align: center; }
        .header { margin-bottom: 20px; }
        table { width: 100%; border-collapse: collapse; margin-bottom: 20px; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
        .total { font-weight: bold; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Payslip</h1>
        <p>PT. Payroll</p>
        <p>Employee: {{.User.Name}} (ID: {{.User.ID}})</p>
        <p>Period: {{.Period.StartDate.Format "02/01/2006"}} - {{.Period.EndDate.Format "02/01/2006"}}</p>
    </div>

    <h2>Attendance Details</h2>
    <table>
        <tr>
            <th>Date</th>
            <th>Check-in</th>
            <th>Check-out</th>
            <th>Duration (Hours)</th>
        </tr>
        {{range .Attendances}}
        <tr>
            <td>{{.CheckIn.Format "02/01/2006"}}</td>
            <td>{{.CheckIn.Format "15:04"}}</td>
            <td>{{if .CheckOut.IsZero}}N/A{{else}}{{.CheckOut.Format "15:04"}}{{end}}</td>
            <td>{{if .CheckOut.IsZero}}N/A{{else}}{{printf "%.2f" (div (sub .CheckOut .CheckIn) 3600)}}{{end}}</td>
        </tr>
        {{else}}
        <tr><td colspan="4">No attendance records</td></tr>
        {{end}}
    </table>

    <h2>Overtime Details</h2>
    <table>
        <tr>
            <th>Date</th>
            <th>Hours</th>
            <th>Rate per Hour</th>
            <th>Total</th>
        </tr>
        {{range .Overtimes}}
        <tr>
            <td>{{.OvertimeDate.Format "02/01/2006"}}</td>
            <td>{{.Hours}}</td>
            <td>{{$.FormatRupiah (mul (div $.User.BaseSalary 160) 2)}}</td>
            <td>{{$.FormatRupiah (mul .Hours (mul (div $.User.BaseSalary 160) 2))}}</td>
        </tr>
        {{else}}
        <tr><td colspan="4">No overtime records</td></tr>
        {{end}}
    </table>

    <h2>Reimbursement Details</h2>
    <table>
        <tr>
            <th>Description</th>
            <th>Amount</th>
        </tr>
        {{range .Reimbursements}}
        <tr>
            <td>{{if .Description}}{{.Description}}{{else}}N/A{{end}}</td>
            <td>{{$.FormatRupiah .Amount}}</td>
        </tr>
        {{else}}
        <tr><td colspan="2">No reimbursement records</td></tr>
        {{end}}
    </table>

    <h2>Summary</h2>
    <table>
        <tr><td>Base Salary</td><td class="total">{{$.FormatRupiah .Payslip.BaseSalary}}</td></tr>
        <tr><td>Overtime Pay</td><td class="total">{{$.FormatRupiah .Payslip.OvertimePay}}</td></tr>
        <tr><td>Reimbursement Pay</td><td class="total">{{$.FormatRupiah .Payslip.ReimbursementPay}}</td></tr>
        <tr><td>Total Take-Home Pay</td><td class="total">{{$.FormatRupiah .Payslip.TotalPay}}</td></tr>
    </table>
</body>
</html>
`

    // Data untuk template
    data := struct {
        User          *User
        Period        *PayrollPeriod
        Payslip       *Payslip
        Attendances   []Attendance
        Overtimes     []Overtime
        Reimbursements []Reimbursement
        FormatRupiah  func(float64) string
    }{
        User:          user,
        Period:        period,
        Payslip:       payslip,
        Attendances:   attendances,
        Overtimes:     overtimes,
        Reimbursements: reimbursements,
        FormatRupiah:  formatRupiah,
    }

    // Parse dan render template
    t, err := template.New("payslip").Funcs(template.FuncMap{
        "div": func(a, b float64) float64 { return a / b },
        "sub": func(a, b time.Time) float64 { return a.Sub(b).Seconds() },
        "mul": func(a, b float64) float64 { return a * b },
    }).Parse(htmlTemplate)
    if err != nil {
        return "", errors.New("failed to parse HTML template")
    }

    var htmlContent bytes.Buffer
    if err := t.Execute(&htmlContent, data); err != nil {
        return "", errors.New("failed to render HTML template")
    }

    return htmlContent.String(), nil
}

func GeneratePDF(htmlContent, outputPath string) error {
    // Buat generator PDF
    pdfg, err := wkhtmltopdf.NewPDFGenerator()
    if err != nil {
        return errors.New("failed to create PDF generator")
    }

    // Tambahkan halaman dari HTML
    pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader([]byte(htmlContent))))
    pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
    pdfg.MarginTop.Set(10)
    pdfg.MarginBottom.Set(10)
    pdfg.MarginLeft.Set(10)
    pdfg.MarginRight.Set(10)

    // Tulis ke file
    err = pdfg.WriteFile(outputPath)
    if err != nil {
        return errors.New("failed to write PDF file")
    }

    return nil
}