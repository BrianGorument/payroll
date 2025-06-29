# Sistem Penggajian (Payroll System)

Selamat datang di repositori **Sistem Penggajian**! Ini adalah aplikasi sederhana yang membantu perusahaan menghitung gaji karyawan berdasarkan absensi, lembur, dan penggantian biaya (reimbursement). Aplikasi ini juga bisa membuat slip gaji dalam format PDF yang mudah dibaca, dengan format uang Rupiah (misalnya, `Rp.1.234.567,89`).

Proyek ini dibuat menggunakan bahasa pemrograman **Go** dan cocok untuk Anda yang ingin belajar cara membuat aplikasi berbasis API atau sistem penggajian.

## Apa yang Bisa Dilakukan Aplikasi Ini?
- **Menghitung Gaji**: Menghitung gaji berdasarkan hari kerja (minimal 8 jam sehari), lembur, dan penggantian biaya.
- **Membuat Slip Gaji**: Menghasilkan file PDF untuk slip gaji dengan rincian absensi, lembur, dan total gaji.
- **Akses Berbasis Peran**: Admin bisa menghitung gaji dan membuat slip gaji untuk karyawan tertentu, sementara karyawan hanya bisa melihat slip gaji mereka sendiri.

## Prasyarat
Untuk menjalankan aplikasi ini, Anda perlu:
1. **Go** (versi 1.22 atau lebih baru). Unduh di [golang.org](https://golang.org/dl/).
2. **Database MySQL** atau database lain yang kompatibel dengan GORM. Anda bisa mengunduh MySQL di [mysql.com](https://www.mysql.com/downloads/).
3. **Git** untuk mengunduh kode dari GitHub. Unduh di [git-scm.com](https://git-scm.com/downloads).
4. Komputer dengan sistem operasi Windows, macOS, atau Linux.
5. Alat untuk menguji API, seperti **Postman** atau **curl** (Postman bisa diunduh di [postman.com](https://www.postman.com/downloads/)).

## Cara Menjalankan Proyek
Ikuti langkah-langkah berikut untuk menjalankan aplikasi di komputer Anda:

### 1. Unduh Kode dari GitHub
1. Buka terminal (Command Prompt di Windows, Terminal di macOS/Linux).
2. Kloning repositori ini:
   ```bash
   git clone https://github.com/<username-anda>/payroll-system.git
   ```
   Ganti `<username-anda>` dengan nama pengguna GitHub Anda.
3. Masuk ke folder proyek:
   ```bash
   cd payroll-system
   ```

### 2. Siapkan Database
1. Buat database MySQL bernama `payroll_db`:
   - Buka MySQL di terminal atau alat seperti phpMyAdmin.
   - Jalankan perintah:
     ```sql
     CREATE DATABASE payroll_db;
     ```
2. Tambahkan data awal (contoh):
   - Masukkan data ke tabel `users`, `payslips`, `attendances`, `overtimes`, dan `reimbursements`. Contoh:
     ```sql
     INSERT INTO users (id, name, salary, role, created_at, updated_at, created_by, updated_by) VALUES
     (2, 'UserA', 13600774.80, 'employee', NOW(), NOW(), 'admin', 'admin'),
     (5, 'UserB', 14419133.79, 'employee', NOW(), NOW(), 'admin', 'admin');

     INSERT INTO payroll_periods (id, start_date, end_date, is_processed) VALUES
     (3, '2025-07-01', '2025-07-31', 0);

     INSERT INTO attendances (user_id, payroll_period_id, check_in, check_out) VALUES
     (2, 3, '2025-07-30 08:00:00', '2025-07-30 16:00:00');

     INSERT INTO overtimes (user_id, payroll_period_id, hours, overtime_date) VALUES
     (2, 3, 4, '2025-07-30');

     INSERT INTO reimbursements (user_id, payroll_period_id, amount, description) VALUES
     (2, 3, 1500.00, 'Travel expense');
     ```

### 3. Instal Dependensi
1. Pastikan Anda berada di folder proyek (`payroll-system`).
2. Jalankan perintah untuk mengunduh dependensi:
   ```bash
   go mod tidy
   ```
   Ini akan mengunduh library seperti `github.com/go-pdf/fpdf` untuk membuat PDF dan `golang.org/x/text` untuk format Rupiah.

### 4. Konfigurasi Database
1. Buat file konfigurasi (misalnya, `config.yaml`) untuk menghubungkan aplikasi ke database:
   ```yaml
   database:
     host: localhost
     port: 3306
     user: root
     password: <password-anda>
     name: payroll_db
   ```
   Ganti `<password-anda>` dengan kata sandi MySQL Anda.
2. Pastikan aplikasi Anda membaca konfigurasi ini (biasanya di `main.go` atau file serupa).

### 5. Jalankan Aplikasi
1. Jalankan aplikasi dengan perintah:
   ```bash
   go run main.go
   ```
   Ganti `main.go` dengan file utama proyek Anda jika bernama lain.
2. Aplikasi akan berjalan di `http://localhost:8080` (atau port lain sesuai konfigurasi).

### 6. Uji API
Gunakan Postman atau curl untuk menguji API:
- **Menghitung Gaji**:
  - Kirim permintaan POST ke `http://localhost:8080/v1/payslips/run`:
    ```json
    {
        "payroll_period_id": "3"
    }
    ```
  - Anda akan mendapatkan respons JSON dengan rincian gaji, misalnya:
    ```json
    [
        {
            "id": "1",
            "user_id": "2",
            "payroll_period_id": "3",
            "base_salary": "Rp.13.600.775,00",
            "salary_base_on_attended": "Rp.3.400.194,00",
            "overtime_pay": "Rp.340.019,00",
            "reimbursement_pay": "Rp.1.500,00",
            "total_pay": "Rp.3.741.713,00"
        }
    ]
    ```
- **Membuat Slip Gaji PDF**:
  - Kirim permintaan POST ke `http://localhost:8080/v1/payslips/generate`:
    ```json
    {
        "payroll_period_id": "3",
        "user_id": "2"
    }
    ```
  - File PDF akan disimpan di folder `tmp/payslips/` (misalnya, `payslip_2_3.pdf`).

### 7. Lihat Slip Gaji
- Buka folder `tmp/payslips/` di komputer Anda.
- Buka file PDF (misalnya, `payslip_2_3.pdf`) untuk melihat slip gaji dengan rincian:
  - Nama karyawan dan periode gaji.
  - Tabel absensi (tanggal, check-in, check-out, durasi).
  - Tabel lembur (tanggal, jam, tarif, total).
  - Tabel penggantian biaya (deskripsi, jumlah).
  - Ringkasan gaji dalam format Rupiah (misalnya, `Rp.3.741.713,00`).
  - Tanda tangan "HR Department".

## Struktur Direktori
Berikut adalah struktur folder proyek ini:
```
payroll-system/
├── go.mod                # File konfigurasi Go module
├── src/
│   ├── payslips/
│   │   ├── dto.go        # Definisi struktur data untuk API
│   │   ├── pdf.go        # Logika pembuatan PDF slip gaji
│   │   ├── service.go    # Logika utama untuk menghitung dan membuat slip gaji
│   │   ├── utils.go      # Fungsi pembantu (format Rupiah, dll.)
│   │   └── ...
│   ├── users/
│   │   ├── user.go       # Definisi data karyawan
│   │   └── ...
├── tmp/payslips/         # Folder untuk menyimpan file PDF
└── main.go               # File utama untuk menjalankan aplikasi
```

## Contoh Slip Gaji
Slip gaji dalam PDF akan terlihat seperti ini:
```
PAYSLIP
PT. Payroll Indonesia
Employee: UserA (ID: 2)
Period: 01/07/2025 - 31/07/2025

Attendance Details
| Date       | Check-in | Check-out       | Duration (Hours) |
|------------|----------|-----------------|------------------|
| 30/07/2025 | 08:00    | 30/07/2025 16:00 | 8.00             |

Overtime Details
| Date       | Hours | Rate per Hour    | Total            |
|------------|-------|------------------|------------------|
| 30/07/2025 | 4     | Rp.170.010,00    | Rp.680.040,00    |

Reimbursement Details
| Description      | Amount        |
|------------------|---------------|
| Travel expense   | Rp.1.500,00   |

Summary
| Item                     | Amount           |
|--------------------------|------------------|
| Base Salary (Full Month) | Rp.13.600.775,00 |
| Salary Based on Attendance | Rp.3.400.194,00 |
| Overtime Pay             | Rp.340.019,00    |
| Reimbursement Pay        | Rp.1.500,00      |
| Total Take-Home Pay      | Rp.3.741.713,00  |

Generated on: 29/06/2025
[Right-aligned] HR Department
```

## Bantuan
Jika Anda mengalami kesulitan:
- Periksa pesan error di terminal dan cari solusi di [Stack Overflow](https://stackoverflow.com).
- Hubungi saya melalui [GitHub Issues](https://github.com/<username-anda>/payroll-system/issues) atau email di `<email-anda>`.
- Pastikan Go, MySQL, dan dependensi sudah terinstal dengan benar.

Selamat mencoba, dan semoga proyek ini bermanfaat!