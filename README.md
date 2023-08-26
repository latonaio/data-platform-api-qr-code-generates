# data-platform-api-barcode-generates

data-platform-api-barcode-generates は、データ連携基盤において、API でバーコードを生成するマイクロサービスです。  
https://xxx.xxx.io/api/API_BARCODE_SRV/generates/

## 動作環境

data-platform-api-barcode-generates の動作環境は、次の通りです。  
・ OS: LinuxOS （必須）  
・ CPU: ARM/AMD/Intel（いずれか必須）  


## 本レポジトリ が 対応する API サービス
data-platform-api-barcode-generates が対応する APIサービス は、次のものです。

APIサービス URL: https://xxx.xxx.io/api/API_BARCODE_SRV/creates/

## 本レポジトリ に 含まれる API名
data-platform-api-barcode-generates には、次の API をコールするためのリソースが含まれています。  

* A_Header（オーダー - ヘッダデータ）
* A_HeaderPartner（オーダー - ヘッダ取引先データ）
* A_HeaderPartnerPlant（オーダー - ヘッダ取引先プラントデータ）
* A_HeaderPartnerContact（オーダー - ヘッダ取引先コンタクトデータ）
* A_Item（オーダー - 明細データ）
* A_ItemPartner（オーダー - 明細取引先データ）
* A_ItemPartnerPlant（オーダー - 明細取引先プラントデータ）
* A_ItemPricingElement（オーダー - 明細取引先プラントデータ）
* A_ItemScheduleLine（オーダー - 明細納入日程行データ）
* A_Address（オーダー - 住所データ）

## API への 値入力条件 の 初期値
data-platform-api-barcode-generates において、API への値入力条件の初期値は、入力ファイルレイアウトの種別毎に、次の通りとなっています。  

## データ連携基盤のAPIの選択的コール

Latona および AION の データ連携基盤 関連リソースでは、Inputs フォルダ下の sample.json の accepter に取得したいデータの種別（＝APIの種別）を入力し、指定することができます。  
なお、同 accepter にAll(もしくは空白)の値を入力することで、全データ（＝全APIの種別）をまとめて取得することができます。  

* sample.jsonの記載例(1)  

accepter において 下記の例のように、データの種別（＝APIの種別）を指定します。  
ここでは、"Header" が指定されています。    
  
```
	"api_schema": "DPFMOrdersCreates",
	"accepter": ["Header"],
	"order_id": null,
	"deleted": false
```
  
* 全データを取得する際のsample.jsonの記載例(2)  

全データを取得する場合、sample.json は以下のように記載します。  

```
	"api_schema": "DPFMOrdersCreates",
	"accepter": ["All"],
	"order_id": null,
	"deleted": false
```

## 指定されたデータ種別のコール

accepter における データ種別 の指定に基づいて DPFM_API_Caller 内の caller.go で API がコールされます。  
caller.go の func() 毎 の 以下の箇所が、指定された API をコールするソースコードです。  

```
func (c *DPFMAPICaller) AsyncOrderCreates(
	accepter []string,
	input *dpfm_api_input_reader.SDC,
	output *sub_func_complementer.SDC,
	log *logger.Logger,
    
) []error {
	wg := sync.WaitGroup{}
	mtx := sync.Mutex{}
	errs := make([]error, 0, 5)
	exconfAllExist := false

	subFuncFin := make(chan error)
	exconfFin := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()
		var e []error
		exconfAllExist, e = c.configure.Conf(input, output, log)
		if len(e) != 0 {
			mtx.Lock()
			errs = append(errs, e...)
			mtx.Unlock()
			exconfFin <- xerrors.New("exconf error")
			return
		}
		exconfFin <- nil
	}()

	for _, fn := range accepter {
		wg.Add(1)
		switch fn {
		case "Header":
			go c.headerCreate(&wg, &mtx, subFuncFin, log, &errs, input, output)
		case "Item":
			errs = append(errs, xerrors.New("accepter Item is not implement yet"))
		default:
			wg.Done()
		}
	}
```

## Output  
本マイクロサービスでは、[golang-logging-library-for-data-platform](https://github.com/latonaio/golang-logging-library-for-data-platform) により、以下のようなデータがJSON形式で出力されます。  
以下の sample.json の例は オーダー の ヘッダデータ が取得された結果の JSON の例です。  
以下の項目のうち、"OrderID" ～ "PlusMinusFlag" は、/DPFM_API_Output_Formatter/type.go 内 の Type Header {} による出力結果です。"cursor" ～ "time"は、golang-logging-library による 定型フォーマットの出力結果です。  

```
{
  "connection_key": "request",
  "result": true,
  "redis_key": "abcdefg",
  "filepath": "/var/lib/aion/Data/rededge_sdc/abcdef.json",
  "api_status_code": 200,
  "runtime_session_id": "786c60bfab5c480abcb2bd53a738e9ca",
  "business_partner": 101,
  "service_label": "BARCODE",
  "api_type": "generates",
  "message": {
    "Product": {
      "Product": "A3750",
      "ProductStandardID": "1234567890128",
      "BarcodeType": "EAN"
    }
  },
  "api_schema": "Barcode",
  "accepter": [
    "All"
  ],
  "deleted": false,
  "api_processing_result": true,
  "api_processing_error": ""
}
```
