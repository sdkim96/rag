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
	*hss = append(*hss, history)
}
