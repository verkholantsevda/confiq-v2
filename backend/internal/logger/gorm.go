package logger

import gormlogger "gorm.io/gorm/logger"

func Gorm() gormlogger.Interface {
	return &GormLogger{}
}
