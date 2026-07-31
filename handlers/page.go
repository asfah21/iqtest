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

// ShowPrivacy menampilkan halaman kebijakan privasi
func ShowPrivacy(c *gin.Context) {
	helpers.Render(c, http.StatusOK, pages.PrivacyPage())
}

// ShowTerms menampilkan halaman syarat & ketentuan
func ShowTerms(c *gin.Context) {
	helpers.Render(c, http.StatusOK, pages.TermsPage())
}

// ShowPurchasePolicy menampilkan halaman kebijakan pembelian
func ShowPurchasePolicy(c *gin.Context) {
	helpers.Render(c, http.StatusOK, pages.PurchasePolicyPage())
}

// ShowFaq menampilkan halaman FAQ
func ShowFaq(c *gin.Context) {
	helpers.Render(c, http.StatusOK, pages.FaqPage())
}

// ShowContact menampilkan halaman kontak
func ShowContact(c *gin.Context) {
	helpers.Render(c, http.StatusOK, pages.ContactPage())
}

// Show404 menampilkan halaman 404
func Show404(c *gin.Context) {
	helpers.Render(c, http.StatusNotFound, pages.ErrorPage("Halaman yang Anda cari tidak ditemukan."))
}

// Show500 menampilkan halaman error server
func Show500(c *gin.Context) {
	helpers.Render(c, http.StatusInternalServerError, pages.ErrorPage("Terjadi kesalahan pada server. Silakan coba lagi."))
}
