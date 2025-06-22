// This Module defines constants used across the application.

package cst

const (
	RootMessageID = "00000000-0000-0000-0000-000000000000" // Represents the root message ID for new conversations.
)

// API Response Statuses
const (
	Ok    = "ok"    // Represents a successful operation
	Error = "error" // Represents an error in the operation
)

// Errors For API Responses
const (
	UnAuthorizedUserError = "Unauthorized user" // 401
	EntityError           = "Entity error"      // 422

	NotInplementedError = "Not implemented"       // 501
	InternalServerError = "Internal server error" // 500
)

// Messages for API Responses
const (
	NewConversationCreated    = "New conversation created"    // Message when a new conversation is created
	ConversationListRetrieved = "Conversation list retrieved" // Message when conversation list is retrieved
)

// Errors for across the application
const (
	AssertionError = "Assertion error" // Error when type assertion fails
)
