package ideamart

import (
	"log"
	"net/http"
)

type MobileTerminatedUSSDOperation string
type MobileOriginatedUSSDOperation string

// USSD Operation type contants
const (
	MobileTermiatedFinal    MobileTerminatedUSSDOperation = "mt-fin"
	MobileTermiatedContinue MobileTerminatedUSSDOperation = "mt-cont"

	MobileOriginatedInitial  MobileOriginatedUSSDOperation = "mo-init"
	MobileOriginatedContinue MobileOriginatedUSSDOperation = "mo-cont"

	statusCodeSuccess string = "S1000"
)

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

type USSDMobileTerminatedResponse struct {
	Response
	Timestamp string `json:"timestamp"`
	RequestID string `json:"requestId"`
}

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

type USSDSessionStore interface {
	Get(string) *USSDSession
	Save(USSDSession)
}

// The USSD client.
// SessionStore should implement the interface USSDSessionStore.
// The provided inMemorySessionStore can be used for this.
// IncomingMessageHandlerFunc is called to get the response to a USSD message.
type USSDClient struct {
	ApplicationID              string
	Password                   string
	SendEndpoint               string
	RetryCount                 int
	SessionStore               USSDSessionStore
	IncomingMessageHandlerFunc func(address, message string, operation MobileOriginatedUSSDOperation, sessionData map[string]interface{}) (response string, responseType MobileTerminatedUSSDOperation, err error)
}

type USSDSession struct {
	ID            string
	RemoteAddress string
	SessionData   map[string]interface{}
}

func newUSSDSession(id, remoteAddr string) USSDSession {
	return USSDSession{
		ID:            id,
		RemoteAddress: remoteAddr,
		SessionData:   map[string]interface{}{},
	}
}

func (client *USSDClient) sendHandlerResponse(session *USSDSession, ussdReq USSDMobileOriginatedRequest) {
	response, responseType, err := client.IncomingMessageHandlerFunc(session.RemoteAddress, ussdReq.Message, ussdReq.USSDOperation, session.SessionData)
	ussdResp := USSDMobileTerminatedRequest{
		ApplicationID:      client.ApplicationID,
		Password:           client.Password,
		Message:            response,
		SessionID:          session.ID,
		USSDOperation:      responseType,
		DestinationAddress: session.RemoteAddress,
	}
	resp := USSDMobileTerminatedResponse{}
	err = doRequest(client.SendEndpoint, ussdResp, &resp)
	if err != nil {
		log.Print(err)
	}
	if resp.StatusCode != statusCodeSuccess {
		log.Print(apiErrorFromCode(resp.StatusCode), resp)
	}
}

// This method should be attached as the handler to the USSD receiving endpoint of the server.
func (client *USSDClient) HandleIncoming(res http.ResponseWriter, req *http.Request) {
	ussdReq := USSDMobileOriginatedRequest{}
	err := unmarshalRequest(req, &ussdReq)
	var session *USSDSession
	if ussdReq.USSDOperation == MobileOriginatedInitial {
		s := newUSSDSession(ussdReq.SessionID, ussdReq.SourceAddress)
		client.SessionStore.Save(s)
		session = client.SessionStore.Get(ussdReq.SessionID)
	} else {
		session = client.SessionStore.Get(ussdReq.SessionID)
	}
	if err != nil || session == nil {
		sendErrorResponse(res)
		return
	} else {
		sendSuccessResponse(res)
	}
	req.Body.Close()
	go client.sendHandlerResponse(session, ussdReq)
}
