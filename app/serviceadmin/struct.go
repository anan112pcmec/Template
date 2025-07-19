package serviceadmin

type RequestAdmin struct {
	Id     string `json:"id"`
	Nama   string `json:"nama"`
	Tujuan string `json:"tujuan`
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
