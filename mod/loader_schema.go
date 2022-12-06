package mod

type ExternalSchema struct {
	Name             string          `json:"name,omitempty"`
	Enabled          *bool           `json:"enabled,omitempty"`
	Priority         int             `json:"priority,omitempty"`
	Repeat           *int            `json:"repeat,omitempty"`
	Scenario         ExtState        `json:"scenario"`
	Request          ExtReq          `json:"request"`
	Response         *ExtRes         `json:"response,omitempty"`
	RandomResponse   *ExtRandomRes   `json:"random_response,omitempty"`
	SequenceResponse *ExtSequenceRes `json:"sequence_response,omitempty"`
	DelayInMs        int             `json:"delay_in_ms,omitempty"`
}

type ExtReq struct {
	URL          any            `json:"url,omitempty"`
	URLMatch     string         `json:"url_match,omitempty"`
	URLPath      any            `json:"url_path,omitempty"`
	URLPathMatch string         `json:"url_path_match,omitempty"`
	Method       string         `json:"method,omitempty"`
	Query        map[string]any `json:"query,omitempty"`
	Header       map[string]any `json:"header,omitempty"`
	Body         any            `json:"body,omitempty"`
}

type ExtState struct {
	Name          string `json:"name,omitempty"`
	RequiredState string `json:"required_state,omitempty"`
	NewState      string `json:"new_state,omitempty"`
}

type ExtRes struct {
	Status        int               `json:"status,omitempty"`
	Header        map[string]string `json:"header,omitempty"`
	Body          any               `json:"body,omitempty"`
	BodyFile      string            `json:"body_file,omitempty"`
	Template      bool              `json:"template,omitempty"`
	TemplateModel any               `json:"template_model,omitempty"`
}

type ExtRandomRes struct {
	Responses []ExtRes `json:"responses,omitempty"`
}

type ExtSequenceRes struct {
	Responses  []ExtRes `json:"responses,omitempty"`
	AfterEnded *ExtRes  `json:"after_ended,omitempty"`
}
