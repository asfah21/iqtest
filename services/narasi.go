package services

import (
	"fmt"
	"sort"
	"strings"

	"ego/models"
)

// ──────────────────────────────────────────────────────────────
// Narrative Generator — menghasilkan laporan kognitif berdasarkan skor domain
// Per IQTEST.md §8 (Result Interpretation)
// ──────────────────────────────────────────────────────────────

// classifyPercentage returns a qualitative label for a domain percentage score.
func classifyPercentage(pct float64) string {
	switch {
	case pct >= 90:
		return "superior"
	case pct >= 75:
		return "sangat_baik"
	case pct >= 60:
		return "baik"
	case pct >= 40:
		return "cukup"
	case pct >= 25:
		return "perlu_latihan"
	default:
		return "rendah"
	}
}

func classifyLabel(pct float64) string {
	switch classifyPercentage(pct) {
	case "superior":
		return "Sangat Baik"
	case "sangat_baik":
		return "Sangat Baik"
	case "baik":
		return "Baik"
	case "cukup":
		return "Cukup"
	case "perlu_latihan":
		return "Perlu Latihan"
	default:
		return "Rendah"
	}
}

// domainNames maps domain codes to Indonesian names.
var domainNames = map[string]string{
	"MTX": "Penalaran Matriks",
	"SEQ": "Deret Logis",
	"SPA": "Rotasi Spasial",
	"ANL": "Analogi Visual",
}

// domainDescriptions provides detailed descriptions of each domain.
var domainDescriptions = map[string]string{
	"MTX": "Kemampuan mengenali pola dalam informasi visual dan menarik kesimpulan induktif.",
	"SEQ": "Kemampuan memahami urutan logis dan aturan berurutan dalam serangkaian informasi.",
	"SPA": "Kemampuan memvisualisasikan objek dalam ruang dan melakukan mental rotation secara akurat.",
	"ANL": "Kemampuan mengidentifikasi hubungan analogis antar konsep visual dan menerapkannya pada konteks baru.",
}

// domainExercises provides exercise recommendations per domain.
var domainExercises = map[string][]string{
	"MTX": {
		"Latih kemampuan pola dengan puzzle Sudoku dan Nonogram",
		"Mainkan game asah otak seperti logic puzzles dan pattern recognition games",
		"Coba teka-teki Raven's Progressive Matrices untuk latihan lanjutan",
	},
	"SEQ": {
		"Pelajari dasar-dasar pemrograman untuk melatih logika sekuensial",
		"Mainkan game strategi seperti catur yang membutuhkan perencanaan langkah",
		"Latih dengan soal deret angka dan figural secara rutin",
	},
	"SPA": {
		"Mainkan game 3D puzzle seperti Rubik's Cube atau puzzle geometri",
		"Latih mental rotation dengan aplikasi khusus atau game puzzle spasial",
		"Belajar membaca dan membuat peta, diagram teknis, atau sketsa 3D",
	},
	"ANL": {
		"Latih dengan soal analogi verbal dan figural secara bergantian",
		"Biasakan mencari hubungan 'A:B seperti C:?' dalam kehidupan sehari-hari",
		"Pelajari mind mapping untuk melatih kemampuan menghubungkan konsep",
	},
}

// ──────────────────────────────────────────────────────────────
// 6.1 — generateExecutiveSummary
// Input: rawScore, maxScore, domainScores, percentile, estimatedIQ
// Output: string deskripsi performa kognitif
// ──────────────────────────────────────────────────────────────

