package xmlfeeds

import "encoding/xml"

type Shop struct {
	XMLName  xml.Name   `xml:"SHOP"`
	ShopItem []ShopItem `xml:"SHOPITEM"`
}

type ShopItem struct {
	ItemID            string        `xml:"ITEM_ID"`
	ProductName       string        `xml:"PRODUCTNAME"`
	Product           string        `xml:"PRODUCT"`
	Description       string        `xml:"DESCRIPTION"`
	URL               string        `xml:"URL"`
	ImgURL            string        `xml:"IMGURL"`
	ImgURLAlternative []string      `xml:"IMGURL_ALTERNATIVE"`
	PriceVAT          string        `xml:"PRICE_VAT"`
	VAT               string        `xml:"VAT"`
	Manufacturer      string        `xml:"MANUFACTURER"`
	CategoryText      string        `xml:"CATEGORYTEXT"`
	EAN               string        `xml:"EAN"`
	DeliveryDate      int           `xml:"DELIVERY_DATE"`
	Delivery          []Delivery    `xml:"DELIVERY"`
	ItemGroupID       string        `xml:"ITEMGROUP_ID"`
	Accessory         []string      `xml:"ACCESSORY"`
	Gift              string        `xml:"GIFT"`
	ExtendedWarranty  *Warranty     `xml:"EXTENDED_WARRANTY"`
	SpecialService    string        `xml:"SPECIAL_SERVICE"`
	SalesVoucher      *SalesVoucher `xml:"SALES_VOUCHER"`
	Param             []Param       `xml:"PARAM"`
}

type Delivery struct {
	DeliveryID       string  `xml:"DELIVERY_ID"`
	DeliveryPrice    float64 `xml:"DELIVERY_PRICE"`
	DeliveryPriceCOD float64 `xml:"DELIVERY_PRICE_COD"`
}

type Param struct {
	Name  string `xml:"PARAM_NAME"`
	Value string `xml:"VAL"`
}

type Warranty struct {
	Value int    `xml:"VAL"`
	Desc  string `xml:"DESC"`
}

type SalesVoucher struct {
	Code string `xml:"CODE"`
	Desc string `xml:"DESC"`
}
