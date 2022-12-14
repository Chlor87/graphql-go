package repo

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repo[T any] struct {
	DB *gorm.DB
}

func New[T any](db *gorm.DB) (r *Repo[T], err error) {
	var t T
	err = db.AutoMigrate(&t)
	if err != nil {
		return
	}
	r = &Repo[T]{db}
	return
}

func (r *Repo[T]) Create(in *T) error {
	return r.DB.Create(in).Error
}

func (r *Repo[T]) Get(id int) (t *T, err error) {
	err = r.DB.First(&t, id).Error
	return
}

func (r *Repo[T]) List() (res []*T, err error) {
	var t T
	err = r.DB.Model(&t).Order("id asc").Find(&res).Error
	return
}

func (r *Repo[T]) Update(id int, update *T) (*T, error) {
	var t T
	tx := r.DB.Model(&t).
		Where("id = ?", id).
		Clauses(clause.Returning{}).
		Updates(update)

	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &t, nil
}
