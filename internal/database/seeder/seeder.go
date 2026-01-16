package seeder

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/senatroxx/filmix-backend/internal/config"
	"github.com/senatroxx/filmix-backend/internal/database/entities"
	"github.com/senatroxx/filmix-backend/internal/integrations/tmdb"
	"golang.org/x/crypto/bcrypt"
)

type Seeder struct {
	db   *sql.DB
	cfg  *config.Config
	repo *Repository
	tmdb *tmdb.Client
}

func NewSeeder(db *sql.DB, cfg *config.Config) *Seeder {
	return &Seeder{
		db:   db,
		cfg:  cfg,
		repo: NewRepository(db),
		tmdb: tmdb.NewClient(cfg.TmdbApiKey),
	}
}

func (s *Seeder) SeedAll() error {
	ctx := context.Background()

	// Truncate all data first
	if err := s.TruncateAll(ctx); err != nil {
		return fmt.Errorf("failed to truncate data: %w", err)
	}

	if err := s.SeedReference(ctx); err != nil {
		return fmt.Errorf("failed to seed reference data: %w", err)
	}
	if err := s.SeedUsers(ctx); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}
	if err := s.SeedCinemas(ctx); err != nil {
		return fmt.Errorf("failed to seed cinemas: %w", err)
	}
	if err := s.SeedMovies(ctx); err != nil {
		return fmt.Errorf("failed to seed movies: %w", err)
	}
	if err := s.SeedShowtimes(ctx); err != nil {
		return fmt.Errorf("failed to seed showtimes: %w", err)
	}

	return nil
}

func (s *Seeder) TruncateAll(ctx context.Context) error {
	log.Println("Truncating all tables...")

	query := `
		TRUNCATE TABLE 
			transaction_items,
			transactions,
			showtimes,
			seats,
			seat_pricing_overrides,
			seat_pricings,
			studios,
			theaters,
			cinemas,
			genre_movie,
			movies,
			movie_genres,
			movie_ratings,
			movie_statuses,
			users,
			roles,
			seat_type
		CASCADE;
	`

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	log.Println("All tables truncated.")
	return nil
}

func (s *Seeder) SeedReference(ctx context.Context) error {
	log.Println("Seeding reference data...")

	// Movie Statuses
	statuses := []string{"Now Showing", "Coming Soon", "Ended"}
	for _, st := range statuses {
		_ = s.repo.CreateStatus(ctx, &entities.MovieStatus{ID: uuid.New(), Status: st})
	}

	// Movie Ratings
	ratings := []string{"G", "PG", "PG-13", "R", "NC-17"}
	for _, r := range ratings {
		_ = s.repo.CreateRating(ctx, &entities.MovieRating{ID: uuid.New(), Rating: r})
	}

	return nil
}

