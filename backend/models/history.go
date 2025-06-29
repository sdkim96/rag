package models

type History struct {
	UserMessage      string         `json:"user_message"`
	AssistantMessage string         `json:"assistant_message"`
	MetaData         map[string]any `json:"metadata"`
}

type Histories []*History

func (hss Histories) GetHistory() string {
	for _, hs := range hss {
		if hs == nil {
			continue
		}
		if hs.UserMessage != "" && hs.AssistantMessage != "" {
			return hs.UserMessage + " " + hs.AssistantMessage
		}
	}
	return ""
}

func (hss *Histories) SetHistory(history *History) {

	// The reason for using dereferencing
	//
	// If we append elements to a slice beyond current capacity, append() method
	// allocates a new backend array.
	// This means that if we recives slice type not by dereferencing, It could not effect to original slice.
	*hss = append(*hss, history)
}
