package gormx

import (
	"context"
	"portal/internal/infra"

	"gorm.io/gorm"
)

const TxKey = "tx_key"

// WithContext 获取上下文中的事务tx
func WithContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	if v, ok := ctx.Value(TxKey).(*gorm.DB); ok {
		return v
	}
	return db
}

// NewContext 将事务tx设置到上下文中返回新的上下文
func NewContext(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, TxKey, db)
}

// Transaction 事务
func Transaction(ctx context.Context, fc func(tx *gorm.DB) error) error {
	return infra.MysqlDB.WithContext(ctx).Transaction(fc)
}
