package logger

import l "github.com/sirupsen/logrus"

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
