package db

import (
	"errors"
	"fmt"
	"github.com/gocraft/dbr"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/structs"
	"reflect"
)

type Repository interface {
	Find(object interface{}, id string) error
	FindOneBy(object interface{}, fb FindBy) error
	FindBy(objects interface{}, fb FindBy) error
	Create(object interface{}) error
	Update(object interface{}) error
	Delete(object interface{}) error
	Count(object interface{}, fb FindBy) (int, error)
}

type FindBy struct {
	Conditions map[string]interface{}
	Search     map[string]interface{}
	OrderBy    map[string]interface{}
	Limit      uint64
	Offset     uint64
}

type BaseRepository struct {
	Db    *dbr.Session
	Sh    structs.Helper
	Table string
}

func NewRepository(db *dbr.Session, sh structs.Helper, table string) Repository {
	return &BaseRepository{db, sh, table}
}

func (r BaseRepository) Find(object interface{}, id string) error {
	if err := r.IsPointer(object); err != nil {
		return err
	}

	return r.Db.Select("*").From(r.Table).Where("id = ?", id).Limit(1).LoadOne(object)
}

func (r BaseRepository) FindOneBy(object interface{}, fb FindBy) error {
	if err := r.IsPointer(object); err != nil {
		return err
	}

	query, err := r.buildQuery(object, fb, true, true)
	if err != nil {
		return err
	}

	query = query.Limit(1) // ensure limit is 1

	return query.LoadOne(object)
}

func (r BaseRepository) FindBy(objects interface{}, fb FindBy) error {
	if err := r.IsPointer(objects); err != nil {
		return err
	}

	object, err := r.getSliceElementType(objects)
	if err != nil {
		return err
	}

	_, err = r.buildQuery(object, fb, true, true)
	if err != nil {
		return err
	}

	// _, err = query.Load(objects)

	return err
}

func (r BaseRepository) Create(object interface{}) error {
	columns := r.Sh.GetTagValues(object, "db")
	_, err := r.Db.InsertInto(r.Table).Columns(columns...).Record(object).Exec()

	return err
}

func (r BaseRepository) Update(object interface{}) error {
	objectMap := r.Sh.GetMapByTag(object, "structs")

	_, err := r.Db.Update(r.Table).SetMap(objectMap).Where("id = ?", objectMap["id"]).Exec()

	return err
}

func (r BaseRepository) Delete(object interface{}) error {
	objectMap := r.Sh.GetMapByTag(object, "structs")

	_, err := r.Db.DeleteFrom(r.Table).Where("id = ?", objectMap["id"]).Exec()

	return err
}

func (r BaseRepository) Count(object interface{}, fb FindBy) (int, error) {
	if reflect.ValueOf(object).Kind() != reflect.Struct {
		return 0, errors.New("object not a struct")
	}

	query, err := r.buildQuery(object, fb, false, false)
	if err != nil {
		return 0, err
	}
	outerQuery := r.Db.Select("COUNT(*)").From(query.As("count"))

	count := 0
	_, err = outerQuery.Load(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r BaseRepository) IsPointer(object interface{}) error {
	v := reflect.ValueOf(object)

	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer to repository methods")
	}

	return nil
}

func (r BaseRepository) getSliceElementType(objects interface{}) (interface{}, error) {
	v := reflect.ValueOf(objects)

	if v.Elem().Kind() != reflect.Slice {
		return nil, errors.New("must pass a slice to repository getSliceType")
	}

	elem := reflect.TypeOf(objects).Elem().Elem()

	return reflect.New(elem).Elem().Interface(), nil
}

func (r BaseRepository) buildQuery(object interface{}, fb FindBy, addOffset bool, addLimit bool) (*dbr.SelectStmt, error) {
	columnMap, err := r.Sh.GetTagMap(object, "json", "db")
	if err != nil {
		return nil, err
	}

	query := r.Db.Select("*").From(r.Table)

	for f, v := range fb.Conditions {
		query = query.Where(columnMap[f]+" = ?", v)
	}

	for f, v := range fb.Search {
		// check if filter exists in Conditions
		if _, ok := fb.Conditions[f]; ok {
			return nil, fmt.Errorf("property '%s' is already being filtered", f)
		}

		query.WhereCond = append(query.WhereCond, dbr.Like(columnMap[f], fmt.Sprintf("%%%v%%", v)))
	}

	for f, v := range fb.OrderBy {
		query = query.OrderDir(columnMap[f], v == "asc")
	}

	if addOffset && fb.Offset != 0 {
		query = query.Offset(fb.Offset)
	}

	if addLimit && fb.Limit != 0 {
		query = query.Limit(fb.Limit)
	}

	return query, nil
}
