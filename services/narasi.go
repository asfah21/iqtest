package services

import (
	"fmt"
	"strings"
)

// ──────────────────────────────────────────────────────────────
// Narrative Generator — transforms raw Dark Triad scores into
// a premium, personalized relationship assessment report.
//
// Every narrative is dynamically constructed from score combinations
// so no two results read exactly the same.
// ──────────────────────────────────────────────────────────────

// level returns a qualitative label for a 0–100 score.
func level(s int) string {
	switch {
	case s <= 25:
		return "rendah"
	case s <= 50:
		return "sedang"
	case s <= 75:
		return "cukup tinggi"
	default:
		return "tinggi"
	}
}

// levelEn returns English label for internal logic.
func levelEn(s int) string {
	switch {
	case s <= 25:
		return "low"
	case s <= 50:
		return "moderate"
	case s <= 75:
		return "moderately_high"
	default:
		return "high"
	}
}

// ── Executive Summary ────────────────────────────────────────

func generateExecutiveSummary(nama string, n, m, p int) string {
	nl := levelEn(n)
	ml := levelEn(m)
	pl := levelEn(p)

	intro := fmt.Sprintf(`Halo, %s. Terima kasih telah menyelesaikan asesmen ShadowSelf. Laporan ini dirancang untuk membantumu memahami pola-pola kepribadian yang memengaruhi cara kamu menjalin dan menjalani hubungan — bukan sebagai label, melainkan sebagai cermin untuk refleksi diri.`, nama)

	var middle string
	switch {
	case nl == "high" && ml == "high" && pl == "high":
		middle = `Hasil asesmen menunjukkan bahwa kamu memiliki energi kepribadian yang cukup dominan dalam hubungan. Kamu cenderung memiliki kebutuhan yang kuat akan perhatian dan validasi, serta secara aktif membentuk dinamika hubungan sesuai dengan keinginanmu. Keputusan dalam hubungan sering kali diambil secara spontan, mengikuti perasaan sesaat. Ini menciptakan dinamika yang intens — penuh gairah, namun juga berpotensi menimbulkan gesekan jika tidak dikelola dengan kesadaran diri yang baik.`
	case nl == "low" && ml == "low" && pl == "low":
		middle = `Hasil asesmen menunjukkan bahwa kamu memiliki pendekatan yang tenang dan reflektif dalam hubungan. Kamu tidak terlalu bergantung pada validasi eksternal, cenderung membiarkan dinamika hubungan mengalir secara alami, dan mengambil keputusan dengan pertimbangan yang matang. Ini adalah fondasi yang sehat untuk hubungan yang stabil dan saling menghormati. Tantanganmu mungkin lebih pada bagaimana memastikan kebutuhanmu tetap terkomunikasikan dengan jelas.`
	case nl == "high" && ml == "low":
		middle = `Hasil asesmen menunjukkan kombinasi yang menarik: kamu memiliki kebutuhan emosional yang cukup tinggi akan perhatian dan validasi, namun cenderung tidak menggunakan strategi untuk mengarahkan situasi. Artinya, kamu mungkin sangat menginginkan kedekatan dan pengakuan, tetapi lebih memilih untuk menerima dinamika hubungan apa adanya. Ini bisa menciptakan kerentanan — kamu mungkin sering merasa kurang dihargai tanpa secara aktif mengomunikasikannya.`
	case nl == "low" && ml == "high":
		middle = `Kamu menunjukkan pendekatan yang cukup strategis dalam hubungan, namun tanpa kebutuhan emosional yang tinggi akan validasi. Ini menciptakan dinamika yang unik: kamu cenderung memikirkan langkah-langkahmu dalam hubungan, namun tidak terlalu bergantung pada respons pasangan untuk merasa berharga. Kemandirian emosional ini bisa menjadi kekuatan, selama tidak membuat pasangan merasa dijaga jarak.`
	case pl == "high" && nl == "high":
		middle = `Kombinasi antara kebutuhan validasi yang tinggi dan kecenderungan pengambilan keputusan yang spontan menciptakan dinamika hubungan yang ekspresif dan penuh semangat. Kamu cenderung mengejar apa yang kamu inginkan dengan intensitas emosional yang tinggi. Ini bisa sangat menarik pada awalnya, namun penting untuk menyadari bahwa hubungan yang sehat juga membutuhkan ruang untuk refleksi dan konsistensi.`
	case pl == "low" && nl == "low":
		middle = `Kamu menunjukkan kestabilan emosional yang baik dalam hubungan. Keputusan diambil dengan hati-hati, dan kamu tidak terlalu bergantung pada validasi dari pasangan. Ini menciptakan dinamika yang aman dan dapat diprediksi — kualitas yang sangat berharga untuk hubungan jangka panjang. Tantanganmu mungkin adalah memastikan hubungan tetap terasa segar dan tidak jatuh ke dalam rutinitas yang monoton.`
	default:
		middle = fmt.Sprintf(`Hasil asesmen menunjukkan bahwa dalam hubungan, kamu memiliki tingkat kebutuhan validasi emosional yang %s, kecenderungan %s dalam mengarahkan situasi, dan pendekatan yang %s dalam pengambilan keputusan. Kombinasi ini menciptakan dinamika hubungan yang khas — dengan kekuatan dan tantangannya masing-masing — yang akan kita bahas lebih dalam di bagian-bagian berikutnya.`, level(n), level(m), level(p))
	}

	closing := `Penting untuk diingat bahwa asesmen ini bukanlah vonis. Ini adalah titik awal untuk memahami dirimu lebih dalam. Setiap pola kepribadian memiliki konteksnya masing-masing, dan kesadaran adalah langkah pertama menuju hubungan yang lebih sehat dan memuaskan.`

	return strings.TrimSpace(intro + "\n\n" + middle + "\n\n" + closing)
}

