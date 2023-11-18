package driver

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	pb "github.com/aousomran/sqlite-og/gen/proto"
	"github.com/aousomran/sqlite-og/internal/connections"
	"github.com/aousomran/sqlite-og/internal/server"
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
var sqliteDB *sql.DB
var ogDB *sql.DB
var grpcServer *grpc.Server
var connectionManager *connections.Manager

func insertRandomData(db *sql.DB, numRows int) error {
	rand.Seed(time.Now().UnixNano())

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

func queryStudents(db *sql.DB, sql string) ([]Student, error) {
	rows, err := db.Query(sql)
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

func setupDBConnections(t *testing.T) {
	db1, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		require.NoError(t, err)
	}
	sqliteDB = db1

	dsn := fmt.Sprintf("%s/%s", Listener.Addr().String(), databaseName)
	db2, err := sql.Open("sqliteog", dsn)
	if err != nil {
		require.NoError(t, err)
	}
	ogDB = db2
}

func setupGRPCServer(t *testing.T) {
	// setup listener
	var err error
	Listener, err = net.Listen("tcp", "localhost:51515")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	if Listener == nil {
		t.Fatal("could not create listener")
	}

	//  setup grpc server
	grpcServer = grpc.NewServer()

	connectionManager = connections.NewManager()
	ogServer := server.New(connectionManager)
	pb.RegisterSqliteOGServer(grpcServer, ogServer)

	go func() {
		errServe := grpcServer.Serve(Listener)
		if errServe != nil {
			require.NoError(t, errServe)
			return
		}
		t.Logf("grpc server listening on %s", Listener.Addr().String())
	}()
	return
}

func setupDatabase(t *testing.T) {
	// Create the example_table
	_, err := sqliteDB.Exec(`
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
		t.Fatal(err)
	}

	// Insert random data into the example_table for testing
	err = insertRandomData(sqliteDB, 10)
	if err != nil {
		t.Fatal(err)
	}
}

func setupSuite(t *testing.T) func(t *testing.T) {
	// setup globals to be used in tests
	setupGRPCServer(t)
	setupDBConnections(t)
	setupDatabase(t)

	// return teardown function
	return func(t *testing.T) {
		_ = sqliteDB.Close()
		_ = ogDB.Close()
		_ = connectionManager.Close()
		grpcServer.GracefulStop()
		_ = os.Remove(fmt.Sprintf("%s", databaseName))
	}
}

func TestCRUD(t *testing.T) {
	teardown := setupSuite(t)
	defer teardown(t)

	t.Run("selected records are equal", func(t *testing.T) {
		query := `SELECT * FROM example_table LIMIT 10`
		expected, err := queryStudents(sqliteDB, query)
		actual, err2 := queryStudents(ogDB, query)

		require.NoError(t, err)
		require.NoError(t, err2)
		require.NotNil(t, expected)
		require.NotNil(t, actual)
		require.Equal(t, expected, actual)
	})

	t.Run("insert works", func(t *testing.T) {
		query := `SELECT * FROM example_table ORDER BY id DESC LIMIT 1`
		before, err := queryStudents(sqliteDB, query)
		require.NoError(t, err)

		err = insertRandomData(ogDB, 1)
		require.NoError(t, err)

		after, err := queryStudents(sqliteDB, query)
		require.NoError(t, err)

		require.NotEqual(t, before, after)
		require.Greater(t, after[0].ID, before[0].ID)
	})

	t.Run("update works", func(t *testing.T) {
		query := `SELECT * FROM example_table ORDER BY id DESC LIMIT 1`
		before, err := queryStudents(sqliteDB, query)
		require.NoError(t, err)

		name := "ao"
		updateResult, err := ogDB.Exec(`UPDATE example_table SET name=? WHERE id=?`, name, before[0].ID)
		require.NoError(t, err)
		affected, err := updateResult.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, 1, int(affected))

		after, err := queryStudents(sqliteDB, query)
		require.NoError(t, err)

		require.NotEqual(t, before[0].Name, after[0].Name)
		require.Equal(t, name, after[0].Name)
	})

	t.Run("delete works", func(t *testing.T) {
		query := `SELECT * FROM example_table ORDER BY id DESC LIMIT 1`
		before, err := queryStudents(ogDB, query)
		require.NoError(t, err)

		deleteResult, err := ogDB.Exec("delete from example_table where id=?", fmt.Sprintf("%d", before[0].ID))
		require.NoError(t, err)
		affected, err := deleteResult.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, 1, int(affected))

		after, err := queryStudents(sqliteDB, query)
		require.NoError(t, err)

		require.NotEqual(t, before, after)
		require.Less(t, after[0].ID, before[0].ID)
	})

	t.Run("test callbacks", func(t *testing.T) {
		sayHello := func(args ...string) []string {
			return []string{"hello " + args[0]}
		}

		sql.Register("og_custom", &SQLiteOGDriver{
			Funcs: map[string]callbackFunc{
				"say_hello": sayHello,
			},
			CallbacksEnabled: true,
		})

		dsn := fmt.Sprintf("%s/%s", Listener.Addr().String(), databaseName)
		db, err := sql.Open("og_custom", dsn)

		rows, err := db.Query(`select say_hello(?)`, "ao")
		require.NoError(t, err)
		result := ""
		for rows.Next() {
			errScan := rows.Scan(&result)
			require.NoError(t, errScan)
		}
		require.Equal(t, "hello ao", result)
		err = db.Close()
		require.NoError(t, err)
		t.Log("test done")
	})
}
