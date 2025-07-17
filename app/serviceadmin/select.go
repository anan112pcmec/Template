package serviceadmin

import (
	"time"

	"gorm.io/gorm"
)

type BukuInduk struct {
	ID         uint           `gorm:"primaryKey;autoIncrement"`
	Judul      string         `gorm:"column:judul;type:varchar(255);not null"`
	Jenis      string         `gorm:"column:jenis;type:varchar(100);not null"`
	Harga      int64          `gorm:"column:harga;not null"`
	Penulis    string         `gorm:"column:penulis;type:varchar(255);not null"`
	Penerbit   string         `gorm:"column:penerbit;type:varchar(255);not null"`
	Stok       int64          `gorm:"column:stok;not null"`
	Tahun      string         `gorm:"column:tahun;type:varchar(10);not null"`
	ISBN       string         `gorm:"column:isbn;type:varchar(50);unique;not null"`
	Kategori   string         `gorm:"column:kategori;type:varchar(100);not null"`
	Bahasa     string         `gorm:"column:bahasa;type:varchar(50);not null"`
	Deskripsi  string         `gorm:"column:deskripsi;type:text"`
	TujuanAksi string         `gorm:"column:tujuan_aksi;type:varchar(100)"`
	Gambar     []byte         `gorm:"column:gambar;type:bytea"`
	CreatedAt  time.Time      `gorm:"column:created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type BukuChild struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	KodeInduk uint   `gorm:"column:Kode_induk;not null"`
	Judul     string `gorm:"column:judul;type:varchar(255);not null"`
	Jenis     string `gorm:"column:jenis;type:varchar(100);not null"`
	Harga     int64  `gorm:"column:harga;not null"`
	Penulis   string `gorm:"column:penulis;type:varchar(255);not null"`
	Penerbit  string `gorm:"column:penerbit;type:varchar(255);not null"`
	Stok      int64  `gorm:"column:stok;not null"`
	Tahun     string `gorm:"column:tahun;type:varchar(10);not null"`
	ISBN      string `gorm:"column:isbn;type:varchar(50);not null"`
	Kategori  string `gorm:"column:kategori;type:varchar(100);not null"`
	Bahasa    string `gorm:"column:bahasa;type:varchar(50);not null"`
	Status    string `gorm:"column:status;type:varchar(50)"`
	Deskripsi string `gorm:"column:deskripsi;type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