// ── Relationship Profile ─────────────────────────────────────

func generateRelationshipProfile(n, m, p int) string {
	nl := levelEn(n)
	ml := levelEn(m)
	pl := levelEn(p)

	// Cara mencari perhatian dan validasi
	var validasi string
	switch nl {
	case "high":
		validasi = `Kamu memiliki kecenderungan yang cukup kuat untuk mencari perhatian dan validasi dalam hubungan. Kehadiran dan pengakuan pasangan terasa penting bagi rasa percaya dirimu. Kamu mungkin merasa lebih bersemangat ketika mendapat pujian atau perhatian, dan bisa merasa kurang dihargai jika pasangan tidak memberikan respons yang kamu harapkan. Ini bukan sesuatu yang salah — namun penting untuk menyadari bahwa validasi eksternal hanyalah satu dari sekian sumber harga diri.`
	case "moderately_high":
		validasi = `Kamu cukup menghargai perhatian dan validasi dari pasangan, namun tidak sepenuhnya bergantung padanya. Kamu menikmati momen-momen ketika pasangan menunjukkan apresiasi, tetapi kamu juga memiliki sumber harga diri lain yang membuatmu tetap stabil. Keseimbangan ini sehat — selama kamu tetap bisa mengomunikasikan kebutuhanmu tanpa menuntut.`
	case "moderate":
		validasi = `Kamu memiliki kebutuhan validasi yang seimbang. Kamu menikmati perhatian dari pasangan, namun tidak terlalu bergantung padanya untuk merasa berharga. Ini menunjukkan tingkat kemandirian emosional yang cukup baik. Kamu cenderung tidak mencari-cari pengakuan secara berlebihan, dan ini bisa menciptakan dinamika yang santai dalam hubungan.`
	default:
		validasi = `Kamu cenderung tidak terlalu bergantung pada validasi eksternal dalam hubungan. Harga dirimu relatif stabil dan tidak banyak dipengaruhi oleh seberapa banyak perhatian yang kamu terima dari pasangan. Ini adalah fondasi yang kuat untuk hubungan yang sehat — kamu mampu memberikan ruang bagi pasangan tanpa perlu terus-menerus mencari konfirmasi. Tantanganmu mungkin adalah memastikan bahwa pasangan tetap merasa dibutuhkan dan dihargai.`
	}

	// Cara menghadapi konflik
	var konflik string
	switch {
	case ml == "high" && pl == "high":
		konflik = `Dalam menghadapi konflik, kamu cenderung mengambil pendekatan yang aktif dan kadang impulsif. Kamu mungkin mencoba mengarahkan situasi sesuai keinginanmu, dan melakukannya dengan cara yang spontan. Ini bisa membuatmu terlihat dominan dalam konflik, namun juga berisiko membuat keputusan yang kurang dipikirkan matang-matang. Belajar untuk memberi ruang dan menunda respons bisa membantu menciptakan resolusi yang lebih seimbang.`
	case ml == "high" && pl == "low":
		konflik = `Kamu cenderung menghadapi konflik dengan strategi yang terencana. Kamu tidak mudah terbawa emosi dan lebih suka memikirkan langkah-langkah sebelum bertindak. Ini bisa menjadi kekuatan besar dalam menyelesaikan masalah — namun pastikan bahwa pendekatanmu yang terstruktur tidak membuat pasangan merasa sedang "dikelola" daripada diajak berdiskusi setara.`
	case ml == "low" && pl == "high":
		konflik = `Dalam konflik, kamu cenderung spontan dan mengikuti perasaan sesaat. Kamu mungkin kesulitan menahan diri untuk tidak merespons secara impulsif ketika emosi memuncak. Sisi positifnya, kamu autentik dan tidak menyimpan perasaan. Tantangannya adalah melatih diri untuk mengambil jeda sebelum merespons, agar konflik tidak berlarut-larut karena kata-kata yang terucap tanpa pertimbangan.`
	default:
		konflik = `Kamu cenderung menghadapi konflik dengan kepala dingin. Kamu tidak mudah terprovokasi dan tidak merasa perlu untuk selalu mengontrol jalannya diskusi. Ini adalah pendekatan yang dewasa dan konstruktif. Pasanganmu mungkin merasa aman karena kamu bisa menjadi penengah yang tenang. Pastikan saja bahwa ketenanganmu tidak disalahartikan sebagai ketidakpedulian.`
	}

	// Cara mengambil keputusan
	var keputusan string
	switch pl {
	case "high":
		keputusan = `Keputusan dalam hubungan — mulai dari hal kecil hingga keputusan besar — sering kamu ambil secara spontan. Kamu mengandalkan intuisi dan perasaan saat itu juga. Ini membuat hubungan terasa dinamis dan tidak membosankan. Namun, untuk keputusan-keputusan penting yang berdampak jangka panjang, melatih diri untuk melibatkan pertimbangan yang lebih matang bisa membantu menghindari penyesalan di kemudian hari.`
	case "moderately_high":
		keputusan = `Kamu cukup spontan dalam mengambil keputusan, namun masih menyisakan ruang untuk pertimbangan. Ada keseimbangan antara mengikuti intuisi dan memikirkan konsekuensi. Ini adalah posisi yang fleksibel — kamu bisa bergerak cepat saat diperlukan, namun juga bisa melambat saat situasi menuntut kehati-hatian.`
	case "moderate":
		keputusan = `Kamu mengambil keputusan dengan pertimbangan yang cukup matang. Tidak terlalu lambat hingga kehilangan momen, namun juga tidak terlalu cepat hingga ceroboh. Pendekatan ini menciptakan stabilitas dalam hubungan — pasangan bisa mengandalkan keputusanmu karena biasanya sudah melalui proses pikir yang memadai.`
	default:
		keputusan = `Kamu cenderung berhati-hati dalam mengambil keputusan, terutama yang menyangkut hubungan. Kamu butuh waktu untuk mempertimbangkan berbagai sudut pandang sebelum bertindak. Ini adalah kualitas yang berharga — keputusanmu biasanya matang dan tidak terburu-buru. Tantangannya adalah memastikan bahwa proses berpikir yang panjang tidak membuat pasangan merasa frustrasi atau tidak diutamakan.`
	}

	// Cara memengaruhi pasangan
	var pengaruh string
	switch {
	case ml == "high" && n >= 50:
		pengaruh = `Kamu memiliki pengaruh yang cukup kuat terhadap pasangan, dan kamu sadar akan hal itu. Kamu cenderung menggunakan kombinasi antara kedekatan emosional dan strategi untuk mendapatkan apa yang kamu inginkan dalam hubungan. Ini bisa efektif, namun penting untuk memastikan bahwa pengaruhmu tidak membuat pasangan merasa tertekan atau kehilangan otonomi. Hubungan yang sehat adalah tentang saling memengaruhi, bukan mengendalikan.`
	case ml == "high" && n < 50:
		pengaruh = `Kamu cenderung memengaruhi pasangan melalui pendekatan yang lebih strategis daripada emosional. Kamu mungkin tidak terlalu ekspresif dalam menunjukkan kebutuhan, namun kamu punya cara untuk membuat situasi berjalan sesuai keinginanmu. Kekuatanmu terletak pada kemampuan membaca situasi. Tantangannya adalah memastikan pasangan tidak merasa bahwa mereka sedang "dikelola" tanpa disadari.`
	case ml == "low" && n >= 50:
		pengaruh = `Pengaruhmu terhadap pasangan lebih banyak berasal dari kedekatan emosional daripada strategi. Kamu cenderung memengaruhi dengan cara menunjukkan kebutuhan dan perasaanmu secara langsung. Ini autentik dan bisa menciptakan kedekatan yang dalam. Namun, penting juga untuk belajar bahwa tidak semua hal harus berjalan sesuai keinginanmu — dan itu tidak mengurangi nilai dirimu.`
	default:
		pengaruh = `Kamu cenderung tidak terlalu aktif dalam memengaruhi pasangan. Kamu lebih suka membiarkan dinamika berjalan alami dan tidak merasa perlu untuk mengarahkan situasi. Ini bisa menciptakan hubungan yang santai dan penuh penerimaan. Pastikan bahwa sikap ini tidak membuat kebutuhanmu sendiri terabaikan — komunikasi tetaplah penting.`
	}

	// Tingkat keterlibatan emosional
	var emosional string
	switch {
	case n >= 60 && p >= 60:
		emosional = `Keterlibatan emosionalmu dalam hubungan cukup intens. Kamu memiliki kebutuhan emosional yang tinggi dan cenderung mengekspresikannya secara spontan. Ini menciptakan hubungan yang penuh gairah dan dramatis — menarik, namun juga menguras energi. Menemukan keseimbangan antara intensitas dan ketenangan bisa menjadi kunci untuk hubungan yang lebih berkelanjutan.`
	case n >= 60 && p < 40:
		emosional = `Keterlibatan emosionalmu cukup dalam, namun terkendali. Kamu memiliki kebutuhan untuk terhubung secara emosional, namun kamu juga mampu menahan diri. Ini adalah kombinasi yang baik — kamu bisa hadir secara emosional tanpa kehilangan kendali. Pasangan mungkin merasa bahwa kamu adalah tempat yang aman untuk berbagi perasaan.`
	case n < 40 && p >= 60:
		emosional = `Keterlibatan emosionalmu dalam hubungan cenderung lebih pada tindakan daripada kata-kata. Kamu mungkin tidak terlalu ekspresif secara emosional, namun kamu menunjukkan perhatian melalui keputusan dan tindakan spontan. Pasangan perlu memahami bahwa caramu terlibat mungkin berbeda dari yang mereka harapkan — dan itu tidak berarti kamu tidak peduli.`
	default:
		emosional = `Keterlibatan emosionalmu dalam hubungan cukup stabil dan seimbang. Kamu mampu hadir untuk pasangan tanpa kehilangan dirimu sendiri. Ini adalah fondasi yang sehat untuk hubungan jangka panjang. Kamu cenderung tidak terlalu reaktif secara emosional, namun juga tidak dingin — sebuah keseimbangan yang ideal.`
	}

	sections := []string{
		"**Mencari Perhatian dan Validasi**\n" + validasi,
		"**Menghadapi Konflik**\n" + konflik,
		"**Mengambil Keputusan**\n" + keputusan,
		"**Memengaruhi Pasangan**\n" + pengaruh,
		"**Tingkat Keterlibatan Emosional**\n" + emosional,
	}
	return strings.Join(sections, "\n\n")
}

