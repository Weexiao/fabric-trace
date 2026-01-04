package chaincode

/*
定义用户结构体
用户ID
用户类型
实名认证信息哈希,包括用户注册的姓名、身份证号、手机号、注册平台同意协议签名的哈希
工业产品列表
*/
type User struct {
	UserID       string               `json:"userID"`
	UserType     string               `json:"userType"`
	RealInfoHash string               `json:"realInfoHash"`
	ProductList  []*IndustrialProduct `json:"productList"`
}

/*
定义工业产品结构体
溯源码
原料供应商输入
制造商输入
物流承运商输入
经销商输入
*/
type IndustrialProduct struct {
	TraceabilityCode  string            `json:"traceabilityCode"`
	RawSupplierInput  RawSupplierInput  `json:"rawSupplierInput"`
	ManufacturerInput ManufacturerInput `json:"manufacturerInput"`
	CarrierInput      CarrierInput      `json:"carrierInput"`
	DealerInput       DealerInput       `json:"dealerInput"`
}

// HistoryQueryResult structure used for handling result of history query
type HistoryQueryResult struct {
	Record    *IndustrialProduct `json:"record"`
	TxId      string             `json:"txId"`
	Timestamp string             `json:"timestamp"`
	IsDelete  bool               `json:"isDelete"`
}

/*
原料供应商（原料提供商）
工业产品的溯源码，一物一码，主打高端市场（自动生成）
产品名称
原料产地
原料到货时间
原料生产时间
原料供应商名称
*/
type RawSupplierInput struct {
	ProductName    string `json:"productName"`
	RawOrigin      string `json:"rawOrigin"`
	ArrivalTime    string `json:"arrivalTime"`
	ProductionTime string `json:"productionTime"`
	SupplierName   string `json:"supplierName"`
	Img            string `json:"img"`
	Txid           string `json:"txid"`
	Timestamp      string `json:"timestamp"`
}

/*
制造商输入
商品名称
生产批次
出厂时间（可以防止黑心商家修改时间）
制造商名称与厂址
联系电话
*/
type ManufacturerInput struct {
	ProductName        string `json:"productName"`
	ProductionBatch    string `json:"productionBatch"`
	FactoryTime        string `json:"factoryTime"`
	FactoryNameAddress string `json:"factoryNameAddress"`
	ContactPhone       string `json:"contactPhone"`
	Img                string `json:"img"`
	Txid               string `json:"txid"`
	Timestamp          string `json:"timestamp"`
}

/*
物流承运商输入
姓名
年龄
电话
车牌号
运输记录
*/
type CarrierInput struct {
	Name            string `json:"name"`
	Age             string `json:"age"`
	Phone           string `json:"phone"`
	PlateNumber     string `json:"plateNumber"`
	TransportRecord string `json:"transportRecord"`
	Img             string `json:"img"`
	Txid            string `json:"txid"`
	Timestamp       string `json:"timestamp"`
}

/*
经销商输入
存入时间
销售时间
经销商名称
经销商位置
经销商电话
*/
type DealerInput struct {
	StoreTime      string `json:"storeTime"`
	SellTime       string `json:"sellTime"`
	DealerName     string `json:"dealerName"`
	DealerLocation string `json:"dealerLocation"`
	DealerPhone    string `json:"dealerPhone"`
	Img            string `json:"img"`
	Txid           string `json:"txid"`
	Timestamp      string `json:"timestamp"`
}

// FileRole defines who uploaded/owns a file for role-based permissions.
type FileRole string

const (
	FileRoleRawSupplier  FileRole = "raw_supplier"
	FileRoleManufacturer FileRole = "manufacturer"
	FileRoleCarrier      FileRole = "carrier"
	FileRoleDealer       FileRole = "dealer"
)

// FileManifest captures off-chain file metadata stored on-chain.
type FileManifest struct {
	TraceabilityCode string   `json:"traceabilityCode"`
	FileID           string   `json:"fileID"`
	CID              string   `json:"cid"`
	Hash             string   `json:"hash"`
	Mime             string   `json:"mime"`
	Size             int64    `json:"size"`
	Encrypted        bool     `json:"encrypted"`
	KeyVersion       string   `json:"keyVersion"`
	Role             FileRole `json:"role"`
	Uploader         string   `json:"uploader"`
	Timestamp        string   `json:"timestamp"`
}
