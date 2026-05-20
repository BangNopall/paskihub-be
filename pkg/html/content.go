package html_content

func GetEmailVerifHTML(link string) string {
	return getPaskiHubActionEmailHTML(
		link,
		"Verifikasi akun PaskiHub kamu",
		`Terima kasih sudah mendaftar di <strong style="font-weight: 700;">PaskiHub</strong>.<br>Sebelum mulai menggunakan layanan kami, silakan verifikasi akun kamu terlebih dahulu. Klik tombol di bawah ini untuk verifikasi akun.`,
		"Verifikasi Akun Saya",
		"Anda menerima email ini karena Anda baru saja mendaftar akun PaskiHub.",
		"766",
		"148",
		"168",
	)
}

func GetEmailForgotPassword(link string) string {
	return getPaskiHubActionEmailHTML(
		link,
		"Ganti password akun PaskiHub",
		`Kami menerima permintaan untuk mengganti password akun <strong style="font-weight: 700;">PaskiHub</strong> kamu. Klik tombol di bawah ini untuk mengganti password`,
		"Reset Password",
		"Anda menerima email ini karena ada permintaan reset password akun PaskiHub.",
		"746",
		"148",
		"168",
	)
}

func getPaskiHubActionEmailHTML(link, title, description, buttonText, footerNote, height, contentTop, footerSpacer string) string {
	return `
<!doctype html>
<html lang="id">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="x-ua-compatible" content="ie=edge">
    <title>PaskiHub</title>
</head>
<body style="margin: 0; padding: 0; background-color: #f5f7fb;">
    <table role="presentation" width="100%" cellspacing="0" cellpadding="0" border="0" style="margin: 0; padding: 0; border-collapse: collapse; background-color: #f5f7fb;">
        <tr>
            <td align="center" style="padding: 0;">
                <table role="presentation" width="600" cellspacing="0" cellpadding="0" border="0" style="width: 600px; max-width: 600px; height: ` + height + `px; border-collapse: collapse; background-color: #ffffff;">
                    <tr>
                        <td style="height: ` + contentTop + `px; line-height: ` + contentTop + `px; font-size: 0;">&nbsp;</td>
                    </tr>
                    <tr>
                        <td align="center" style="padding: 0 72px;">
                            <table role="presentation" width="456" cellspacing="0" cellpadding="0" border="0" style="width: 456px; border-collapse: collapse;">
                                <tr>
                                    <td align="center" style="padding: 0; color: #0f172a; font-family: Montserrat, Arial, Helvetica, sans-serif; font-size: 30px; font-weight: 700; line-height: 36px;">
                                        PaskiHub
                                    </td>
                                </tr>
                                <tr>
                                    <td style="height: 48px; line-height: 48px; font-size: 0;">&nbsp;</td>
                                </tr>
                                <tr>
                                    <td align="center" style="padding: 0; color: #0f172a; font-family: Montserrat, Arial, Helvetica, sans-serif; font-size: 24px; font-weight: 500; line-height: 36px;">
                                        ` + title + `
                                    </td>
                                </tr>
                                <tr>
                                    <td style="height: 32px; line-height: 32px; font-size: 0;">&nbsp;</td>
                                </tr>
                                <tr>
                                    <td align="center" style="padding: 0; color: #737373; font-family: Poppins, Arial, Helvetica, sans-serif; font-size: 14px; font-weight: 400; line-height: 20px;">
                                        ` + description + `
                                    </td>
                                </tr>
                                <tr>
                                    <td style="height: 48px; line-height: 48px; font-size: 0;">&nbsp;</td>
                                </tr>
                                <tr>
                                    <td align="center" bgcolor="#ff6f59" style="height: 48px; padding: 0; background-color: #ff6f59; border-radius: 50px;">
                                        <a href="` + link + `" target="_blank" style="display: block; padding: 14px 28px; color: #ffffff; font-family: Poppins, Arial, Helvetica, sans-serif; font-size: 14px; font-weight: 700; line-height: 20px; text-align: center; text-decoration: none; border-radius: 50px;">
                                            ` + buttonText + `
                                        </a>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                    <tr>
                        <td style="height: ` + footerSpacer + `px; line-height: ` + footerSpacer + `px; font-size: 0;">&nbsp;</td>
                    </tr>
                    <tr>
                        <td align="center" bgcolor="#cdeaff" style="height: 144px; padding: 0 72px; background-color: #cdeaff;">
                            <table role="presentation" width="456" cellspacing="0" cellpadding="0" border="0" style="width: 456px; border-collapse: collapse;">
                                <tr>
                                    <td style="height: 48px; line-height: 48px; font-size: 0;">&nbsp;</td>
                                </tr>
                                <tr>
                                    <td align="center" style="padding: 0; color: #404040; font-family: Poppins, Arial, Helvetica, sans-serif; font-size: 12px; font-weight: 500; line-height: 16px;">
                                        &copy; 2026 PaskiHub. All rights reserved.
                                    </td>
                                </tr>
                                <tr>
                                    <td style="height: 8px; line-height: 8px; font-size: 0;">&nbsp;</td>
                                </tr>
                                <tr>
                                    <td align="center" style="padding: 0; color: #404040; font-family: Poppins, Arial, Helvetica, sans-serif; font-size: 12px; font-weight: 400; line-height: 16px;">
                                        ` + footerNote + `
                                    </td>
                                </tr>
                                <tr>
                                    <td style="height: 48px; line-height: 48px; font-size: 0;">&nbsp;</td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`
}

