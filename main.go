package main

import (
	"log"
	"net/http"
	"net/mail"
    "github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"fmt"
    "math/rand"
    "time"
)

type VerifyOTPRequest struct {
	OTP string `json:"otp" validate:"required"`
	EMAIL string `json:"email" validate:"required"`
}

type UserRecord struct {
    Id        string                 `json:"id"`
    Created   string                 `json:"created"`
    Updated   string                 `json:"updated"`
    CollectionId string              `json:"collectionId"`
    Data      map[string]interface{} `json:"data"`
}

func generateOTP() string {
    // Seed the random number generator to ensure different results each time
    rand.Seed(time.Now().UnixNano())

    // Generate a random number between 100000 and 999999
    otp := rand.Intn(900000) + 100000

    // Convert the OTP to a string
    return fmt.Sprintf("%06d", otp)
}

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/verify-otp", func(c echo.Context) error {
			var req VerifyOTPRequest
			
			// Bind the request JSON to the struct
			if err := c.Bind(&req); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
			}
	
			// Retrieve the user record by email
            userRecord, err := app.Dao().FindFirstRecordByData("users", "email", req.EMAIL)
            if err != nil {
                return c.JSON(http.StatusOK, map[string]string{"message": "","error": "User not found","code": "202"})
            }

            // Compare the received OTP with the stored OTP
            storedOTP := userRecord.GetString("otp")
            if storedOTP != req.OTP {
                return c.JSON(http.StatusOK, map[string]string{"message": "","error": "Invalid OTP","code": "201"})
            }

            // OTP verified successfully
            return c.JSON(http.StatusOK, map[string]string{"message": "OTP verified successfully", "code": "200"})
       
		}, /* optional middlewares */)
	
		return nil
	})

	// fires only for "users" collections
	app.OnRecordAfterCreateRequest("users").Add(func(e *core.RecordCreateEvent) error {

		otp := generateOTP() 

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
			HTML: `
				<!DOCTYPE html>
				<html lang="en">
				<head>
					<meta charset="UTF-8">
					<meta name="viewport" content="width=device-width, initial-scale=1.0">
					<title>OTP for Login</title>
					<style>
						body {
							font-family: Arial, sans-serif;
							background-color: #f4f4f4;
							margin: 0;
							padding: 0;
						}
						.container {
							width: 100%;
							max-width: 600px;
							margin: 0 auto;
							background-color: #ffffff;
							padding: 20px;
							box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
						}
						.header {
							text-align: center;
							padding: 10px 0;
						}
						.header img {
							width: 150px;
						}
						.content {
							padding: 20px;
							text-align: center;
						}
						.content h1 {
							color: #333333;
						}
						.content p {
							color: #666666;
							line-height: 1.5;
						}
						.otp {
							display: inline-block;
							padding: 10px 20px;
							font-size: 18px;
							color: #000000;
							background-color: #ffcc00;
							border-radius: 5px;
							margin-top: 20px;
							text-decoration: none;
						}
						.footer {
							text-align: center;
							padding: 10px;
							font-size: 12px;
							color: #999999;
						}
					</style>
				</head>
				<body>
					<div class="container">
					
						<div class="content">
							<h1>Welcome to GreatApe</h1>
							<p>Hello, welcome to GreatApe. Here is your One-Time Password (OTP) to signup:</p>
							<div class="otp">` + otp + `</div>
						</div>
						
					</div>
				</body>
				</html>
			`,
			// bcc, cc, attachments and custom headers are also supported...
		}

		return app.NewMailClient().Send(message)
	})

	//fires only for "users" collections
	app.OnRecordAuthRequest("users").Add(func(e *core.RecordAuthEvent) error {

		otp := generateOTP() 

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
			Subject: "OTP for Login to GreatApe",
			HTML: `
				<!DOCTYPE html>
				<html lang="en">
				<head>
					<meta charset="UTF-8">
					<meta name="viewport" content="width=device-width, initial-scale=1.0">
					<title>OTP for Login</title>
					<style>
						body {
							font-family: Arial, sans-serif;
							background-color: #f4f4f4;
							margin: 0;
							padding: 0;
						}
						.container {
							width: 100%;
							max-width: 600px;
							margin: 0 auto;
							background-color: #ffffff;
							padding: 20px;
							box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
						}
						.header {
							text-align: center;
							padding: 10px 0;
						}
						.header img {
							width: 150px;
						}
						.content {
							padding: 20px;
							text-align: center;
						}
						.content h1 {
							color: #333333;
						}
						.content p {
							color: #666666;
							line-height: 1.5;
						}
						.otp {
							display: inline-block;
							padding: 10px 20px;
							font-size: 18px;
							color: #000000;
							background-color: #ffcc00;
							border-radius: 5px;
							margin-top: 20px;
							text-decoration: none;
						}
						.footer {
							text-align: center;
							padding: 10px;
							font-size: 12px;
							color: #999999;
						}
					</style>
				</head>
				<body>
					<div class="container">
					
						<div class="content">
							<h1>Welcome to GreatApe</h1>
							<p>Hello, welcome to GreatApe. Here is your One-Time Password (OTP) to login:</p>
							<div class="otp">` + otp + `</div>
						</div>
						
					</div>
				</body>
				</html>
			`,
			// bcc, cc, attachments and custom headers are also supported...
		}

		return app.NewMailClient().Send(message)
    })

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