// ── Kekuatan Potensial ───────────────────────────────────────

func generateKekuatan(n, m, p int) []string {
	var items []string

	nl := levelEn(n)
	ml := levelEn(m)
	pl := levelEn(p)

	// Based on low scores (strengths in moderation)
	if nl == "low" || nl == "moderate" {
		items = append(items, "Kemandirian emosional yang baik — kamu tidak bergantung pada validasi eksternal untuk merasa berharga, yang menciptakan fondasi hubungan yang lebih stabil.")
	}
	if ml == "low" || ml == "moderate" {
		items = append(items, "Pendekatan yang autentik dan transparan dalam hubungan — kamu tidak merasa perlu menggunakan strategi untuk mendapatkan apa yang kamu inginkan.")
	}
	if pl == "low" || pl == "moderate" {
		items = append(items, "Kemampuan mengambil keputusan dengan pertimbangan matang — kamu cenderung tidak bertindak impulsif, sehingga keputusanmu lebih dapat diandalkan.")
	}

	// Based on high scores (reframed as strengths)
	if nl == "high" || nl == "moderately_high" {
		items = append(items, "Kepekaan emosional yang tinggi — kamu mampu merasakan dan merespons kebutuhan emosional dalam hubungan, yang bisa menciptakan kedekatan yang dalam dengan pasangan.")
	}
	if ml == "high" || ml == "moderately_high" {
		items = append(items, "Kemampuan navigasi sosial yang baik — kamu pandai membaca dinamika hubungan dan menyesuaikan pendekatanmu sesuai situasi.")
	}
	if pl == "high" || pl == "moderately_high" {
		items = append(items, "Keberanian dan spontanitas — kamu tidak takut mengambil risiko dalam hubungan, yang bisa membawa pengalaman baru dan kesegaran dalam dinamika bersama pasangan.")
	}

	// Combination-based strengths
	if n >= 40 && n <= 60 {
		items = append(items, "Keseimbangan antara kebutuhan validasi dan kemandirian — kamu mampu menerima kasih sayang tanpa kehilangan identitas dirimu.")
	}
	if m <= 40 && p <= 40 {
		items = append(items, "Ketenangan dan stabilitas emosional — kamu menjadi tempat yang aman bagi pasangan untuk menjadi diri mereka sendiri tanpa tekanan.")
	}
	if n >= 50 && m <= 50 {
		items = append(items, "Keterbukaan emosional yang diimbangi dengan penerimaan — kamu mampu mengekspresikan kebutuhan tanpa memaksakan kehendak.")
	}

	// Ensure we have 3-5 items
	if len(items) < 3 {
		items = append(items, "Kesadaran diri untuk mengikuti asesmen ini — langkah pertama menuju hubungan yang lebih sehat adalah kemauan untuk memahami diri sendiri.")
	}
	if len(items) < 3 {
		items = append(items, "Potensi untuk tumbuh — setiap pola kepribadian bisa dikelola dan diarahkan dengan kesadaran dan latihan.")
	}

	if len(items) > 5 {
		items = items[:5]
	}

	return items
}

