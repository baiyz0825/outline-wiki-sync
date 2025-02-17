// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/baiyz0825/outline-wiki-sync/model"
)

func newOutlineWikiCollectionMapping(db *gorm.DB, opts ...gen.DOOption) outlineWikiCollectionMapping {
	_outlineWikiCollectionMapping := outlineWikiCollectionMapping{}

	_outlineWikiCollectionMapping.outlineWikiCollectionMappingDo.UseDB(db, opts...)
	_outlineWikiCollectionMapping.outlineWikiCollectionMappingDo.UseModel(&model.OutlineWikiCollectionMapping{})

	tableName := _outlineWikiCollectionMapping.outlineWikiCollectionMappingDo.TableName()
	_outlineWikiCollectionMapping.ALL = field.NewAsterisk(tableName)
	_outlineWikiCollectionMapping.Id = field.NewUint(tableName, "id")
	_outlineWikiCollectionMapping.CollectionId = field.NewString(tableName, "collection_id")
	_outlineWikiCollectionMapping.CurrentId = field.NewString(tableName, "current_id")
	_outlineWikiCollectionMapping.ParentId = field.NewString(tableName, "parent_id")
	_outlineWikiCollectionMapping.CollectionPath = field.NewString(tableName, "collection_path")
	_outlineWikiCollectionMapping.CollectionName = field.NewString(tableName, "collection_name")
	_outlineWikiCollectionMapping.RealCollection = field.NewBool(tableName, "real_collection")
	_outlineWikiCollectionMapping.Sync = field.NewBool(tableName, "sync")
	_outlineWikiCollectionMapping.CreatedAt = field.NewTime(tableName, "created_at")
	_outlineWikiCollectionMapping.UpdatedAt = field.NewTime(tableName, "updated_at")
	_outlineWikiCollectionMapping.Deleted = field.NewField(tableName, "deleted")

	_outlineWikiCollectionMapping.fillFieldMap()

	return _outlineWikiCollectionMapping
}

type outlineWikiCollectionMapping struct {
	outlineWikiCollectionMappingDo outlineWikiCollectionMappingDo

	ALL            field.Asterisk
	Id             field.Uint
	CollectionId   field.String
	CurrentId      field.String
	ParentId       field.String
	CollectionPath field.String
	CollectionName field.String
	RealCollection field.Bool
	Sync           field.Bool
	CreatedAt      field.Time
	UpdatedAt      field.Time
	Deleted        field.Field

	fieldMap map[string]field.Expr
}

func (o outlineWikiCollectionMapping) Table(newTableName string) *outlineWikiCollectionMapping {
	o.outlineWikiCollectionMappingDo.UseTable(newTableName)
	return o.updateTableName(newTableName)
}

func (o outlineWikiCollectionMapping) As(alias string) *outlineWikiCollectionMapping {
	o.outlineWikiCollectionMappingDo.DO = *(o.outlineWikiCollectionMappingDo.As(alias).(*gen.DO))
	return o.updateTableName(alias)
}

func (o *outlineWikiCollectionMapping) updateTableName(table string) *outlineWikiCollectionMapping {
	o.ALL = field.NewAsterisk(table)
	o.Id = field.NewUint(table, "id")
	o.CollectionId = field.NewString(table, "collection_id")
	o.CurrentId = field.NewString(table, "current_id")
	o.ParentId = field.NewString(table, "parent_id")
	o.CollectionPath = field.NewString(table, "collection_path")
	o.CollectionName = field.NewString(table, "collection_name")
	o.RealCollection = field.NewBool(table, "real_collection")
	o.Sync = field.NewBool(table, "sync")
	o.CreatedAt = field.NewTime(table, "created_at")
	o.UpdatedAt = field.NewTime(table, "updated_at")
	o.Deleted = field.NewField(table, "deleted")

	o.fillFieldMap()

	return o
}

func (o *outlineWikiCollectionMapping) WithContext(ctx context.Context) IOutlineWikiCollectionMappingDo {
	return o.outlineWikiCollectionMappingDo.WithContext(ctx)
}

func (o outlineWikiCollectionMapping) TableName() string {
	return o.outlineWikiCollectionMappingDo.TableName()
}

func (o outlineWikiCollectionMapping) Alias() string { return o.outlineWikiCollectionMappingDo.Alias() }

func (o *outlineWikiCollectionMapping) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := o.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (o *outlineWikiCollectionMapping) fillFieldMap() {
	o.fieldMap = make(map[string]field.Expr, 11)
	o.fieldMap["id"] = o.Id
	o.fieldMap["collection_id"] = o.CollectionId
	o.fieldMap["current_id"] = o.CurrentId
	o.fieldMap["parent_id"] = o.ParentId
	o.fieldMap["collection_path"] = o.CollectionPath
	o.fieldMap["collection_name"] = o.CollectionName
	o.fieldMap["real_collection"] = o.RealCollection
	o.fieldMap["sync"] = o.Sync
	o.fieldMap["created_at"] = o.CreatedAt
	o.fieldMap["updated_at"] = o.UpdatedAt
	o.fieldMap["deleted"] = o.Deleted
}