func GetBasicEmail(content string) string {
	htmlBody := `
<div style="margin: 0; padding: 0;">
    <div align="center" style="display: flex; justify-content: center;">
        <div style="background-color: #000025; background-image: url('https://lh3.googleusercontent.com/d/1TUWZ5C94WYx9E9QqV84mlQYqmVXw8Nrx'); object-fit: cover; height: 100%; max-width: 34rem; margin: auto; position: relative;">
            <div style="padding: 1.5rem 0;">
                <a style="text-align: center;" href="https://hology.ub.ac.id/">
                    <img src="https://lh3.googleusercontent.com/d/1nJs3eWbLXeXMNCv920BsnLtzhqF91zD_" alt="">
                </a>
            </div>
            <div style="background: linear-gradient(to right, #EE6418, #F1C117); height: 1px; margin-bottom: 6em;"></div>
            <div style="padding-left: 4rem; padding-right: 4rem; margin-bottom: 4vh;">
                <div style="color: #fff; text-align: left; margin-bottom: 1.25rem; min-width: 420px;">
                    <p style="overflow-wrap: break-word; white-space: pre-line;">
                        ` + content + `
                    </p>
                </div>
                <div style="text-align: left; color: white;">
                    Best regards, Hology Admin
                </div>
            </div>        
            <div>
                <p style="color: white;">Our Social Media</p>
                <div>
                    <a style="margin: 0 2px; color: white; text-align: center; display: inline;" href="https://www.tiktok.com/@hology_ub" target="_blank">
                        <img src="https://lh3.googleusercontent.com/d/1N_PBErdQFzqAVKK8o-xxEJTZmiVV2MJ8" alt="" width="24em" height="24em">
                    </a>  
                    <a style="margin: 0 2px; color: white; text-align: center; display: inline;" href="https://x.com/HOLOGY_UB" target="_blank">
                        <img src="https://lh3.googleusercontent.com/d/1T_cmDou8t61dZGyNbDctzhbm7rQVBoFp" alt="" width="24em" height="24em">
                    </a>                
                    <a style="margin: 0 2px; color: white; text-align: center; display: inline;" href="https://www.youtube.com/@hologyub9984" target="_blank">
                        <img src="https://lh3.googleusercontent.com/d/1NStY8W1qwZM1IDeAoQHUnMJYtu7wjdU6" alt="" width="24em" height="24em">
                    </a>
                    <a style="margin: 0 2px; color: white; text-align: center; display: inline;" href="https://www.instagram.com/hology_ub/" target="_blank">
                        <img src="https://lh3.googleusercontent.com/d/1Qcvfyi27CtwgtMoUgttN4PBbGcLpQ-Sj" alt="" width="24em" height="24em">
                    </a>
                    <a style="margin: 0 2px; color: white; text-align: center;  display: inline;" href="https://www.linkedin.com/company/hology-ub/" target="_blank">
                        <img src="https://lh3.googleusercontent.com/d/1cRijj4t_PFUrn1pG5hhCoEDzhl-JZJeJ" alt="" width="24em" height="24em">
                    </a>
                </div>
            </div>
            <div>
                <p style="color: white;">Contact Person</p>
                <div>
                    <a href="https://wa.me/628563462687" style="margin: 0 2px; text-align: center; display: inline;" target="_blank">
                        <img src="https://lh3.googleusercontent.com/d/1BRWkFWl46ldYQduzQGWGQzEqrWDLMjM_" alt="" width="24em">
                    </a>
                    <a href="https://wa.me/62895611807906" style="margin: 0 2px; text-align: center; display: inline;" target="_blank">
                        <img src="https://lh3.googleusercontent.com/d/1BRWkFWl46ldYQduzQGWGQzEqrWDLMjM_" alt="" width="24em">
                    </a>                    
                </div>
                <div>

                </div>
            </div>
            <div style="margin: 2em 0 0 0;">
                <div style="background: linear-gradient(to right, #EE6418, #F1C117); height: 1px; margin-bottom: 2rem;"></div>
                <p style="color: white; padding-bottom: 1rem;">&copy; HOLOGY 8.0 All Right Reserved</p>
            </div>
        </div>
    </div>
</div>
    `

	return htmlBody
}
