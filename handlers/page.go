package handlers

import (
	"net/http"

	"ego/helpers"
	"ego/templ/pages"

	"github.com/gin-gonic/gin"
)

// ShowHome menampilkan halaman utama
func ShowHome(c *gin.Context) {
	helpers.Render(c, http.StatusOK, pages.IndexPage())
}

// ShowQuiz menampilkan halaman kuesioner
func ShowQuiz(c *gin.Context) {
	helpers.Render(c, http.StatusOK, pages.QuizPage())
}

// ShowTentang menampilkan halaman tentang kami
func ShowTentang(c *gin.Context) {
	helpers.Render(c, http.StatusOK, pages.TentangPage())
}

// Show404 menampilkan halaman 404
func Show404(c *gin.Context) {
	helpers.Render(c, http.StatusNotFound, pages.ErrorPage("Halaman yang Anda cari tidak ditemukan."))
}

// Show500 menampilkan halaman error server
func Show500(c *gin.Context) {
	helpers.Render(c, http.StatusInternalServerError, pages.ErrorPage("Terjadi kesalahan pada server. Silakan coba lagi."))
}
