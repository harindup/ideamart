package ideamart

/*
	The SMS client.
	Only sending is currently supported.
*/

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	// SMS delivery status codes from Ideamart
	SMSStatusSent          = "SENT"
	SMSStatusDelivered     = "DELIVERED"
	SMSStatusExpired       = "EXPIRED"
	SMSStatusDeleted       = "DELETED"
	SMSStatusUndeliverable = "UNDELIVERABLE"
	SMSStatusAccepted      = "ACCEPTED"
	SMSStatusUnknown       = "UNKNOWN"
	SMSStatusRejected      = "REJECTED"

	smsTimestampFormat = "0601021504"
)

// The SMS client.
// DeliveryStatusCallback is called to notify a delivery.
type SMSClient struct {
	ApplicationID          string
	Password               string
	SendEndpoint           string
	RetryCount             int
	MaxAddressCount        int
	DeliveryStatusCallback func(messageId, address, status string, timestamp time.Time)
}

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

type SMSDestinationResponse struct {
	Address      string `json:"address"`
	Timestamp    string `json:"timeStamp"`
	MessageID    string `json:"messageId"`
	StatusCode   string `json:"statusCode"`
	StatusDetail string `json:"statusDetail"`
	Sent         bool   `json:"-"`
	Error        *Error `json:"-"`
}

type SMSSendResponse struct {
	StatusCode           string                   `json:"statusCode"`
	StatusDetail         string                   `json:"statusDetail"`
	RequestID            string                   `json:"requestId"`
	Version              string                   `json:"version"`
	DestinationResponses []SMSDestinationResponse `json:"destinationResponses"`
}

type SMSDeliveryReport struct {
	DestinationAddress string    `json:"destinationAddress"`
	RawTimestamp       string    `json:"timeStamp"`
	MessageID          string    `json:"requestId"`
	DeliveryStatus     string    `json:"deliveryStatus"`
	Timestamp          time.Time `json:"-"`
}

func parseSMSTimestamp(value string) time.Time {
	t, err := time.ParseInLocation(smsTimestampFormat, value, timestampLocation)
	if err != nil {
		log.Print("Error in parsing SMS timestamp: ", value, err)
	}
	return t
}

func splitAddrSlice(addresses []string, maxSize int) [][]string {
	slices := [][]string{}
	c := len(addresses) / maxSize
	for i := 0; i <= c; i++ {
		f := i * maxSize
		l := f + maxSize
		if i == c {
			if f == len(addresses) {
				break
			} else {
				l = len(addresses)
			}
		}
		slices = append(slices, addresses[f:l])
	}
	return slices
}

func (request *SMSSendRequest) sendWithRetries(endpoint string, retryCount int) ([]SMSDestinationResponse, error) {
	resp := SMSSendResponse{}
	for c := 0; c < retryCount; c++ {
		err := doRequest(endpoint, *request, &resp)
		if err != nil && err != ErrInvalidJSON {
			return []SMSDestinationResponse{}, err
		}
		if isSuccessCode(resp.StatusCode) {
			formatDestinationResponses(resp.DestinationResponses)
			return resp.DestinationResponses, nil
		}
		if isErrorCode(resp.StatusCode) {
			apiErr := apiErrorFromCode(resp.StatusCode)
			if !(apiErr == ErrTempSysErr || apiErr == ErrMsgDelivFailed) {
				return resp.DestinationResponses, ErrSendingFailed
			}
		}
	}
	return resp.DestinationResponses, ErrSendingFailed
}

func formatDestinationResponses(responses []SMSDestinationResponse) {
	for i := range responses {
		responses[i].Sent = isSuccessCode(responses[i].StatusCode)
		if responses[i].Sent {
			//responses[i].Timestamp = parseSMSTimestamp(responses[i].rawTimestamp)
		} else {
			e := apiErrorFromCode(responses[i].StatusCode)
			responses[i].Error = &e
		}
	}
}

func (client *SMSClient) sendSMS(sms SMSSendRequest, recipients []string) (destResps []SMSDestinationResponse, failures []string, err error) {
	destResps = []SMSDestinationResponse{}
	failures = []string{}
	addressBlocks := splitAddrSlice(recipients, client.MaxAddressCount)
	for _, block := range addressBlocks {
		sms.DestinationAddresses = block
		d, err := sms.sendWithRetries(client.SendEndpoint, client.RetryCount)
		if err != nil {
			failures = append(failures, block...)
		}
		for _, r := range d {
			if r.Sent {
				destResps = append(destResps, r)
			} else {
				failures = append(failures, r.Address)
			}
		}
	}
	if len(failures) == len(recipients) {
		return destResps, failures, ErrSendingFailed
	}
	return destResps, failures, nil
}

func (client *SMSClient) SendTextMessage(message string, recipients []string, chargingAmount float32, requestDeliveryReports bool) (destResps []SMSDestinationResponse, failures []string, err error) {
	smsReq := SMSSendRequest{
		ApplicationID: client.ApplicationID,
		Password:      client.Password,
		Message:       message,
	}
	if chargingAmount > 0 {
		a := fmt.Sprint(chargingAmount)
		smsReq.ChargingAmount = &a
	}
	if requestDeliveryReports {
		d := "1"
		smsReq.DeliveryStatusRequest = &d
	}
	return client.sendSMS(smsReq, recipients)
}

// This method should be attached as the handler for the delivery report endpoint.
func (client *SMSClient) HandleDeliveryReport(res http.ResponseWriter, req *http.Request) {
	report := SMSDeliveryReport{}
	err := unmarshalRequest(req, &report)
	if err != nil {
		sendErrorResponse(res)
	} else {
		sendSuccessResponse(res)
	}
	req.Body.Close()
	report.Timestamp = parseSMSTimestamp(report.RawTimestamp)
	go client.DeliveryStatusCallback(report.MessageID, report.DestinationAddress, report.DeliveryStatus, report.Timestamp)
}
