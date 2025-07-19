package serviceadmin

type RequestAdmin struct {
	Id     string `json:"id"`
	Nama   string `json:"nama"`
	Tujuan string `json:"tujuan`
}

type UserRequest struct {
	Nama       string `json:"nama"`
	Password   string `json:"password"`
	Favorit    string `json:"favorit"`
	KreditSkor int    `json:"kreditskor"`
	Email      string `json:"email"`
	Alamat     string `json:"alamat"`
	Status     string `json:"status"`
	Bergabung  string `json:"bergabung"`
}

type BukuBaruRequest struct {
	ID           string `json:"id"`
	Tujuan       string `json:"tujuan"`
	Judul        string `json:"judul"`
	Jenis        string `json:"jenis"`
	Harga        string `json:"harga"`
	Penulis      string `json:"penulis"`
	Penerbit     string `json:"penerbit"`
	Stok         string `json:"stok"`
	Tahun        string `json:"tahun"`
	ISBN         string `json:"ISBN"`
	Kategori     string `json:"kategori"`
	Bahasa       string `json:"bahasa"`
	Deskripsi    string `json:"deskripsi"`
	GambarBase64 string `json:"gambarBase64"`
}

type Peminjamanbuku struct {
	ID           int    `gorm:"primaryKey;column:id"`
	Judul        string `gorm:"column:judul"`
	Tanggal      string `gorm:"column:tanggal"`
	Kategori     string `gorm:"column:kategori"`
	Status       string `gorm:"column:status"`
	Kodebuku     int64  `gorm:"column:kodebuku"`
	ISBN         string `gorm:"column:ISBN"`
	Iduser       int64  `gorm:"column:iduser"`
	Namapeminjam string `gorm:"column:namapeminjam"`
}