func generateExecutiveSummary(nama string, rawScore, maxScore float64, domainScores map[string]models.DomainScore, percentile float64, estimatedIQ *float64) string {
	pct := rawScore / maxScore * 100
	var overallLabel string
	switch {
	case pct >= 80:
		overallLabel = "sangat baik"
	case pct >= 60:
		overallLabel = "baik"
	case pct >= 40:
		overallLabel = "cukup"
	default:
		overallLabel = "perlu pengembangan"
	}

	intro := fmt.Sprintf(`Halo, %s. Terima kasih telah menyelesaikan asesmen kemampuan kognitif ShadowSelf. Laporan ini memberikan gambaran tentang profil kognitifmu berdasarkan 4 domain:
- Penalaran Matriks (MTX): %s
- Deret Logis (SEQ): %s
- Rotasi Spasial (SPA): %s
- Analogi Visual (ANL): %s`,
		nama,
		classifyLabel(domainScores["MTX"].Percentage),
		classifyLabel(domainScores["SEQ"].Percentage),
		classifyLabel(domainScores["SPA"].Percentage),
		classifyLabel(domainScores["ANL"].Percentage),
	)

	middle := fmt.Sprintf(`Secara keseluruhan, kemampuan kognitifmu berada pada kategori %s dengan skor mentah %.1f dari maksimal %.1f. Pencapaianmu berada pada persentil %.1f, artinya kamu memiliki skor yang sama atau lebih tinggi dari %.1f%% peserta lainnya di platform ini.`, overallLabel, rawScore, maxScore, percentile, percentile)

	var iqNote string
	if estimatedIQ != nil {
		iqNote = fmt.Sprintf(` Estimasi skor IQ berdasarkan data normatif yang tersedia adalah %.0f.`, *estimatedIQ)
	} else {
		iqNote = ` Estimasi IQ belum tersedia karena data normatif masih dalam tahap pengumpulan (target: 1.000 partisipan).`
	}

	closing := `Penting untuk diingat bahwa hasil ini bersifat indikatif dan merupakan potret kemampuan kognitif pada saat tes dilakukan. Kemampuan kognitif dapat berkembang dengan latihan dan pembelajaran berkelanjutan.`

	return strings.TrimSpace(intro + "\n\n" + middle + iqNote + "\n\n" + closing)
}

// ──────────────────────────────────────────────────────────────
// 6.2 — generateKekuatan
// Input: domainScores map[string]DomainScore
// Output: []string — kekuatan berdasarkan domain dengan skor tertinggi
// ──────────────────────────────────────────────────────────────

func generateKekuatan(domainScores map[string]models.DomainScore) []string {
	type domainEntry struct {
		code string
		pct  float64
	}

	var entries []domainEntry
	for code, ds := range domainScores {
		entries = append(entries, domainEntry{code, ds.Percentage})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].pct > entries[j].pct
	})

	var items []string

	// Top 1-2 domains as strengths
	for i, e := range entries {
		if i >= 2 {
			break
		}
		if e.pct >= 40 {
			items = append(items, fmt.Sprintf(`**%s** (%.0f%%) — %s Kemampuan ini termasuk dalam kategori %s dan menjadi salah satu kekuatan kognitif terkuatmu. %s`,
				domainNames[e.code], e.pct, domainDescriptions[e.code], classifyLabel(e.pct), getStrengthSentence(e.code, e.pct)))
		}
	}

	// If all domains are low, provide general encouragement
	if len(items) == 0 {
		items = append(items, "Konsistensi adalah kunci — setiap domain kognitif dapat dilatih dan ditingkatkan dengan pendekatan yang tepat.")
		items = append(items, "Kesadaran diri untuk mengikuti asesmen ini adalah langkah pertama menuju pengembangan kognitif yang terarah.")
	}

	// Ensure 3-5 items
	if len(items) < 3 {
		items = append(items, "Keseimbangan antar domain menunjukkan fleksibilitas kognitif — kemampuan untuk beralih antara berbagai jenis pemikiran adalah aset berharga.")
	}
	if len(items) < 3 {
		items = append(items, "Potensi untuk tumbuh — setiap sesi latihan kognitif memperkuat koneksi neural dan meningkatkan performa.")
	}

	if len(items) > 5 {
		items = items[:5]
	}

	return items
}

func getStrengthSentence(code string, pct float64) string {
	switch classifyPercentage(pct) {
	case "superior", "sangat_baik":
		switch code {
		case "MTX":
			return "Kamu unggul dalam mendeteksi pola kompleks dan menarik kesimpulan dari informasi visual."
		case "SEQ":
			return "Kamu memiliki kemampuan luar biasa dalam memahami urutan dan aturan logis."
		case "SPA":
			return "Visualisasi spasialmu sangat tajam — kamu dapat memutar dan memanipulasi objek mental dengan akurat."
		case "ANL":
			return "Kemampuan analogis yang kuat membantumu melihat hubungan abstrak antar konsep."
		}
	case "baik":
		switch code {
		case "MTX":
			return "Kemampuan penalaran pola yang solid — kamu dapat mengidentifikasi keteraturan dalam informasi visual."
		case "SEQ":
			return "Pemahaman logika sekuensial yang baik — kamu mampu mengikuti dan memprediksi urutan."
		case "SPA":
			return "Kemampuan spasial yang memadai — kamu dapat melakukan orientasi dan rotasi mental dasar."
		case "ANL":
			return "Penalaran analogis yang baik — kamu mampu mentransfer hubungan dari satu konteks ke konteks lain."
		}
	default:
		return "Dengan latihan yang konsisten, kemampuan ini dapat ditingkatkan secara signifikan."
	}
	return ""
}

