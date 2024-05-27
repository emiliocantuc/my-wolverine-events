// Handles everything related to the database

package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// The database connection
type DB struct {
	*sql.DB
}

// Connect to the SQLite database
func ConnectDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// User

func (db *DB) UserExists(email string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (db *DB) AddUser(email string) (int64, error) {

	stmt, err := db.Prepare("INSERT INTO users (email) VALUES (?)")
	if err != nil {
		return 0, fmt.Errorf("could not prepare statement: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(email)
	if err != nil {
		return 0, fmt.Errorf("could not execute statement: %v", err)
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("could not retrieve last insert ID: %v", err)
	}

	return userId, nil
}

func (db *DB) DeleteUser(userID int64) error {

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}

	_, err = tx.Exec("DELETE FROM preferences WHERE user_id = ?", userID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not delete preferences: %v", err)
	}

	_, err = tx.Exec("DELETE FROM users WHERE user_id = ?", userID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not delete user: %v", err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("could not commit transaction: %v", err)
	}
	return nil
}

// Preferences
type Preferences struct {
	EducationLevel  string
	Program         string
	CampusLocation  string
	Interests       string
	IncSeminars     bool
	IncSports       bool
	IncSocial       bool
	SendEmail       bool
	KeywordsToAvoid string
}

func (db *DB) GetPreferences(userId int64) (Preferences, error) {

	var prefs Preferences
	var incSeminars, incSports, incSocial, sendEmail int

	query := `
        SELECT educationLevel, program, campusLocation, interests,
               incSeminars, incSports, incSocial, sendEmail, keywordsToAvoid
        FROM preferences
        WHERE user_id = ?
    `

	row := db.QueryRow(query, userId)

	err := row.Scan(&prefs.EducationLevel, &prefs.Program, &prefs.CampusLocation, &prefs.Interests,
		&incSeminars, &incSports, &incSocial, &sendEmail, &prefs.KeywordsToAvoid)
	if err != nil {
		if err == sql.ErrNoRows {
			return Preferences{}, fmt.Errorf("no preferences found for user_id %d", userId)
		}
		return Preferences{}, fmt.Errorf("could not retrieve preferences: %v", err)
	}

	prefs.IncSeminars = incSeminars == 1
	prefs.IncSports = incSports == 1
	prefs.IncSocial = incSocial == 1
	prefs.SendEmail = sendEmail == 1

	return prefs, nil
}

func (db *DB) UpdatePreference(userId int64, preferenceName, preferenceValue string) error {

	// TODO validate preferenceName
	query := fmt.Sprintf(`
        INSERT INTO preferences (user_id, %s)
        VALUES (?, ?)
        ON CONFLICT(user_id)
        DO UPDATE SET %s = excluded.%s
    `, preferenceName, preferenceName, preferenceName)

	_, err := db.Exec(query, userId, preferenceValue)
	if err != nil {
		return fmt.Errorf("could not update preference %s: %v", preferenceName, err)
	}

	return nil
}

// Events

// Need for templeting
type EventCard struct {
	Id           int
	Title        string
	Subtitle     string
	EventType    string
	Description  string
	StartDate    string
	VoteDiff     int
	CalendarLink string
	PermaLink    string
	BuildingName string
	LoggedIn     bool
}

func formatEvent(event EventCard) EventCard {
	truncate := func(s string, max int) string {
		if len(s) > max {
			return s[:max] + "..."
		}
		return s
	}

	formatTime := func(startDate string) string {
		t, err := time.Parse(time.RFC3339, startDate)
		if err != nil {
			return ""
		}
		return t.Format("3:04 PM")
	}

	cleanString := func(s, delimiter, toRemove string) string {
		if idx := strings.Index(s, delimiter); idx != -1 {
			s = s[:idx]
		}
		return strings.TrimSpace(strings.Replace(s, toRemove, "", -1))
	}

	// Truncate title and description
	event.Title = truncate(event.Title, 75)
	event.Description = truncate(event.Description, 300)

	// Format EventType and BuildingName
	eventType := cleanString(event.EventType, "/", "")
	buildingName := event.BuildingName //cleanString(event.BuildingName, "", "location")

	// Format Start Time
	startTime := formatTime(event.StartDate)

	// Construct Subtitle
	var subtitleParts []string
	if eventType != "" {
		subtitleParts = append(subtitleParts, eventType)
	}
	if startTime != "" {
		subtitleParts = append(subtitleParts, startTime)
	}
	if buildingName != "" {
		subtitleParts = append(subtitleParts, buildingName)
	}
	event.Subtitle = strings.Join(subtitleParts, " | ")

	return event
}

// GetMaxNweek fetches the latest nweek from the statistics table
func (db *DB) GetMaxNweek() (int, error) {
	query := `SELECT MAX(nweek) FROM statistics`
	var maxNweek int
	err := db.QueryRow(query).Scan(&maxNweek)
	if err != nil {
		return 0, err
	}
	return maxNweek, nil
}

func unmarshallEvents(rows *sql.Rows) ([]EventCard, error) {
	var events []EventCard
	for rows.Next() {
		var event EventCard
		var voteDiff sql.NullInt64
		err := rows.Scan(&event.Id, &event.Title, &event.EventType, &event.Description, &event.StartDate, &voteDiff, &event.CalendarLink, &event.PermaLink, &event.BuildingName)
		if err != nil {
			return nil, err
		}
		if voteDiff.Valid {
			event.VoteDiff = int(voteDiff.Int64)
		} else {
			event.VoteDiff = 0
		}
		events = append(events, formatEvent(event))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

// GetTopEvents fetches the events with the highest up-down votes for the latest nweek
func (db *DB) GetTopEvents(n int) ([]EventCard, error) {

	// TODO put a cache in front of this ?

	maxNweek, err := db.GetMaxNweek()
	if err != nil {
		return nil, err
	}

	query := `
        SELECT 
            e.event_id AS Id,
            e.title AS Title,
			e.type AS EventType,
            e.event_description AS Description,
            e.event_date AS StartDate,
            COALESCE(SUM(CASE WHEN v.vote_type = 'U' THEN 1 ELSE 0 END), 0) - 
            COALESCE(SUM(CASE WHEN v.vote_type = 'D' THEN 1 ELSE 0 END), 0) AS VoteDiff,
            e.gcal_link AS CalendarLink,
            e.permalink AS PermaLink,
            e.building_name AS BuildingName
        FROM 
            events e
        LEFT JOIN 
            votes v ON e.event_id = v.event_id
        WHERE
            e.nweek = ?
        GROUP BY 
            e.event_id, e.title, e.event_description, e.event_date, e.gcal_link, e.permalink, e.building_name
        ORDER BY 
            VoteDiff DESC
        LIMIT ?
    `

	rows, err := db.Query(query, maxNweek, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events, err := unmarshallEvents(rows)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (db *DB) GetRecommendedEvents(userId int64) ([]EventCard, error) {
	query := `
        SELECT e.event_id, e.title, e.type, e.event_description, e.event_date, 
               COALESCE(SUM(CASE WHEN v.vote_type = 'U' THEN 1 WHEN v.vote_type = 'D' THEN -1 ELSE 0 END), 0) as vote_diff,
               e.gcal_link, e.permalink, e.building_name
        FROM recommended_events re
        JOIN events e ON re.event_id = e.event_id
        LEFT JOIN votes v ON e.event_id = v.event_id
        WHERE re.user_id = ?
        GROUP BY e.event_id
        ORDER BY vote_diff DESC, e.event_date DESC;
    `

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("could not query recommended events: %v", err)
	}
	defer rows.Close()

	events, err := unmarshallEvents(rows)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// Vote
func (db *DB) Vote(userId int64, eventId int, voteType string) error {

	if voteType != "U" && voteType != "D" && voteType != "C" {
		return fmt.Errorf("invalid vote type: %s", voteType)
	}

	query := `
        INSERT INTO votes (user_id, event_id, vote_type, voted_at)
        VALUES (?, ?, ?, ?)
        ON CONFLICT(user_id, event_id)
        DO UPDATE SET vote_type = excluded.vote_type, voted_at = excluded.voted_at
    `
	_, err := db.Exec(query, userId, eventId, voteType, time.Now())
	return err
}
