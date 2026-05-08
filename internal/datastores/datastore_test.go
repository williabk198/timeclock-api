package datastores

import "github.com/DATA-DOG/go-sqlmock"

type queryExpectationsFunc func(sqlmock.Sqlmock)

func (qef queryExpectationsFunc) SetExpectactions(s sqlmock.Sqlmock) {
	qef(s)
}
