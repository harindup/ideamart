Go Ideamart Client
==================

A client library for the [Dialog Ideamart APIs](http://www.ideamart.lk/idea-pro) written in [Go](https://golang.org/).

Supported APIs
--------------
 * USSD
 * SMS
 * CaaS
 * Subscription

Additional Features
-------------------
* USSD session handler with support for custom sessions stores.
* An in-memory USSD session store with built-in garbage collection.

LICENSE
-------

### The MIT License (MIT)

Copyright © 2015 R. A. Harindu Perera

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

Usage
-----

### Import

--
    import "github.com/harindup/ideamart"

### Endpoint constants

```go
const (
    USSDSendEndpointLocal             = "http://localhost:7000/ussd/send"
	USSDSendEndpointLive              = "https://api.dialog.lk/ussd/send"
	SMSSendEndpointLocal              = "http://localhost:7000/sms/send"
	SMSSendEndpointLive               = "https://api.dialog.lk/sms/send"
	SubscriptionEndpointLocal         = "http://localhost:7000/subscription/send"
	SubscriptionEndpointLive          = "https://api.dialog.lk/subscription/send"
	SubscriptionBaseSizeEndpointLocal = "http://localhost:7000/subscription/query-base"
	SubscriptionBaseSizeEndpointLive  = "https://api.dialog.lk/subscription/query-base"
	SubscriptionStatusEndpointLocal   = "http://localhost:7000/subscription/getStatus"
	SubscriptionStatusEndpointLive    = "https://api.dialog.lk/subscription/getStatus"
	CaaSQueryBalanceEndpointLocal     = "http://localhost:7000/caas/get/balance"
	CaaSQueryBalanceEndpointLive      = "https://api.dialog.lk/caas/get/balance"
)
```

### Error types
```go
const (
	TypeAPIError     = "ApiError"
	TypeClientError  = "ClientError"
	TypeUnknownError = "UnknownError"
)
```

### SMS delivery status codes from Ideamart
```go
const (
	SMSStatusSent          = "SENT"
	SMSStatusDelivered     = "DELIVERED"
	SMSStatusExpired       = "EXPIRED"
	SMSStatusDeleted       = "DELETED"
	SMSStatusUndeliverable = "UNDELIVERABLE"
	SMSStatusAccepted      = "ACCEPTED"
	SMSStatusUnknown       = "UNKNOWN"
	SMSStatusRejected      = "REJECTED"
)
```

###The Subscriber status codes returned by Ideamart API
```go
const (
	SubscriberStatusRegistered    = "REGISTERED"
	SubscriberStatusUnregistered  = "UNREGISTERED"
	SubscriberStatusPendingCharge = "PENDING CHARGE"
)
```

```go
const (
	MobileTermiatedFinal    MobileTerminatedUSSDOperation = "mt-fin"
	MobileTermiatedContinue MobileTerminatedUSSDOperation = "mt-cont"

	MobileOriginatedInitial  MobileOriginatedUSSDOperation = "mo-init"
	MobileOriginatedContinue MobileOriginatedUSSDOperation = "mo-cont"

	StatusCodeSuccess string = "S1000"
)
```

```go
const (
	CaaSMobileAccount = "MobileAccount"
)
```

### Error codes returned by the Ideamart API
```go
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
```

#### func  NewInMemorySessionStore

```go
func NewInMemorySessionStore(maxSize int) inMemorySessionStore
```
Returns a properly initialized in-memory sessions store. maxSize is the maximum
number of sessions it will store before discarding the "oldest" session.

#### func  NewSMSQueue

```go
func NewSMSQueue(client *SMSClient, capacity, messagesPerSecond, maxRetryCount int, sendCallback func(smsMessage, recipient, smsMessageId string)) smsQueue
```
Initializes and returns a new SMS queue. Make sure that an application has only
one queue if request throttling should be properly functional.

#### type CaaSBalanceRequest

```go
type CaaSBalanceRequest struct {
	ApplicationID         string `json:"applicationId"`
	Password              string `json:"password"`
	SubscriberID          string `json:"subscriberId"`
	PaymentInstrumentName string `json:"paymentInstrumentName"`
}
```


#### type CaaSBalanceResponse

```go
type CaaSBalanceResponse struct {
	StatusCode        string `json:"statusCode"`
	StatusDetail      string `json:"statusDetail"`
	ChargeableBalance string `json:"chargeableBalance"`
	AccountStatus     string `json:"accountStatus"`
	AccountType       string `json:"accountType"`
}
```


#### type CaaSClient

```go
type CaaSClient struct {
	ApplicationID       string
	Password            string
	BalanceEndpoint     string
	DirectDebitEndpoint string
}
```

CaasS Client.

#### func (*CaaSClient) DirectDebit

```go
func (client *CaaSClient) DirectDebit(subscriberId, paymentInstrumentName, externalTrxId string, amount float64) (string, time.Time, error)
```

#### func (*CaaSClient) GetBalance

```go
func (client *CaaSClient) GetBalance(subscriberId, paymentInstrumentName string) (float64, error)
```

#### type CaaSDirectDebitRequest

```go
type CaaSDirectDebitRequest struct {
	ApplicationID         string `json:"applicationId"`
	Password              string `json:"password"`
	SubscriberID          string `json:"subscriberId"`
	PaymentInstrumentName string `json:"paymentInstrumentName"`
	ExternalTransactionID string `json:"externalTrxId"`
	Amount                string `json:"amount"`
}
```


#### type CaaSDirectDebitResponse

```go
type CaaSDirectDebitResponse struct {
	StatusCode            string `json:"statusCode"`
	Timestamp             string `json:"timeStamp"`
	ShortDescription      string `json:"shortDescription"`
	StatusDetail          string `json:"statusDetail"`
	ExternalTransactionID string `json:"externalTrxId"`
	LongDescription       string `json:"longDescription"`
	InternalTransactionID string `json:"internalTrxId"`
}
```


#### type Error

```go
type Error struct {
	Type        string
	Code        string
	Description string
	Retryable   bool
}
```

Ideamart Error type. Type is either ApiError, ClientError or UnknownError.
Retryable is a boolean value which gives whether the same request can be retied
after the error.

#### func (Error) Error

```go
func (e Error) Error() string
```

#### type MobileOriginatedUSSDOperation

```go
type MobileOriginatedUSSDOperation string
```


#### type MobileTerminatedUSSDOperation

```go
type MobileTerminatedUSSDOperation string
```


#### type Response

```go
type Response struct {
	StatusCode   string `json:"statusCode"`
	StatusDetail string `json:"statusDetail"`
}
```


```go
var (
	SuccessResponse Response = Response{"S1000", "Success"}
)
```

#### func (Response) Error

```go
func (r Response) Error() *Error
```

#### func (Response) Success

```go
func (r Response) Success() bool
```

#### type SMSClient

```go
type SMSClient struct {
	ApplicationID          string
	Password               string
	SendEndpoint           string
	RetryCount             int
	MaxAddressCount        int
	DeliveryStatusCallback func(messageId, address, status string, timestamp time.Time)
}
```

The SMS client. DeliveryStatusCallback is called to notify a delivery.

#### func (*SMSClient) HandleDeliveryReport

```go
func (client *SMSClient) HandleDeliveryReport(res http.ResponseWriter, req *http.Request)
```
This method should be attached as the handler for the delivery report endpoint.

#### func (*SMSClient) SendTextMessage

```go
func (client *SMSClient) SendTextMessage(message string, recipients []string, chargingAmount float32, requestDeliveryReports bool) (destResps []SMSDestinationResponse, failures []string, err error)
```

#### type SMSDeliveryReport

```go
type SMSDeliveryReport struct {
	DestinationAddress string    `json:"destinationAddress"`
	RawTimestamp       string    `json:"timeStamp"`
	MessageID          string    `json:"requestId"`
	DeliveryStatus     string    `json:"deliveryStatus"`
	Timestamp          time.Time `json:"-"`
}
```


#### type SMSDestinationResponse

```go
type SMSDestinationResponse struct {
	Address      string `json:"address"`
	Timestamp    string `json:"timeStamp"`
	MessageID    string `json:"messageId"`
	StatusCode   string `json:"statusCode"`
	StatusDetail string `json:"statusDetail"`
	Sent         bool   `json:"-"`
	Error        *Error `json:"-"`
}
```


#### type SMSSendRequest

```go
type SMSSendRequest struct {
	ApplicationID         string   `json:"applicationId"`
	Password              string   `json:"password"`
	ChargingAmount        *string  `json:"chargingAmount,omitempty"`
	SourceAddress         string   `json:"sourceAddress"`
	DeliveryStatusRequest *string  `json:"deliveryStatusRequest,omitempty"`
	Message               string   `json:"message"`
	DestinationAddresses  []string `json:"destinationAddresses"`
	Version               *string  `json:"version,omitempty"`
	Encoding              *string  `json:"encoding,omitempty"`
}
```


#### type SMSSendResponse

```go
type SMSSendResponse struct {
	StatusCode           string                   `json:"statusCode"`
	StatusDetail         string                   `json:"statusDetail"`
	RequestID            string                   `json:"requestId"`
	Version              string                   `json:"version"`
	DestinationResponses []SMSDestinationResponse `json:"destinationResponses"`
}
```


#### type StatusCode

```go
type StatusCode string
```


#### type SubscriptionBaseSizeRequest

```go
type SubscriptionBaseSizeRequest struct {
	ApplicationID string `json:"applicationId"`
	Password      string `json:"password"`
}
```


#### type SubscriptionBaseSizeResponse

```go
type SubscriptionBaseSizeResponse struct {
	Version      string `json:"version"`
	StatusCode   string `json:"statusCode"`
	StatusDetail string `json:"statusDetail"`
	BaseSize     string `json:"baseSize"`
}
```


#### type SubscriptionClient

```go
type SubscriptionClient struct {
	ApplicationID              string
	Password                   string
	SubscriptionEndpoint       string
	StatusQueryEndpoint        string
	BaseSizeEndpoint           string
	SubscriptionStatusCallback func(subscriberId, status string, timestamp time.Time)
}
```

The Subscription service client. SubscriptionStatusCallback is called to handle
Subscription notifications.

#### func (*SubscriptionClient) GetBaseSize

```go
func (client *SubscriptionClient) GetBaseSize() (int, error)
```

#### func (*SubscriptionClient) GetStatus

```go
func (client *SubscriptionClient) GetStatus(subscriberId string) (string, error)
```

#### func (*SubscriptionClient) HandleSubscriptionNotification

```go
func (client *SubscriptionClient) HandleSubscriptionNotification(res http.ResponseWriter, req *http.Request)
```
This method should be attached as the subscription notification endpoint
handler.

#### func (*SubscriptionClient) Subscribe

```go
func (client *SubscriptionClient) Subscribe(subscriberId string) (string, error)
```

#### func (*SubscriptionClient) Unsubscribe

```go
func (client *SubscriptionClient) Unsubscribe(subscriberId string) (string, error)
```

#### type SubscriptionNotification

```go
type SubscriptionNotification struct {
	ApplicationID string `json:"applicationId"`
	Frequency     string `json:"frequency"`
	Status        string `json:"status"`
	SubscriberID  string `json:"subscriberId"`
	Version       string `json:"version"`
	Timestamp     string `json:"timeStamp"`
}
```


#### type SubscriptionRequest

```go
type SubscriptionRequest struct {
	ApplicationID string `json:"applicationId"`
	Password      string `json:"password"`
	SubscriberID  string `json:"subscriberId"`
	Action        string `json:"action"`
	Version       string `json:"version"`
}
```


#### type SubscriptionResponse

```go
type SubscriptionResponse struct {
	Version            string `json:"version"`
	RequestID          string `json:"requestId"`
	StatusCode         string `json:"statusCode"`
	StatusDetail       string `json:"statusDetail"`
	SubscriptionStatus string `json:"subscriptionStatus"`
}
```


#### type SubscriptionStatusRequest

```go
type SubscriptionStatusRequest struct {
	ApplicationID string `json:"applicationId"`
	Password      string `json:"password"`
	SubscriberID  string `json:"subscriberId"`
}
```


#### type SubscriptionStatusResponse

```go
type SubscriptionStatusResponse struct {
	Version            string `json:"version"`
	StatusCode         string `json:"statusCode"`
	StatusDetail       string `json:"statusDetail"`
	SubscriptionStatus string `json:"subscriptionStatus"`
}
```


#### type USSDClient

```go
type USSDClient struct {
	ApplicationID              string
	Password                   string
	SendEndpoint               string
	RetryCount                 int
	SessionStore               USSDSessionStore
	IncomingMessageHandlerFunc func(address, message string, operation MobileOriginatedUSSDOperation, sessionData map[string]interface{}) (response string, responseType MobileTerminatedUSSDOperation, err error)
}
```

The USSD client. SessionStore should implement the interface USSDSessionStore.
The provided inMemorySessionStore can be used for this.
IncomingMessageHandlerFunc is called to get the response to a USSD message.

#### func (*USSDClient) HandleIncoming

```go
func (client *USSDClient) HandleIncoming(res http.ResponseWriter, req *http.Request)
```
This method should be attached as the handler to the USSD receiving endpoint of
the server.

#### type USSDMobileOriginatedRequest

```go
type USSDMobileOriginatedRequest struct {
	ApplicationID string                        `json:"applicationId"`
	SourceAddress string                        `json:"sourceAddress"`
	USSDOperation MobileOriginatedUSSDOperation `json:"ussdOperation"`
	Message       string                        `json:"message"`
	//RequestID     string                        `json:"requestId"`
	SessionID string `json:"sessionId"`
	Encoding  string `json:"encoding"`
	Version   string `json:"version"`
}
```


#### type USSDMobileTerminatedRequest

```go
type USSDMobileTerminatedRequest struct {
	ApplicationID      string                        `json:"applicationId"`
	Password           string                        `json:"password"`
	USSDOperation      MobileTerminatedUSSDOperation `json:"ussdOperation"`
	Message            string                        `json:"message"`
	SessionID          string                        `json:"sessionId"`
	DestinationAddress string                        `json:"destinationAddress"`
	Version            *string                       `json:"version,omitempty"`
	Encoding           *string                       `json:"encoding,omitempty"`
}
```


#### type USSDMobileTerminatedResponse

```go
type USSDMobileTerminatedResponse struct {
	Response
	Timestamp string `json:"timestamp"`
	RequestID string `json:"requestId"`
}
```


#### type USSDSession

```go
type USSDSession struct {
	ID            string
	RemoteAddress string
	SessionData   map[string]interface{}
}
```


#### type USSDSessionStore

```go
type USSDSessionStore interface {
	Get(string) *USSDSession
	Save(USSDSession)
}
```

