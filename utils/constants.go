package utils

type Role string

const (
	User  Role = "USER"
	Admin Role = "ADMIN"
)

// Mongo db collections
type Collection string

const (
	Users Collection = "users"
)

const (
	Database = "test-db"

	//Fiber Ctx Keys
	InfoLogger  = "info-logger"
	ErrorLogger = "err-logger"
	MongoClient = "mongo-client"

	// Constant Keys
	Subject      = "subject"
	Expiration   = "exp"
	Scope        = "scope"
	Roles        = "roles"
	TempToken    = "temp-token"
	AccessToken  = "access-token"
	RefreshToken = "refresh-token"
	JWTToken     = "jwtToken"
)
