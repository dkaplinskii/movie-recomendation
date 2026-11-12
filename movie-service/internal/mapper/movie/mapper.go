package movie

import (
	"fmt"
	"strconv"

	"studentgit.kata.academy/movie-recommendation-platform-20/movie-service/internal/domain"
	"studentgit.kata.academy/movie-recommendation-platform-20/movie-service/internal/errors"
)

func MapToMovie(film FilmResponseDTO, videos VideosResponseDTO) (domain.Movie, error) {
	movie := domain.Movie{
		KinopoiskID: film.KinopoiskID,
		Description: film.Description,
		PosterURL:   film.PosterURL,
	}

	switch {
	case film.NameRu != "":
		movie.Title = film.NameRu
	case film.NameEn != "":
		movie.Title = film.NameEn
	case film.NameOriginal != "":
		movie.Title = film.NameOriginal
	default:
		return domain.Movie{}, errors.ErrInvalidTitle
	}

	for _, filmCountry := range film.Countries {
		if filmCountry.Country == "" {
			continue
		}

		movie.Countries = append(movie.Countries, filmCountry.Country)
	}

	for _, filmGenre := range film.Genres {
		if filmGenre.Genre == "" {
			continue
		}

		movie.Genres = append(movie.Genres, filmGenre.Genre)
	}

	if film.Year != "" {
		year, err := strconv.Atoi(film.Year)
		if err != nil {
			return domain.Movie{}, fmt.Errorf("%w: %q", errors.ErrInvalidYear, film.Year)
		}

		movie.Year = year
	}

	for _, video := range videos.Items {
		if video.URL == "" {
			continue
		}

		movie.TrailerURL = video.URL
		break
	}

	return movie, nil
}
