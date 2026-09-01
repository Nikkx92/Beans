package storage

const (
	WidthImg     = "w_img"
	HeightImg    = "h_img"
	ITopImg      = "i_top_img"
	IBottomImg   = "i_down_img"
	ILeftImg     = "i_left_img"
	IRightImg    = "i_right_img"
	Kegel        = "kegel"
	ITopKegel    = "i_top_txt"
	IBottomKegel = "i_down_txt"
	ILeftKegel   = "i_left_txt"
	IRightKegel  = "i_right_txt"
	Value        = "value"
)

var ParamGroups = map[string][]string{
	"img": {
		WidthImg, HeightImg, ITopImg, IBottomImg, ILeftImg, IRightImg,
	},
	"text": {
		Kegel, ITopKegel, ILeftKegel, IRightKegel, IBottomKegel,
	},
	"single": {
		Value,
	},
	"icon": {
		WidthImg, HeightImg, ITopImg, IBottomImg, ILeftImg, IRightImg,
		Kegel, ITopKegel, ILeftKegel, IRightKegel, IBottomKegel,
	},
}

type FieldConfig struct {
	Title     string
	Group     string
	IsMobile  bool
	IsDesktop bool
}

var SettingItemsConfig = []FieldConfig{
	{Title: "страна", Group: "icon", IsDesktop: true, IsMobile: true},
	{Title: "сорт", Group: "text", IsDesktop: true, IsMobile: true},
	{Title: "купаж", Group: "text", IsDesktop: true, IsMobile: true},
	{Title: "грейд", Group: "icon", IsDesktop: true, IsMobile: true},
	{Title: "локация", Group: "icon", IsDesktop: true, IsMobile: true},
	{Title: "высота", Group: "icon", IsDesktop: true, IsMobile: true},
	{Title: "сушка", Group: "icon", IsDesktop: true, IsMobile: true},
	{Title: "размер", Group: "icon", IsDesktop: true, IsMobile: true},
	{Title: "вкус", Group: "icon", IsDesktop: true, IsMobile: true},
	{Title: "цена", Group: "text", IsDesktop: true, IsMobile: true},
	{Title: "кол-во", Group: "text", IsDesktop: true, IsMobile: true},
	{Title: "наличие", Group: "text", IsDesktop: true, IsMobile: true},
	{Title: "limited", Group: "img", IsDesktop: true, IsMobile: true},
	{Title: "h_header", Group: "single", IsDesktop: true, IsMobile: true},
	{Title: "w_column1", Group: "single", IsDesktop: true, IsMobile: true},
	{Title: "w_column2", Group: "single", IsDesktop: true, IsMobile: true},
	{Title: "w_column3", Group: "single", IsDesktop: true, IsMobile: true},
	{Title: "w_column4", Group: "single", IsDesktop: true, IsMobile: false},
	{Title: "cell_size", Group: "single", IsDesktop: true, IsMobile: true},
}

type StateCoffee struct {
	Items            []Items           `json:"items"`
	Parameters       []FieldParameters `json:"desc_settings"`
	MobileParameters []FieldParameters `json:"mobile_settings"`
}
type Items struct {
	Country string   `json:"страна"`
	Sorts   []Fields `json:"sorts"`
}

type FieldParameters struct {
	Title      string          `json:"title"`
	Parameters map[string]Type `json:"parameters"`
}

type Fields struct {
	Sort        string `json:"сорт"`
	Blend       string `json:"купаж"`
	Grade       string `json:"грейд"`
	Limited     bool   `json:"limited"`
	Location    string `json:"локация"`
	HeightAbove string `json:"высота"`
	Drying      string `json:"сушка"`
	CornSize    string `json:"размер"`
	Taste       string `json:"вкус"`
	Price       string `json:"цена"`
	Quantity    string `json:"кол-во"`
	StockDesc   string `json:"наличие"`
	InStock     bool   `json:"inStock"`
}

type Type struct {
	Num int
	Str string
}
