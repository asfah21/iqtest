package handlers

import (
	"net/http"

	"ego/helpers"
	"ego/templ/pages"

	"github.com/gin-gonic/gin"
)

// SetupRoutes mendaftarkan semua route ke Gin engine
func SetupRoutes(r *gin.Engine) {
	// 0. Recovery Middleware untuk 5xx — harus PALING ATAS
	r.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				helpers.Render(c, http.StatusInternalServerError, pages.ErrorPage("Terjadi kesalahan pada server. Silakan coba lagi."))
				c.Abort()
			}
		}()
		c.Next()
	})

	// 1. Halaman Utama (HEAD juga penting untuk curl -I dan health check OpenResty)
	r.GET("/", ShowHome)
	r.HEAD("/", ShowHome)

	// 2. Halaman Kuesioner
	r.GET("/quiz", ShowQuiz)
	r.HEAD("/quiz", ShowQuiz)

	// 3. Proses Jawaban
	r.POST("/submit-tes", SubmitTest)

	// 4. Paywall
	r.GET("/paywall/:id", ShowPaywall)
	r.HEAD("/paywall/:id", ShowPaywall)

	// 5. Konfirmasi Pembayaran
	r.POST("/konfirmasi-bayar/:id", KonfirmasiBayar)

	// 6. Hasil Premium (hanya jika PAID)
	r.GET("/hasil/:id", ShowResult)
	r.HEAD("/hasil/:id", ShowResult)

	// 7. Halaman Informasi
	r.GET("/tentang", ShowTentang)
	r.HEAD("/tentang", ShowTentang)

	// 8. Admin Routes
	r.GET("/admin/login", ShowLogin)
	r.HEAD("/admin/login", ShowLogin)
	r.POST("/admin/login", LoginProcess)
	r.GET("/admin/dashboard", ShowDashboard)
	r.HEAD("/admin/dashboard", ShowDashboard)
	r.GET("/admin/user/:id", ShowUserDetail)
	r.HEAD("/admin/user/:id", ShowUserDetail)
	r.GET("/admin/logout", LogoutProcess)

	// 9. Handle 404 — harus di PALING AKHIR
	r.NoRoute(Show404)
}
