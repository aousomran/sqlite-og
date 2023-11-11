package driver

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/require"
	"log"
	"math/rand"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Student struct {
	ID               int
	Name             string
	Age              int
	Height           float64
	IsStudent        int
	BirthDate        string
	RegistrationTime string
	//Notes            []byte
}

const databaseName = "tester.db"

var filePath = fmt.Sprintf("../../%s", databaseName)

func insertRandomData(db *sql.DB, numRows int) error {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Define a function to generate random dates within a specified range
	randomDate := func(min, max time.Time) time.Time {
		delta := max.Sub(min)
		randomDelta := time.Duration(rand.Int63n(int64(delta)))
		return min.Add(randomDelta)
	}

	// Insert random data into the example_table
	for i := 0; i < numRows; i++ {
		name := fmt.Sprintf("Person%d", i+1)
		age := rand.Intn(50) + 20
		height := rand.Float64()*50 + 150
		isStudent := rand.Intn(2)
		birthDate := randomDate(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), time.Now())
		registrationTime := time.Now().Add(-time.Duration(rand.Intn(365 * 24 * int(time.Hour))))
		//notes := make([]byte, 8)
		//rand.Read(notes)

		_, err := db.Exec(`
			INSERT INTO example_table (name, age, height, is_student, birth_date, registration_time)
			VALUES (?, ?, ?, ?, ?, ?)
		`, name, age, height, isStudent, birthDate.Format("2006-01-02"), registrationTime.Format("2006-01-02 15:04:05"))
		if err != nil {
			return err
		}
	}

	return nil
}

func setupSuite(t *testing.T) func(t *testing.T) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Create the example_table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS example_table (
			id INTEGER PRIMARY KEY,
			name TEXT,
			age INTEGER,
			height REAL,
			is_student INTEGER,
			birth_date DATE,
			registration_time TIMESTAMP
			--notes BLOB
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Insert random data into the example_table for testing
	err = insertRandomData(db, 2)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("database setup for testing completed.")
	return func(t *testing.T) {
		_ = db.Close()
	}
}

func queryStudents(db *sql.DB) ([]Student, error) {
	query := `SELECT * FROM example_table LIMIT 10`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []Student

	for rows.Next() {
		var student Student
		errScan := rows.Scan(
			&student.ID, &student.Name, &student.Age,
			&student.Height, &student.IsStudent, &student.BirthDate,
			&student.RegistrationTime)
		//&student.Notes)
		if errScan != nil {
			return nil, errScan
		}
		result = append(result, student)
	}
	return result, nil
}

func TestQuery(t *testing.T) {
	teardown := setupSuite(t)
	defer teardown(t)

	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	students, err := queryStudents(db)
	require.NoError(t, err)
	require.NotNil(t, students)

	db2, err := sql.Open("sqliteog", fmt.Sprintf("localhost:50051/%s", databaseName))
	defer db2.Close()

	students2, err2 := queryStudents(db2)
	require.NoError(t, err2)
	require.NotNil(t, students2)

	require.Equal(t, students, students2)

}
