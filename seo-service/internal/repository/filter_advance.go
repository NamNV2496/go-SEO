package repository

// import (
// 	"context"

// 	"gorm.io/gorm"
// )

// type QueryOption struct {
// 	Limit   int
// 	Offset  int
// 	OrderBy string
// }
// type WhereCondition struct {
// 	AndOpertor bool
// 	OrOpertor  bool
// 	Statement  string
// 	Value      any
// }

// type Filter[T any] struct {
// 	db *gorm.DB
// }

// func NewFilter[T any](
// 	database IDatabase,
// ) *Filter[T] {
// 	return &Filter[T]{
// 		db: database.GetDB(),
// 	}
// }

// func (_self *Filter[T]) Query(ctx context.Context, opt QueryOption, model T, conditions ...WhereCondition) ([]T, error) {
// 	db := _self.db.WithContext(ctx).Model(&model)
// 	if len(conditions) > 0 {
// 		for _, condition := range conditions {
// 			if condition.AndOpertor {
// 				db = db.Where(condition.Statement, condition.Value)
// 			} else if condition.OrOpertor {
// 				db = db.Or(condition.Statement, condition.Value)
// 			}
// 		}
// 	}
// 	var results []T
// 	err := db.
// 		Order(opt.OrderBy).
// 		Offset(opt.Offset).
// 		Limit(opt.Limit).
// 		Find(&results).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return results, nil
// }