func (s *Seeder) SeedUsers(ctx context.Context) error {
	log.Println("Seeding users...")

	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	query := `INSERT INTO users (id, name, email, password, role_id) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`

	roleAdminID := uuid.New()
	roleUserID := uuid.New()

	roleQuery := `INSERT INTO roles (id, name) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	s.db.ExecContext(ctx, roleQuery, roleAdminID, "admin")
	s.db.ExecContext(ctx, roleQuery, roleUserID, "user")

	// Admin
	s.db.ExecContext(ctx, query, uuid.New(), "Admin User", "admin@filmix.com", string(hash), roleAdminID)
	// User
	s.db.ExecContext(ctx, query, uuid.New(), "Normal User", "user@filmix.com", string(hash), roleUserID)

	return nil
}

func (s *Seeder) SeedCinemas(ctx context.Context) error {
	log.Println("Seeding cinemas...")

	// 1. Cinema
	cinemaName := "Filmix Central Park"
	cinema, err := s.repo.GetCinemaByName(ctx, cinemaName)
	if err != nil {
		// If not found or error, try create
		cinemaID := uuid.New()
		errCreate := s.repo.CreateCinema(ctx, &entities.Cinema{
			ID:      cinemaID,
			Name:    cinemaName,
			LogoURL: "https://example.com/logo.png",
		})
		if errCreate != nil {
			// If create fails and we couldn't find it, it's a real error
			return fmt.Errorf("failed to create cinema: %w", errCreate)
		}
		cinema = &entities.Cinema{ID: cinemaID}
	}

	// 2. Seat Types
	stdType, _ := s.repo.GetSeatTypeByName(ctx, "Standard")
	if stdType == nil {
		stdID := uuid.New()
		err = s.repo.CreateSeatType(ctx, &entities.SeatType{ID: stdID, Name: "Standard", CinemaID: cinema.ID})
		if err != nil {
			log.Printf("Failed to create Standard seat type: %v", err)
		}
		stdType = &entities.SeatType{ID: stdID}
	}

	vipType, _ := s.repo.GetSeatTypeByName(ctx, "VIP")
	if vipType == nil {
		vipID := uuid.New()
		err = s.repo.CreateSeatType(ctx, &entities.SeatType{ID: vipID, Name: "VIP", CinemaID: cinema.ID})
		if err != nil {
			log.Printf("Failed to create VIP seat type: %v", err)
		}
		vipType = &entities.SeatType{ID: vipID}
	}

	// 3. Theater
	theaterName := "Theater 1"
	theater, err := s.repo.GetTheaterByName(ctx, cinema.ID, theaterName)
	if err != nil {
		theaterID := uuid.New()
		errCreate := s.repo.CreateTheater(ctx, &entities.Theater{
			ID:        theaterID,
			CinemaID:  cinema.ID,
			Name:      theaterName,
			Address:   "Jl. Letjen S. Parman No.28",
			Latitude:  -6.175392,
			Longitude: 106.827153,
		})
		if errCreate != nil {
			return fmt.Errorf("failed to create theater: %w", errCreate)
		}
		theater = &entities.Theater{ID: theaterID}
	}

	// 4. Studio
	studioName := "Studio 1 (IMAX)"
	studio, err := s.repo.GetStudioByName(ctx, theater.ID, studioName)
	if err != nil {
		studioID := uuid.New()
		errCreate := s.repo.CreateStudio(ctx, &entities.Studio{ID: studioID, TheaterID: theater.ID, Name: studioName})
		if errCreate != nil {
			return fmt.Errorf("failed to create studio: %w", errCreate)
		}
		studio = &entities.Studio{ID: studioID}

		// Seats (10x10) - Only seed if studio is new
		for row := 0; row < 10; row++ {
			rowChar := string(rune('A' + row))
			for num := 1; num <= 10; num++ {
				typeID := stdType.ID
				if row > 7 {
					typeID = vipType.ID
				}
				_ = s.repo.CreateSeat(ctx, &entities.Seat{
					ID:         uuid.New(),
					StudioID:   studio.ID,
					Row:        rowChar,
					Number:     num,
					SeatTypeID: typeID,
					Active:     true,
				})
			}
		}

	}

	// Seat Pricing
	// Check if exists for this theater
	var pricingCount int
	s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM seat_pricings WHERE theater_id = $1", theater.ID).Scan(&pricingCount)

	if pricingCount == 0 {
		log.Println("Creating seat pricing...")
		err = s.repo.CreateSeatPricing(ctx, &entities.SeatPricing{
			ID:         uuid.New(),
			Price:      50000,
			DayType:    "weekday",
			SeatTypeID: stdType.ID,
			TheaterID:  theater.ID,
		})
		if err != nil {
			log.Printf("Failed to create seat pricing (weekday): %v", err)
		}

		s.repo.CreateSeatPricing(ctx, &entities.SeatPricing{
			ID:         uuid.New(),
			Price:      75000,
			DayType:    "weekend",
			SeatTypeID: stdType.ID,
			TheaterID:  theater.ID,
		})
	}

	return nil
}

func (s *Seeder) SeedMovies(ctx context.Context) error {
	log.Println("Seeding movies from TMDB...")

	genres, err := s.tmdb.GetGenres()
	if err != nil {
		log.Printf("Warning: failed to fetch genres: %v", err)
	}

	// Seed genres
	for _, name := range genres {
		_ = s.repo.UpsertGenre(ctx, &entities.MovieGenre{ID: uuid.New(), Genre: name})
	}

	rawMovies, err := s.tmdb.FetchNowPlayingRaw()
	if err != nil {
		return err
	}

	// Get all statuses and ratings for random selection
	statusNames := []string{"Now Showing", "Coming Soon", "Ended"}
	ratingNames := []string{"G", "PG", "PG-13", "R", "NC-17"}

	// Fetch all status IDs
	var statuses []*entities.MovieStatus
	for _, name := range statusNames {
		st, _ := s.repo.GetStatusByName(ctx, name)
		if st != nil {
			statuses = append(statuses, st)
		}
	}

	// Fetch all rating IDs
	var ratings []*entities.MovieRating
	for _, name := range ratingNames {
		rt, _ := s.repo.GetRatingByName(ctx, name)
		if rt != nil {
			ratings = append(ratings, rt)
		}
	}

	// Fallback if none found
	if len(statuses) == 0 {
		sid := uuid.New()
		s.repo.CreateStatus(ctx, &entities.MovieStatus{ID: sid, Status: "Now Showing"})
		statuses = append(statuses, &entities.MovieStatus{ID: sid})
	}
	if len(ratings) == 0 {
		rid := uuid.New()
		s.repo.CreateRating(ctx, &entities.MovieRating{ID: rid, Rating: "PG-13"})
		ratings = append(ratings, &entities.MovieRating{ID: rid})
	}

	for _, m := range rawMovies {
		// Randomly select status and rating
		randomStatus := statuses[rand.Intn(len(statuses))]
		randomRating := ratings[rand.Intn(len(ratings))]

		movieID := uuid.New()
		err := s.repo.CreateMovie(ctx, &entities.Movie{
			ID:            movieID,
			Title:         m.Title,
			Tagline:       "", // TMDB list doesn't have tagline, detail does.
			Overview:      m.Overview,
			PosterURL:     tmdb.ImageBase + m.PosterPath,
			BackdropURL:   tmdb.ImageBase + m.BackdropPath,
			TrailerURL:    "",  // Requires another call
			Duration:      120, // List doesn't have runtime
			Popularity:    int(m.Popularity),
			MovieStatusID: randomStatus.ID,
			MovieRatingID: randomRating.ID,
		})
		if err != nil {
			log.Printf("Failed to seed movie %s: %v", m.Title, err)
		}
	}
	return nil
}

func (s *Seeder) SeedShowtimes(ctx context.Context) error {
	log.Println("Seeding showtimes...")

	var studioID uuid.UUID
	err := s.db.QueryRowContext(ctx, "SELECT id FROM studios LIMIT 1").Scan(&studioID)
	if err != nil {
		return fmt.Errorf("no studio found")
	}

	rows, err := s.db.QueryContext(ctx, "SELECT id FROM movies")
	if err != nil {
		return err
	}
	defer rows.Close()

	var movieIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err == nil {
			movieIDs = append(movieIDs, id)
		}
	}

	// Seed 3 showtimes per movie for today
	var pricingID uuid.UUID
	err = s.db.QueryRowContext(ctx, "SELECT id FROM seat_pricings LIMIT 1").Scan(&pricingID)
	if err != nil {
		return fmt.Errorf("no seat pricing found, seeding failed")
	}

	for _, mid := range movieIDs {
		for i := 0; i < 3; i++ {
			start := time.Now().Add(time.Duration(i*3) * time.Hour)
			_ = s.repo.CreateShowtime(ctx, &entities.Showtime{
				ID:            uuid.New(),
				MovieID:       mid,
				StudioID:      studioID,
				Time:          start,
				ExpiredAt:     start.Add(2 * time.Hour),
				SeatPricingID: pricingID,
				Status:        true,
			})
		}
	}
	return nil
}
