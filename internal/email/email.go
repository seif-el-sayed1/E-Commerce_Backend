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

func AdminForgotPasswordEmail(token, emailAddress string) error {
	forgotLink := config.Env.AdminForgotPass
	loginLink := config.Env.AdminLogin
	appName := config.Env.AppName

	html := GenerateHTML(TemplateOptions{
		EmailTitle:    "Reset your admin account password",
		EmailSubTitle: "Tap the button below to reset your account password.",
		BtnText:       "Reset Password",
		BtnLink:       forgotLink + token,
		BelowText:     "You can login from here:",
		BelowLink:     loginLink,
		FooterNote:    "You are receiving this email because a request to reset the password for your " + appName + " admin account has been initiated. If you did not initiate this action, please disregard this message.",
	})

	return Send(SendOptions{
		Email:   emailAddress,
		Subject: appName + " reset admin account password",
		HTML:    html,
	})
}

func UserVerificationEmail(code, emailAddress string) error {
	appName := config.Env.AppName

	html := GenerateHTML(TemplateOptions{
		EmailTitle:    "Verify Your user Account",
		EmailSubTitle: "Use the code below to verify your email address.",
		BtnText:       code,
		FooterNote:    "You received this email because you have registered on " + appName + ". If you did not initiate this action, please ignore this email.",
	})

	return Send(SendOptions{
		Email:   emailAddress,
		Subject: appName + " account verification",
		HTML:    html,
	})
}
