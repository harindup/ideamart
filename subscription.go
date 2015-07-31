package ideamart

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	SubscriberStatusRegistered    = "REGISTERED"
	SubscriberStatusUnregistered  = "UNREGISTERED"
	SubscriberStatusPendingCharge = "PENDING CHARGE"
	version                       = "1.0"
	subscribeAction               = "1"
	unsubscribeAction             = "0"
	subscriptionTimestampFormat   = "20060102150405"
)

type SubscriptionClient struct {
	ApplicationID              string
	Password                   string
	SubscriptionEndpoint       string
	StatusQueryEndpoint        string
	BaseSizeEndpoint           string
	SubscriptionStatusCallback func(subscriberId, status string, timestamp time.Time)
}

type SubscriptionRequest struct {
	ApplicationID string `json:"applicationId"`
	Password      string `json:"password"`
	SubscriberID  string `json:"subscriberId"`
	Action        string `json:"action"`
	Version       string `json:"version"`
}

type SubscriptionResponse struct {
	Version            string `json:"version"`
	RequestID          string `json:"requestId"`
	StatusCode         string `json:"statusCode"`
	StatusDetail       string `json:"statusDetail"`
	SubscriptionStatus string `json:"subscriptionStatus"`
}

type SubscriptionStatusRequest struct {
	ApplicationID string `json:"applicationId"`
	Password      string `json:"password"`
	SubscriberID  string `json:"subscriberId"`
}

type SubscriptionStatusResponse struct {
	Version            string `json:"version"`
	StatusCode         string `json:"statusCode"`
	StatusDetail       string `json:"statusDetail"`
	SubscriptionStatus string `json:"subscriptionStatus"`
}

type SubscriptionBaseSizeRequest struct {
	ApplicationID string `json:"applicationId"`
	Password      string `json:"password"`
}

type SubscriptionBaseSizeResponse struct {
	Version      string `json:"version"`
	StatusCode   string `json:"statusCode"`
	StatusDetail string `json:"statusDetail"`
	BaseSize     string `json:"baseSize"`
}

type SubscriptionNotification struct {
	ApplicationID string `json:"applicationId"`
	Frequency     string `json:"frequency"`
	Status        string `json:"status"`
	SubscriberID  string `json:"subscriberId"`
	Version       string `json:"version"`
	Timestamp     string `json:"timeStamp"`
}

func parseSubscriptionTimestamp(value string) time.Time {
	t, err := time.ParseInLocation(subscriptionTimestampFormat, value, timestampLocation)
	if err != nil {
		log.Print("Error in parsing SMS timestamp: ", value, err)
	}
	return t
}

func (client *SubscriptionClient) sendSubscriptionRequest(subscriberId, action string) (string, error) {
	req := SubscriptionRequest{
		ApplicationID: client.ApplicationID,
		Password:      client.Password,
		SubscriberID:  subscriberId,
		Version:       version,
		Action:        action,
	}
	res := SubscriptionResponse{}
	err := doRequest(client.SubscriptionEndpoint, req, &res)
	if err != nil {
		return "", err
	}
	log.Print(res)
	return res.SubscriptionStatus, nil
}

func (client *SubscriptionClient) Subscribe(subscriberId string) (string, error) {
	return client.sendSubscriptionRequest(subscriberId, subscribeAction)
}

func (client *SubscriptionClient) Unsubscribe(subscriberId string) (string, error) {
	return client.sendSubscriptionRequest(subscriberId, unsubscribeAction)
}

func (client *SubscriptionClient) GetBaseSize() (int, error) {
	req := SubscriptionBaseSizeRequest{
		ApplicationID: client.ApplicationID,
		Password:      client.Password,
	}
	res := SubscriptionBaseSizeResponse{}
	err := doRequest(client.BaseSizeEndpoint, req, &res)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res.BaseSize)
}

func (client *SubscriptionClient) GetStatus(subscriberId string) (string, error) {
	req := SubscriptionStatusRequest{
		ApplicationID: client.ApplicationID,
		Password:      client.Password,
		SubscriberID:  subscriberId,
	}
	res := SubscriptionStatusResponse{}
	err := doRequest(client.StatusQueryEndpoint, req, &res)
	if err != nil {
		return "", err
	}
	return res.SubscriptionStatus, nil
}

func (client *SubscriptionClient) HandleSubscriptionNotification(res http.ResponseWriter, req *http.Request) {
	notification := SubscriptionNotification{}
	err := unmarshalRequest(req, &notification)
	if err != nil {
		sendErrorResponse(res)
	} else {
		sendSuccessResponse(res)
	}
	req.Body.Close()
	go client.SubscriptionStatusCallback(notification.SubscriberID, notification.Status, parseSubscriptionTimestamp(notification.Timestamp))
}
