package dpfm_api_output_formatter

import "data-platform-api-barcode-generates/DPFM_API_Caller/requests"

func ConvertToProduct(product *Product) *requests.Product {
	pm := &requests.Product{}

	pm = &requests.Product{
		Product:           product.Product,
		ProductStandardID: product.ProductStandardID,
		BarcodeType:       product.BarcodeType,
	}

	return pm
}
