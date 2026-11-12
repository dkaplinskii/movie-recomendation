package stub

import "recommendation-service/internal/domain"

func defaultReviews(userID int64) []domain.Review {
	return []domain.Review{
		{UserID: userID, MovieID: 1, Rating: 10},
		{UserID: userID, MovieID: 2, Rating: 8},
	}
}

func defaultReviewedMovies() []domain.Movie {
	return []domain.Movie{
		{
			ID:      1,
			Title:   "Reviewed Drama",
			Genres:  []string{"drama"},
			Country: "USA",
			Year:    2020,
		},
		{
			ID:      2,
			Title:   "Reviewed Thriller",
			Genres:  []string{"thriller"},
			Country: "USA",
			Year:    2020,
		},
	}
}

func defaultCandidateMovies() []domain.Movie {
	return []domain.Movie{
		{
			ID:          10,
			Title:       "The Shawshank Redemption",
			Description: "Two imprisoned men bond over years, finding hope and friendship.",
			Genres:      []string{"drama"},
			PosterURL:   "https://example.com/posters/shawshank.jpg",
			Year:        1994,
			Country:     "USA",
			TrailerURL:  "https://example.com/trailers/shawshank",
		},
		{
			ID:          11,
			Title:       "Inception",
			Description: "A thief enters dreams to steal secrets and plant an idea.",
			Genres:      []string{"sci-fi", "action", "thriller"},
			PosterURL:   "https://example.com/posters/inception.jpg",
			Year:        2010,
			Country:     "USA",
			TrailerURL:  "https://example.com/trailers/inception",
		},
		{
			ID:          12,
			Title:       "The Dark Knight",
			Description: "Batman faces the Joker in Gotham City.",
			Genres:      []string{"action", "drama", "thriller"},
			PosterURL:   "https://example.com/posters/dark-knight.jpg",
			Year:        2008,
			Country:     "USA",
			TrailerURL:  "https://example.com/trailers/dark-knight",
		},
		{
			ID:          13,
			Title:       "Interstellar",
			Description: "Explorers travel through a wormhole in space.",
			Genres:      []string{"sci-fi", "drama"},
			PosterURL:   "https://example.com/posters/interstellar.jpg",
			Year:        2014,
			Country:     "USA",
			TrailerURL:  "https://example.com/trailers/interstellar",
		},
		{
			ID:          14,
			Title:       "Pulp Fiction",
			Description: "Interconnected stories of crime in Los Angeles.",
			Genres:      []string{"crime", "drama"},
			PosterURL:   "https://example.com/posters/pulp-fiction.jpg",
			Year:        1994,
			Country:     "USA",
			TrailerURL:  "https://example.com/trailers/pulp-fiction",
		},
		{
			ID:          15,
			Title:       "The Matrix",
			Description: "A hacker discovers the true nature of reality.",
			Genres:      []string{"sci-fi", "action"},
			PosterURL:   "https://example.com/posters/matrix.jpg",
			Year:        1999,
			Country:     "USA",
			TrailerURL:  "https://example.com/trailers/matrix",
		},
		{
			ID:          16,
			Title:       "Forrest Gump",
			Description: "The life story of a kind-hearted man from Alabama.",
			Genres:      []string{"drama", "romance"},
			PosterURL:   "https://example.com/posters/forrest-gump.jpg",
			Year:        1994,
			Country:     "USA",
			TrailerURL:  "https://example.com/trailers/forrest-gump",
		},
		{
			ID:          17,
			Title:       "Gladiator",
			Description: "A Roman general seeks revenge as a gladiator.",
			Genres:      []string{"action", "drama"},
			PosterURL:   "https://example.com/posters/gladiator.jpg",
			Year:        2000,
			Country:     "USA",
			TrailerURL:  "https://example.com/trailers/gladiator",
		},
		{
			ID:          18,
			Title:       "The Prestige",
			Description: "Two magicians engage in a bitter rivalry.",
			Genres:      []string{"drama", "mystery", "thriller"},
			PosterURL:   "https://example.com/posters/prestige.jpg",
			Year:        2006,
			Country:     "USA",
			TrailerURL:  "https://example.com/trailers/prestige",
		},
		{
			ID:          19,
			Title:       "Se7en",
			Description: "Two detectives hunt a serial killer using the seven deadly sins.",
			Genres:      []string{"crime", "drama", "thriller"},
			PosterURL:   "https://example.com/posters/se7en.jpg",
			Year:        1995,
			Country:     "USA",
			TrailerURL:  "https://example.com/trailers/se7en",
		},
	}
}
