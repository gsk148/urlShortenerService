package api

// ShortenRequest model for /api/shorten request
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse model for /api/shorten response
type ShortenResponse struct {
	Result string `json:"result"`
}

// BatchShortenRequestItem model for batch request
type BatchShortenRequestItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchShortenResponseItem model for batch response
type BatchShortenResponseItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// Statistic model for statistic response
type Statistic struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

// URLInfo model for url info
type URLInfo struct {
	UUID          string `json:"uuid,omitempty"`
	UserID        string `json:"userID,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
	IsDeleted     bool   `json:"is_deleted,omitempty"`
}

// ShortenedData model for url info
type ShortenedData struct {
	UserID      string `json:"userID"`
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	IsDeleted   bool   `json:"is_deleted"`
}
