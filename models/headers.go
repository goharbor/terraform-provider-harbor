package models

type ResponseHeaders struct {
	Connection    []string `json:"Connection,omitempty"`
	ContentLength []string `json:"Content-Length,omitempty"`
	ContentType   []string `json:"Content-Type,omitempty"`
	Date          []string `json:"Date,omitempty"`
	Location      []string `json:"Location,omitempty"`
	Server        []string `json:"Server,omitempty"`
	SetCookie     []string `json:"Set-Cookie,omitempty"`
	XRequestID    []string `json:"X-Request-Id,omitempty"`
}
