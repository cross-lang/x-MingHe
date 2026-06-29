package sms

import "testing"

func TestSms(t *testing.T) {
	err := SendSMS([]string{"+8618512869479"}, "2581380", []string{"3306"}, "明河", "1401074343", "yourSecretId", "yourSecretKey")
	t.Log(err)
}