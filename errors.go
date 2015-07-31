package ideamart

import "fmt"

type Error struct {
	Type        string
	Code        string
	Description string
	Retryable   bool
}

func (e Error) Error() string {
	return fmt.Sprintf("%s %s: %s", e.Type, e.Code, e.Description)
}

const (
	TypeAPIError     = "ApiError"
	TypeClientError  = "ClientError"
	TypeUnknownError = "UnknownError"
)

var (
	ErrAuthFailed           = Error{TypeAPIError, "E1313", "Authentication failed", false}
	ErrIPNotProvisioned     = Error{TypeAPIError, "E1303", "IP address from which this request originated is not provisioned", false}
	ErrReqInvalid           = Error{TypeAPIError, "E1312", "Request is invalid", false}
	ErrSMSServUnavailable   = Error{TypeAPIError, "E1309", "Requested SMS service is not allowed for this application", false}
	ErrMOTermNotAllowed     = Error{TypeAPIError, "E1311", "Mobile terminated SMS messages have not been enabled", false}
	ErrSMSServNotFound      = Error{TypeAPIError, "E1315", "Cannot find the requested service SMS or it is not active", false}
	ErrMSISDNInvalid        = Error{TypeAPIError, "E1317", "MSISDN in request is invalid or not allowed", false}
	ErrChargeOpNotAllowed   = Error{TypeAPIError, "E1328", "Charging operation requested is not allowed", false}
	ErrReqFailedToAllDest   = Error{TypeAPIError, "E1341", "Request failed; errors occurred while sending the request to all destinations", false}
	ErrSMSMsgTooLong        = Error{TypeAPIError, "E1334", "SMS message length too long", false}
	ErrSMSAdvertTooLong     = Error{TypeAPIError, "E1335", "SMS advertisement message length too long", false}
	ErrDuplicateReq         = Error{TypeAPIError, "E1337", "Duplicate request", false}
	ErrUnexpected           = Error{TypeAPIError, "E1601", "System experienced an unexpected error", false}
	ErrMSISDNBlacklisted    = Error{TypeAPIError, "E1342", "MSISDN black listed", false}
	ErrMSISDNNotWhitelisted = Error{TypeAPIError, "E1343", "MSISDN not white listed", false}
	ErrAddrFormatInvalid    = Error{TypeAPIError, "E1325", "Address format invalid", false}
	ErrSrcAddrNotAllowed    = Error{TypeAPIError, "E1331", "Surce address is not allowed", false}
	ErrPermChrg             = Error{TypeAPIError, "E1308", "Permanent charging error", false}

	ErrTrxnLimExceededPerSec = Error{TypeAPIError, "E1318", "Transaction limit per second has exceeded. Please throttle requests not to exceed the transaction limit. Contact Idea Mart admin to increase the traffic limit.", true}
	ErrTrxnLimExceededPerDay = Error{TypeAPIError, "E1319", "Transaction limit for today is exceeded. Please try again tomorrow or contact Idea Mart admin to increase the transaction per day limit.", true}
	ErrBalInsufficient       = Error{TypeAPIError, "E1326", "Insufficient balance.", true}
	ErrMsgDelivFailed        = Error{TypeAPIError, "E1602", "Message delivery failed. Please retry.", true}
	ErrTempSysErr            = Error{TypeAPIError, "E1603", "Temporary System Error occurred while delivering your request.", true}

	ErrInvalidJSON   = Error{TypeClientError, "", "Ideamart API sent invalid JSON", true}
	ErrSendingFailed = Error{TypeClientError, "", "Sending message failed after retries", false}
)

var apiErrMap = map[string]Error{}

func apiErrorFromCode(code string) Error {
	if apiErrMap[code].Code != "" {
		return apiErrMap[code]
	}
	return Error{TypeAPIError, code, "Unknown API Error", false}
}

func init() {
	apiErrMap["E1313"] = ErrAuthFailed
	apiErrMap["E1303"] = ErrIPNotProvisioned
	apiErrMap["E1312"] = ErrReqInvalid
	apiErrMap["E1309"] = ErrSMSServUnavailable
	apiErrMap["E1311"] = ErrMOTermNotAllowed
	apiErrMap["E1315"] = ErrSMSServNotFound
	apiErrMap["E1317"] = ErrMSISDNInvalid
	apiErrMap["E1328"] = ErrChargeOpNotAllowed
	apiErrMap["E1341"] = ErrReqFailedToAllDest
	apiErrMap["E1334"] = ErrSMSMsgTooLong
	apiErrMap["E1335"] = ErrSMSAdvertTooLong
	apiErrMap["E1337"] = ErrDuplicateReq
	apiErrMap["E1601"] = ErrUnexpected
	apiErrMap["E1342"] = ErrMSISDNBlacklisted
	apiErrMap["E1343"] = ErrMSISDNNotWhitelisted
	apiErrMap["E1325"] = ErrAddrFormatInvalid
	apiErrMap["E1331"] = ErrSrcAddrNotAllowed
	apiErrMap["E1308"] = ErrPermChrg

	apiErrMap["E1318"] = ErrTrxnLimExceededPerSec
	apiErrMap["E1319"] = ErrTrxnLimExceededPerDay
	apiErrMap["E1326"] = ErrBalInsufficient
	apiErrMap["E1602"] = ErrMsgDelivFailed
	apiErrMap["E1603"] = ErrTempSysErr
}
