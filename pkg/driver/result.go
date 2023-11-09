package driver

import (
	"fmt"
	pb "gitlab.com/aous.omr/sqlite-og/gen/proto"
)

type Result struct {
	pbr *pb.ExecuteResult
}

func resultFromPB(pbResult *pb.ExecuteResult) (*Result, error) {
	if pbResult == nil {
		return nil, fmt.Errorf("empty pbResult")
	}
	return &Result{
		pbResult,
	}, nil
}

func (r *Result) LastInsertId() (int64, error) {
	return r.pbr.GetLastInsertId(), nil
}

func (r *Result) RowsAffected() (int64, error) {
	return r.pbr.GetAffectedRows(), nil
}