func (o outlineWikiCollectionMapping) clone(db *gorm.DB) outlineWikiCollectionMapping {
	o.outlineWikiCollectionMappingDo.ReplaceConnPool(db.Statement.ConnPool)
	return o
}

func (o outlineWikiCollectionMapping) replaceDB(db *gorm.DB) outlineWikiCollectionMapping {
	o.outlineWikiCollectionMappingDo.ReplaceDB(db)
	return o
}

type outlineWikiCollectionMappingDo struct{ gen.DO }

type IOutlineWikiCollectionMappingDo interface {
	gen.SubQuery
	Debug() IOutlineWikiCollectionMappingDo
	WithContext(ctx context.Context) IOutlineWikiCollectionMappingDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IOutlineWikiCollectionMappingDo
	WriteDB() IOutlineWikiCollectionMappingDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IOutlineWikiCollectionMappingDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IOutlineWikiCollectionMappingDo
	Not(conds ...gen.Condition) IOutlineWikiCollectionMappingDo
	Or(conds ...gen.Condition) IOutlineWikiCollectionMappingDo
	Select(conds ...field.Expr) IOutlineWikiCollectionMappingDo
	Where(conds ...gen.Condition) IOutlineWikiCollectionMappingDo
	Order(conds ...field.Expr) IOutlineWikiCollectionMappingDo
	Distinct(cols ...field.Expr) IOutlineWikiCollectionMappingDo
	Omit(cols ...field.Expr) IOutlineWikiCollectionMappingDo
	Join(table schema.Tabler, on ...field.Expr) IOutlineWikiCollectionMappingDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IOutlineWikiCollectionMappingDo
	RightJoin(table schema.Tabler, on ...field.Expr) IOutlineWikiCollectionMappingDo
	Group(cols ...field.Expr) IOutlineWikiCollectionMappingDo
	Having(conds ...gen.Condition) IOutlineWikiCollectionMappingDo
	Limit(limit int) IOutlineWikiCollectionMappingDo
	Offset(offset int) IOutlineWikiCollectionMappingDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IOutlineWikiCollectionMappingDo
	Unscoped() IOutlineWikiCollectionMappingDo
	Create(values ...*model.OutlineWikiCollectionMapping) error
	CreateInBatches(values []*model.OutlineWikiCollectionMapping, batchSize int) error
	Save(values ...*model.OutlineWikiCollectionMapping) error
	First() (*model.OutlineWikiCollectionMapping, error)
	Take() (*model.OutlineWikiCollectionMapping, error)
	Last() (*model.OutlineWikiCollectionMapping, error)
	Find() ([]*model.OutlineWikiCollectionMapping, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.OutlineWikiCollectionMapping, err error)
	FindInBatches(result *[]*model.OutlineWikiCollectionMapping, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.OutlineWikiCollectionMapping) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IOutlineWikiCollectionMappingDo
	Assign(attrs ...field.AssignExpr) IOutlineWikiCollectionMappingDo
	Joins(fields ...field.RelationField) IOutlineWikiCollectionMappingDo
	Preload(fields ...field.RelationField) IOutlineWikiCollectionMappingDo
	FirstOrInit() (*model.OutlineWikiCollectionMapping, error)
	FirstOrCreate() (*model.OutlineWikiCollectionMapping, error)
	FindByPage(offset int, limit int) (result []*model.OutlineWikiCollectionMapping, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IOutlineWikiCollectionMappingDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (o outlineWikiCollectionMappingDo) Debug() IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Debug())
}

func (o outlineWikiCollectionMappingDo) WithContext(ctx context.Context) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.WithContext(ctx))
}

func (o outlineWikiCollectionMappingDo) ReadDB() IOutlineWikiCollectionMappingDo {
	return o.Clauses(dbresolver.Read)
}

func (o outlineWikiCollectionMappingDo) WriteDB() IOutlineWikiCollectionMappingDo {
	return o.Clauses(dbresolver.Write)
}

func (o outlineWikiCollectionMappingDo) Session(config *gorm.Session) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Session(config))
}

func (o outlineWikiCollectionMappingDo) Clauses(conds ...clause.Expression) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Clauses(conds...))
}

func (o outlineWikiCollectionMappingDo) Returning(value interface{}, columns ...string) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Returning(value, columns...))
}

func (o outlineWikiCollectionMappingDo) Not(conds ...gen.Condition) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Not(conds...))
}

func (o outlineWikiCollectionMappingDo) Or(conds ...gen.Condition) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Or(conds...))
}

func (o outlineWikiCollectionMappingDo) Select(conds ...field.Expr) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Select(conds...))
}

func (o outlineWikiCollectionMappingDo) Where(conds ...gen.Condition) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Where(conds...))
}

func (o outlineWikiCollectionMappingDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IOutlineWikiCollectionMappingDo {
	return o.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (o outlineWikiCollectionMappingDo) Order(conds ...field.Expr) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Order(conds...))
}

