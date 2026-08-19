package email

import (
	"seif-el-sayed1/E-Commerce_Backend.git/internal/config"
)

func AdminVerificationEmail(token, emailAddress string) error {
	adminURL := config.Env.AdminUrl
	verifyLink := config.Env.AdminVerify
	loginLink := config.Env.AdminLogin
	appName := config.Env.AppName

	html := GenerateHTML(TemplateOptions{
		Link:          adminURL,
		EmailTitle:    "Verify Your Admin Account",
		EmailSubTitle: "Tap the button below to verify your email address.",
		BtnText:       "Verify Account",
		BtnLink:       verifyLink + token,
		BelowText:     "You can login from here:",
		BelowLink:     loginLink,
		FooterNote:    "You received this email because you were added as an admin on " + appName + ". If you did not initiate this action, please ignore this email.",
		FooterLink:    appName,
	})

	return Send(SendOptions{
		Email:   emailAddress,
		Subject: appName + " admin account verification",
		HTML:    html,
	})
}
