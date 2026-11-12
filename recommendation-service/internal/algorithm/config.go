package algorithm

type Config struct {
	LatestReviewsLimit int
	RecommendationSize int
	MinPreferenceScore float64
	MaxFilterValues    int
	PreferredYearRange int
	GenreWeight        float64
	CountryWeight      float64
	YearWeight         float64
	RatingWeight       float64
	MinYearSampleSize  int
	CandidatePoolSize  int
}
