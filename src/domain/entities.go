package domain

type Player struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Stats  Stats  `json:"stats"`
	Family Family `json:"family"`
}

type Stats struct {
	// Scale 1 to 10
	Agility       int `json:"agility"`
	Intelligence  int `json:"intelligence"`
	ArtisticSkill int `json:"artistic_skill"`
	Communication int `json:"communication"`
	SocialEnergy  int `json:"social_energy"`
}

type Family struct {
	ID int64 `json:"id"`
	//Name string `json:"name"`
}

type Team struct {
	//ID   int64  `json:"id"`
	Name string `json:"name"`
	//Color   string   `json:"color"`
	Members []Player `json:"members"`
	Score   int      `json:"score"`
}

type Game struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	RelevantStats Stats  `json:"relevant_stats"`
}
