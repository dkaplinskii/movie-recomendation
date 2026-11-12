package models

type OutboxEvent struct {
	EventID int64
	UserID  int64
	FilmID  int64
	Review  int32
}
