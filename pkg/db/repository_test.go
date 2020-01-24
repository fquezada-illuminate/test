package db

import (
	"database/sql/driver"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gocraft/dbr"
	"github.com/gocraft/dbr/dialect"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/structs"
	"reflect"
	"regexp"
	"testing"
)

type MockObject struct {
	Id   string `json:"id" db:"id" structs:"id"`
	Name string `json:"name" db:"name"`
}

type MockObjectNoTags struct {
	Id   string
	Name string
}

func TestNewRepository(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()
	conn := &dbr.Connection{DB: db, Dialect: dialect.PostgreSQL, EventReceiver: &dbr.NullEventReceiver{}}
	sess := conn.NewSession(nil)

	sh := structs.Helper{}
	table := "resource"

	repo := NewRepository(sess, sh, table)

	if _, ok := repo.(*BaseRepository); !ok {
		t.Errorf("Exepected %s, got %s", reflect.TypeOf(&BaseRepository{}), reflect.TypeOf(repo))
	}
}

func TestBaseRepository_Find(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	conn := &dbr.Connection{DB: db, Dialect: dialect.PostgreSQL, EventReceiver: &dbr.NullEventReceiver{}}
	sess := conn.NewSession(nil)

	sh := structs.Helper{}
	table := "resource"

	expectedResult := &MockObject{"123", "Test Name"}

	repo := NewRepository(sess, sh, table)

	buff := dbr.NewBuffer()
	sess.
		Select("*").
		From(table).
		Where(fmt.Sprintf("id = '%s'", expectedResult.Id)).
		Build(sess.Dialect, buff)

	rows := sqlmock.
		NewRows([]string{"id", "name"}).
		AddRow(expectedResult.Id, expectedResult.Name)

	mock.ExpectQuery(regexp.QuoteMeta(buff.String())).WillReturnRows(rows)

	// test pointer error
	err := repo.Find(MockObject{}, expectedResult.Id)
	if err == nil {
		t.Error("Expected error and got none")
	}

	// test nil pointer
	err = repo.Find(nil, expectedResult.Id)

	// test response
	mockObj := &MockObject{}
	err = repo.Find(mockObj, expectedResult.Id)

	if err != nil {
		t.Errorf("Expected response, got error: %s", err.Error())
	}

	if !reflect.DeepEqual(mockObj, expectedResult) {
		t.Error("Object did not populate correctly")
	}
}

func TestBaseRepository_FindOneBy(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	conn := &dbr.Connection{DB: db, Dialect: dialect.PostgreSQL, EventReceiver: &dbr.NullEventReceiver{}}
	sess := conn.NewSession(nil)

	sh := structs.Helper{}
	table := "resource"

	repo := NewRepository(sess, sh, table)

	// test pointer error
	err := repo.FindOneBy(MockObject{}, FindBy{})
	if err == nil {
		t.Error("Expected pointer error and got none")
	}

	// test not tag err
	err = repo.FindOneBy(&MockObjectNoTags{}, FindBy{})
	if err == nil {
		t.Error("Expected tag error and got none")
	}

	// test response
	mockObj := &MockObject{}
	expectedResult := &MockObject{"123", "Test Name"}

	columnMap, _ := sh.GetTagMap(mockObj, "json", "db")

	fb := FindBy{
		Conditions: map[string]interface{}{
			"name": expectedResult.Name,
		},
		Search: map[string]interface{}{
			"id": expectedResult.Id,
		},
		OrderBy: map[string]interface{}{
			"name": "asc",
		},
		Limit:  1,
		Offset: 10,
	}

	buff := dbr.NewBuffer()
	query := sess.Select("*").From(table)

	for f, v := range fb.Conditions {
		query = query.Where(fmt.Sprintf(columnMap[f]+" = '%s'", v))
	}

	for f, v := range fb.Search {
		query.WhereCond = append(query.WhereCond, dbr.Like(columnMap[f], fmt.Sprintf("%%%v%%", v)))
	}

	query.OrderDir("name", true).Offset(10).Limit(1).Build(sess.Dialect, buff)

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(expectedResult.Id, expectedResult.Name)

	mock.ExpectQuery(regexp.QuoteMeta(buff.String())).WillReturnRows(rows)

	err = repo.FindOneBy(mockObj, fb)

	if err != nil {
		t.Errorf("Expected response, got error: %s", err.Error())
	}

	if !reflect.DeepEqual(mockObj, expectedResult) {
		t.Error("Object did not populate correctly")
	}

	// test same property on filter and search
	property := "test"
	conditions := map[string]interface{}{}
	conditions[property] = "value"
	search := map[string]interface{}{}
	search[property] = "value2"
	fb = FindBy{
		Conditions: conditions,
		Search:     search,
	}

	err = repo.FindOneBy(mockObj, fb)
	if err.Error() != fmt.Sprintf("property '%s' is already being filtered", property) {
		t.Error("Expected error but got none")
	}
}

