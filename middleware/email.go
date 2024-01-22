package middleware

import (
	"net/smtp"
)

func SendEmail(email, password string) error {

	auth := smtp.PlainAuth("", "balkanlingo@gmail.com", password, "smtp.gmail.com")

	to := []string{email}
	// Create the email message
	message := []byte("To: " + email + "\r\n" +
		"Subject: Password reset\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		`<html>
			<head></head>
			<body style="font-family: 'Arial', sans-serif;
				background-color: #f0f0f0; padding: 10%">
				<div class="container" style="max-width: 600px;
					margin: 20px auto;
					padding: 20px;
					background-color: #fff;
					border-radius: 10px;
					box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);">
					<h1 style="color: #333;
						text-align: center;">Vaša jednokratna lozinka</h1>
					<p style="color: #666;">
						<strong>Email:</strong> ` + email + `
					</p>
					<img style="display:none">
					<p style="color: #666;">
						<strong>Jednokratna lozinka:</strong> ` + password + `
					</p>
					<div class="button-container" style="text-align: center;">
						<a href="https://balkanlingo.online/login" class="button" style="display: inline-block;
							padding: 10px 20px;
							font-size: 16px;
							text-align: center;
							text-decoration: none;
							background-color: #4CAF50;
							color: #fff;
							border-radius: 5px;
							cursor: pointer;">Prijavi se</a>
					</div>
					<hr>
					<p style="color: #666;">
						<strong>Molimo vas da se što prije prijavite u sustav i promijenite lozinku.</strong>
					</p>
					<h6 style="color: #4CA;">
						<i>Balkan Lingo, FER projekt</i>
					<h6>
				</div>
			</body>
		</html>
		`)

	err := smtp.SendMail("smtp.gmail.com:587", auth, "balkanlingo@gmail.com", to, message)

	return err

}
