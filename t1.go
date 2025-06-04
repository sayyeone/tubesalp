package main

import (
	"fmt"
	"os"
	"strings"
)

// KONSTANTA
const NMAX = 100
const MAX_SUBJECTS = 3
const MAX_TRYOUTS = 3

// TIPE BENTUKAN 1: userLes
type userLes struct {
	nama     string
	status   string
	username string
	email    string
	password string
}

// TIPE BENTUKAN 2: statusLogin
type statusLogin struct {
	statusUser    string
	apakahLogin   bool
	userYangLogin userLes
}

// TIPE BENTUKAN 3: mataPelajaran
type mataPelajaran struct {
	nama        string
	nilaiTryOut [MAX_TRYOUTS]float64
	nilaiRata   float64
}

// TIPE BENTUKAN 4: jadwalLes
type jadwalLes struct {
	day  string
	time string
}

// TIPE BENTUKAN 5: siswa
type siswa struct {
	id            int
	nama          string
	kelas         string
	username      string
	status        string
	userLes       userLes
	mataPelajaran [MAX_SUBJECTS]mataPelajaran
	jadwalLes     [MAX_SUBJECTS]jadwalLes
	telepon       string
	email         string
	tanggalLahir  string
	catatan       string
	totalNilai    float64
	ranking       int
}

// TIPE BENTUKAN 6: admin
type admin struct {
	nama     string
	username string
	status   string
	email    string
	password string
}

// ARRAY GLOBAL 1: Data Siswa
var globalSiswa [NMAX]siswa
var nSiswa int

// ARRAY GLOBAL 2: Data Admin
var globalAdmin [NMAX]admin
var nAdmin int

// ARRAY GLOBAL 3: Data Login Status
var globalLogin [1]statusLogin

// ============================================================================
// FUNGSI-FUNGSI UTILITAS (FUNCTIONS)
// ============================================================================

// FUNGSI 1: Normalisasi username untuk konsistensi
// I.S.: username terdefinisi
// F.S.: mengembalikan username dalam lowercase dan tanpa spasi
func normalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

// FUNGSI 2: Validasi format email
// I.S.: email terdefinisi
// F.S.: mengembalikan true jika email valid, false jika tidak
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".") && len(email) > 5
}

// FUNGSI 3: Validasi format tanggal (dd/mm/yyyy)
// I.S.: date terdefinisi
// F.S.: mengembalikan true jika format valid, false jika tidak
func isValidDate(date string) bool {
	if len(date) != 10 {
		return false
	}
	if strings.Count(date, "/") != 2 {
		return false
	}
	parts := strings.Split(date, "/")
	return len(parts) == 3 && len(parts[0]) == 2 && len(parts[1]) == 2 && len(parts[2]) == 4
}

// FUNGSI 4: Validasi kekuatan password
// I.S.: password terdefinisi
// F.S.: mengembalikan true jika password cukup kuat, false jika tidak
func isValidPassword(password string) bool {
	return len(password) >= 6
}

// FUNGSI 5: Hitung rata-rata nilai mata pelajaran
// I.S.: mata pelajaran terdefinisi dengan nilai tryout
// F.S.: mengembalikan nilai rata-rata dari tryout yang ada
func hitungRataMataPelajaran(mp mataPelajaran) float64 {
	var total float64
	var count int
	for i := 0; i < MAX_TRYOUTS; i++ {
		if mp.nilaiTryOut[i] > 0 {
			total += mp.nilaiTryOut[i]
			count++
		}
	}
	if count > 0 {
		return total / float64(count)
	}
	return 0.0
}

// FUNGSI 6: Hitung total nilai siswa
// I.S.: siswa terdefinisi dengan mata pelajaran
// F.S.: mengembalikan rata-rata dari semua mata pelajaran
func hitungTotalNilaiSiswa(s siswa) float64 {
	var total float64
	var count int
	for i := 0; i < MAX_SUBJECTS; i++ {
		if s.mataPelajaran[i].nilaiRata > 0 {
			total += s.mataPelajaran[i].nilaiRata
			count++
		}
	}
	if count > 0 {
		return total / float64(count)
	}
	return 0.0
}

// FUNGSI 7: Pencarian Linear - Cari siswa berdasarkan username
// I.S.: username terdefinisi
// F.S.: mengembalikan index siswa jika ditemukan, -1 jika tidak
func linearSearchSiswaByUsername(username string) int {
	username = normalizeUsername(username)
	for i := 0; i < nSiswa; i++ {
		if normalizeUsername(globalSiswa[i].username) == username {
			return i
		}
	}
	return -1
}

