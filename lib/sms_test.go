package lib

import "testing"

func TestSms(t *testing.T) {
	var s = NewAliSms()
	err := s.Send("15900000001", `{"customer":"test"}`, "SMS_71390007")

	if err != nil {
		t.Errorf("call error: %v", err)
	}
}
