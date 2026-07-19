package handlers

import (
	"net/http"

	"ego/helpers"
	"ego/services"
	"ego/templ/pages"
	"ego/templ/types"

	"github.com/gin-gonic/gin"
)

// Admin credentials (hardcoded)
const (
	adminUsername = "admin"
	adminPassword = "admin360"
)

// ShowLogin menampilkan halaman login admin
func ShowLogin(c *gin.Context) {
	helpers.Render(c, http.StatusOK, pages.LoginPage(""))
}

// LoginProcess memproses login admin
func LoginProcess(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username != adminUsername || password != adminPassword {
		helpers.Render(c, http.StatusOK, pages.LoginPage("Username atau password salah"))
		return
	}

	// Set session cookie (simple token-based)
	c.SetCookie("admin_token", "authenticated", 3600*2, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/admin/dashboard")
}

// ShowDashboard menampilkan dashboard admin dengan data user
func ShowDashboard(c *gin.Context) {
	// Cek autentikasi
	token, err := c.Cookie("admin_token")
	if err != nil || token != "authenticated" {
		c.Redirect(http.StatusSeeOther, "/admin/login")
		return
	}

	users, err := services.GetAllUsers()
	if err != nil {
		helpers.Render(c, http.StatusInternalServerError, pages.ErrorPage("Gagal mengambil data pengguna."))
		return
	}

	// Hitung statistik
	totalUser := len(users)
	sudahBayar := 0
	belumBayar := 0
	totalPendapatan := 0

	var rows []types.DashboardUserRow
	for _, u := range users {
		isPaid := u.StatusPembayaran == "paid"
		if isPaid {
			sudahBayar++
			totalPendapatan += 14900
		} else {
			belumBayar++
		}
		rows = append(rows, types.DashboardUserRow{
			ID:         u.ID,
			Nama:       u.Nama,
			Email:      u.Email,
			SudahBayar: isPaid,
			MBTITipe:   u.MBTITipe,
			SkorEI:     u.SkorEI,
			SkorSN:     u.SkorSN,
			SkorTF:     u.SkorTF,
			SkorJP:     u.SkorJP,
			Dibuat:     "-",
		})
	}

	dashData := types.DashboardPageData{
		Users:           rows,
		TotalUser:       totalUser,
		SudahBayar:      sudahBayar,
		BelumBayar:      belumBayar,
		TotalPendapatan: totalPendapatan,
	}

	helpers.Render(c, http.StatusOK, pages.DashboardPage(dashData))
}

// ShowUserDetail menampilkan detail user untuk admin
func ShowUserDetail(c *gin.Context) {
	// Cek autentikasi
	token, err := c.Cookie("admin_token")
	if err != nil || token != "authenticated" {
		c.Redirect(http.StatusSeeOther, "/admin/login")
		return
	}

	id := c.Param("id")
	user, err := services.GetUserByID(id)
	if err != nil {
		helpers.Render(c, http.StatusNotFound, pages.ErrorPage("User tidak ditemukan."))
		return
	}

	helpers.Render(c, http.StatusOK, pages.UserDetailPage(user))
}

// LogoutProcess menghapus session admin
func LogoutProcess(c *gin.Context) {
	c.SetCookie("admin_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/admin/login")
}
