package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 定义合约结构体
type SmartContract struct {
	contractapi.Contract
}

// 注册用户
func (s *SmartContract) RegisterUser(ctx contractapi.TransactionContextInterface, userID string, userType string, realInfoHash string) error {
	user := User{
		UserID:       userID,
		UserType:     userType,
		RealInfoHash: realInfoHash,
		ProductList:  []*IndustrialProduct{},
	}
	userAsBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(userID, userAsBytes)
	if err != nil {
		return err
	}
	return nil
}

// 工业产品上链，传入用户ID、工业产品上链信息
// arg7 为可选的压缩证据 JSON（CompressionEvidence），为空时表示未压缩
func (s *SmartContract) Uplink(ctx contractapi.TransactionContextInterface, userID string, traceabilityCode string, arg1 string, arg2 string, arg3 string, arg4 string, arg5 string, arg6 string, arg7 string) (string, error) {
	// 解析压缩证据（可选）
	var compEvidence *CompressionEvidence
	if arg7 != "" {
		compEvidence = &CompressionEvidence{}
		if err := json.Unmarshal([]byte(arg7), compEvidence); err != nil {
			return "", fmt.Errorf("failed to unmarshal compression evidence: %v", err)
		}
	}

	// 获取用户类型
	userType, err := s.GetUserType(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user type: %v", err)
	}

	// 通过溯源码获取工业产品的上链信息
	productAsBytes, err := ctx.GetStub().GetState(traceabilityCode)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	// 将工业产品的信息转换为结构体
	var product IndustrialProduct
	if productAsBytes != nil {
		err = json.Unmarshal(productAsBytes, &product)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal product: %v", err)
		}
	}

	//获取时间戳,修正时区
	txtime, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("failed to read TxTimestamp: %v", err)
	}
	timeLocation, _ := time.LoadLocation("Asia/Shanghai") // 选择你所在的时区
	txTimeStr := time.Unix(txtime.Seconds, 0).In(timeLocation).Format("2006-01-02 15:04:05")

	// 获取交易ID
	txid := ctx.GetStub().GetTxID()
	// 给工业产品信息中加上溯源码
	product.TraceabilityCode = traceabilityCode
	// 不同用户类型的上链的参数不一致
	switch userType {
	// 原料供应商
	case "原料供应商":
		// 原料提供商输入
		product.RawSupplierInput.ProductName = arg1
		product.RawSupplierInput.RawOrigin = arg2
		product.RawSupplierInput.ArrivalTime = arg3
		product.RawSupplierInput.ProductionTime = arg4
		product.RawSupplierInput.SupplierName = arg5
		product.RawSupplierInput.Img = arg6
		product.RawSupplierInput.Txid = txid
		product.RawSupplierInput.Timestamp = txTimeStr
		product.RawSupplierInput.CompressionEvidence = compEvidence
	// 制造商
	case "制造商":
		product.ManufacturerInput.ProductName = arg1
		product.ManufacturerInput.ProductionBatch = arg2
		product.ManufacturerInput.FactoryTime = arg3
		product.ManufacturerInput.FactoryNameAddress = arg4
		product.ManufacturerInput.ContactPhone = arg5
		product.ManufacturerInput.Img = arg6
		product.ManufacturerInput.Txid = txid
		product.ManufacturerInput.Timestamp = txTimeStr
		product.ManufacturerInput.CompressionEvidence = compEvidence
	// 物流承运商
	case "物流承运商":
		product.CarrierInput.Name = arg1
		product.CarrierInput.Age = arg2
		product.CarrierInput.Phone = arg3
		product.CarrierInput.PlateNumber = arg4
		product.CarrierInput.TransportRecord = arg5
		product.CarrierInput.Img = arg6
		product.CarrierInput.Txid = txid
		product.CarrierInput.Timestamp = txTimeStr
		product.CarrierInput.CompressionEvidence = compEvidence
	// 经销商
	case "经销商":
		product.DealerInput.StoreTime = arg1
		product.DealerInput.SellTime = arg2
		product.DealerInput.DealerName = arg3
		product.DealerInput.DealerLocation = arg4
		product.DealerInput.DealerPhone = arg5
		product.DealerInput.Img = arg6
		product.DealerInput.Txid = txid
		product.DealerInput.Timestamp = txTimeStr
		product.DealerInput.CompressionEvidence = compEvidence

	}

	//将工业产品的信息转换为json格式
	productAsBytes, err = json.Marshal(product)
	if err != nil {
		return "", fmt.Errorf("failed to marshal product: %v", err)
	}
	//将工业产品的信息存入区块链
	err = ctx.GetStub().PutState(traceabilityCode, productAsBytes)
	if err != nil {
		return "", fmt.Errorf("failed to put product: %v", err)
	}
	//将工业产品存入用户的工业产品列表
	err = s.AddIndustrialProduct(ctx, userID, &product)
	if err != nil {
		return "", fmt.Errorf("failed to add product to user: %v", err)

	}

	return txid, nil
}

