package payslips

import (
	"regexp"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"github.com/go-pdf/fpdf"
)


type TerbilangConverter struct {
	dasar  []string
	angka  []int
	satuan []string
	reBelas *regexp.Regexp // Regex untuk "satu puluh X" menjadi "X belas"
	reSe    *regexp.Regexp // Regex untuk "satu [unit]" menjadi "se[unit]"
}


func FormatRupiah(value float64) string {
    p := message.NewPrinter(language.Indonesian)
    rounded := float64(int(value*100+0.5)) / 100
    return p.Sprintf("Rp.%.2f", rounded)
}



// NewTerbilangConverter membuat instance baru dari TerbilangConverter.
func NewTerbilangConverter() *TerbilangConverter {
	return &TerbilangConverter{
		// Indeks 0 kosong agar angka dasar sesuai dengan indeks 1-9
		dasar:  []string{"", "satu", "dua", "tiga", "empat", "lima", "enam", "tujuh", "delapan", "sembilan"},
		// Urutan dari terbesar ke terkecil: milyar, juta, ribu, ratus, puluh, satuan
		angka:  []int{1000000000, 1000000, 1000, 100, 10, 1},
		satuan: []string{"milyar", "juta", "ribu", "ratus", "puluh", ""},
		// Kompilasi regex dengan opsi case-insensitive (?i)
		reBelas: regexp.MustCompile(`(?i)satu puluh (\w+)`),
		reSe:    regexp.MustCompile(`(?i)satu (ribu|ratus|puluh|belas)`),
	}
}

func (tc *TerbilangConverter) Convert(n int) string {
	if n == 0 {
		return "nol"
	}

	if n < 0 {
		return "minus " + tc.Convert(-n) // Tangani angka negatif
	}

	var result string
	tempN := n // Gunakan variabel sementara untuk angka yang akan diproses

	// Iterasi melalui array `angka` dan `satuan` dari unit terbesar ke terkecil
	for i := 0; i < len(tc.angka); i++ {
		divider := tc.angka[i]
		unit := tc.satuan[i]

		// Hindari pembagian oleh nol jika divider menjadi 0 (untuk kasus satuan kosong di akhir)
		if divider == 0 {
			continue
		}

		count := tempN / divider // Hitung berapa kali 'divider' masuk ke 'tempN'

		if count > 0 {
			// Jika 'count' lebih dari atau sama dengan 10, panggil fungsi Convert secara rekursif
			// untuk mengubah 'count' itu sendiri ke terbilang (misal: "seratus" untuk 100)
			if count >= 10 {
				result += tc.Convert(count) + " " + unit + " "
			} else { // Jika 'count' antara 1 sampai 9
				result += tc.dasar[count] + " " + unit + " "
			}
		}
		tempN %= divider // Kurangi 'tempN' dengan bagian yang sudah diproses
	}

	// Hapus spasi di awal dan akhir string
	result = strings.TrimSpace(result)

	//Mengubah "satu puluh [kata]" menjadi "[kata] belas"
	result = tc.reBelas.ReplaceAllString(result, "${1} belas")

    //Mengubah "satu [ribu|ratus|puluh|belas]" menjadi "se[ribu|ratus|puluh|belas]"
	result = tc.reSe.ReplaceAllString(result, "se${1}")

	return result
}

func GetTerbilang(n int) string {
	converter := NewTerbilangConverter()
	return converter.Convert(n)
}

func TransformString(input string) string {
	transformed := strings.ReplaceAll(input, " ", "_")
	transformed = strings.ReplaceAll(transformed, "/", "|")
	return transformed
}

func GetBulan(month int) string {
	switch month {
	case 1:
		return "Januari"
	case 2:
		return "Februari"
	case 3:
		return "Maret"
	case 4:
		return "April"
	case 5:
		return "Mei"
	case 6:
		return "Juni"
	case 7:
		return "Juli"
	case 8:
		return "Agustus"
	case 9:
		return "September"
	case 10:
		return "Oktober"
	case 11:
		return "November"
	case 12:
		return "Desember"
	default:
		return "" // Mengembalikan string kosong jika 	bulan tidak valid
	}
}

func GetHari(Day string) string {
	switch Day {
	case "Monday":
		return "Senin"
	case "Tuesday":
		return "Selasa"
	case "Wednesday":
		return "Rabu"
	case "Thursday":
		return "Kamis"
	case "Friday":
		return "Jumat"
	case "Saturday":
		return "Sabtu"
	case "Sunday":
		return "Minggu"
	}
	return ""
}


// CalculateMultiCellHeight calculates the height needed for a MultiCell based on text content
func CalculateMultiCellHeight(pdf *fpdf.Fpdf, text string, width, lineHeight float64) float64 {
    if text == "" {
        return lineHeight
    }
    words := strings.Split(text, " ")
    line := ""
    count := 0
    for _, word := range words {
        if pdf.GetStringWidth(line+word+" ") < width-2 {
            line += word + " "
        } else {
            count++
            line = word + " "
        }
    }
    if line != "" {
        count++
    }
    return float64(count) * lineHeight
}

// Max returns the maximum value from a slice of float64
func Max(values ...float64) float64 {
    maxVal := values[0]
    for _, v := range values[1:] {
        if v > maxVal {
            maxVal = v
        }
    }
    return maxVal
}