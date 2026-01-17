package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/senatroxx/filmix-backend/internal/database/entities"
)

type ISeatRepository interface {
	FindByStudioID(ctx context.Context, studioID uuid.UUID) ([]entities.Seat, error)
	FindBookedSeatIDs(ctx context.Context, showtimeID uuid.UUID) ([]uuid.UUID, error)
}

type SeatRepository struct {
	db *sql.DB
}

func NewSeatRepository(db *sql.DB) ISeatRepository {
	return &SeatRepository{db: db}
}

func (r *SeatRepository) FindByStudioID(ctx context.Context, studioID uuid.UUID) ([]entities.Seat, error) {
	query := `
		SELECT s.id, s.row, s.number, s.active, s.studio_id, s.seat_type_id,
		       st.id, st.name
		FROM seats s
		JOIN seat_type st ON s.seat_type_id = st.id
		WHERE s.studio_id = $1 AND s.active = true
		ORDER BY s.row, s.number
	`

	rows, err := r.db.QueryContext(ctx, query, studioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []entities.Seat
	for rows.Next() {
		var seat entities.Seat
		var seatType entities.SeatType

		err := rows.Scan(
			&seat.ID, &seat.Row, &seat.Number, &seat.Active, &seat.StudioID, &seat.SeatTypeID,
			&seatType.ID, &seatType.Name,
		)
		if err != nil {
			return nil, err
		}

		seat.SeatType = &seatType
		seats = append(seats, seat)
	}

	return seats, nil
}

func (r *SeatRepository) FindBookedSeatIDs(ctx context.Context, showtimeID uuid.UUID) ([]uuid.UUID, error) {
	query := `
		SELECT ti.seat_id
		FROM transaction_items ti
		JOIN transactions t ON ti.transaction_id = t.id
		WHERE t.showtime_id = $1 AND t.status IN ('pending', 'paid')
	`

	rows, err := r.db.QueryContext(ctx, query, showtimeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookedIDs []uuid.UUID
	for rows.Next() {
		var seatID uuid.UUID
		if err := rows.Scan(&seatID); err != nil {
			return nil, err
		}
		bookedIDs = append(bookedIDs, seatID)
	}

	return bookedIDs, nil
}
