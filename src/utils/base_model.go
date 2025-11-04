package utils

import (
	"database/sql"
	"fmt"
	"strings"
)

// BaseModel merepresentasikan model dasar yang bisa digunakan ulang
type BaseModel struct {
	DB        *sql.DB
	TableName string
}

// NewBaseModel membuat instance baru dari BaseModel
func NewBaseModel(db *sql.DB, tableName string) *BaseModel {
	return &BaseModel{
		DB:        db,
		TableName: tableName,
	}
}

// Create melakukan insert data ke tabel dan mengembalikan id dari data yang baru
func (b *BaseModel) Create(data map[string]interface{}) (map[string]interface{}, error) {
	//TODO: make() itu digunakan untuk membuat slice, map, atau channel dengan ukuran tertentu. make() memiliki 3 argumen yaitu tipe data, panjang awal, dan kapasitas awal (opsional). Jika kapasitas tidak ditambahkan, otomatis kapasitas akan diambil dari length.
	fields := make([]string, 0, len(data))
	placeholders := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data)) //TODO interface[] adalah tipe data bebas. Ini cocok untuk menampung value yang biasanya bermacam-macam tipe data.

	//TODO: range pada map di Go mengembalikan key dan value untuk setiap elemen.
	//TODO: keyword range variable 1 dan 2 nya bisa beda beda. Jika slice of string yang akan dikembalikan adalah index number dan values. Kalo map yang dikembalikan adalah keys dan values.
	//TODO: Nah, ini bukan loop manual seperti for { ... } atau for i < something { ... }. Tapi ini adalah loop berbasis range, yang otomatis berhenti setelah semua elemen di data habis diiterasi.
	i := 1
	for k, v := range data {
		fields = append(fields, k)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i)) //TODO: fmt.Sprintf() gunanya untuk membuat string dengan format tertentu menggunakan placeholders khusus.
		values = append(values, v)
		i++
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING *",
		b.TableName,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "),
	)

	//TODO: Jika query menggunakan .QueryRow() maka tidak perlu di defer Close() karena hanya mengembalikan satu baris.
	row := b.DB.QueryRow(query, values...)

	// Karena RETURNING * hasilnya bisa dinamis, kita scan ke map pakai Rows.Columns
	columns, err := b.getColumns()
	if err != nil {
		return nil, err
	}

	//TODO: Pennggunaan make disini agar slicenya dapat diakses indeksnya. Karena jika langsung membuat menggunakan slice seperti []interface langsung, tidak akan bisa diakses index nya dan akan error. Karena dengan make() ini si length sama kapasitasnya sudah diatur jadi indeksnya sudah ada duluan meskipun isinya kosong.
	result := make(map[string]interface{})
	columnPointers := make([]interface{}, len(columns))
	columnValues := make([]interface{}, len(columns))

	//TODO: i otomatis berisi 0. itu disebabkan otomatis oleh range.
	for i := range columns {
		columnPointers[i] = &columnValues[i]
	}

	//TODO: Variadic expansion konsepnya sama seperti spread operator.
	if err := row.Scan(columnPointers...); err != nil {
		return nil, err
	}

	for i, col := range columns {
		result[col] = columnValues[i] //TODO: Jika key "col" belum ada di map, Go otomatis membuat key baru dan assign value-nya. Map di Go dinamis → bisa menambah key kapan saja dengan assignment
	}

	return result, nil
}

// FindById mencari satu record berdasarkan id
func (b *BaseModel) FindById(id string) (map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", b.TableName)
	row := b.DB.QueryRow(query, id)

	columns, err := b.getColumns()
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	columnPointers := make([]interface{}, len(columns))
	columnValues := make([]interface{}, len(columns))

	//TODO: Scan() membutuhkan pointer (menunjuk ke alamat asli) dari setiap value yang akan di-scan.). Makanya di code di bawah dibuat columnPointers itu menunjuk ke alaamat asli dari columnValues.
	//TODO: 'range' adalah keyword looping yang bisa dipakai untuk array, slice, map, string, channel, dan langsung built-in bahasa Go
	for i := range columns {
		columnPointers[i] = &columnValues[i]
	}

	if err := row.Scan(columnPointers...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // tidak ditemukan
		}
		return nil, err
	}

	for i, col := range columns {
		result[col] = columnValues[i]
	}

	return result, nil
}

// GetAll mengambil semua record dari tabel
func (b *BaseModel) GetAll() ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s", b.TableName)
	rows, err := b.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := b.getColumns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		columnPointers := make([]interface{}, len(columns))
		columnValues := make([]interface{}, len(columns))
		for i := range columns {
			columnPointers[i] = &columnValues[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[col] = columnValues[i]
		}

		results = append(results, rowMap)
	}

	return results, nil
}

// getColumns mengambil nama kolom dari tabel terkait (dibutuhkan untuk mapping hasil RETURNING *)
func (b *BaseModel) getColumns() ([]string, error) {
	//TODO: select column_name from information_schema.columns itu pasti mengembalikan semua nama kolom perbaris. Bayangin aja select name from users, itu pasti menghasilkan banyak baris berisi nama-nama user dari kolom "name"
	query := fmt.Sprintf(`SELECT column_name FROM information_schema.columns WHERE table_name = '%s'`, b.TableName)
	rows, err := b.DB.Query(query) //TODO: rows isinya adalah iterator. Beda dengan JS yang isinya adalah array of objects. Rows ini perlu "ditarik" perbaris menggunakan .Next().
	if err != nil {
		return nil, err
	}
	//TODO: Jika query menggunakan .Query() maka harus di defer Close() untuk menutup koneksi setelah selesai digunakan.
	defer rows.Close()

	var columns []string
	//TODO: Bisa dikatakan untuk hasil dari .Query() itu akan selalu menggunakan looping "for namaVariable.Next()" untuk menarik / mengiterasi setiap baris agar dapat discan.
	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}
	return columns, nil
}
