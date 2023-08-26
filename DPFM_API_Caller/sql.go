package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-barcode-generates/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-barcode-generates/DPFM_API_Output_Formatter"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	"golang.org/x/xerrors"
)

func (c *DPFMAPICaller) createSqlProcess(
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	product *dpfm_api_output_formatter.Product,
	errs *[]error,
	log *logger.Logger,
) interface{} {
	c.barcodeCreateSql(nil, input, output, product, errs, log)
	response := dpfm_api_output_formatter.ConvertToProduct(product)

	data := dpfm_api_output_formatter.Message{
		Product: dpfm_api_output_formatter.Product{
			Product:           response.Product,
			ProductStandardID: response.ProductStandardID,
			BarcodeType:       response.BarcodeType,
		},
	}

	return data
}

func (c *DPFMAPICaller) barcodeCreateSql(
	ctx context.Context,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	product *dpfm_api_output_formatter.Product,
	errs *[]error,
	log *logger.Logger,
) {
	if ctx == nil {
		ctx = context.Background()
	}
	sessionID := input.RuntimeSessionID
	res, err := c.rmq.SessionKeepRequest(nil, c.conf.RMQ.QueueToSQL()[0], map[string]interface{}{"message": product, "function": "BarcodeGenerate", "runtime_session_id": sessionID})
	if err != nil {
		err = xerrors.Errorf("rmq error: %w", err)
		*errs = append(*errs, err)
		return
	}
	res.Success()

	return
}