// 添加工业产品到用户的列表
func (s *SmartContract) AddIndustrialProduct(ctx contractapi.TransactionContextInterface, userID string, product *IndustrialProduct) error {
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if userBytes == nil {
		return fmt.Errorf("the user %s does not exist", userID)
	}
	// 将返回的结果转换为User结构体
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return err
	}
	// 遍历用户的工业产品列表，检查是否已经存在该工业产品
	for _, existingProduct := range user.ProductList {
		if existingProduct.TraceabilityCode == product.TraceabilityCode {
			return fmt.Errorf("the product with traceability code %s already exists in user %s's product list", product.TraceabilityCode, userID)
		}
	}
	// 如果不存在，则将工业产品添加到用户的工业产品列表中
	user.ProductList = append(user.ProductList, product)
	userAsBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(userID, userAsBytes)
	if err != nil {
		return err
	}
	return nil
}

// 获取用户类型
func (s *SmartContract) GetUserType(ctx contractapi.TransactionContextInterface, userID string) (string, error) {
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if userBytes == nil {
		return "", fmt.Errorf("the user %s does not exist", userID)
	}
	// 将返回的结果转换为User结构体
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return "", err
	}
	return user.UserType, nil
}

// 获取用户信息
func (s *SmartContract) GetUserInfo(ctx contractapi.TransactionContextInterface, userID string) (*User, error) {
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return &User{}, fmt.Errorf("failed to read from world state: %v", err)
	}
	if userBytes == nil {
		return &User{}, fmt.Errorf("the user %s does not exist", userID)
	}
	// 将返回的结果转换为User结构体
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return &User{}, err
	}
	return &user, nil
}

// 获取工业产品的上链信息
func (s *SmartContract) GetIndustrialProductInfo(ctx contractapi.TransactionContextInterface, traceabilityCode string) (*IndustrialProduct, error) {
	productAsBytes, err := ctx.GetStub().GetState(traceabilityCode)
	if err != nil {
		return &IndustrialProduct{}, fmt.Errorf("failed to read from world state: %v", err)
	}

	// 将返回的结果转换为IndustrialProduct结构体
	var product IndustrialProduct
	if productAsBytes != nil {
		err = json.Unmarshal(productAsBytes, &product)
		if err != nil {
			return &IndustrialProduct{}, fmt.Errorf("failed to unmarshal product: %v", err)
		}
		if product.TraceabilityCode != "" {
			return &product, nil
		}
	}
	return &IndustrialProduct{}, fmt.Errorf("the industrial product %s does not exist", traceabilityCode)
}

// 获取用户的工业产品列表
func (s *SmartContract) GetIndustrialProductList(ctx contractapi.TransactionContextInterface, userID string) ([]*IndustrialProduct, error) {
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if userBytes == nil {
		return nil, fmt.Errorf("the user %s does not exist", userID)
	}
	// 将返回的结果转换为User结构体
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, err
	}
	return user.ProductList, nil
}

// 获取所有的工业产品信息
func (s *SmartContract) GetAllIndustrialProductInfo(ctx contractapi.TransactionContextInterface) ([]IndustrialProduct, error) {
	productListAsBytes, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	defer func() { _ = productListAsBytes.Close() }()
	var products []IndustrialProduct
	for productListAsBytes.HasNext() {
		queryResponse, err := productListAsBytes.Next()
		if err != nil {
			return nil, err
		}
		var product IndustrialProduct
		err = json.Unmarshal(queryResponse.Value, &product)
		if err != nil {
			return nil, err
		}
		// 过滤非工业产品的信息
		if product.TraceabilityCode != "" {
			products = append(products, product)
		}
	}
	return products, nil
}

// 获取工业产品上链历史
func (s *SmartContract) GetIndustrialProductHistory(ctx contractapi.TransactionContextInterface, traceabilityCode string) ([]HistoryQueryResult, error) {
	log.Printf("GetAssetHistory: ID %v", traceabilityCode)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(traceabilityCode)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resultsIterator.Close() }()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var product IndustrialProduct
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &product)
			if err != nil {
				return nil, err
			}
		} else {
			product = IndustrialProduct{
				TraceabilityCode: traceabilityCode,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}
		// 指定目标时区
		targetLocation, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			return nil, err
		}

		// 将时间戳转换到目标时区
		timestamp = timestamp.In(targetLocation)
		// 格式化时间戳为指定格式的字符串
		formattedTime := timestamp.Format("2006-01-02 15:04:05")

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: formattedTime,
			Record:    &product,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}

const maxFileSizeBytes int64 = 5 * 1024 * 1024 * 1024

