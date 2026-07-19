package handlers

import (
	"net/http"
	"time"

	"ego/helpers"
	"ego/middleware"
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

	// 2. API Questions (tanpa correctOption — per IQTEST.md §11.2)
	r.GET("/api/questions", GetQuestions)

	// 3. Halaman Kuesioner
	r.GET("/quiz", ShowQuiz)
	r.HEAD("/quiz", ShowQuiz)

	// 4. Proses Jawaban (dengan rate limiting: max 5/jam per IP — per IQTEST.md §9.2)
	r.POST("/submit-tes", middleware.RateLimitMiddleware(5, time.Hour), SubmitTest)

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
