package errorKit

import (
	"fmt"

	"google.golang.org/grpc/status"
)

const (
	StatusBoInternalError = 999
	StatusBoInvalidParams = 10002
	StatusBoIpNotAllowed  = 10003
)

var (
	ErrBoDisabled             = status.Errorf(StatusBoInternalError, "bo disabled")
	ErrForbidden              = status.Errorf(StatusBoInternalError, "forbidden")
	ErrUnauthorized           = status.Errorf(StatusBoInternalError, "unauthorized")
	ErrInvalidToken           = status.Errorf(StatusBoInternalError, "invalid token")
	ErrInvalidParams          = status.Errorf(StatusBoInvalidParams, "invalid params")
	ErrGetHeader              = status.Errorf(StatusBoInternalError, "get header error")
	ErrRecordNotFound         = status.Errorf(StatusBoInternalError, "record not found")
	ErrMessageQueue           = status.Errorf(StatusBoInternalError, "message queue error")
	ErrDatabase               = status.Errorf(StatusBoInternalError, "database error")
	ErrCurrencyNotSupport     = status.Errorf(StatusBoInternalError, "currency not support")
	ErrProviderLineNotSupport = status.Errorf(StatusBoInternalError, "provider line not support")
	ErrCommissionRate         = status.Errorf(StatusBoInternalError, "invalid commission rate")
	ErrCurrencyIsRequired     = status.Errorf(StatusBoInternalError, "currency is required")
	ErrPermissionEmpty        = status.Errorf(StatusBoInternalError, "permission should not empty")

	ErrCreditReqQuery        = status.Errorf(StatusBoInternalError, "failed to find credit req")
	ErrCreditReqUpdated      = status.Errorf(StatusBoInternalError, "failed to update credit req")
	ErrCreditReqNotUpdated   = status.Errorf(StatusBoInternalError, "credit req not updated")
	ErrCreditReqCreated      = status.Errorf(StatusBoInternalError, "credit req created")
	ErrUserNotUpdated        = status.Errorf(StatusBoInternalError, "user not updated")
	ErrRoleNotFound          = status.Errorf(StatusBoInternalError, "role not found")
	ErrAgentNodFound         = status.Errorf(StatusBoInternalError, "agent not found")
	ErrAgentDuplicated       = status.Errorf(StatusBoInternalError, "agent duplicated")
	ErrAgentCodeDuplicated   = status.Errorf(StatusBoInternalError, "agent code duplicated")
	ErrUserNotFound          = status.Errorf(StatusBoInternalError, "user not found")
	ErrUserExists            = status.Errorf(StatusBoInternalError, "user exists")
	ErrInvalidPassword       = status.Errorf(StatusBoInternalError, "invalid password")
	ErrInvalidOTP            = status.Errorf(StatusBoInternalError, "invalid otp")
	ErrUserUpdated           = status.Errorf(StatusBoInternalError, "user updated failed")
	ErrUserCreated           = status.Errorf(StatusBoInternalError, "user created failed")
	ErrWalletCreated         = status.Errorf(StatusBoInternalError, "wallet created failed")
	ErrRoleCreated           = status.Errorf(StatusBoInternalError, "role created failed")
	ErrChannelCreated        = status.Errorf(StatusBoInternalError, "channel created failed")
	ErrMemberCreated         = status.Errorf(StatusBoInternalError, "member created failed")
	ErrSalesNotFound         = status.Errorf(StatusBoInternalError, "sales not found")
	ErrWhiteIPExists         = status.Errorf(StatusBoInternalError, "white ip exists")
	ErrWhiteIPCreated        = status.Errorf(StatusBoInternalError, "white ip created failed")
	ErrWhiteIPNotUpdated     = status.Errorf(StatusBoInternalError, "white ip not updated")
	ErrCreateConversionRates = status.Errorf(StatusBoInternalError, "conversion rates created failed")

	ErrGetCache    = status.Errorf(StatusBoInternalError, "failed to get cache")
	ErrDelCache    = status.Errorf(StatusBoInternalError, "failed to del cache")
	ErrGetGameList = status.Errorf(StatusBoInternalError, "failed to get game list")

	// auth
	ErrUserNameEmpty  = status.Errorf(StatusBoInvalidParams, "username should not empty")
	ErrPasswordEmpty  = status.Errorf(StatusBoInvalidParams, "password should not empty")
	ErrAgentCodeEmpty = status.Errorf(StatusBoInvalidParams, "agent code should not empty")
	ErrOTPIsDisable   = status.Errorf(StatusBoInvalidParams, "otp is disable")
	ErrInvalidIP      = func(ip string) error {
		return status.Errorf(StatusBoIpNotAllowed, "IP is not allowed. (%s)", ip)
	}
	ErrInvalidSign = status.Errorf(StatusBoInternalError, "invalid sign")

	// platform
	ErrInsufficientAgBalance = fmt.Errorf("insufficient agent balance")
	ErrPlayerNotFound        = fmt.Errorf("player not found")
	ErrProviderLineNotFound  = fmt.Errorf("provider line not found")
	ErrProviderNotFound      = fmt.Errorf("provider not found")
	ErrPasswordIncorrect     = fmt.Errorf("password incorrect")
	ErrWagerAlreadyExists    = fmt.Errorf("wager already exists")
	ErrWagerAlreadySettled   = fmt.Errorf("wager already settled")
	ErrTxAlreadyExists       = fmt.Errorf("transaction already exists")
	ErrTxAlreadySettled      = fmt.Errorf("transaction already settled")
	ErrInvalidOperation      = fmt.Errorf("invalid operation")
	ErrGameNotSupported      = fmt.Errorf("game not supported")
	ErrInvalidGameCode       = fmt.Errorf("invalid game code")
	ErrWagerNotFound         = fmt.Errorf("wager not found")
	ErrTxNotFound            = fmt.Errorf("transaction not found")
	ErrSetOTPFailed          = fmt.Errorf("set otp failed")
	ErrInvalidCurrency       = fmt.Errorf("invalid currency")
	ErrBetAlreadyExists      = fmt.Errorf("bet transaction already exists")
	ErrSettleAlreadyExists   = fmt.Errorf("settle transaction already exists")
	ErrCancelAlreadyExists   = fmt.Errorf("cancel transaction already exists")
	ErrInsufficientBalance   = fmt.Errorf("insufficient balance")
	ErrInvalidSignature      = fmt.Errorf("invalid signature")
	ErrInvalidParameters     = fmt.Errorf("invalid parameters")
	ErrSessionExpired        = fmt.Errorf("session expired")
	ErrGetPlatformHeader     = fmt.Errorf("get header error")
	ErrInvalidProviderToken  = fmt.Errorf("invalid provider token")
	ErrInvalidPlayerSession  = fmt.Errorf("invalid player session")
	ErrGameNotFound          = fmt.Errorf("game not found")
	ErrLaunchGameFailed      = fmt.Errorf("launch game failed")

	// provider
	ErrFailedToCreatePlayer     = fmt.Errorf("failed to create player")
	ErrFalledToToSigninProvider = fmt.Errorf("failed to signin provider")
	ErrFailedToGetGameList      = fmt.Errorf("failed to get game list")
	ErrFailedToGetGameURL       = fmt.Errorf("failed to get game url")
)
