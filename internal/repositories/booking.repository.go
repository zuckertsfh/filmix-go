package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/senatroxx/filmix-backend/internal/database/entities"
)

type IBookingRepository interface {
	Create(ctx context.Context, tx *entities.Transaction, items []entities.TransactionItem) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Transaction, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entities.Transaction, error)
	CheckSeatsAvailable(ctx context.Context, showtimeID uuid.UUID, seatIDs []uuid.UUID) (bool, error)
}

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) IBookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(ctx context.Context, tx *entities.Transaction, items []entities.TransactionItem) error {
	dbTx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer dbTx.Rollback()

	txQuery := `
		INSERT INTO transactions (id, status, external_ref, invoice_number, amount, expired_at, payment_method_id, showtime_id, theater_id, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = dbTx.ExecContext(ctx, txQuery,
		tx.ID, tx.Status, tx.ExternalRef, tx.InvoiceNumber, tx.Amount, tx.ExpiredAt,
		tx.PaymentMethodID, tx.ShowtimeID, tx.TheaterID, tx.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	itemQuery := `
		INSERT INTO transaction_items (id, price, transaction_id, seat_id, seat_type_id)
		VALUES ($1, $2, $3, $4, $5)
	`
	for _, item := range items {
		_, err = dbTx.ExecContext(ctx, itemQuery,
			item.ID, item.Price, item.TransactionID, item.SeatID, item.SeatTypeID,
		)
		if err != nil {
			return fmt.Errorf("failed to insert transaction item: %w", err)
		}
	}

	return dbTx.Commit()
}

func (r *BookingRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Transaction, error) {
	query := `
		SELECT 
			t.id, t.status, t.external_ref, t.invoice_number, t.amount, t.expired_at, t.paid_at,
			t.payment_method_id, t.showtime_id, t.theater_id, t.user_id,
			s.id, s.time, s.movie_id,
			m.id, m.title, m.poster_url,
			th.id, th.name
		FROM transactions t
		JOIN showtimes s ON t.showtime_id = s.id
		JOIN movies m ON s.movie_id = m.id
		JOIN theaters th ON t.theater_id = th.id
		WHERE t.id = $1
	`

	var tx entities.Transaction
	var showtime entities.Showtime
	var movie entities.Movie
	var theater entities.Theater

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tx.ID, &tx.Status, &tx.ExternalRef, &tx.InvoiceNumber, &tx.Amount, &tx.ExpiredAt, &tx.PaidAt,
		&tx.PaymentMethodID, &tx.ShowtimeID, &tx.TheaterID, &tx.UserID,
		&showtime.ID, &showtime.Time, &showtime.MovieID,
		&movie.ID, &movie.Title, &movie.PosterURL,
		&theater.ID, &theater.Name,
	)
	if err != nil {
		return nil, err
	}

	showtime.Movie = &movie
	tx.Showtime = &showtime
	tx.Theater = &theater

	itemQuery := `
		SELECT ti.id, ti.price, ti.seat_id, ti.seat_type_id,
		       s.id, s.row, s.number,
		       st.id, st.name
		FROM transaction_items ti
		JOIN seats s ON ti.seat_id = s.id
		JOIN seat_type st ON ti.seat_type_id = st.id
		WHERE ti.transaction_id = $1
	`

	rows, err := r.db.QueryContext(ctx, itemQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entities.TransactionItem
		var seat entities.Seat
		var seatType entities.SeatType

		err := rows.Scan(
			&item.ID, &item.Price, &item.SeatID, &item.SeatTypeID,
			&seat.ID, &seat.Row, &seat.Number,
			&seatType.ID, &seatType.Name,
		)
		if err != nil {
			return nil, err
		}

		seat.SeatType = &seatType
		item.Seat = &seat
		tx.Items = append(tx.Items, item)
	}

	return &tx, nil
}

func (r *BookingRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entities.Transaction, error) {
	query := `
		SELECT 
			t.id, t.status, t.amount, t.expired_at, t.paid_at,
			s.id, s.time,
			m.id, m.title, m.poster_url,
			th.id, th.name
		FROM transactions t
		JOIN showtimes s ON t.showtime_id = s.id
		JOIN movies m ON s.movie_id = m.id
		JOIN theaters th ON t.theater_id = th.id
		WHERE t.user_id = $1
		ORDER BY t.expired_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []entities.Transaction
	for rows.Next() {
		var tx entities.Transaction
		var showtime entities.Showtime
		var movie entities.Movie
		var theater entities.Theater

		err := rows.Scan(
			&tx.ID, &tx.Status, &tx.Amount, &tx.ExpiredAt, &tx.PaidAt,
			&showtime.ID, &showtime.Time,
			&movie.ID, &movie.Title, &movie.PosterURL,
			&theater.ID, &theater.Name,
		)
		if err != nil {
			return nil, err
		}

		showtime.Movie = &movie
		tx.Showtime = &showtime
		tx.Theater = &theater
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func (r *BookingRepository) CheckSeatsAvailable(ctx context.Context, showtimeID uuid.UUID, seatIDs []uuid.UUID) (bool, error) {
	if len(seatIDs) == 0 {
		return false, nil
	}

	// Check if any of the seats are already booked for this showtime
	query := `
		SELECT COUNT(*) FROM transaction_items ti
		JOIN transactions t ON ti.transaction_id = t.id
		WHERE t.showtime_id = $1 
		AND t.status IN ('pending', 'paid')
		AND t.expired_at > $2
		AND ti.seat_id = ANY($3)
	`

	// Convert UUID slice to string slice for pq.Array
	seatIDStrings := make([]string, len(seatIDs))
	for i, id := range seatIDs {
		seatIDStrings[i] = id.String()
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, showtimeID, time.Now(), pq.Array(seatIDStrings)).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}