// ── Area yang Perlu Diperhatikan ─────────────────────────────

func generateAreaPerhatian(n, m, p int) []string {
	var items []string

	nl := levelEn(n)
	ml := levelEn(m)
	pl := levelEn(p)

	if nl == "high" || nl == "moderately_high" {
		items = append(items, "Ketergantungan pada validasi eksternal — penting untuk mengembangkan sumber harga diri yang tidak hanya berasal dari pengakuan pasangan, agar hubungan tidak menjadi beban emosional.")
	}
	if ml == "high" || ml == "moderately_high" {
		items = append(items, "Kecenderungan untuk mengarahkan situasi secara berlebihan — cobalah sesekali melepaskan kendali dan membiarkan dinamika berjalan alami, agar pasangan merasa setara dalam hubungan.")
	}
	if pl == "high" || pl == "moderately_high" {
		items = append(items, "Impulsivitas dalam pengambilan keputusan — melatih diri untuk mengambil jeda sebelum bertindak bisa membantu menghindari keputusan yang mungkin disesali di kemudian hari.")
	}

	// Combination-based concerns
	if n >= 60 && m >= 60 {
		items = append(items, "Kombinasi antara kebutuhan validasi yang tinggi dan kecenderungan strategis bisa menciptakan dinamika yang rumit — pastikan bahwa hubungan tidak berubah menjadi arena negosiasi yang melelahkan.")
	}
	if n >= 60 && p >= 60 {
		items = append(items, "Intensitas emosional yang tinggi perlu dikelola dengan baik agar tidak menimbulkan kelelahan dalam hubungan. Belajar untuk menenangkan diri sebelum merespons bisa sangat membantu.")
	}
	if m >= 60 && p >= 60 {
		items = append(items, "Kombinasi antara pendekatan strategis dan impulsif bisa membuat pasangan merasa kewalahan. Cobalah untuk lebih transparan tentang niat dan perasaanmu.")
	}
	if n <= 25 && m <= 25 {
		items = append(items, "Kecenderungan untuk terlalu pasif dalam hubungan — pastikan bahwa kebutuhan dan perasaanmu tetap terkomunikasikan, meskipun kamu tidak merasa perlu untuk mendominasi.")
	}
	if p <= 25 && m >= 60 {
		items = append(items, "Kombinasi antara kehati-hatian yang tinggi dan pendekatan strategis bisa membuatmu tampak terlalu terkontrol. Jangan takut untuk menunjukkan kerentanan sesekali.")
	}

	// Ensure we have 3-5 items
	if len(items) < 3 {
		items = append(items, "Komunikasi asertif — pastikan bahwa kamu mampu menyampaikan kebutuhan dan perasaanmu dengan jelas, tanpa menunggu pasangan menebaknya.")
	}
	if len(items) < 3 {
		items = append(items, "Keseimbangan antara memberi dan menerima — hubungan yang sehat membutuhkan kedua pihak untuk saling mendukung dan mendengarkan.")
	}

	if len(items) > 5 {
		items = items[:5]
	}

	return items
}

