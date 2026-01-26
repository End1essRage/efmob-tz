package logger

import l "github.com/sirupsen/logrus"

// хук для отладки
type DebugHook struct {
	Pkg  string
	Func string
}

func (h *DebugHook) Levels() []l.Level {
	return l.AllLevels
}

func (h *DebugHook) Fire(entry *l.Entry) error {
	entry.Data["package"] = h.Pkg
	entry.Data["method"] = h.Func
	return nil
}

// хук для метаинфы о сервисе
type MetaHook struct {
	Name string
}

func (h *MetaHook) Levels() []l.Level {
	return l.AllLevels
}

func (h *MetaHook) Fire(entry *l.Entry) error {
	entry.Data["service"] = h.Name
	return nil
}
