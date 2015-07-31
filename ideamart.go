package ideamart

import "time"

var timestampLocation *time.Location

func init() {
	timestampLocation, _ = time.LoadLocation("Asia/Colombo")
	if timestampLocation == nil {
		panic("Failed to load time zone info. Please check TZData.")
	}
}
