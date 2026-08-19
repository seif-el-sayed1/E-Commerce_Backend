// internal/email/template.go
package email

import (
	"fmt"
	"seif-el-sayed1/E-Commerce_Backend.git/internal/config"
)

type TemplateOptions struct {
	Link            string
	Logo            string
	BackgroundColor string
	PrimaryColor    string
	SecondaryColor  string
	EmailTitle      string
	EmailSubTitle   string
	BtnText         string
	BtnLink         string
	BelowText       string
	BelowLink       string
	FooterNote      string
	FooterLink      string
}

func GenerateHTML(opts TemplateOptions) string {
	if opts.BackgroundColor == "" {
		opts.BackgroundColor = "#F4F6F8"
	}
	if opts.PrimaryColor == "" {
		opts.PrimaryColor = "#0F766E"
	}
	if opts.SecondaryColor == "" {
		opts.SecondaryColor = "#ffffff"
	}

	buttonHTML := ""
	if opts.BtnText != "" || opts.BtnLink != "" {
		if opts.BtnLink != "" {
			buttonHTML = fmt.Sprintf(`
				<div class="button-wrapper">
					<a href="%s" class="button" target="_blank">%s</a>
				</div>`, opts.BtnLink, opts.BtnText)
		} else {
			buttonHTML = fmt.Sprintf(`
				<div class="button-wrapper">
					<div class="button">%s</div>
				</div>`, opts.BtnText)
		}
	}

	belowHTML := ""
	if opts.BelowText != "" || opts.BelowLink != "" {
		belowLinkHTML := ""
		if opts.BelowLink != "" {
			belowLinkHTML = fmt.Sprintf(`<a href="%s" target="_blank">%s</a>`, opts.BelowLink, opts.BelowLink)
		}
		belowHTML = fmt.Sprintf(`
			<div class="below">
				<p>%s</p>
				%s
			</div>`, opts.BelowText, belowLinkHTML)
	}

	footerLinkHTML := ""
	if opts.FooterLink != "" {
		footerLinkHTML = fmt.Sprintf(`<a href="%s" target="_blank">%s</a>`, opts.FooterLink, opts.FooterLink)
	}

	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="ar" dir="rtl">
	<head>
		<meta charset="UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>Email</title>
		<style>
			body { margin: 0; padding: 0; background: %s; font-family: 'Arial', 'Tahoma', sans-serif; }
			.container { width: 100%%; max-width: 600px; margin: 30px auto; background: %s; border-radius: 16px; overflow: hidden; box-shadow: 0 10px 30px rgba(0,0,0,0.1); }
			.header { background: %s; padding: 40px 30px 25px; text-align: center; }
			.header h1 { margin: 0; color: white; font-size: 28px; font-weight: bold; }
			.content { padding: 35px 30px; color: #374151; font-size: 16px; line-height: 28px; }
			.button-wrapper { text-align: center; padding: 20px 30px 40px; }
			.button { background: %s; color: #fff !important; padding: 16px 36px; font-size: 17px; border-radius: 12px; text-decoration: none; display: inline-block; font-weight: bold; }
			.below { padding: 0 30px 30px; font-size: 15.5px; line-height: 26px; color: #4B5563; }
			.footer { text-align: center; padding: 25px 30px; background: #F9FAFB; font-size: 13px; color: #9CA3AF; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header"><h1>%s</h1></div>
			<div class="content"><p>%s</p></div>
			%s
			%s
			<div class="thank-you"><p>Thank you<br><strong>%s</strong></p></div>
			<div class="footer"><p>%s</p>%s</div>
		</div>
	</body>
	</html>
	`, opts.BackgroundColor, opts.SecondaryColor, opts.PrimaryColor, opts.PrimaryColor,
		opts.EmailTitle, opts.EmailSubTitle, buttonHTML, belowHTML,
		config.Env.AppName, opts.FooterNote, footerLinkHTML)
}
