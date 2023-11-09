package driver

import (
	"context"
	"database/sql/driver"
	"errors"
	"github.com/stretchr/testify/require"
	pb "gitlab.com/aous.omr/sqlite-og/gen/proto"
	"gitlab.com/aous.omr/sqlite-og/mocks"
	"go.uber.org/mock/gomock"
	"golang.org/x/exp/rand"
	"testing"
)

const testCnxId = "test-id"

// Connection

func TestSQLiteOGConn_Ping(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	ogclient := mocks.NewMockSqliteOGClient(ctrl)
	ogclient.EXPECT().Ping(gomock.Eq(ctx), gomock.Eq(&pb.Empty{})).Return(&pb.Empty{}, nil).Times(1)
	conn := &SQLiteOGConn{
		ID:       testCnxId,
		GRPCConn: nil,
		OGClient: ogclient,
	}
	err := conn.Ping(ctx)
	require.NoError(t, err)
}

func TestSQLiteOGConn_ExecContext(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	ogclient := mocks.NewMockSqliteOGClient(ctrl)
	conn := &SQLiteOGConn{
		ID:       testCnxId,
		GRPCConn: nil,
		OGClient: ogclient,
	}
	// doesn't matter what the sql is
	sql := "insert into test_table values (?,?)"

	t.Run("execute with no args", func(t *testing.T) {
		ogclient.EXPECT().Execute(gomock.Eq(ctx), gomock.Eq(&pb.Statement{
			Sql:    sql,
			Params: []string{},
			CnxId:  testCnxId,
		})).Return(&pb.ExecuteResult{
			LastInsertId: 1,
			AffectedRows: 1,
		}, nil).Times(1)

		res, err := conn.ExecContext(ctx, sql, nil)
		require.NotNil(t, res)
		require.NoError(t, err)

		lastID, err := res.LastInsertId()
		require.NoError(t, err)
		require.Equal(t, int64(1), lastID)

		affected, err := res.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)
	})

	t.Run("execute with args", func(t *testing.T) {
		args := []driver.NamedValue{
			{
				Name:    "",
				Ordinal: 1,
				Value:   "one",
			},
			{
				Name:    "",
				Ordinal: 2,
				Value:   2,
			},
		}
		params := []string{"one", "2"}

		ogclient.EXPECT().Execute(gomock.Eq(ctx), gomock.Eq(&pb.Statement{
			Sql:    sql,
			Params: params,
			CnxId:  testCnxId,
		})).Return(&pb.ExecuteResult{
			LastInsertId: 1,
			AffectedRows: 1,
		}, nil).Times(1)

		res, err := conn.ExecContext(ctx, sql, args)
		require.NotNil(t, res)
		require.NoError(t, err)

		lastID, err := res.LastInsertId()
		require.NoError(t, err)
		require.Equal(t, int64(1), lastID)

		affected, err := res.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)
	})

	t.Run("execute errors", func(t *testing.T) {
		ogclient.EXPECT().Execute(gomock.Eq(ctx), gomock.Any()).Return(nil, errors.New("an error")).Times(1)

		res, err := conn.ExecContext(ctx, sql, nil)
		require.Nil(t, res)
		require.Error(t, err)
		require.Equal(t, "an error", err.Error())
	})
}

func TestSQLiteOGConn_QueryContext(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	ogclient := mocks.NewMockSqliteOGClient(ctrl)
	conn := &SQLiteOGConn{
		ID:       testCnxId,
		GRPCConn: nil,
		OGClient: ogclient,
	}
	sql := "select col1,col2 from test_table"
	pbResult := &pb.QueryResult{
		Columns:     []string{"col1", "col2"},
		ColumnTypes: []string{"TEXT", "INTEGER"},
		Rows: []*pb.Row{
			{
				Fields: []string{
					"row1field1", "row1field2",
				},
			},
			{
				Fields: []string{
					"row2field1", "row2field2",
				},
			},
		},
	}

	t.Run("query with no args", func(t *testing.T) {
		ogclient.EXPECT().Query(gomock.Eq(ctx), gomock.Eq(&pb.Statement{
			Sql:    sql,
			Params: []string{},
			CnxId:  testCnxId,
		})).Return(pbResult, nil).Times(1)

		rows, err := conn.QueryContext(ctx, sql, nil)
		require.NotNil(t, rows)
		require.NoError(t, err)
		require.Equal(t, pbResult.Columns, rows.Columns())
	})

	t.Run("query with args", func(t *testing.T) {
		args := []driver.NamedValue{
			{
				Name:    "",
				Ordinal: 2,
				Value:   2,
			},
			{
				Name:    "",
				Ordinal: 1,
				Value:   "one",
			},
		}
		params := []string{"one", "2"}

		ogclient.EXPECT().Query(gomock.Eq(ctx), gomock.Eq(&pb.Statement{
			Sql:    sql,
			Params: params,
			CnxId:  testCnxId,
		})).Return(pbResult, nil).Times(1)

		res, err := conn.QueryContext(ctx, sql, args)
		require.NotNil(t, res)
		require.NoError(t, err)
	})

	t.Run("query errors", func(t *testing.T) {

		ogclient.EXPECT().Query(gomock.Eq(ctx), gomock.Any()).Return(nil, errors.New("an error")).Times(1)

		res, err := conn.QueryContext(ctx, sql, nil)
		require.Nil(t, res)
		require.Error(t, err)
		require.Equal(t, "an error", err.Error())
	})
}

// Result

func TestResult(t *testing.T) {
	id := rand.Int63n(int64(100))
	affected := rand.Int63n(int64(100))
	pbResult := &pb.ExecuteResult{
		LastInsertId: id,
		AffectedRows: affected,
	}
	expected := &Result{pbResult}

	t.Run("resultFromPb no error", func(t *testing.T) {
		res, err := resultFromPB(pbResult)
		require.NoError(t, err)
		require.Equal(t, expected, res)

		lastInsertId, err1 := res.LastInsertId()
		affectedRows, err2 := res.RowsAffected()

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.Equal(t, id, lastInsertId)
		require.Equal(t, affected, affectedRows)
	})

	t.Run("resultFromPb errors when pb is nil", func(t *testing.T) {
		res, err := resultFromPB(nil)
		require.Error(t, err)
		require.Nil(t, res)
	})

}

// Rows

func TestRows(t *testing.T) {

}
