package serviceuser

type RequestUser struct {
	Tujuan           string   `json:"tujuan"`
	Berdasarkan      string   `json:"berdasarkan"`
	BukanBuku        []string `json:"bukanbuku"`
	Pencarian        string   `json:"pencarian"`
	Favorit          []string `json:"favoritnya"`
	IdUser           string   `json:"iduser"`
	NamaBuku         string   `json:"namabuku"`
	ISBnBuku         string   `json:"isbn"`
	NamaUser         string   `json:"namauser"`
	JenisContributor string   `json:"jeniskontributor"`
	IdBuku           string   `json:"idbuku"`
}
