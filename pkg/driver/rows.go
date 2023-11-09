package driver

import (
	"database/sql/driver"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"time"

	pb "github.com/aousomran/sqlite-og/gen/proto"
)

type Rows struct {
	mu     sync.RWMutex
	pbr    *pb.QueryResult
	index  int
	closed bool
}

func rowsFromPB(pbResult *pb.QueryResult) (*Rows, error) {
	if pbResult == nil {
		return nil, fmt.Errorf("empty pbResult")
	}
	return &Rows{
		mu:    sync.RWMutex{},
		pbr:   pbResult,
		index: 0,
	}, nil
}

func (r *Rows) Columns() []string {
	return r.pbr.GetColumns()
}

func (r *Rows) Close() error {
	r.closed = true
	return nil
}

func (r *Rows) Next(dest []driver.Value) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// if we've reached the end of the result set return io.EOF
	if r.index == len(r.pbr.GetRows()) {
		return io.EOF
	}
	for k, v := range r.pbr.GetRows()[r.index].Fields {
		dest[k] = v
		typeName := r.ColumnTypeDatabaseTypeName(k)
		if typeName == "DATETIME" || typeName == "DATE" {
			t, err := time.Parse(time.RFC3339, v)

			fmt.Printf("before: %s, after: %+v | %s", v, t, t.Format(time.DateTime))
			if err != nil {
				fmt.Printf("error parsing time: %s", err.Error())
			}
			dest[k] = t
		}
	}
	// increment the index so that the subsequent call to .Next()
	// points to the correct result set
	r.index++
	return nil
}

func (r *Rows) ColumnTypeDatabaseTypeName(index int) string {
	return strings.ToUpper(strings.Split(r.pbr.ColumnTypes[index], "(")[0])
}

func (r *Rows) ColumnTypeScanType(index int) reflect.Type {
	// TODO: what if columns are nullable ??
	switch r.ColumnTypeDatabaseTypeName(index) {
	case "INTEGER", "INT", "TINYINT", "BIGINT":
		return reflect.TypeOf(int64(0))
	case "REAL", "FLOAT", "DOUBLE":
		return reflect.TypeOf(float64(0))
	case "TEXT", "BLOB", "CHAR", "VARCHAR":
		return reflect.TypeOf("")
	case "BOOL":
		return reflect.TypeOf(true)
	case "TIME", "DATETIME", "DATE":
		return reflect.TypeOf(time.Now())
	default:
		return reflect.TypeOf("")
	}
}
