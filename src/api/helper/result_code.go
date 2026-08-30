package helper

type ResultCode int

const (
	Success         = 0
	ValidationError = 40001
	AuthError       = 40101
	ForbiddenError  = 40301
	NotFoundError   = 40401
	LimiterError    = 42901
	OtpLimiterError = 42902
	CustomRecovery  = 50001
	InternalError   = 50002
)
