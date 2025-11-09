package http

type EnvelopeConfig struct {
	// mapping the error code to http status code, default is false,
	// and the http code will be always be 200 when the server returns.
	MappingCode    bool   `json:"mappingCode"`
	ErrorWrapped   bool   `json:"errorWrapped"`   // there is an error field in the enveloped response
	SuccessCode    string `json:"successCode"`    // the success code value when error not wrapped, the default is "200"
	SuccessMessage string `json:"successMessage"` // the default is "OK"
}
