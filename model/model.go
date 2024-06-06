package model

type Credentials struct {
	Hostname string
	Port     string
	Username string
	Password string
	DBName   string
}

type IMDBMovie struct {
	ID               string `bson:"_id"`
	Title            string `bson:"title"`
	ReleaseDate      string `bson:"relesase_date"`
	Genres           string `bson:"genres"`
	OriginalLanguage string `bson:"original_language"`
	Overview         string `bson:"popularity" type:"desc"`
	VoteCount        string `bson:"vote_count"`
	VoteAverage      string `bson:"vote_average"`
}