func TestBaseRepository_FindBy(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	conn := &dbr.Connection{DB: db, Dialect: dialect.PostgreSQL, EventReceiver: &dbr.NullEventReceiver{}}
	sess := conn.NewSession(nil)

	sh := structs.Helper{}
	table := "resource"

	repo := NewRepository(sess, sh, table)

	// test pointer error
	err := repo.FindBy([]MockObject{}, FindBy{})
	if err == nil {
		t.Error("Expected pointer error and got none")
	}

	// test not slice error
	err = repo.FindBy(&MockObject{}, FindBy{})
	if err == nil {
		t.Error("Expected slice")
	}

	// test not tag err
	err = repo.FindBy(&[]MockObjectNoTags{}, FindBy{})
	if err == nil {
		t.Error("Expected tag error and got none")
	}

	// test response
	mockObj := &[]MockObject{}
	expectedResult := []MockObject{
		{"123", "Test Name"},
	}

	columnMap, _ := sh.GetTagMap(MockObject{}, "json", "db")

	fb := FindBy{
		Conditions: map[string]interface{}{
			"name": expectedResult[0].Name,
		},
		Search: map[string]interface{}{
			"id": expectedResult[0].Id,
		},
	}

	buff := dbr.NewBuffer()
	query := sess.Select("*").From(table)

	for f, v := range fb.Conditions {
		query = query.Where(fmt.Sprintf(columnMap[f]+" = '%s'", v))
	}

	for f, v := range fb.Search {
		query.WhereCond = append(query.WhereCond, dbr.Like(columnMap[f], fmt.Sprintf("%%%v%%", v)))
	}

	query.Build(sess.Dialect, buff)

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(expectedResult[0].Id, expectedResult[0].Name)

	mock.ExpectQuery(regexp.QuoteMeta(buff.String())).WillReturnRows(rows)

	err = repo.FindBy(mockObj, fb)

	if err != nil {
		t.Errorf("Expected response, got error: %s", err.Error())
	}

	if !reflect.DeepEqual(mockObj, &expectedResult) {
		t.Error("Object did not populate correctly")
	}

	// test same property on filter and search
	property := "test"
	conditions := map[string]interface{}{}
	conditions[property] = "value"
	search := map[string]interface{}{}
	search[property] = "value2"
	fb = FindBy{
		Conditions: conditions,
		Search:     search,
	}

	err = repo.FindBy(mockObj, fb)
	if err.Error() != fmt.Sprintf("property '%s' is already being filtered", property) {
		t.Error("Expected error but got none")
	}
}

