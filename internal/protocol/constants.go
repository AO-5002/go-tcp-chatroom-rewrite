package protocol

// Validation limits
const (
	MinUsernameLength = 1
	MaxUsernameLength = 32
	MaxPortLength     = 4
	MinMessageLength  = 1
	MaxMessageLength  = 512
	MaxUserListLength = 10
	MaxProtocolLength = 1024
)

// Error codes (from the spec's protocol table)
type ErrorCode int

const (
	UsernameTaken   ErrorCode = iota // 0
	InvalidUsername                  // 1
	UnknownUserPMSG                  // 2
	InvalidMessage                   // 3
)

// Transport
const (
	ServerPort         = 8080
	ClientPort         = 3000
	MaxUserConnections = 10
)