// ──────────────────────────────────────────────────────────────
// 6.3 — generateAreaPerhatian
// Input: domainScores map[string]DomainScore
// Output: []string — area pengembangan berdasarkan domain dengan skor terendah
// ──────────────────────────────────────────────────────────────

func generateAreaPerhatian(domainScores map[string]models.DomainScore) []string {
	type domainEntry struct {
		code string
		pct  float64
	}

	var entries []domainEntry
	for code, ds := range domainScores {
		entries = append(entries, domainEntry{code, ds.Percentage})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].pct < entries[j].pct
	})

	var items []string

	for i, e := range entries {
		if i >= 2 {
			break
		}
		if e.pct < 60 {
			items = append(items, fmt.Sprintf(`**%s** (%.0f%%) — Meskipun masih dalam tahap pengembangan, area ini memiliki potensi besar untuk ditingkatkan. Fokus pada latihan spesifik dapat membantu memperkuat kemampuan ini.`,
				domainNames[e.code], e.pct))
		}
	}

	if len(items) == 0 {
		items = append(items, "Semua domain berada pada tingkat yang memadai — fokus pada pemeliharaan dan pengayaan di area yang paling sering digunakan dalam aktivitas sehari-hari.")
	}

	if len(items) < 2 {
		items = append(items, "Keseimbangan antar domain adalah kunci — pastikan untuk tidak mengabaikan area yang lebih jarang digunakan.")
	}

	if len(items) < 2 {
		items = append(items, "Variasi dalam latihan kognitif membantu menjaga otak tetap adaptif dan responsif terhadap tantangan baru.")
	}

	if len(items) > 5 {
		items = items[:5]
	}

	return items
}

// ──────────────────────────────────────────────────────────────
// 6.4 — generateRekomendasi
// Input: domainScores
// Output: []string — saran latihan spesifik per domain rendah
// ──────────────────────────────────────────────────────────────

func generateRekomendasi(domainScores map[string]models.DomainScore) []string {
	var items []string

	// Sort from lowest to highest
	type domainEntry struct {
		code string
		pct  float64
	}
	var entries []domainEntry
	for code, ds := range domainScores {
		entries = append(entries, domainEntry{code, ds.Percentage})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].pct < entries[j].pct
	})

	for _, e := range entries {
		if e.pct < 60 {
			items = append(items, fmt.Sprintf(`**%s** (%.0f%%) — %s`, domainNames[e.code], e.pct, strings.Join(domainExercises[e.code], " ")))
		}
	}

	if len(items) == 0 {
		items = append(items, "Pertahankan kebiasaan baikmu — tantang dirimu dengan masalah yang lebih kompleks untuk terus mengembangkan kemampuan kognitif.")
	}

	if len(items) > 5 {
		items = items[:5]
	}

	return items
}

// ──────────────────────────────────────────────────────────────
// 6.5 — GenerateAllNarratives
// Memanggil semua fungsi narasi dan mengembalikan output terstruktur
// ──────────────────────────────────────────────────────────────

func GenerateAllNarratives(nama string, rawScore, maxScore float64, domainScores map[string]models.DomainScore, percentile float64, estimatedIQ *float64) (executiveSummary string, kekuatan, areaPerhatian, rekomendasi []string) {
	executiveSummary = generateExecutiveSummary(nama, rawScore, maxScore, domainScores, percentile, estimatedIQ)
	kekuatan = generateKekuatan(domainScores)
	areaPerhatian = generateAreaPerhatian(domainScores)
	rekomendasi = generateRekomendasi(domainScores)
	return
}
