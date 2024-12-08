package model

type File struct {
	Announce  string `json:"announce"`
	CreatedBy string `json:"created by"`
	Info      Info   `json:"info"`
}
type Info struct {
	Length int64 `json:"length"`
	Name string `json:"name"`
	PieceLength int64 `json:"piece length"`
	Pieces string `json:"pieces"`
}