// ── Relationship Insight ─────────────────────────────────────

func generateRelationshipInsight(n, m, p int) string {
	nl := levelEn(n)
	ml := levelEn(m)
	pl := levelEn(p)

	// Build insight based on score combinations
	type pattern struct {
		condition bool
		insight   string
	}

	patterns := []pattern{
		{
			nl == "high" && ml == "high" && pl == "high",
			`Pola "Intensitas Total" — Ketiga dimensi berada pada tingkat yang tinggi, menciptakan kepribadian yang sangat dominan dalam hubungan. Kamu cenderung menjadi pusat dari dinamika hubungan: mencari perhatian, mengarahkan situasi, dan bertindak spontan. Pasangan mungkin merasa tertarik pada energimu yang besar, namun juga bisa merasa kewalahan. Kuncinya adalah mengembangkan kesadaran diri untuk mengenali kapan intensitasmu membantu dan kapan justru menghambat keintiman.`,
		},
		{
			nl == "low" && ml == "low" && pl == "low",
			`Pola "Ketenangan Reflektif" — Ketiga dimensi berada pada tingkat yang rendah, menunjukkan pendekatan yang tenang, reflektif, dan tidak reaktif dalam hubungan. Kamu tidak mudah terpengaruh oleh dinamika eksternal dan cenderung membiarkan hubungan mengalir secara alami. Ini adalah fondasi yang stabil, namun perlu diingat bahwa hubungan juga membutuhkan gairah dan inisiatif. Tantanganmu adalah memastikan bahwa ketenanganmu tidak disalahartikan sebagai ketidakpedulian.`,
		},
		{
			nl == "high" && ml == "high" && pl != "high",
			`Pola "Arsitek Hubungan" — Kamu memiliki kebutuhan validasi yang tinggi dan cenderung menggunakan strategi untuk mendapatkannya, namun tanpa impulsivitas yang berlebihan. Ini membuatmu seperti "arsitek" dalam hubungan — kamu tahu apa yang kamu inginkan dan punya rencana untuk mencapainya. Ini bisa efektif, namun pastikan bahwa pasangan merasa diajak membangun bersama, bukan sekadar mengikuti cetak birumu.`,
		},
		{
			nl == "high" && pl == "high" && ml != "high",
			`Pola "Gairah Spontan" — Kebutuhan validasi yang tinggi dikombinasikan dengan spontanitas menciptakan dinamika yang penuh gairah dan ekspresif. Kamu mengejar apa yang kamu inginkan dengan antusiasme yang menular. Hubungan bersamamu terasa hidup dan tidak membosankan. Tantangannya adalah menjaga konsistensi — gairah yang besar kadang diikuti oleh kejenuhan yang cepat jika tidak diimbangi dengan komitmen yang stabil.`,
		},
		{
			ml == "high" && pl == "high" && nl != "high",
			`Pola "Pemain Catur" — Kamu cenderung strategis dan spontan secara bersamaan — kombinasi yang jarang. Kamu bisa membaca situasi dengan cepat dan mengambil keputusan taktis dalam sekejap. Dalam hubungan, ini membuatmu selalu selangkah lebih maju. Namun, pasangan mungkin merasa bahwa mereka sedang berhadapan dengan lawan, bukan kekasih. Cobalah untuk sesekali melepaskan "permainan" dan menikmati momen apa adanya.`,
		},
		{
			nl == "low" && ml == "high" && pl == "low",
			`Pola "Pengamat Strategis" — Kamu tidak terlalu bergantung pada validasi emosional, namun kamu cenderung mengamati dan mengarahkan situasi dengan cara yang terukur. Ini membuatmu objektif dan tidak mudah goyah. Dalam hubungan, kamu bisa menjadi penasihat yang baik dan penengah yang netral. Tantanganmu adalah menunjukkan kehangatan emosional — pasangan mungkin merasa bahwa kamu terlalu rasional dan kurang ekspresif secara perasaan.`,
		},
		{
			nl == "high" && ml == "low" && pl == "low",
			`Pola "Hati yang Terbuka" — Kebutuhan validasi yang tinggi tanpa kecenderungan strategis atau impulsif menunjukkan bahwa kamu adalah pribadi yang emosional namun terkendali. Kamu menginginkan perhatian dan kedekatan, namun kamu mengejarnya dengan cara yang tulus dan tidak memaksa. Ini adalah kombinasi yang menarik — kamu terlihat rentan namun juga aman. Pastikan bahwa kamu tidak menahan diri terlalu banyak sehingga kebutuhanmu tidak terpenuhi.`,
		},
		{
			nl == "low" && ml == "low" && pl == "high",
			`Pola "Petualang Bebas" — Kamu tidak terikat oleh kebutuhan validasi atau strategi, namun kamu bertindak spontan. Ini menciptakan kepribadian yang bebas dan tidak terduga. Dalam hubungan, kamu membawa kesegaran dan petualangan. Pasangan tidak akan bosan bersamamu. Tantangannya adalah komitmen dan konsistensi — pastikan bahwa spontanitasmu tidak membuat pasangan merasa tidak aman atau tidak diutamakan.`,
		},
	}

	// Find matching pattern
	for _, p := range patterns {
		if p.condition {
			return p.insight
		}
	}

	// Fallback: generic insight based on dominant trait
	dominant := "narsisme"
	domScore := n
	if m > domScore {
		dominant = "machiavellian"
		domScore = m
	}
	if p > domScore {
		dominant = "psikopati"
		domScore = p
	}

	switch dominant {
	case "narsisme":
		return fmt.Sprintf(`Dimensi yang paling menonjol dalam hasil asesmenmu adalah kebutuhan validasi emosional (skor %d). Ini menunjukkan bahwa dalam hubungan, kamu sangat menghargai perhatian, pengakuan, dan respons emosional dari pasangan. Dinamika hubunganmu kemungkinan besar sangat dipengaruhi oleh seberapa baik kebutuhan ini terpenuhi. Kesadaran akan pola ini adalah langkah penting — bukan untuk mengubah siapa dirimu, tetapi untuk memahami mengapa kamu bereaksi dengan cara tertentu dalam situasi tertentu.`, n)
	case "machiavellian":
		return fmt.Sprintf(`Dimensi yang paling menonjol adalah kecenderungan mengarahkan situasi (skor %d). Kamu cenderung mendekati hubungan dengan kesadaran strategis — kamu memperhatikan pola, membaca situasi, dan menyesuaikan pendekatanmu. Ini bukan sesuatu yang negatif; banyak orang sukses dalam hubungan karena mereka mampu memahami dinamika yang terjadi. Yang perlu diperhatikan adalah apakah pendekatan ini membuatmu sulit untuk benar-benar melepaskan kendali dan menikmati momen tanpa perhitungan.`, m)
	default:
		return fmt.Sprintf(`Dimensi yang paling menonjol adalah pengambilan keputusan impulsif (skor %d). Kamu adalah pribadi yang bertindak berdasarkan dorongan hati dan perasaan sesaat. Dalam hubungan, ini bisa berarti kamu sangat autentik — apa yang kamu rasakan, itulah yang kamu lakukan. Namun, pola ini juga bisa menimbulkan ketidakstabilan jika tidak diimbangi dengan refleksi. Belajar untuk menciptakan jeda antara perasaan dan tindakan bisa menjadi keterampilan yang sangat berharga.`, p)
	}
}

