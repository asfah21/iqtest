package handlers

import (
	"net/http"

	"ego/helpers"
	"ego/models"
	"ego/services"
	"ego/templ/pages"
	"ego/templ/types"

	"github.com/gin-gonic/gin"
)

// GetQuestions — GET /api/questions
// Mengembalikan 20 soal tanpa correctOption (per IQTEST.md §11.2)
func GetQuestions(c *gin.Context) {
	questions := services.GetQuestions()

	type questionResponse struct {
		ID           string            `json:"id"`
		QuestionCode string            `json:"question_code"`
		Domain       string            `json:"domain"`
		ImageURL     string            `json:"image_url"`
		Options      map[string]string `json:"options"`
	}

	var resp []questionResponse
	for _, q := range questions {
		resp = append(resp, questionResponse{
			ID:           q.ID,
			QuestionCode: q.QuestionCode,
			Domain:       q.Domain,
			ImageURL:     q.ImageURL,
			Options: map[string]string{
				"A": q.OptionImages[0],
				"B": q.OptionImages[1],
				"C": q.OptionImages[2],
				"D": q.OptionImages[3],
			},
		})
	}

	c.JSON(http.StatusOK, resp)
}

// SubmitTest — POST /submit-tes
// Menerima payload JSON { nama, email, answers, tab_switch_count }
func SubmitTest(c *gin.Context) {
	var req struct {
		Nama    string `json:"nama"`
		Email   string `json:"email"`
		Answers []struct {
			QuestionCode   string  `json:"question_code"`
			SelectedOption *string `json:"selected_option"`
			TimeTakenMs    int     `json:"time_taken_ms"`
			TimedOut       bool    `json:"timed_out"`
		} `json:"answers"`
		TabSwitchCount int `json:"tab_switch_count"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data tidak valid: " + err.Error(),
		})
		return
	}

	// Convert to map[string]string for ProcessQuizAnswers
	rawAnswers := make(map[string]string)
	for _, a := range req.Answers {
		if a.SelectedOption != nil {
			rawAnswers[a.QuestionCode] = *a.SelectedOption
		}
	}

	sessionID, err := services.ProcessQuizAnswers(req.Email, req.Nama, rawAnswers, req.TabSwitchCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menyimpan data tes: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": sessionID,
	})
}

// ShowPaywall menampilkan halaman pembayaran
func ShowPaywall(c *gin.Context) {
	id := c.Param("id")
	data, err := services.GetPaywallData(id)
	if err != nil {
		c.String(http.StatusNotFound, "User tidak ditemukan")
		return
	}

	helpers.Render(c, http.StatusOK, pages.PaywallPage(*data))
}

// ShowResult menampilkan hasil kuis (hanya jika sudah bayar)
func ShowResult(c *gin.Context) {
	id := c.Param("id")
	result, err := services.GetQuizResult(id)
	if err != nil {
		c.String(http.StatusNotFound, "Data tidak ditemukan")
		return
	}
	if result == nil {
		c.Redirect(http.StatusSeeOther, "/paywall/"+id+"?error=belum_bayar")
		return
	}

	hasilData := quizResultToHasilData(result)
	helpers.Render(c, http.StatusOK, pages.HasilPage(hasilData))
}

// KonfirmasiBayar memproses konfirmasi pembayaran dari user
func KonfirmasiBayar(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		NamaPengirim string `json:"nama_pengirim"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data tidak valid",
		})
		return
	}

	err := services.ConfirmPayment(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal memproses pembayaran: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"id":      id,
	})
}

// quizResultToHasilData converts the service-layer QuizResult to the template data type.
func quizResultToHasilData(r *models.QuizResult) types.HasilPageData {
	domainViews := make(map[string]types.DomainScoreView)
	domainLabels := map[string]string{
		"MTX": "Penalaran Matriks",
		"SEQ": "Deret Logis",
		"SPA": "Rotasi Spasial",
		"ANL": "Analogi Visual",
	}
	for k, ds := range r.DomainScores {
		domainViews[k] = types.DomainScoreView{
			Domain:     k,
			Percentage: ds.Percentage,
			Label:      domainLabels[k],
		}
	}

	percentile := 0.0
	if r.Percentile != nil {
		percentile = *r.Percentile
	}

	return types.HasilPageData{
		Nama:             r.Nama,
		RawScore:         r.RawScore,
		MaxPossible:      r.MaxPossible,
		Percentile:       percentile,
		EstimatedIQ:      r.EstimatedIQ,
		DomainScores:     domainViews,
		AvgResponseMs:    r.AvgResponseMs,
		IsReliable:       r.IsReliable,
		ExecutiveSummary: r.ExecutiveSummary,
		Kekuatan:         r.Kekuatan,
		AreaPerhatian:    r.AreaPerhatian,
	}
}
