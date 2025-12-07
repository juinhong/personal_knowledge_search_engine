package models

type Note struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
	Type    string   `json:"type"`
}