// ── Compatibility Notes ──────────────────────────────────────

func generateCompatibilityNotes(n, m, p int) string {
	nl := levelEn(n)
	ml := levelEn(m)
	pl := levelEn(p)

	var cocok strings.Builder
	cocok.WriteString("**Tipe Pasangan yang Biasanya Cocok**\n")

	switch {
	case nl == "high":
		cocok.WriteString("Kamu cenderung cocok dengan pasangan yang hangat, ekspresif, dan mampu memberikan perhatian yang kamu butuhkan tanpa merasa kelelahan. Pasangan yang memiliki kepercayaan diri yang stabil dan tidak mudah terancam oleh kebutuhan emosionalmu bisa menjadi pendamping yang ideal. Mereka yang mampu memberikan validasi secara alami — tanpa diminta — akan membuatmu merasa aman dan dihargai.\n\n")
	case nl == "low":
		cocok.WriteString("Kemandirian emosionalmu membuatmu cocok dengan berbagai tipe kepribadian. Kamu tidak terlalu menuntut secara emosional, sehingga pasangan yang lebih pendiam atau introvert pun bisa merasa nyaman. Pasangan yang menghargai ruang pribadi dan tidak terlalu bergantung pada validasi eksternal akan menemukan dinamika yang selaras denganmu.\n\n")
	default:
		cocok.WriteString("Kamu cocok dengan pasangan yang memiliki keseimbangan antara kehangatan dan kemandirian — seseorang yang bisa memberikan perhatian tanpa membuatmu merasa terkekang. Pasangan yang komunikatif dan mampu membaca kebutuhanmu tanpa perlu diinstruksikan akan menciptakan dinamika yang harmonis.\n\n")
	}

	switch {
	case ml == "high":
		cocok.WriteString("Dengan kecenderungan strategismu, kamu membutuhkan pasangan yang cerdas secara emosional dan mampu mengikuti alur berpikirmu. Pasangan yang tidak mudah dimanipulasi dan bisa menandingi ketajaman analisismu akan membuat hubungan terasa menantang secara intelektual. Hindari pasangan yang terlalu naif atau mudah dipengaruhi — kamu mungkin secara tidak sadar mengambil keuntungan dari dinamika tersebut.\n\n")
	case ml == "low":
		cocok.WriteString("Pendekatanmu yang lugas dan tanpa strategi membuatmu cocok dengan pasangan yang juga menghargai kejujuran dan transparansi. Hubungan yang tidak rumit dan penuh keterbukaan adalah lingkungan yang paling nyaman bagimu. Pasangan yang tidak suka permainan dalam hubungan akan sangat menghargai ketulusanmu.\n\n")
	default:
		cocok.WriteString("Kamu membutuhkan pasangan yang bisa diajak berdiskusi secara terbuka dan tidak defensif. Komunikasi yang jujur dan langsung adalah kunci — pasangan yang bisa menerima masukan tanpa merasa diserang akan menciptakan hubungan yang sehat bersamamu.\n\n")
	}

	switch {
	case pl == "high":
		cocok.WriteString("Spontanitasmu membutuhkan pasangan yang fleksibel dan tidak terlalu kaku dengan rencana. Pasangan yang bisa mengikuti ritme impulsifmu dan menikmati kejutan akan cocok dengan gayamu. Namun, pastikan pasanganmu juga memiliki cukup stabilitas untuk menjadi jangkar saat keputusan impulsifmu membutuhkan koreksi.\n\n")
	case pl == "low":
		cocok.WriteString("Kamu membutuhkan pasangan yang sabar dan tidak terburu-buru dalam mengambil keputusan. Pasangan yang menghargai proses dan tidak memaksamu untuk bertindak sebelum siap akan membuatmu merasa aman. Hindari pasangan yang terlalu impulsif — perbedaan tempo pengambilan keputusan bisa menjadi sumber konflik.\n\n")
	default:
		cocok.WriteString("Keseimbanganmu dalam pengambilan keputusan membuatmu cocok dengan berbagai tipe. Kamu bisa menyesuaikan diri dengan pasangan yang lebih cepat maupun lebih lambat dalam memutuskan sesuatu. Fleksibilitas ini adalah aset berharga dalam hubungan.\n\n")
	}

	cocok.WriteString("**Tipe Dinamika yang Berpotensi Menimbulkan Konflik**\n")

	conflicts := []string{}
	if nl == "high" {
		conflicts = append(conflicts, "Pasangan yang cenderung dingin secara emosional atau sulit mengekspresikan perasaan. Kebutuhan validasimu mungkin tidak terpenuhi, menyebabkan frustrasi di kedua pihak.")
	}
	if ml == "high" {
		conflicts = append(conflicts, "Pasangan yang sangat dominan dan juga suka mengarahkan situasi. Benturan kehendak bisa menciptakan hubungan yang kompetitif daripada kolaboratif.")
	}
	if pl == "high" {
		conflicts = append(conflicts, "Pasangan yang sangat kaku dan membutuhkan kepastian dalam segala hal. Perbedaan tempo dan gaya pengambilan keputusan bisa menjadi sumber gesekan yang terus-menerus.")
	}
	if nl == "low" && ml == "low" {
		conflicts = append(conflicts, "Pasangan yang sangat membutuhkan drama dan konfrontasi emosional. Ketenanganmu mungkin disalahartikan sebagai ketidakpedulian, sementara kebutuhan mereka akan intensitas bisa terasa menguras energi.")
	}
	if len(conflicts) < 2 {
		conflicts = append(conflicts, "Pasangan yang tidak mau berkompromi atau sulit menerima sudut pandang yang berbeda. Hubungan yang sehat membutuhkan fleksibilitas dari kedua pihak.")
	}
	if len(conflicts) < 2 {
		conflicts = append(conflicts, "Dinamika di mana salah satu pihak selalu mengejar dan pihak lain selalu menjauh. Ketidakseimbangan dalam tingkat keterlibatan emosional bisa menimbulkan kelelahan hubungan.")
	}

	for i, c := range conflicts {
		cocok.WriteString(fmt.Sprintf("%d. %s\n", i+1, c))
	}

	return cocok.String()
}

