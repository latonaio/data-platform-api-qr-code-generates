package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-barcode-generates/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-barcode-generates/DPFM_API_Output_Formatter"
	"data-platform-api-barcode-generates/config"
	"errors"
	"fmt"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/code39"
	"image/png"
	"os"
	"strconv"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/ean"
	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client-for-data-platform"
)

const (
	TypeEAN13   = "EAN"
	TypeCode39  = "Code39"
	TypeCode128 = "Code128"
)

const FileExtension = "png"

type DPFMAPICaller struct {
	ctx  context.Context
	conf *config.Conf
	rmq  *rabbitmq.RabbitmqClient
}

type GenerateFileInfo struct {
	Width         int
	Height        int
	MountPath     string
	FileName      string
	FileExtension string
	FilePath      string
	BarcodeType   string
}

func NewDPFMAPICaller(
	conf *config.Conf, rmq *rabbitmq.RabbitmqClient,
) *DPFMAPICaller {
	return &DPFMAPICaller{
		ctx:  context.Background(),
		conf: conf,
		rmq:  rmq,
	}
}

func (c *DPFMAPICaller) AsyncBarcodeGenerates(
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	log *logger.Logger,
	errs *[]error,
	conf *config.Conf,
) (interface{}, *[]error) {
	product := input.Product.Product
	productStandardID := input.Product.ProductStandardID
	barcodeType := input.Product.BarcodeType

	imageInfo := conf.Image.ImageInfo()

	err := barcodeGenerates(productStandardID, barcodeType, imageInfo, conf)
	if err != nil {
		*errs = append(*errs, err)
		return nil, errs
	}

	response := c.createSqlProcess(input, output, &dpfm_api_output_formatter.Product{
		Product:           product,
		ProductStandardID: productStandardID,
		BarcodeType:       barcodeType,
	}, errs, log)

	return response, nil
}

func barcodeGenerates(
	productStandardID string,
	barcodeType string,
	imageInfo *config.Image,
	config *config.Conf,
) error {
	width, err := strconv.Atoi(imageInfo.Width)
	if err != nil {
		return err
	}

	height, err := strconv.Atoi(imageInfo.Height)
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/%s/%s.%s",
		config.MountPath,
		barcodeType,
		productStandardID,
		FileExtension,
	)

	generateFileInfo := GenerateFileInfo{
		Width:       width,
		Height:      height,
		MountPath:   config.MountPath,
		FileName:    productStandardID,
		FilePath:    filePath,
		BarcodeType: barcodeType,
	}

	encodeData, err := func(barcodeType string, productStandardID string) (barcode.BarcodeIntCS, error) {
		if barcodeType == TypeEAN13 {
			return ean.Encode(productStandardID)
		}

		if barcodeType == TypeCode39 {
			return code39.Encode(productStandardID, true, true)
		}

		if barcodeType == TypeCode128 {
			return code128.Encode(productStandardID)
		}

		return nil, errors.New(fmt.Sprintf("Error: %s", "Invalid barcode type"))
	}(barcodeType, productStandardID)
	if err != nil {
		return err
	}

	file, err := generateBarcode(encodeData, generateFileInfo)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	return err
}

func generateBarcode(encodeData barcode.BarcodeIntCS, generateFileInfo GenerateFileInfo) (*os.File, error) {
	directoryPath := fmt.Sprintf("%s/%s", generateFileInfo.MountPath, generateFileInfo.BarcodeType)

	err := os.MkdirAll(directoryPath, 0777)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(generateFileInfo.FilePath)
	if err != nil {
		return nil, err
	}

	barcodeScale, err := barcode.Scale(encodeData, generateFileInfo.Width, generateFileInfo.Height)

	err = png.Encode(file, barcodeScale)
	if err != nil {
		return nil, err
	}

	return file, nil
}
