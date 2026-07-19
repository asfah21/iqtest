package models

// QuestionDef — metadata satu soal pilihan ganda bergambar (per IQTEST.md §4.1)
type QuestionDef struct {
	ID             string    // UUID dari database
	QuestionCode   string    // e.g., "Q_MTX_001"
	Domain         string    // "MTX" | "SEQ" | "SPA" | "ANL"
	ImageURL       string    // gambar soal utama
	OptionImages   [4]string // URL gambar opsi A, B, C, D
	CorrectOption  string    // "A" | "B" | "C" | "D" (HANYA di server!)
	Difficulty     string    // "easy" | "medium" | "hard" | "very_hard"
	Weight         float64   // 1.0 / 1.5 / 2.0 / 2.5 sesuai kesulitan
	PValue         *float64  // nullable — dikalibrasi dari data uji coba
	Discrimination *float64  // nullable — dikalibrasi dari data uji coba
}
