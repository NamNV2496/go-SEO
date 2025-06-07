package entity

const (
	// template url
	UrlTemplateBuySell  = "mua-ban"
	UrlTemplateExchange = "trao-doi"
	UrlTemplateBuy      = "mua"
	UrlTemplateSell     = "ban"
	UrlTemplateRent     = "cho-thue"
	// kind seo data
	UrlKindLocation = "localtion"
	UrlKindProduct  = "product"
	UrlKindCategory = "category"
	UrlKindBrand    = "brand"
	UrlKindYear     = "year"
	UrlKindMonth    = "month"
)

type DynamicParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