// PutFileManifest writes off-chain file metadata onto the ledger with role/size checks.
func (s *SmartContract) PutFileManifest(ctx contractapi.TransactionContextInterface, userID string, manifestJSON string) (string, error) {
	var manifest FileManifest
	if err := json.Unmarshal([]byte(manifestJSON), &manifest); err != nil {
		return "", fmt.Errorf("invalid manifest json: %v", err)
	}
	if manifest.TraceabilityCode == "" || manifest.FileID == "" || manifest.CID == "" || manifest.Hash == "" {
		return "", errors.New("traceabilityCode, fileID, cid, hash are required")
	}
	if manifest.Size <= 0 || manifest.Size > maxFileSizeBytes {
		return "", fmt.Errorf("file size invalid or exceeds limit %d bytes", maxFileSizeBytes)
	}
	userType, err := s.GetUserType(ctx, userID)
	if err != nil {
		return "", err
	}
	role, err := mapUserTypeToFileRole(userType)
	if err != nil {
		return "", err
	}
	if manifest.Role != "" && manifest.Role != role {
		return "", fmt.Errorf("role mismatch: manifest role %s not allowed for user type %s", manifest.Role, userType)
	}
	manifest.Role = role
	if manifest.SourceHash == "" {
		manifest.SourceHash = manifest.Hash
	}
	if manifest.Encrypted && manifest.KeyVersion == "" {
		return "", errors.New("keyVersion required when encrypted is true")
	}

	// ensure product exists
	_, err = s.GetIndustrialProductInfo(ctx, manifest.TraceabilityCode)
	if err != nil {
		return "", fmt.Errorf("traceability code not found: %v", err)
	}

	fileKey := "file:" + manifest.FileID
	existing, err := ctx.GetStub().GetState(fileKey)
	if err != nil {
		return "", fmt.Errorf("failed to read manifest: %v", err)
	}
	if existing != nil {
		return "", fmt.Errorf("file manifest %s already exists", manifest.FileID)
	}

	// stamp uploader and timestamp from transaction context
	manifest.Uploader = userID
	txtime, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("failed to get tx timestamp: %v", err)
	}
	timeLocation, _ := time.LoadLocation("Asia/Shanghai")
	manifest.Timestamp = time.Unix(txtime.Seconds, 0).In(timeLocation).Format("2006-01-02 15:04:05")

	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal manifest: %v", err)
	}

	if err := ctx.GetStub().PutState(fileKey, manifestBytes); err != nil {
		return "", fmt.Errorf("failed to put manifest: %v", err)
	}

	// create composite key for traceabilityCode -> fileID index
	indexKey, err := ctx.GetStub().CreateCompositeKey("trace~file", []string{manifest.TraceabilityCode, manifest.FileID})
	if err != nil {
		return "", fmt.Errorf("failed to create index key: %v", err)
	}
	if err := ctx.GetStub().PutState(indexKey, []byte{0}); err != nil {
		return "", fmt.Errorf("failed to put index: %v", err)
	}

	return manifest.FileID, nil
}

// GetFileManifest returns manifest by fileID.
func (s *SmartContract) GetFileManifest(ctx contractapi.TransactionContextInterface, fileID string) (*FileManifest, error) {
	fileKey := "file:" + fileID
	data, err := ctx.GetStub().GetState(fileKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %v", err)
	}
	if data == nil {
		return nil, fmt.Errorf("file manifest %s not found", fileID)
	}
	var manifest FileManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to unmarshal manifest: %v", err)
	}
	return &manifest, nil
}

// GetFileManifestsByTrace lists manifests associated with a traceability code.
func (s *SmartContract) GetFileManifestsByTrace(ctx contractapi.TransactionContextInterface, traceabilityCode string) ([]FileManifest, error) {
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey("trace~file", []string{traceabilityCode})
	if err != nil {
		return nil, fmt.Errorf("failed to query index: %v", err)
	}
	defer func() { _ = iter.Close() }()

	var manifests []FileManifest
	for iter.HasNext() {
		resp, err := iter.Next()
		if err != nil {
			return nil, err
		}
		_, parts, err := ctx.GetStub().SplitCompositeKey(resp.Key)
		if err != nil || len(parts) != 2 {
			continue
		}
		fileID := parts[1]
		m, err := s.GetFileManifest(ctx, fileID)
		if err != nil {
			return nil, err
		}
		manifests = append(manifests, *m)
	}
	return manifests, nil
}

func mapUserTypeToFileRole(userType string) (FileRole, error) {
	switch userType {
	case "原料供应商":
		return FileRoleRawSupplier, nil
	case "制造商":
		return FileRoleManufacturer, nil
	case "物流承运商":
		return FileRoleCarrier, nil
	case "经销商":
		return FileRoleDealer, nil
	default:
		return "", fmt.Errorf("unsupported user type %s", userType)
	}
}
