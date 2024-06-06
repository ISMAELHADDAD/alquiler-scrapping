package database

type Entry struct {
	Title string `json:"title"`
	Price int    `json:"price"`
	Url   string `json:"url"`
}

type Entries = []Entry
