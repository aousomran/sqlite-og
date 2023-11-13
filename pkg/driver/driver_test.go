package driver

import (
	"database/sql"
	"fmt"
	pb "github.com/aousomran/sqlite-og/gen/proto"
	"github.com/aousomran/sqlite-og/internal/connections"
	"github.com/aousomran/sqlite-og/internal/server"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"os"
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
}

const databaseName = "tester.db"

var Listener net.Listener = nil

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
	db, err := sql.Open("sqlite3", databaseName)
	defer func() {
		_ = db.Close()
	}()
	if err != nil {
		t.Fatal(err)
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

	// setup listener & grpc server
	Listener, err = net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	if Listener == nil {
		t.Fatal("could not create listener")
	}
	t.Logf("listener started, address %v", Listener.Addr().String())

	s := grpc.NewServer()

	manager := connections.NewManager()
	srv := server.New(manager)
	pb.RegisterSqliteOGServer(s, srv)

	go func() {
		err = s.Serve(Listener)
		if err != nil {
			t.Errorf("grpc serve error %v", err)
			return
		}
		t.Logf("grpc server listening on %s", Listener.Addr().String())
	}()

	// return teardown function
	return func(t *testing.T) {
		_ = manager.Close()
		s.GracefulStop()
		_ = Listener.Close()
		errRemove := os.Remove(fmt.Sprintf("%s", databaseName))
		if errRemove != nil {
			t.Logf("error removing datafile %s", errRemove)
		}
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

	db, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	students, err := queryStudents(db)
	require.NoError(t, err)
	require.NotNil(t, students)

	dsn := fmt.Sprintf("%s/%s", Listener.Addr().String(), databaseName)
	db2, err := sql.Open("sqliteog", dsn)
	defer db2.Close()

	students2, err2 := queryStudents(db2)
	require.NoError(t, err2)
	require.NotNil(t, students2)
}
