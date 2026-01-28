package http

import "encoding/json"

type NullableStringUpdate struct {
	present bool
	null    bool
	value   string
}

func (n *NullableStringUpdate) UnmarshalJSON(data []byte) error {
	n.present = true

	// null
	if string(data) == "null" {
		n.null = true
		return nil
	}

	// string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// "" → тоже считаем очисткой
	if s == "" {
		n.null = true
		return nil
	}

	n.value = s
	return nil
}

func (n NullableStringUpdate) IsSet() bool {
	return n.present
}

func (n NullableStringUpdate) IsNull() bool {
	return n.present && n.null
}

func (n NullableStringUpdate) Value() string {
	return n.value
}