func (o outlineWikiCollectionMappingDo) Distinct(cols ...field.Expr) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Distinct(cols...))
}

func (o outlineWikiCollectionMappingDo) Omit(cols ...field.Expr) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Omit(cols...))
}

func (o outlineWikiCollectionMappingDo) Join(table schema.Tabler, on ...field.Expr) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Join(table, on...))
}

func (o outlineWikiCollectionMappingDo) LeftJoin(table schema.Tabler, on ...field.Expr) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.LeftJoin(table, on...))
}

func (o outlineWikiCollectionMappingDo) RightJoin(table schema.Tabler, on ...field.Expr) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.RightJoin(table, on...))
}

func (o outlineWikiCollectionMappingDo) Group(cols ...field.Expr) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Group(cols...))
}

func (o outlineWikiCollectionMappingDo) Having(conds ...gen.Condition) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Having(conds...))
}

func (o outlineWikiCollectionMappingDo) Limit(limit int) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Limit(limit))
}

func (o outlineWikiCollectionMappingDo) Offset(offset int) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Offset(offset))
}

func (o outlineWikiCollectionMappingDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Scopes(funcs...))
}

func (o outlineWikiCollectionMappingDo) Unscoped() IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Unscoped())
}

func (o outlineWikiCollectionMappingDo) Create(values ...*model.OutlineWikiCollectionMapping) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Create(values)
}

func (o outlineWikiCollectionMappingDo) CreateInBatches(values []*model.OutlineWikiCollectionMapping, batchSize int) error {
	return o.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (o outlineWikiCollectionMappingDo) Save(values ...*model.OutlineWikiCollectionMapping) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Save(values)
}

func (o outlineWikiCollectionMappingDo) First() (*model.OutlineWikiCollectionMapping, error) {
	if result, err := o.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.OutlineWikiCollectionMapping), nil
	}
}

func (o outlineWikiCollectionMappingDo) Take() (*model.OutlineWikiCollectionMapping, error) {
	if result, err := o.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.OutlineWikiCollectionMapping), nil
	}
}

func (o outlineWikiCollectionMappingDo) Last() (*model.OutlineWikiCollectionMapping, error) {
	if result, err := o.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.OutlineWikiCollectionMapping), nil
	}
}

func (o outlineWikiCollectionMappingDo) Find() ([]*model.OutlineWikiCollectionMapping, error) {
	result, err := o.DO.Find()
	return result.([]*model.OutlineWikiCollectionMapping), err
}

func (o outlineWikiCollectionMappingDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.OutlineWikiCollectionMapping, err error) {
	buf := make([]*model.OutlineWikiCollectionMapping, 0, batchSize)
	err = o.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (o outlineWikiCollectionMappingDo) FindInBatches(result *[]*model.OutlineWikiCollectionMapping, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return o.DO.FindInBatches(result, batchSize, fc)
}

func (o outlineWikiCollectionMappingDo) Attrs(attrs ...field.AssignExpr) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Attrs(attrs...))
}

func (o outlineWikiCollectionMappingDo) Assign(attrs ...field.AssignExpr) IOutlineWikiCollectionMappingDo {
	return o.withDO(o.DO.Assign(attrs...))
}

func (o outlineWikiCollectionMappingDo) Joins(fields ...field.RelationField) IOutlineWikiCollectionMappingDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Joins(_f))
	}
	return &o
}

func (o outlineWikiCollectionMappingDo) Preload(fields ...field.RelationField) IOutlineWikiCollectionMappingDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Preload(_f))
	}
	return &o
}

func (o outlineWikiCollectionMappingDo) FirstOrInit() (*model.OutlineWikiCollectionMapping, error) {
	if result, err := o.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.OutlineWikiCollectionMapping), nil
	}
}

func (o outlineWikiCollectionMappingDo) FirstOrCreate() (*model.OutlineWikiCollectionMapping, error) {
	if result, err := o.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.OutlineWikiCollectionMapping), nil
	}
}

func (o outlineWikiCollectionMappingDo) FindByPage(offset int, limit int) (result []*model.OutlineWikiCollectionMapping, count int64, err error) {
	result, err = o.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = o.Offset(-1).Limit(-1).Count()
	return
}

func (o outlineWikiCollectionMappingDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = o.Count()
	if err != nil {
		return
	}

	err = o.Offset(offset).Limit(limit).Scan(result)
	return
}

func (o outlineWikiCollectionMappingDo) Scan(result interface{}) (err error) {
	return o.DO.Scan(result)
}

func (o outlineWikiCollectionMappingDo) Delete(models ...*model.OutlineWikiCollectionMapping) (result gen.ResultInfo, err error) {
	return o.DO.Delete(models)
}

func (o *outlineWikiCollectionMappingDo) withDO(do gen.Dao) *outlineWikiCollectionMappingDo {
	o.DO = *do.(*gen.DO)
	return o
}
