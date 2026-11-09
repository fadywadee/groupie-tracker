package data

// Artist يمثل هيكل بيانات الفنان الأساسي
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate *int     `json:"creationDate"` // استخدام مؤشر لتجنب القيمة 0 الخاطئة
	FirstAlbum   *string  `json:"firstAlbum"`
	Locations    string   `json:"locations"`    // رابط API الفرعي للأماكن
	ConcertDates string   `json:"concertDates"` // رابط API الفرعي للتواريخ
	Relations    string   `json:"relations"`    // رابط API الفرعي للعلاقات
}

// Location يمثل أماكن الحفلات
type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

// Date يمثل تواريخ الحفلات
type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// Relation يربط الفنان بالتواريخ والأماكن معاً
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
