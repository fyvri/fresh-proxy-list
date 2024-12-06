package entity

import "encoding/json"

type Source struct {
	Method    string `json:"method"`
	Category  string `json:"category"`
	URL       string `json:"url"`
	IsChecked bool   `json:"is_checked"`
}

func (s *Source) UnmarshalJSON(data []byte) error {
	type Alias Source
	alias := &struct {
		IsChecked *bool `json:"is_checked"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	if alias.IsChecked == nil {
		s.IsChecked = true
	} else {
		s.IsChecked = *alias.IsChecked
	}

	return nil
}