// ── Reflection Questions ─────────────────────────────────────

func generateReflectionQuestions(n, m, p int) []string {
	nl := levelEn(n)
	ml := levelEn(m)
	pl := levelEn(p)

	questions := []string{}

	// Always include 1-2 universal questions
	questions = append(questions, "Dari hasil asesmen ini, adakah satu pola yang paling kamu kenali dalam dirimu? Bagaimana pola itu muncul dalam hubungan terakhirmu?")
	questions = append(questions, "Ketika membaca bagian kekuatan dan area yang perlu diperhatikan, mana yang paling menantang untuk kamu akui? Mengapa?")

	// Trait-specific questions
	switch {
	case nl == "high" || nl == "moderately_high":
		questions = append(questions, "Seberapa sering kamu merasa kurang dihargai dalam hubungan? Apakah perasaan itu muncul karena pasangan benar-benar tidak memberikan perhatian, atau karena ekspektasimu yang mungkin perlu disesuaikan?")
	case nl == "low":
		questions = append(questions, "Apakah kamu merasa nyaman dengan tingkat kemandirian emosionalmu saat ini, atau adakah bagian dirimu yang sebenarnya ingin lebih terhubung secara emosional dengan pasangan?")
	}

	switch {
	case ml == "high" || ml == "moderately_high":
		questions = append(questions, "Dalam situasi apa kamu merasa perlu untuk mengarahkan jalannya hubungan? Apa yang kamu takutkan terjadi jika kamu melepaskan kendali?")
	case ml == "low":
		questions = append(questions, "Apakah kamu pernah merasa kebutuhanmu terabaikan karena kamu terlalu membiarkan dinamika berjalan alami? Bagaimana kamu bisa lebih aktif tanpa merasa memaksa?")
	}

	switch {
	case pl == "high" || pl == "moderately_high":
		questions = append(questions, "Pernahkah kamu menyesali keputusan yang diambil secara impulsif dalam hubungan? Apa yang bisa kamu lakukan untuk menciptakan jeda sebelum bertindak?")
	case pl == "low":
		questions = append(questions, "Apakah kehati-hatianmu dalam mengambil keputusan pernah membuatmu kehilangan momen penting dalam hubungan? Bagaimana menyeimbangkan antara pertimbangan dan spontanitas?")
	}

	// Ensure we have exactly 3 questions
	if len(questions) > 3 {
		// Keep the first universal question, then pick 2 most relevant trait-specific ones
		selected := []string{questions[0]}
		// Pick the last 2 from the middle (skip the second universal)
		for i := len(questions) - 1; i >= 1 && len(selected) < 3; i-- {
			if i != 1 { // skip the second universal
				selected = append(selected, questions[i])
			}
		}
		questions = selected
	}

	return questions
}

// ── Public API ───────────────────────────────────────────────

// GenerateAllNarratives generates all narrative sections for a given set of scores.
func GenerateAllNarratives(nama string, narsisme, machiavellian, psikopati int) (executiveSummary, relationshipProfile string, kekuatan, areaPerhatian []string, relationshipInsight, compatibilityNotes string, reflectionQuestions []string) {
	executiveSummary = generateExecutiveSummary(nama, narsisme, machiavellian, psikopati)
	relationshipProfile = generateRelationshipProfile(narsisme, machiavellian, psikopati)
	kekuatan = generateKekuatan(narsisme, machiavellian, psikopati)
	areaPerhatian = generateAreaPerhatian(narsisme, machiavellian, psikopati)
	relationshipInsight = generateRelationshipInsight(narsisme, machiavellian, psikopati)
	compatibilityNotes = generateCompatibilityNotes(narsisme, machiavellian, psikopati)
	reflectionQuestions = generateReflectionQuestions(narsisme, machiavellian, psikopati)
	return
}
