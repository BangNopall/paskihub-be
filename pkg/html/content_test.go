package html_content

import (
	"strings"
	"testing"
)

func TestGetEmailVerifHTMLUsesPaskiHubVerificationTemplate(t *testing.T) {
	link := "https://paskihub.example/verify?token=abc"

	htmlBody := GetEmailVerifHTML(link)

	assertContains(t, htmlBody, "PaskiHub")
	assertContains(t, htmlBody, "Verifikasi akun PaskiHub kamu")
	assertContains(t, htmlBody, "Terima kasih sudah mendaftar di")
	assertContains(t, htmlBody, "Verifikasi Akun Saya")
	assertContains(t, htmlBody, `href="`+link+`"`)
	assertContains(t, htmlBody, "2026 PaskiHub. All rights reserved.")
	assertContains(t, htmlBody, "Anda menerima email ini karena Anda baru saja mendaftar akun PaskiHub.")
	assertNotContains(t, htmlBody, "HOLOGY")
}

func TestGetEmailForgotPasswordUsesPaskiHubResetTemplate(t *testing.T) {
	link := "https://paskihub.example/reset?token=abc"

	htmlBody := GetEmailForgotPassword(link)

	assertContains(t, htmlBody, "PaskiHub")
	assertContains(t, htmlBody, "Ganti password akun PaskiHub")
	assertContains(t, htmlBody, "Kami menerima permintaan untuk mengganti password akun")
	assertContains(t, htmlBody, "Reset Password")
	assertContains(t, htmlBody, `href="`+link+`"`)
	assertContains(t, htmlBody, "2026 PaskiHub. All rights reserved.")
	assertContains(t, htmlBody, "Anda menerima email ini karena ada permintaan reset password akun PaskiHub.")
	assertNotContains(t, htmlBody, "HOLOGY")
}

func assertContains(t *testing.T, body, want string) {
	t.Helper()
	if !strings.Contains(body, want) {
		t.Fatalf("expected email HTML to contain %q", want)
	}
}

func assertNotContains(t *testing.T, body, unwanted string) {
	t.Helper()
	if strings.Contains(body, unwanted) {
		t.Fatalf("expected email HTML not to contain %q", unwanted)
	}
}
