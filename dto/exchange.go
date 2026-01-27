package dto

// ExchangeCreateRequest is used to create an exchange record.
type ExchangeCreateRequest struct {
	MQMExchangeCode    string  `json:"mqmExchangeCode" binding:"required,max=10"`
	ClearExchangeCode  string  `json:"clearExchangeCode" binding:"required,max=10"`
	GlobexExchangeCode string  `json:"globexExchangeCode" binding:"required,max=10"`
	Description        *string `json:"description" binding:"omitempty,max=3000"`
	SegType            string  `json:"segType" binding:"required,max=10"`
}

// ExchangeUpdateRequest updates mutable exchange fields.
type ExchangeUpdateRequest struct {
	ClearExchangeCode  string  `json:"clearExchangeCode" binding:"required,max=10"`
	GlobexExchangeCode string  `json:"globexExchangeCode" binding:"required,max=10"`
	Description        *string `json:"description" binding:"omitempty,max=3000"`
	SegType            string  `json:"segType" binding:"required,max=10"`
}

// ExchangeResponse returns exchange data to clients.
type ExchangeResponse struct {
	MQMExchangeCode    string  `json:"mqmExchangeCode"`
	ClearExchangeCode  string  `json:"clearExchangeCode"`
	GlobexExchangeCode string  `json:"globexExchangeCode"`
	Description        *string `json:"description,omitempty"`
	SegType            string  `json:"segType"`
}

// ExchangeListResponse wraps paginated exchanges.
type ExchangeListResponse struct {
	Total int64              `json:"total"`
	Items []ExchangeResponse `json:"items"`
}
