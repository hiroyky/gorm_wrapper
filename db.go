package gormwrapper

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DB interface {
	Create(values any) error
	UpdateColumns(m any, updateValues map[string]any) error
	Delete(m any) error
	First(dst any) error
	Find(dst any) error
	Count(dst int64) error
	Raw(sql string, values ...any) error
	Exec(sql string, values ...any) error
	Where(query any, args ...any) DB
	OrderBy(column string, direction OrderDirection) DB
	Limit(limit int) DB
	Offset(offset int) DB
	Clauses(cons ...clause.Expression) DB
}

func NewDB(conn *gorm.DB) DB {
	return &db{
		conn: conn,
	}
}

type db struct {
	conn *gorm.DB
}

func (d *db) Create(values any) error {
	return d.conn.Create(values).Error
}

func (d *db) UpdateColumns(m any, updateValues map[string]any) error {
	return d.conn.Model(m).UpdateColumns(updateValues).Error
}

func (d *db) Delete(m any) error {
	return d.conn.Delete(m).Error
}

func (d *db) First(dst any) error {
	return d.conn.First(&dst).Error
}

func (d *db) Find(dst any) error {
	return d.conn.Find(&dst).Error
}

func (d *db) Count(m Model, dst int64) error {
	return d.conn.Model(m).Count(&dst).Error
}

func (d *db) Raw(sql string, values ...any) error {
	return d.conn.Raw(sql, values...).Error
}

func (d *db) Exec(sql string, values ...any) error {
	return d.conn.Exec(sql, values...).Error
}

func (d *db) Where(query any, args ...any) DB {
	conn := d.conn.Where(query, args...)
	return d.cloneQuery(conn)
}

func (d *db) OrderBy(column string, direction OrderDirection) DB {
	conn := d.conn.Order(fmt.Sprintf("`%s` %s", column, direction))
	return d.cloneQuery(conn)
}

func (d *db) Limit(limit int) DB {
	conn := d.conn.Limit(limit)
	return d.cloneQuery(conn)
}

func (d *db) Offset(offset int) DB {
	conn := d.conn.Offset(offset)
	return d.cloneQuery(conn)
}

func (d *db) Clauses(cons ...clause.Expression) DB {
	conn := d.conn.Clauses(cons...)
	return d.cloneQuery(conn)
}

func (d *db) cloneQuery(conn *gorm.DB) *db {
	return &db{conn: conn}
}
