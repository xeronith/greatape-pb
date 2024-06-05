package hooks

import (
	"greateape-pb/utility"
	"log"
	"net/mail"
	"strings"

	_ "embed"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

//go:embed templates/otp_signup.html
var SignupOTPTemplate string

func RecordAfterCreateRequestHook(e *core.RecordCreateEvent) error {

	otp := utility.GenerateOTP()

	// Update the user record with the OTP
	e.Record.Set("otp", otp)
	err := app.Dao().SaveRecord(e.Record)
	if err != nil {
		log.Println("Error updating record:", err)
		return err
	}

	// Send the email
	message := &mailer.Message{
		From: mail.Address{
			Address: app.Settings().Meta.SenderAddress,
			Name:    app.Settings().Meta.SenderName,
		},
		To:      []mail.Address{{Address: e.Record.Email()}},
		Subject: "OTP for Signup to GreatApe",
		HTML:    strings.ReplaceAll(SignupOTPTemplate, "{{ OTP }}", otp),
	}

	return app.NewMailClient().Send(message)
}