// FUNGSI 8: Pencarian Binary - Cari siswa berdasarkan ID (array harus terurut)
// I.S.: targetID terdefinisi, array siswa terurut berdasarkan ID
// F.S.: mengembalikan index siswa jika ditemukan, -1 jika tidak
func binarySearchSiswaByID(targetID int) int {
	if targetID <= 0 || nSiswa == 0 {
		return -1
	}

	left := 0
	right := nSiswa - 1

	for left <= right {
		mid := left + (right-left)/2
		if globalSiswa[mid].id == targetID {
			return mid
		}
		if globalSiswa[mid].id < targetID {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}

// FUNGSI 9: Cek apakah username sudah ada
// I.S.: username terdefinisi
// F.S.: mengembalikan kode unik jika username ditemukan, 0 jika tidak
func cekUsernameExists(username string) int {
	username = normalizeUsername(username)

	// Cek di admin (return 1000000 + index)
	for i := 0; i < nAdmin; i++ {
		if normalizeUsername(globalAdmin[i].username) == username {
			return 1000000 + i
		}
	}

	// Cek di siswa (return 2000000 + index)
	for i := 0; i < nSiswa; i++ {
		if normalizeUsername(globalSiswa[i].username) == username {
			return 2000000 + i
		}
	}

	return 0
}

// ============================================================================
// PROSEDUR-PROSEDUR UTILITAS (PROCEDURES)
// ============================================================================

// PROSEDUR 1: Input string dengan spasi (berakhir dengan "STOP")
// I.S.: pointer ke string terdefinisi
// F.S.: string terisi dengan input user hingga kata "STOP"
func inputDenganSpasi(result *string) {
	var kata string
	var kalimat string
	wordCount := 0
	const maxWords = 50

	fmt.Scan(&kata)
	for kata != "STOP" && wordCount < maxWords {
		if kalimat != "" {
			kalimat += " "
		}
		kalimat += kata
		wordCount++
		fmt.Scan(&kata)
	}
	*result = strings.TrimSpace(kalimat)
}

// PROSEDUR 2: Selection Sort siswa berdasarkan ID (ascending)
// I.S.: array siswa terdefinisi
// F.S.: array siswa terurut berdasarkan ID secara ascending
func selectionSortSiswaByID() {
	for i := 0; i < nSiswa-1; i++ {
		minIdx := i
		for j := i + 1; j < nSiswa; j++ {
			if globalSiswa[j].id < globalSiswa[minIdx].id {
				minIdx = j
			}
		}
		if minIdx != i {
			temp := globalSiswa[i]
			globalSiswa[i] = globalSiswa[minIdx]
			globalSiswa[minIdx] = temp
		}
	}
}

// PROSEDUR 3: Insertion Sort siswa berdasarkan total nilai (descending untuk ranking)
// I.S.: array siswa terdefinisi
// F.S.: array siswa terurut berdasarkan total nilai secara descending dan ranking di-update
func insertionSortSiswaByNilai() {
	// Update total nilai semua siswa
	for i := 0; i < nSiswa; i++ {
		for j := 0; j < MAX_SUBJECTS; j++ {
			globalSiswa[i].mataPelajaran[j].nilaiRata = hitungRataMataPelajaran(globalSiswa[i].mataPelajaran[j])
		}
		globalSiswa[i].totalNilai = hitungTotalNilaiSiswa(globalSiswa[i])
	}

	// Insertion sort descending
	for i := 1; i < nSiswa; i++ {
		key := globalSiswa[i]
		j := i - 1
		for j >= 0 && globalSiswa[j].totalNilai < key.totalNilai {
			globalSiswa[j+1] = globalSiswa[j]
			j--
		}
		globalSiswa[j+1] = key
	}

	// Update ranking
	for i := 0; i < nSiswa; i++ {
		globalSiswa[i].ranking = i + 1
	}
}

// PROSEDUR 5: Tampilkan daftar username siswa
// I.S.: array siswa terdefinisi
// F.S.: daftar username dan nama siswa ditampilkan
func tampilkanDaftarSiswa() {
	if nSiswa == 0 {
		fmt.Println("❌ Belum ada siswa yang terdaftar.")
		return
	}

	fmt.Println("\nDaftar Siswa Terdaftar:")
	fmt.Println("========================")
	for i := 0; i < nSiswa; i++ {
		fmt.Printf("%d. Username: %-15s | Nama: %s\n",
			i+1, globalSiswa[i].username, globalSiswa[i].nama)
	}
	fmt.Println("========================")
}

// PROSEDUR 6: Registrasi user baru
// I.S.: -
// F.S.: user baru terdaftar dalam sistem (admin atau siswa)
func registrasi() {
	var inputStatus string
	var usernameBaru string
	var emailBaru string
	var passwordBaru string
	var namaBaru string

	fmt.Println("\n============ REGISTRASI ============")
	fmt.Print("Apa status anda? (Admin/Siswa): ")
	fmt.Scan(&inputStatus)

	// Validasi status
	for inputStatus != "Admin" && inputStatus != "Siswa" {
		fmt.Print("❌ Status salah! Masukkan 'Admin' atau 'Siswa': ")
		fmt.Scan(&inputStatus)
	}

	if inputStatus == "Admin" {
		// Cek kapasitas
		if nAdmin >= NMAX {
			fmt.Printf("❌ Error: Maksimal %d admin sudah tercapai!\n", NMAX)
			return
		}

		// Input nama
		fmt.Print("Masukkan nama lengkap Anda (akhiri dengan 'STOP'): ")
		inputDenganSpasi(&namaBaru)

		// Validasi nama tidak kosong
		for strings.TrimSpace(namaBaru) == "" {
			fmt.Print("❌ Nama tidak boleh kosong! Masukkan nama (akhiri dengan 'STOP'): ")
			inputDenganSpasi(&namaBaru)
		}

		// Input username
		fmt.Print("Masukkan username Anda: ")
		fmt.Scan(&usernameBaru)
		usernameBaru = normalizeUsername(usernameBaru)

		// Validasi username unik
		for cekUsernameExists(usernameBaru) != 0 {
			fmt.Print("❌ Username sudah ada. Masukkan username lain: ")
			fmt.Scan(&usernameBaru)
			usernameBaru = normalizeUsername(usernameBaru)
		}

		// Input email
		fmt.Print("Masukkan email Anda: ")
		fmt.Scan(&emailBaru)

		// Validasi email
		for !isValidEmail(emailBaru) {
			fmt.Print("❌ Format email tidak valid! Masukkan email yang benar: ")
			fmt.Scan(&emailBaru)
		}

		// Input password
		fmt.Print("Masukkan password Anda (minimal 6 karakter): ")
		fmt.Scan(&passwordBaru)

		// Validasi password
		for !isValidPassword(passwordBaru) {
			fmt.Print("❌ Password terlalu pendek! Masukkan password minimal 6 karakter: ")
			fmt.Scan(&passwordBaru)
		}

		// Simpan data admin
		globalAdmin[nAdmin].nama = namaBaru
		globalAdmin[nAdmin].username = usernameBaru
		globalAdmin[nAdmin].email = emailBaru
		globalAdmin[nAdmin].password = passwordBaru
		globalAdmin[nAdmin].status = "Admin"
		nAdmin++

	} else { // Siswa
		// Cek kapasitas
		if nSiswa >= NMAX {
			fmt.Printf("❌ Error: Maksimal %d siswa sudah tercapai!\n", NMAX)
			return
		}

		// Input nama
		fmt.Print("Masukkan nama lengkap anda (akhiri dengan 'STOP'): ")
		inputDenganSpasi(&namaBaru)

		// Validasi nama tidak kosong
		for strings.TrimSpace(namaBaru) == "" {
			fmt.Print("❌ Nama tidak boleh kosong! Masukkan nama (akhiri dengan 'STOP'): ")
			inputDenganSpasi(&namaBaru)
		}

		// Input username
		fmt.Print("Masukkan username Anda: ")
		fmt.Scan(&usernameBaru)
		usernameBaru = normalizeUsername(usernameBaru)

		// Validasi username unik
		for cekUsernameExists(usernameBaru) != 0 {
			fmt.Print("❌ Username sudah ada. Masukkan username lain: ")
			fmt.Scan(&usernameBaru)
			usernameBaru = normalizeUsername(usernameBaru)
		}

		// Input email
		fmt.Print("Masukkan email Anda: ")
		fmt.Scan(&emailBaru)

		// Validasi email
		for !isValidEmail(emailBaru) {
			fmt.Print("❌ Format email tidak valid! Masukkan email yang benar: ")
			fmt.Scan(&emailBaru)
		}

		// Input password
		fmt.Print("Masukkan password Anda (minimal 6 karakter): ")
		fmt.Scan(&passwordBaru)

		// Validasi password
		for !isValidPassword(passwordBaru) {
			fmt.Print("❌ Password terlalu pendek! Masukkan password minimal 6 karakter: ")
			fmt.Scan(&passwordBaru)
		}

		// Simpan data siswa
		globalSiswa[nSiswa].nama = namaBaru
		globalSiswa[nSiswa].username = usernameBaru
		globalSiswa[nSiswa].email = emailBaru
		globalSiswa[nSiswa].status = "Siswa"
		globalSiswa[nSiswa].userLes.nama = namaBaru
		globalSiswa[nSiswa].userLes.username = usernameBaru
		globalSiswa[nSiswa].userLes.email = emailBaru
		globalSiswa[nSiswa].userLes.password = passwordBaru
		globalSiswa[nSiswa].userLes.status = "Siswa"
		nSiswa++
	}

	fmt.Println("✅ Registrasi berhasil! Silahkan login sekarang.")
}

// PROSEDUR 7: Login user
// I.S.: user sudah registrasi
// F.S.: user berhasil login dan status login terupdate
func login() {
	var inputUsername, inputPassword string
	var userCode int
	var loginBerhasil bool
	const maxAttempts = 3
	var attempts int

	fmt.Println("\n============= LOG IN ===============")

	// Input username dengan retry
	attempts = 0
	loginBerhasil = false
	for attempts < maxAttempts && !loginBerhasil {
		fmt.Print("Masukkan username anda: ")
		fmt.Scan(&inputUsername)
		inputUsername = normalizeUsername(inputUsername)
		userCode = cekUsernameExists(inputUsername)

		if userCode != 0 {
			loginBerhasil = true
		} else {
			attempts++
			fmt.Printf("❌ Username tidak ditemukan (percobaan %d/%d)\n", attempts, maxAttempts)
		}
	}

	if !loginBerhasil {
		fmt.Println("❌ Maksimal percobaan username tercapai!")
		fmt.Print("Ingin registrasi? (Y/N): ")
		var choice string
		fmt.Scan(&choice)
		if choice == "Y" || choice == "y" {
			registrasi()
		}
		return
	}

	// Input password dengan retry
	attempts = 0
	loginBerhasil = false
	for attempts < maxAttempts && !loginBerhasil {
		fmt.Print("Masukkan password anda: ")
		fmt.Scan(&inputPassword)

		// Cek password berdasarkan tipe user
		if userCode >= 2000000 { // Siswa
			studentIdx := userCode - 2000000
			if inputPassword == globalSiswa[studentIdx].userLes.password {
				// Setup login data untuk siswa
				globalLogin[0].statusUser = "Siswa"
				globalLogin[0].apakahLogin = true
				globalLogin[0].userYangLogin = globalSiswa[studentIdx].userLes
				loginBerhasil = true
			}
		} else { // Admin
			adminIdx := userCode - 1000000
			if inputPassword == globalAdmin[adminIdx].password {
				// Setup login data untuk admin
				globalLogin[0].statusUser = "Admin"
				globalLogin[0].apakahLogin = true
				globalLogin[0].userYangLogin = userLes{
					nama:     globalAdmin[adminIdx].nama,
					status:   globalAdmin[adminIdx].status,
					username: globalAdmin[adminIdx].username,
					email:    globalAdmin[adminIdx].email,
					password: globalAdmin[adminIdx].password,
				}
				loginBerhasil = true
			}
		}

		if !loginBerhasil {
			attempts++
			fmt.Printf("❌ Password salah (percobaan %d/%d)\n", attempts, maxAttempts)
		}
	}

	if !loginBerhasil {
		fmt.Println("❌ Terlalu banyak percobaan password salah!")
		return
	}

	fmt.Println("✅ Login berhasil!")
	fmt.Printf("Selamat datang %s pada apps Brained\n", globalLogin[0].userYangLogin.nama)
	fmt.Printf("Status anda adalah %s\n", globalLogin[0].userYangLogin.status)

	// Redirect ke menu yang sesuai
	if globalLogin[0].statusUser == "Admin" {
		menuAdmin()
	} else if globalLogin[0].statusUser == "Siswa" {
		menuSiswa()
	}
}

// PROSEDUR 8: Input data siswa (untuk admin)
// I.S.: admin sudah login, siswa sudah terdaftar
// F.S.: data detail siswa terinput (ID, kelas, mata pelajaran, jadwal, dll)
func inputDataSiswa() {
	if nSiswa == 0 {
		fmt.Println("❌ Belum ada siswa yang terdaftar!")
		return
	}

	var username string
	var siswaIdx int
	var idBaru int
	var kelasBaru, teleponBaru, tanggalLahirBaru string
	var mataPelajaranBaru string
	var jadwalHari, jadwalJam string

	tampilkanDaftarSiswa()

	fmt.Print("Masukkan username siswa yang akan diinput datanya: ")
	fmt.Scan(&username)

	siswaIdx = linearSearchSiswaByUsername(username)
	if siswaIdx == -1 {
		fmt.Println("❌ Username siswa tidak ditemukan!")
		return
	}

	fmt.Printf("✅ Siswa ditemukan: %s\n", globalSiswa[siswaIdx].nama)

	// Input ID dengan validasi
	fmt.Print("Masukkan ID siswa: ")
	fmt.Scan(&idBaru)

	// Cek ID unik
	idSudahAda := false
	for i := 0; i < nSiswa; i++ {
		if i != siswaIdx && globalSiswa[i].id == idBaru {
			idSudahAda = true
		}
	}

	for idSudahAda {
		fmt.Print("❌ ID sudah digunakan! Masukkan ID lain: ")
		fmt.Scan(&idBaru)
		idSudahAda = false
		for i := 0; i < nSiswa; i++ {
			if i != siswaIdx && globalSiswa[i].id == idBaru {
				idSudahAda = true
			}
		}
	}

	// Input data lainnya
	fmt.Print("Masukkan kelas les siswa: ")
	fmt.Scan(&kelasBaru)

	fmt.Print("Masukkan nomor telepon siswa: ")
	fmt.Scan(&teleponBaru)

	fmt.Print("Masukkan tanggal lahir siswa (dd/mm/yyyy): ")
	fmt.Scan(&tanggalLahirBaru)

	// Validasi tanggal
	for !isValidDate(tanggalLahirBaru) {
		fmt.Print("❌ Format tanggal salah! Gunakan format dd/mm/yyyy: ")
		fmt.Scan(&tanggalLahirBaru)
	}

	// Input mata pelajaran
	for i := 0; i < MAX_SUBJECTS; i++ {
		fmt.Printf("Masukkan mata pelajaran ke-%d: ", i+1)
		fmt.Scan(&mataPelajaranBaru)
		globalSiswa[siswaIdx].mataPelajaran[i].nama = mataPelajaranBaru
	}

	// Input jadwal les
	for i := 0; i < MAX_SUBJECTS; i++ {
		if globalSiswa[siswaIdx].mataPelajaran[i].nama != "" {
			fmt.Printf("Masukkan jadwal les untuk %s:\n", globalSiswa[siswaIdx].mataPelajaran[i].nama)
			fmt.Print("Hari: ")
			fmt.Scan(&jadwalHari)
			fmt.Print("Jam: ")
			fmt.Scan(&jadwalJam)
			globalSiswa[siswaIdx].jadwalLes[i].day = jadwalHari
			globalSiswa[siswaIdx].jadwalLes[i].time = jadwalJam
		}
	}

	// Simpan data
	globalSiswa[siswaIdx].id = idBaru
	globalSiswa[siswaIdx].kelas = kelasBaru
	globalSiswa[siswaIdx].telepon = teleponBaru
	globalSiswa[siswaIdx].tanggalLahir = tanggalLahirBaru

	fmt.Println("✅ Data siswa telah diperbarui.")
}

// PROSEDUR 9: Input nilai try out
// I.S.: siswa sudah memiliki data dan mata pelajaran
// F.S.: nilai try out siswa terinput dan rata-rata dihitung
func inputNilaiTryOut() {
	if nSiswa == 0 {
		fmt.Println("❌ Belum ada siswa yang terdaftar!")
		return
	}

	var idSiswa int
	var siswaIdx int
	var nilai float64
	const maxRetry = 3
	var attempts int
	var found bool

	// Urutkan dulu untuk binary search
	selectionSortSiswaByID()

	// Input ID dengan retry
	attempts = 0
	found = false
	for attempts < maxRetry && !found {
		fmt.Printf("Masukkan ID siswa (percobaan %d/%d): ", attempts+1, maxRetry)
		fmt.Scan(&idSiswa)

		siswaIdx = binarySearchSiswaByID(idSiswa)
		if siswaIdx != -1 {
			found = true
		} else {
			attempts++
			fmt.Printf("❌ Siswa dengan ID %d tidak ditemukan!\n", idSiswa)
		}
	}

	if !found {
		fmt.Println("❌ Maksimal percobaan tercapai!")
		return
	}

	fmt.Printf("✅ Siswa ditemukan: %s\n", globalSiswa[siswaIdx].nama)

	// Cek mata pelajaran tersedia
	mataPelajaranAda := false
	for i := 0; i < MAX_SUBJECTS; i++ {
		if globalSiswa[siswaIdx].mataPelajaran[i].nama != "" {
			mataPelajaranAda = true
		}
	}

	if !mataPelajaranAda {
		fmt.Println("❌ Siswa belum memiliki mata pelajaran! Input data siswa terlebih dahulu.")
		return
	}

	// Input nilai untuk setiap mata pelajaran
	for i := 0; i < MAX_SUBJECTS; i++ {
		if globalSiswa[siswaIdx].mataPelajaran[i].nama != "" {
			fmt.Printf("\n=== Input nilai %s ===\n", globalSiswa[siswaIdx].mataPelajaran[i].nama)

			for j := 0; j < MAX_TRYOUTS; j++ {
				fmt.Printf("Masukkan nilai try out ke-%d (0-100): ", j+1)
				fmt.Scan(&nilai)

				// Validasi nilai
				for nilai < 0 || nilai > 100 {
					fmt.Print("❌ Nilai harus antara 0-100! Masukkan ulang: ")
					fmt.Scan(&nilai)
				}

				globalSiswa[siswaIdx].mataPelajaran[i].nilaiTryOut[j] = nilai
			}

			// Hitung rata-rata
			globalSiswa[siswaIdx].mataPelajaran[i].nilaiRata =
				hitungRataMataPelajaran(globalSiswa[siswaIdx].mataPelajaran[i])

			fmt.Printf("✅ Nilai rata-rata %s: %.2f\n",
				globalSiswa[siswaIdx].mataPelajaran[i].nama,
				globalSiswa[siswaIdx].mataPelajaran[i].nilaiRata)
		}
	}

	// Update total nilai dan ranking
	globalSiswa[siswaIdx].totalNilai = hitungTotalNilaiSiswa(globalSiswa[siswaIdx])
	insertionSortSiswaByNilai()

	fmt.Printf("✅ Total nilai siswa: %.2f\n", globalSiswa[siswaIdx].totalNilai)
	fmt.Println("✅ Input nilai try out selesai!")
}

// PROSEDUR 10: Tampilkan data siswa (diurutkan berdasarkan ID)
// I.S.: array siswa terdefinisi
// F.S.: data siswa ditampilkan dalam bentuk tabel terurut berdasarkan ID

// Menampilkan data siswa
func tampilkanDataSiswa() {
	if nSiswa == 0 {
		fmt.Println("❌ Belum ada data siswa yang diinputkan.")
		return
	}

	selectionSortSiswaByID()

	fmt.Println("\n================= DATA SISWA =================")
	fmt.Println("| ID  | Nama Siswa       | Kelas  | Nomor Telepon       | Tanggal Lahir |")
	fmt.Println("===================================================")

	for i := 0; i < nSiswa; i++ {
		fmt.Printf("| %-3d | %-16s | %-6s | %-18s | %-12s |\n",
			globalSiswa[i].id,
			globalSiswa[i].nama,
			globalSiswa[i].kelas,
			globalSiswa[i].telepon,
			globalSiswa[i].tanggalLahir)
	}

	fmt.Println("===================================================")
}

// Fungsi untuk mengedit data siswa
func editDataSiswa() {
	if nSiswa == 0 {
		fmt.Println("❌ Belum ada siswa yang terdaftar!")
		return
	}

	var idSiswa int
	var siswaIdx int
	const maxRetry = 3

	// Pre-sort untuk binary search
	selectionSortSiswaByID()

	// Input ID siswa dengan retry
	attempt := 0
	for attempt < maxRetry {
		fmt.Printf("Masukkan ID siswa yang ingin diedit (percobaan %d/%d): ", attempt+1, maxRetry)
		fmt.Scan(&idSiswa)

		siswaIdx = binarySearchSiswaByID(idSiswa)

		if siswaIdx != -1 {
			break
		}

		attempt++
		fmt.Printf("❌ Siswa dengan ID %d tidak ditemukan!\n", idSiswa)

		if attempt >= maxRetry {
			fmt.Println("❌ Maksimal percobaan tercapai. Kembali ke menu...")
			return
		}
	}

	fmt.Println("✅ Siswa ditemukan:", globalSiswa[siswaIdx].nama)
	fmt.Println("(Kosongkan field untuk tidak mengubah)")

	var input string

	// Edit kelas les
	fmt.Print("Masukkan kelas les baru: ")
	fmt.Scan(&input)
	if strings.TrimSpace(input) != "" && input != "-" {
		globalSiswa[siswaIdx].kelas = input
	}

	// Edit nomor telepon
	fmt.Print("Masukkan nomor telepon baru: ")
	fmt.Scan(&input)
	if strings.TrimSpace(input) != "" && input != "-" {
		globalSiswa[siswaIdx].telepon = input
	}

	// Edit tanggal lahir
	fmt.Print("Masukkan tanggal lahir baru (dd/mm/yyyy): ")
	fmt.Scan(&input)
	if strings.TrimSpace(input) != "" && input != "-" {
		globalSiswa[siswaIdx].tanggalLahir = input
	}

	// Edit catatan
	fmt.Print("Masukkan catatan baru: ")
	fmt.Scan(&input)
	if strings.TrimSpace(input) != "" && input != "-" {
		globalSiswa[siswaIdx].catatan = input
	}

	fmt.Println("✅ Data siswa berhasil diperbarui!")
}

// Fungsi untuk menghapus data siswa
func hapusSiswa() {
	if nSiswa == 0 {
		fmt.Println("❌ Belum ada siswa yang terdaftar!")
		return
	}

	var idSiswa int
	var siswaIdx int
	var konfirmasi string

	// Pre-sort untuk binary search
	selectionSortSiswaByID()

	// Input ID siswa yang ingin dihapus
	fmt.Print("Masukkan ID siswa yang ingin dihapus: ")
	fmt.Scan(&idSiswa)

	// Pencarian siswa berdasarkan ID
	siswaIdx = binarySearchSiswaByID(idSiswa)

	if siswaIdx == -1 {
		fmt.Printf("❌ Siswa dengan ID %d tidak ditemukan!\n", idSiswa)
		return
	}

	// Konfirmasi penghapusan siswa
	fmt.Printf("Apakah Anda yakin ingin menghapus siswa %s (ID: %d)? (Y/N): ",
		globalSiswa[siswaIdx].nama, idSiswa)
	fmt.Scan(&konfirmasi)

	if konfirmasi == "Y" || konfirmasi == "y" {
		// Menggeser elemen array untuk menghapus siswa
		for i := siswaIdx; i < nSiswa-1; i++ {
			globalSiswa[i] = globalSiswa[i+1]
		}
		nSiswa--

		// Menghitung ulang ranking
		insertionSortSiswaByNilai()

		fmt.Println("✅ Siswa berhasil dihapus!")
	} else {
		fmt.Println("❌ Penghapusan dibatalkan.")
	}
}

// Fungsi untuk menu siswa
func menuSiswa() {
	var pilihan int

	for {
		fmt.Println("\n============ MENU SISWA ============")
		fmt.Println("1. Lihat Data Diri")
		fmt.Println("2. Lihat Nilai Try Out")
		fmt.Println("3. Lihat Ranking")
		fmt.Println("4. Lihat Jadwal Les")
		fmt.Println("5. Kembali ke Menu Utama")
		fmt.Println("6. Keluar")
		fmt.Print("Pilih menu (1-6): ")

		fmt.Scan(&pilihan)

		switch pilihan {
		case 1:
			lihatDataDiri()
		case 2:
			lihatNilaiTryOut()
		case 3:
			lihatRanking()
		case 4:
			lihatJadwalLes()
		case 5:
			globalLogin[0].apakahLogin = false
			globalLogin[0].statusUser = ""
			globalLogin[0].userYangLogin = userLes{}
			fmt.Println("✅ Berhasil logout!")
			return
		case 6:
			fmt.Println("✅ Terima kasih telah menggunakan Brained!")
			fmt.Println("   Sampai jumpa lagi! 👋")
			os.Exit(0)
		default:
			fmt.Println("❌ Pilihan tidak valid. Silakan pilih 1-6.")
		}
	}
}

// Fungsi untuk melihat data diri siswa
func lihatDataDiri() {
	username := normalizeUsername(globalLogin[0].userYangLogin.username)
	siswaIdx := linearSearchSiswaByUsername(username)

	if siswaIdx == -1 {
		fmt.Println("❌ Data siswa tidak ditemukan!")
		return
	}

	siswa := globalSiswa[siswaIdx]

	fmt.Println("\n=============== DATA DIRI ===============")
	fmt.Printf("ID Siswa      : %d\n", siswa.id)
	fmt.Printf("Nama          : %s\n", siswa.nama)
	fmt.Printf("Username      : %s\n", siswa.username)
	fmt.Printf("Email         : %s\n", siswa.email)
	fmt.Printf("Kelas         : %s\n", siswa.kelas)
	fmt.Printf("Telepon       : %s\n", siswa.telepon)
	fmt.Printf("Tanggal Lahir : %s\n", siswa.tanggalLahir)
	fmt.Printf("Catatan       : %s\n", siswa.catatan)
	fmt.Printf("Total Nilai   : %.2f\n", siswa.totalNilai)
	fmt.Printf("Ranking       : %d\n", siswa.ranking)
	fmt.Println("==========================================")
}

// Fungsi untuk melihat nilai try out siswa
func lihatNilaiTryOut() {
	username := normalizeUsername(globalLogin[0].userYangLogin.username)
	siswaIdx := linearSearchSiswaByUsername(username)

	if siswaIdx == -1 {
		fmt.Println("❌ Data siswa tidak ditemukan!")
		return
	}

	siswa := globalSiswa[siswaIdx]

	fmt.Println("\n=============== NILAI TRY OUT ===============")
	fmt.Printf("Nama Siswa: %s\n", siswa.nama)
	fmt.Println("============================================")

	hasGrades := false
	for i := 0; i < 3; i++ {
		if siswa.mataPelajaran[i].nama != "" {
			hasGrades = true
			fmt.Printf("\nMata Pelajaran: %s\n", siswa.mataPelajaran[i].nama)
			fmt.Printf("Nilai Try Out 1: %.2f\n", siswa.mataPelajaran[i].nilaiTryOut[0])
			fmt.Printf("Nilai Try Out 2: %.2f\n", siswa.mataPelajaran[i].nilaiTryOut[1])
			fmt.Printf("Nilai Try Out 3: %.2f\n", siswa.mataPelajaran[i].nilaiTryOut[2])
			fmt.Printf("Nilai Rata-rata: %.2f\n", siswa.mataPelajaran[i].nilaiRata)
			fmt.Println("----------------------------")
		}
	}

	if !hasGrades {
		fmt.Println("❌ Belum ada nilai try out yang diinput!")
	} else {
		fmt.Printf("\nTotal Nilai Keseluruhan: %.2f\n", siswa.totalNilai)
	}
	fmt.Println("============================================")
}

// Fungsi untuk melihat ranking siswa
func lihatRanking() {
	if nSiswa == 0 {
		fmt.Println("❌ Belum ada data siswa!")
		return
	}

	// Update ranking
	insertionSortSiswaByNilai()

	fmt.Println("\n=============== RANKING SISWA ===============")
	fmt.Println("| Rank | Nama Siswa       | Total Nilai |")
	fmt.Println("===========================================")

	for i := 0; i < nSiswa; i++ {
		if globalSiswa[i].totalNilai > 0 {
			fmt.Printf("| %-4d | %-16s | %-11.2f |\n",
				globalSiswa[i].ranking,
				globalSiswa[i].nama,
				globalSiswa[i].totalNilai)
		}
	}
	fmt.Println("===========================================")

	// Highlight current student's ranking
	username := normalizeUsername(globalLogin[0].userYangLogin.username)
	siswaIdx := linearSearchSiswaByUsername(username)

	if siswaIdx != -1 && globalSiswa[siswaIdx].totalNilai > 0 {
		fmt.Printf("\n🏆 Ranking Anda: %d dari %d siswa\n",
			globalSiswa[siswaIdx].ranking, nSiswa)
	}
}

// Fungsi untuk melihat jadwal les siswa
func lihatJadwalLes() {
	username := normalizeUsername(globalLogin[0].userYangLogin.username)
	siswaIdx := linearSearchSiswaByUsername(username)

	if siswaIdx == -1 {
		fmt.Println("❌ Data siswa tidak ditemukan!")
		return
	}

	siswa := globalSiswa[siswaIdx]

	fmt.Println("\n=============== JADWAL LES ===============")
	fmt.Printf("Nama Siswa: %s\n", siswa.nama)
	fmt.Println("==========================================")

	hasSchedule := false
	for i := 0; i < 3; i++ {
		if siswa.mataPelajaran[i].nama != "" && siswa.jadwalLes[i].day != "" {
			hasSchedule = true
			fmt.Printf("Mata Pelajaran: %s\n", siswa.mataPelajaran[i].nama)
			fmt.Printf("Hari          : %s\n", siswa.jadwalLes[i].day)
			fmt.Printf("Jam           : %s\n", siswa.jadwalLes[i].time)
			fmt.Println("----------------------------")
		}
	}

	if !hasSchedule {
		fmt.Println("❌ Belum ada jadwal les yang diatur!")
	}
	fmt.Println("==========================================")
}

func menuAdmin() {
	var pilihan int

	for {
		fmt.Println("\n============ MENU ADMIN ============")
		fmt.Println("1. Input Data Siswa")
		fmt.Println("2. Input Nilai Try Out")
		fmt.Println("3. Tampilkan Data Siswa")
		fmt.Println("4. Edit Data Siswa")
		fmt.Println("5. Hapus Siswa")
		fmt.Println("6. Lihat Ranking Siswa")
		fmt.Println("7. Kembali ke Menu Utama")
		fmt.Println("8. Keluar")
		fmt.Print("Pilih menu (1-8): ")

		fmt.Scan(&pilihan)

		switch pilihan {
		case 1:
			inputDataSiswa()
		case 2:
			inputNilaiTryOut() // INI YANG HILANG!
		case 3:
			tampilkanDataSiswa()
		case 4:
			editDataSiswa()
		case 5:
			hapusSiswa()
		case 6:
			lihatRanking()
		case 7:
			globalLogin[0].apakahLogin = false
			globalLogin[0].statusUser = ""
			globalLogin[0].userYangLogin = userLes{}
			fmt.Println("✅ Berhasil logout!")
			return
		case 8:
			fmt.Println("✅ Terima kasih telah menggunakan Brained!")
			fmt.Println("   Sampai jumpa lagi! 👋")
			os.Exit(0)
		default:
			fmt.Println("❌ Pilihan tidak valid. Silakan pilih 1-8.")
		}
	}
}

// Fungsi untuk menu utama
func menuUtama() {
	var pilihan int

	for {
		fmt.Println("\n========== SELAMAT DATANG DI BRAINED ==========")
		fmt.Println("           Platform Bimbel Terbaik")
		fmt.Println("===============================================")
		fmt.Println("1. Login")
		fmt.Println("2. Register")
		fmt.Println("3. Keluar")
		fmt.Println("===============================================")
		fmt.Print("Pilih menu (1-3): ")

		fmt.Scan(&pilihan)

		switch pilihan {
		case 1:
			login()
		case 2:
			registrasi()
		case 3:
			fmt.Println("✅ Terima kasih telah menggunakan Brained!")
			fmt.Println("   Sampai jumpa lagi! 👋")
			return
		default:
			fmt.Println("❌ Pilihan tidak valid. Silakan pilih 1-3.")
		}
	}
}

// Main function
func main() {
	// Initialize global data
	nAdmin = 0
	nSiswa = 0
	globalLogin[0].apakahLogin = false
	globalLogin[0].statusUser = ""

	// Start the main menu
	menuUtama()
}
