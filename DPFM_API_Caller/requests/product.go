package requests

type Product struct {
	Product           string `json:"Product"`
	ProductStandardID string `json:"ProductStandardID"`
	BarcodeType       string `json:"BarcodeType"`
}
