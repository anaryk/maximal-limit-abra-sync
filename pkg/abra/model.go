package abra

type ContactData struct {
	Name       string `json:"name"`
	Street     string `json:"street"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile"`
}

type APIResponse struct {
	Winstrom struct {
		Version string `json:"@version"`
		Success string `json:"success"`
		Stats   struct {
			Created any `json:"created"`
			Updated any `json:"updated"`
			Deleted any `json:"deleted"`
			Skipped any `json:"skipped"`
			Failed  any `json:"failed"`
		} `json:"stats"`
		Results []struct {
			ID     string `json:"id,omitempty"`
			Ref    string `json:"ref,omitempty"`
			Errors []struct {
				Message     string `json:"message"`
				For         string `json:"for"`
				Path        string `json:"path"`
				Code        string `json:"code"`
				MessageCode string `json:"messageCode"`
			} `json:"errors,omitempty"`
		} `json:"results"`
	} `json:"winstrom"`
}

type APIResponseContacts struct {
	Winstrom struct {
		Version string `json:"@version"`
		Adresar []struct {
			ID         string `json:"id"`
			LastUpdate string `json:"lastUpdate"`
			Kod        string `json:"kod"`
			Nazev      string `json:"nazev"`
			IC         string `json:"ic"`
			DIC        string `json:"dic"`
			Ulice      string `json:"ulice"`
			Mesto      string `json:"mesto"`
			PSC        string `json:"psc"`
			Stat       string `json:"stat"`
			StatRef    string `json:"stat@ref"`
			StatShowAs string `json:"stat@showAs"`
		} `json:"adresar"`
	} `json:"winstrom"`
}

type InvoiceRequest struct {
	Winstrom struct {
		FakturaVydana []FakturaVydana `json:"faktura-vydana"`
	} `json:"winstrom"`
}

type FakturaVydana struct {
	Kod              string           `json:"kod"`
	DatVyst          string           `json:"datVyst"`
	DatSplat         string           `json:"datSplat"`
	StavUhrady       string           `json:"stavUhrK"`
	Polozky          []FakturaPolozka `json:"polozkyFaktury"`
	IdFirmy          string           `json:"firma"`
	TypFaktury       string           `json:"typDokl"`
	AccountingType   string           `json:"typUcOp"`
	FormaUhradyCislo string           `json:"formaUhradyCis"`
}

type FakturaPolozka struct {
	Popis   string  `json:"nazev"`
	Pocet   float64 `json:"mnozMj"`
	CenaKus float64 `json:"cenaMj"`
}

type SaleReceiptItem struct {
	Nazev  string  `json:"nazev,omitempty"`
	CenaMj float64 `json:"cenaMj,omitempty"`
	MnozMj float64 `json:"mnozMj"`
	Cenik  string  `json:"cenik,omitempty"`
	Kod    string  `json:"kod,omitempty"`
	Sklad  string  `json:"sklad,omitempty"`
}

type SaleProdejka struct {
	Mena           string            `json:"mena"`
	TypDokl        string            `json:"typDokl"`
	TypUcOp        string            `json:"typUcOp"`
	PrimUcet       string            `json:"primUcet"`
	ProtiUcet      string            `json:"protiUcet"`
	FormaUhradyCis string            `json:"formaUhradyCis"`
	PolozkyDokladu []SaleReceiptItem `json:"polozkyDokladu"`
}

type SaleReceipt struct {
	Winstrom struct {
		Version  string       `json:"@version"`
		Prodejka SaleProdejka `json:"prodejka"`
	} `json:"winstrom"`
}

type Cenik struct {
	ID            string `json:"id"`
	Kod           string `json:"kod"`
	Nazev         string `json:"nazev"`
	Popis         string `json:"popis"`
	EanKod        string `json:"eanKod"`
	NakupCena     string `json:"nakupCena"`
	CenJednotka   string `json:"cenJednotka"`
	Skladove      string `json:"skladove"`
	ExportNaEshop string `json:"exportNaEshop"`
	EvidExpir     string `json:"evidExpir"`
	ProdejKasa    string `json:"prodejKasa"`
	SkupZboz      string `json:"skupZboz"`
	Mj1           string `json:"mj1"`
	Dodavatel     string `json:"dodavatel"`
	BaseCode      string `json:"popisA"`
	Kategorie     string `json:"popisB,omitempty"`
	ObrazekURL    string `json:"popisC"`
}

type CenikWrapper struct {
	Winstrom CenikWrapperWinstrom `json:"winstrom"`
}

type CenikWrapperWinstrom struct {
	Version string  `json:"@version"`
	Cenik   []Cenik `json:"cenik"`
}
