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