func TestBaseRepository_Count(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	conn := &dbr.Connection{DB: db, Dialect: dialect.PostgreSQL, EventReceiver: &dbr.NullEventReceiver{}}
	sess := conn.NewSession(nil)

	sh := structs.Helper{}
	table := "resource"

	repo := NewRepository(sess, sh, table)

	// test pointer error
	_, err := repo.Count([]MockObject{}, FindBy{})
	if err == nil {
		t.Error("Expected not struct error and got none")
	}

	// test response
	mockObj := MockObject{}

	columnMap, _ := sh.GetTagMap(MockObject{}, "json", "db")

	fb := FindBy{
		Conditions: map[string]interface{}{
			"name": "a",
		},
		Search: map[string]interface{}{
			"id": "123",
		},
	}

	buff := dbr.NewBuffer()
	query := sess.Select("*").From(table)

	// add where conditions (exact match)
	for f, v := range fb.Conditions {
		query = query.Where(fmt.Sprintf("%s = '%s'", columnMap[f], v.(string)))
	}

	// add where conditions with wild card
	for f, v := range fb.Search {
		query.WhereCond = append(query.WhereCond, dbr.Like(columnMap[f], fmt.Sprintf("%%%v%%", v)))
	}

	buff.WriteString("SELECT COUNT(*) FROM (")
	query.Build(sess.Dialect, buff)
	buff.WriteString(") AS \"count\"")

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(10)

	mock.ExpectQuery(regexp.QuoteMeta(buff.String())).WillReturnRows(rows)

	count, err := repo.Count(mockObj, fb)

	if err != nil {
		t.Errorf("Expected response, got error: %s", err.Error())
	}

	if count != 10 {
		t.Errorf("Unexpected count value was returned")
	}

	count, err = repo.Count(MockObject{}, FindBy{})
	if count != 0 && err == nil {
		t.Errorf("Expected error, got none")
	}

	// test same property on filter and search
	property := "test"
	conditions := map[string]interface{}{}
	conditions[property] = "value"
	search := map[string]interface{}{}
	search[property] = "value2"
	fb = FindBy{
		Conditions: conditions,
		Search:     search,
	}

	_, err = repo.Count(mockObj, fb)
	if err.Error() != fmt.Sprintf("property '%s' is already being filtered", property) {
		t.Error("Expected error but got none")
	}
}

func TestBaseRepository_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	conn := &dbr.Connection{DB: db, Dialect: dialect.PostgreSQL, EventReceiver: &dbr.NullEventReceiver{}}
	sess := conn.NewSession(nil)

	sh := structs.Helper{}
	table := "resource"

	repo := NewRepository(sess, sh, table)

	mockObj := MockObject{Id: "123"}

	buff := dbr.NewBuffer()
	_ = sess.DeleteFrom(table).Where(fmt.Sprintf("id = '%s'", mockObj.Id)).Build(sess.Dialect, buff)

	expectedResult := sqlmock.NewResult(1, 1)
	mock.ExpectExec(regexp.QuoteMeta(buff.String())).WillReturnResult(expectedResult)

	err := repo.Delete(mockObj)
	if err != nil {
		t.Errorf("Did not expect error and got: %s", err)
	}
}

func TestBaseRepository_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	conn := &dbr.Connection{DB: db, Dialect: dialect.PostgreSQL, EventReceiver: &dbr.NullEventReceiver{}}
	sess := conn.NewSession(nil)

	sh := structs.Helper{}
	table := "resource"

	repo := NewRepository(sess, sh, table)

	mockObj := MockObject{Id: "123", Name: "test"}

	buff := dbr.NewBuffer()
	_ = sess.InsertInto(table).Build(sess.Dialect, buff)

	expectedResult := sqlmock.NewResult(1, 1)
	mock.ExpectExec(regexp.QuoteMeta(buff.String())).WillReturnResult(driver.Result(expectedResult))

	err := repo.Create(mockObj)
	if err != nil {
		t.Errorf("Did not expect error and got: %s", err)
	}
}

func TestBaseRepository_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	conn := &dbr.Connection{DB: db, Dialect: dialect.PostgreSQL, EventReceiver: &dbr.NullEventReceiver{}}
	sess := conn.NewSession(nil)

	sh := structs.Helper{}
	table := "resource"

	repo := NewRepository(sess, sh, table)

	mockObj := MockObject{Id: "123", Name: "test"}

	buff := dbr.NewBuffer()
	_ = sess.InsertInto(table).Build(sess.Dialect, buff)

	expectedResult := sqlmock.NewResult(1, 1)
	mock.ExpectExec(regexp.QuoteMeta(buff.String())).WillReturnResult(driver.Result(expectedResult))

	err := repo.Update(mockObj)
	if err != nil {
		t.Errorf("Did not expect error and got: %s", err)
	}
}
