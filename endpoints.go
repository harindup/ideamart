package ideamart

// Endpoint constants
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
