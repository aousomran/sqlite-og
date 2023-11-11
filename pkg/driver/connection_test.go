package driver

import (
	"context"
	"database/sql/driver"
	"fmt"
	pb "github.com/aousomran/sqlite-og/gen/proto"
	"github.com/aousomran/sqlite-og/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

const testConnectionId = "fc11e35e-932d-48eb-9c53-1f9e20c8d514"

func TestSQLiteOGConn_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockSqliteOGClient(ctrl)
	ctx := context.Background()
	conn := &SQLiteOGConn{
		ID:       testConnectionId,
		GRPCConn: nil,
		OGClient: client,
	}
	client.EXPECT().Ping(gomock.Any(), &pb.Empty{}).Return(&pb.Empty{}, nil).Times(1)
	err := conn.Ping(ctx)
	assert.NoError(t, err)
}

func TestSQLiteOGConn_ResetSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockSqliteOGClient(ctrl)
	ctx := context.Background()
	conn := &SQLiteOGConn{
		ID:       testConnectionId,
		GRPCConn: nil,
		OGClient: client,
	}
	cnxId := &pb.ConnectionId{Id: testConnectionId}
	client.EXPECT().ResetSession(gomock.Any(), cnxId).Return(cnxId, nil).Times(1)
	err := conn.ResetSession(ctx)
	assert.NoError(t, err)
}

func TestSQLiteOGConn_IsValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockSqliteOGClient(ctrl)
	conn := &SQLiteOGConn{
		ID:       testConnectionId,
		GRPCConn: nil,
		OGClient: client,
	}
	cnxId := &pb.ConnectionId{Id: testConnectionId}
	client.EXPECT().IsValid(gomock.Any(), cnxId).Return(&pb.Empty{}, nil).Times(1)
	isValid := conn.IsValid()
	assert.True(t, isValid)
	client.EXPECT().IsValid(gomock.Any(), cnxId).Return(&pb.Empty{}, fmt.Errorf("an error")).Times(1)
	isValid = conn.IsValid()
	assert.False(t, isValid)
}

func TestSQLiteOGConn_Prepare(t *testing.T) {
	conn := &SQLiteOGConn{}
	t.Run("sql with 2 qmark returns a statement", func(t *testing.T) {
		sql := `select * from mytable where mycolumn=? and myothercolumn=?`
		stmt, err := conn.Prepare(sql)
		assert.NoError(t, err)
		assert.NotNil(t, stmt)
		assert.Equal(t, 2, stmt.NumInput())
	})

	t.Run("sql with no qmark returns a statement", func(t *testing.T) {
		sql := `select * from mytable`
		stmt, err := conn.Prepare(sql)
		assert.NoError(t, err)
		assert.NotNil(t, stmt)
		assert.Equal(t, 0, stmt.NumInput())
	})

	t.Run("empty query returns an error", func(t *testing.T) {
		stmt, err := conn.Prepare("")
		assert.Error(t, err)
		assert.Nil(t, stmt)
	})
}

func TestSQLiteOGConn_QueryContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	client := mocks.NewMockSqliteOGClient(ctrl)
	conn := &SQLiteOGConn{
		ID:       testConnectionId,
		GRPCConn: nil,
		OGClient: client,
	}
	sql := `select * from mytable where mycolumn=? and myothercolumn=?`
	args := []driver.NamedValue{
		{
			Ordinal: 1,
			Value:   "first",
		},
	}

	client.EXPECT().Query(gomock.Any(), gomock.Any()).Return(&pb.QueryResult{
		Columns:     []string{"col1", "col2"},
		ColumnTypes: []string{"NUMBER", "TEXT"},
		Rows: []*pb.Row{
			{
				Fields: []string{"1", "two"},
			},
			{
				Fields: []string{"3", "four"},
			},
		},
	}, nil).Times(1)
	rows, err := conn.QueryContext(ctx, sql, args)
	assert.NoError(t, err)
	assert.NotNil(t, rows)

}

func TestSQLiteOGConn_ExecContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockSqliteOGClient(ctrl)
	conn := &SQLiteOGConn{
		ID:       testConnectionId,
		GRPCConn: nil,
		OGClient: client,
	}
	ctx := context.Background()
	sql := `update mytable set mycolumn=? where myothercolumn=?`
	args := []driver.NamedValue{
		{
			Ordinal: 1,
			Value:   1,
		},
		{
			Ordinal: 2,
			Value:   "two",
		},
	}
	client.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(&pb.ExecuteResult{
		LastInsertId: 100,
		AffectedRows: 1,
	}, nil).Times(1)
	result, err := conn.ExecContext(ctx, sql, args)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	lastInsertId, err2 := result.LastInsertId()
	assert.NoError(t, err2)
	rowsAffected, err3 := result.RowsAffected()
	assert.NoError(t, err3)
	assert.Equal(t, int64(100), lastInsertId)
	assert.Equal(t, int64(1), rowsAffected)
}

func Test_namedValuesToParams(t *testing.T) {
	t.Run("namedValues are in order", func(t *testing.T) {
		namedValues := []driver.NamedValue{
			{
				Ordinal: 1,
				Value:   "one",
			},
			{
				Ordinal: 2,
				Value:   2,
			},
			{
				Ordinal: 3,
				Value:   "three",
			},
		}
		expected := []string{"one", "2", "three"}
		actual, err := namedValuesToParams(namedValues)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, expected, actual)
	})

	t.Run("named values are not in order", func(t *testing.T) {
		namedValues := []driver.NamedValue{
			{
				Ordinal: 3,
				Value:   "three",
			},
			{
				Ordinal: 1,
				Value:   "one",
			},
			{
				Ordinal: 2,
				Value:   2,
			},
		}
		expected := []string{"one", "2", "three"}
		actual, err := namedValuesToParams(namedValues)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, expected, actual)
	})

	t.Run("one of the named values has an invalid ordinal", func(t *testing.T) {
		namedValues := []driver.NamedValue{
			{
				Ordinal: 3,
				Value:   "three",
			},
			{
				Ordinal: 0,
				Value:   "one",
			},
			{
				Ordinal: 2,
				Value:   2,
			},
		}
		expected := []string{"one", "2", "three"}
		actual, err := namedValuesToParams(namedValues)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.NotEqual(t, expected, actual)
	})
}
