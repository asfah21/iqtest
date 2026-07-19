package types

// DashboardUserRow represents a single user row in the admin dashboard table.
type DashboardUserRow struct {
	ID         string
	Nama       string
	Email      string
	SudahBayar bool
	MBTITipe   string
	SkorEI     int
	SkorSN     int
	SkorTF     int
	SkorJP     int
	Dibuat     string
}

// DashboardPageData is the data required by the admin dashboard page.
type DashboardPageData struct {
	Users           []DashboardUserRow
	TotalUser       int
	SudahBayar      int
	BelumBayar      int
	TotalPendapatan int
}
