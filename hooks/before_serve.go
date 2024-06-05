package hooks

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
)

func BeforeServeHook(e *core.ServeEvent) error {
	e.Router.POST("/verify-otp", func(c echo.Context) error {
		var req struct {
			OTP   string `json:"otp" validate:"required"`
			EMAIL string `json:"email" validate:"required"`
		}

		// Bind the request JSON to the struct
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// Retrieve the user record by email
		userRecord, err := app.Dao().FindFirstRecordByData("users", "email", req.EMAIL)
		if err != nil {
			return c.JSON(http.StatusOK, map[string]string{"message": "", "error": "User not found", "code": "202"})
		}

		// Compare the received OTP with the stored OTP
		storedOTP := userRecord.GetString("otp")
		if storedOTP != req.OTP {
			return c.JSON(http.StatusOK, map[string]string{"message": "", "error": "Invalid OTP", "code": "201"})
		}

		// OTP verified successfully
		return c.JSON(http.StatusOK, map[string]string{"message": "OTP verified successfully", "code": "200"})

	})

	return nil
}
