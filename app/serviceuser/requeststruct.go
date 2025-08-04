package serviceuser

type RequestUser struct {
	Tujuan      string   `json:"tujuan"`
	Berdasarkan string   `json:"berdasarkan"`
	BukanBuku   []string `json:"bukanbuku"`
	Pencarian   string   `json:"pencarian"`
	Favorit     []string `json:"favoritnya"`
}
