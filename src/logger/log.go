package logger

import "github.com/mohsenHa/messenger/logger/loggerentity"

type Log struct {
	cat   loggerentity.Category
	sub   loggerentity.SubCategory
	msg   string
	extra map[loggerentity.ExtraKey]interface{}
}

func NewLog(msg string) *Log {
	return &Log{
		cat:   loggerentity.CategoryNotDefined,
		sub:   loggerentity.SubCategoryNotDefined,
		msg:   msg,
		extra: map[loggerentity.ExtraKey]interface{}{},
	}
}

func (l *Log) WithCategory(cat loggerentity.Category) *Log {
	l.cat = cat
	return l
}

func (l *Log) WithSubCategory(sub loggerentity.SubCategory) *Log {
	l.sub = sub
	return l
}
func (l *Log) With(key loggerentity.ExtraKey, value interface{}) *Log {
	l.extra[key] = value
	return l
}

func (l *Log) Debug() {
	L().Debug(l.cat, l.sub, l.msg, l.extra)
}

func (l *Log) Info() {
	L().Info(l.cat, l.sub, l.msg, l.extra)
}

func (l *Log) Warn() {
	L().Warn(l.cat, l.sub, l.msg, l.extra)
}

func (l *Log) Error() {
	L().Error(l.cat, l.sub, l.msg, l.extra)
}

func (l *Log) Fatal() {
	L().Fatal(l.cat, l.sub, l.msg, l.extra)
}
