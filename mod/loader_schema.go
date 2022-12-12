package mod

type ExtMock struct {
	Name             string                   `json:"name,omitempty"`
	Enabled          *bool                    `json:"enabled,omitempty"`
	Priority         int                      `json:"priority,omitempty"`
	Repeat           *int                     `json:"repeat,omitempty"`
	Scenario         ExtMockScenario          `json:"scenario,omitempty"`
	Request          ExtMockRequest           `json:"request"`
	Response         *ExtMockResponse         `json:"response,omitempty"`
	RandomResponse   *ExtMockRandomResponse   `json:"random_response,omitempty"`
	SequenceResponse *ExtMockSequenceResponse `json:"sequence_response,omitempty"`
	DelayInMs        int                      `json:"delay_in_ms,omitempty"`

	Raw map[string]any `json:"-"`
}

type ExtMockRequest struct {
	URL          any            `json:"url,omitempty"`
	URLMatch     string         `json:"url_match,omitempty"`
	URLPath      any            `json:"url_path,omitempty"`
	URLPathMatch string         `json:"url_path_match,omitempty"`
	Method       any            `json:"method,omitempty"`
	Query        map[string]any `json:"query,omitempty"`
	Header       map[string]any `json:"header,omitempty"`
	Body         any            `json:"body,omitempty"`
	Form         map[string]any `json:"form,omitempty"`
}

type ExtMockScenario struct {
	Name          string `json:"name,omitempty"`
	RequiredState string `json:"required_state,omitempty"`
	NewState      string `json:"new_state,omitempty"`
}

type ExtMockResponse struct {
	Status        int               `json:"status,omitempty"`
	Header        map[string]string `json:"header,omitempty"`
	Body          any               `json:"body,omitempty"`
	BodyFile      string            `json:"body_file,omitempty"`
	Template      bool              `json:"template,omitempty"`
	TemplateModel any               `json:"template_model,omitempty"`
}

type ExtMockRandomResponse struct {
	Responses []ExtMockResponse `json:"responses,omitempty"`
}

type ExtMockSequenceResponse struct {
	Responses  []ExtMockResponse `json:"responses,omitempty"`
	AfterEnded *ExtMockResponse  `json:"after_ended,omitempty"`
}
