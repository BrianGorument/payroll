package response

const (
	// Response Code
	RCSuccess      = "00"
	RCDataNotFound = "01"
	RCServerError  = "99"

	// Descriptions
	DescriptionSuccess = "SUCCESS"
	DescriptionFailed  = "FAILED"

	// Messages
	DataSuccess      = "Data retrieved successfully"
	DataSuccessDesc  = "The requested data has been successfully retrieved"
	DataNotFound     = "Data not found"
	DataNotFoundDesc = "The requested data could not be found"
	ServerError      = "Internal server error"
	ServerErrorDesc  = "An unexpected error occurred on the server"
)
