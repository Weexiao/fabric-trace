package controller

import (
	"backend/pkg"
	"backend/pkg/storage"
	"backend/service"
	"backend/settings"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 保存了工业产品上链与查询的函数

func Uplink(c *gin.Context) {
	// 与userID不一样，取ID从第二位作为溯源码，即18位，雪花ID的结构如下（转化为10进制共19位）：
	// +--------------------------------------------------------------------------+
	// | 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID |
	// +--------------------------------------------------------------------------+
	traceCode := pkg.GenerateID()[1:]
	args := buildArgs(c, traceCode)
	if args == nil {
		return
	}
	res, err := pkg.ChaincodeInvoke("Uplink", args)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "uplink failed" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":             200,
		"message":          "uplink success",
		"txid":             res,
		"traceabilityCode": args[1],
	})
}

// 获取工业产品的上链信息
func GetIndustrialProductInfo(c *gin.Context) {
	traceCode := readTraceabilityCode(c)
	// DEBUG: print what we actually received/normalized
	fmt.Printf("[GetIndustrialProductInfo] remote=%s content-type=%s rawForm=%q rawQuery=%q normalizedTraceCode=%q\n",
		c.ClientIP(),
		c.GetHeader("Content-Type"),
		c.PostForm("traceabilityCode"),
		c.Query("traceabilityCode"),
		traceCode,
	)
	fmt.Printf("[GetIndustrialProductInfo] traceCode=%q len=%d\n", traceCode, len(traceCode))

	if traceCode == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请检查溯源码是否正确!!",
		})
		return
	}

	res, err := pkg.ChaincodeQuery("GetIndustrialProductInfo", traceCode)
	if err != nil {
		msg := err.Error()
		fmt.Printf("[GetIndustrialProductInfo] chaincode query failed traceCode=%q err=%s\n", traceCode, msg)
		if strings.Contains(msg, "does not exist") {
			c.JSON(http.StatusOK, gin.H{
				"code":    404,
				"message": "该溯源码不存在，请确认后重试",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "查询失败：" + msg,
		})
		return
	}

	// Enrich product tree with off-chain image hash (derived from stored img filename)
	productEnriched := enrichProductWithImgHash(res)

	// Query file manifests and aggregate source/compressed hashes for frontend trace view.
	manifestsRaw, mErr := pkg.ChaincodeQuery("GetFileManifestsByTrace", traceCode)
	fileHashesByRole := map[string][]string{
		"raw_supplier": {},
		"manufacturer": {},
		"carrier":      {},
		"dealer":       {},
	}
	fileHashEntriesByRole := map[string][]gin.H{
		"raw_supplier": {},
		"manufacturer": {},
		"carrier":      {},
		"dealer":       {},
	}

	if mErr == nil && strings.TrimSpace(manifestsRaw) != "" {
		trim := strings.TrimSpace(manifestsRaw)
		var manifests []storage.Manifest
		if err := json.Unmarshal([]byte(trim), &manifests); err == nil {
			seenLegacy := map[string]map[string]bool{}
			seenEntries := map[string]map[string]bool{}
			for k := range fileHashesByRole {
				seenLegacy[k] = map[string]bool{}
				seenEntries[k] = map[string]bool{}
			}

			for _, m := range manifests {
				role := strings.TrimSpace(m.Role)
				if _, exists := fileHashesByRole[role]; !exists {
					continue
				}
				sourceHash := strings.TrimSpace(m.SourceHash)
				if sourceHash == "" {
					sourceHash = strings.TrimSpace(m.Hash)
				}
				compressedHash := strings.TrimSpace(m.CompressedHash)
				compressedBits := m.CompressedBits
				if len(compressedBits) == 0 && compressedHash != "" {
					compressedBits = hashHexToBits01(compressedHash)
				}
				if sourceHash == "" {
					continue
				}

				if !seenLegacy[role][sourceHash] {
					seenLegacy[role][sourceHash] = true
					fileHashesByRole[role] = append(fileHashesByRole[role], sourceHash)
				}

				pairKey := sourceHash + "|" + compressedHash
				if seenEntries[role][pairKey] {
					continue
				}
				seenEntries[role][pairKey] = true
				fileHashEntriesByRole[role] = append(fileHashEntriesByRole[role], gin.H{
					"sourceHash":     sourceHash,
					"compressedHash": compressedHash,
					"compressedBits": compressedBits,
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "query success",
		"data": gin.H{
			"product":               productEnriched,
			"fileHashesByRole":      fileHashesByRole,
			"fileHashEntriesByRole": fileHashEntriesByRole,
		},
	})
}

func hashHexToBits01(hashHex string) []int {
	if strings.TrimSpace(hashHex) == "" {
		return nil
	}
	b, err := hex.DecodeString(strings.TrimSpace(hashHex))
	if err != nil || len(b) == 0 {
		return nil
	}
	bits := make([]int, 0, len(b)*8)
	for _, by := range b {
		for i := 7; i >= 0; i-- {
			if (by>>uint(i))&1 == 1 {
				bits = append(bits, 1)
			} else {
				bits = append(bits, 0)
			}
		}
	}
	return bits
}

// enrichProductWithImgHash adds `imgHash` fields (sha256 hex) under each role input if an img filename is present.
// It does NOT modify chain data; this is only for query display.
func enrichProductWithImgHash(productJSON string) any {
	trim := strings.TrimSpace(productJSON)
	if trim == "" {
		return productJSON
	}
	if !(strings.HasPrefix(trim, "{") || strings.HasPrefix(trim, "[")) {
		return productJSON
	}

	var obj map[string]any
	if err := json.Unmarshal([]byte(trim), &obj); err != nil {
		return productJSON
	}

	setHash := func(section string) {
		secAny, ok := obj[section]
		if !ok {
			return
		}
		sec, ok := secAny.(map[string]any)
		if !ok {
			return
		}
		imgAny, ok := sec["img"]
		if !ok {
			return
		}
		img, _ := imgAny.(string)
		img = strings.TrimSpace(img)
		if img == "" {
			return
		}
		base := filepath.Base(img)
		ext := filepath.Ext(base)
		hash := strings.TrimSuffix(base, ext)
		if len(hash) == 64 {
			if _, err := hex.DecodeString(hash); err == nil {
				sec["imgHash"] = hash
				return
			}
		}
		sec["imgHash"] = hash
	}

	setHash("rawSupplierInput")
	setHash("manufacturerInput")
	setHash("carrierInput")
	setHash("dealerInput")

	return obj
}

// 获取所有的工业产品信息
func GetAllIndustrialProductInfo(c *gin.Context) {
	res, err := pkg.ChaincodeQuery("GetAllIndustrialProductInfo", "")
	if err != nil {
		c.JSON(200, gin.H{
			"message": "query failed" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}

// 获取工业产品上链历史
func GetIndustrialProductHistory(c *gin.Context) {
	traceCode := readTraceabilityCode(c)
	if traceCode == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请检查溯源码是否正确!!",
		})
		return
	}
	res, err := pkg.ChaincodeQuery("GetIndustrialProductHistory", traceCode)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "does not exist") {
			c.JSON(http.StatusOK, gin.H{
				"code":    404,
				"message": "该溯源码不存在，请确认后重试",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "query failed" + msg,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}

// 获取用户的工业产品ID列表
func GetIndustrialProductList(c *gin.Context) {
	userID, _ := c.Get("userID")
	res, err := pkg.ChaincodeQuery("GetIndustrialProductList", fmt.Sprint(userID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "query failed: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}

func buildArgs(c *gin.Context, traceCode string) []string {
	var args []string
	userID, _ := c.Get("userID")
	userType, _ := pkg.ChaincodeQuery("GetUserType", userID.(string))
	args = append(args, userID.(string))
	if userType == "原料供应商" {
		args = append(args, traceCode)
	} else {
		formTraceCode := readTraceabilityCode(c)
		res, err := pkg.ChaincodeQuery("GetIndustrialProductInfo", formTraceCode)
		if res == "" || err != nil || formTraceCode == "" {
			c.JSON(200, gin.H{
				"code":    400,
				"message": "请检查溯源码是否正确!!",
			})
			return nil
		}
		args = append(args, formTraceCode)
	}

	arg1 := c.PostForm("arg1")
	arg2 := c.PostForm("arg2")
	arg3 := c.PostForm("arg3")
	arg4 := c.PostForm("arg4")
	arg5 := c.PostForm("arg5")
	args = append(args, arg1, arg2, arg3, arg4, arg5)

	// arg6: 图片文件名
	file, _ := c.FormFile("file")
	if file != nil {
		if err := c.SaveUploadedFile(file, "files/uploadfiles/"+file.Filename); err != nil {
			c.JSON(500, gin.H{
				"message": "upload failed: " + err.Error(),
			})
			return nil
		}
		fileSHA256, _ := pkg.CalculateFileSHA256("files/uploadfiles/" + file.Filename)
		ext := pkg.GetFileExt(file.Filename)
		if err := c.SaveUploadedFile(file, fmt.Sprintf("files/images/%s.%s", fileSHA256, ext)); err != nil {
			c.JSON(500, gin.H{
				"message": "save image failed: " + err.Error(),
			})
			return nil
		}
		_ = os.Remove("files/uploadfiles/" + file.Filename)
		args = append(args, fmt.Sprintf("%s.%s", fileSHA256, ext))
	} else {
		args = append(args, "") // arg6 为空
	}

	// arg7: 压缩证据 JSON（如果启用压缩，则对业务字段做压缩并生成 evidence）
	evidenceJSON := ""
	if settings.Cfg.Compression.Enabled {
		compressor := service.NewCompressor(settings.Cfg.Compression.Algorithm)
		payload := fmt.Sprintf("%s|%s|%s|%s|%s", arg1, arg2, arg3, arg4, arg5)
		_, evidence, err := compressor.Compress([]byte(payload))
		if err != nil {
			log.Printf("[buildArgs] compression failed, skipping evidence: %v", err)
		} else {
			eb, _ := json.Marshal(evidence)
			evidenceJSON = string(eb)
		}
	}
	args = append(args, evidenceJSON)

	return args
}

func GetImg(c *gin.Context) {
	filename := c.Param("filename")
	filePath := fmt.Sprintf("files/images/%s", filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(404, gin.H{
			"message": "file not found",
		})
		return
	}
	c.File(filePath)
}

// UplinkCompressed 接受前端压缩后的 JSON 负载进行上链。
// 请求体格式: { "compressedPayload": "<base64>", "traceabilityCode": "...(可选)" }
// 前端将业务字段 JSON 经 Gzip 压缩 + Base64 编码后发送，后端解压还原后走正常上链流程。
// 压缩证据（CompressionEvidence）自动生成并作为 arg7 传入链码。
func UplinkCompressed(c *gin.Context) {
	var req struct {
		CompressedPayload string `json:"compressedPayload" binding:"required"`
		TraceabilityCode  string `json:"traceabilityCode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request: " + err.Error(),
		})
		return
	}

	// 解压还原
	compressor := service.NewCompressor(settings.Cfg.Compression.Algorithm)
	originalData, err := compressor.Decompress(req.CompressedPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "decompress failed: " + err.Error(),
		})
		return
	}

	// 解析业务字段
	var fields struct {
		Arg1 string `json:"arg1"`
		Arg2 string `json:"arg2"`
		Arg3 string `json:"arg3"`
		Arg4 string `json:"arg4"`
		Arg5 string `json:"arg5"`
	}
	if err := json.Unmarshal(originalData, &fields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid decompressed payload: " + err.Error(),
		})
		return
	}

	// 生成压缩证据
	compressedBytes, _ := pkg.Base64Decode(req.CompressedPayload)
	evidence := pkg.BuildCompressionEvidence(compressor.Algorithm(), originalData, compressedBytes)
	evidenceBytes, _ := json.Marshal(evidence)

	// 构建链码参数
	userID, _ := c.Get("userID")
	userIDStr := fmt.Sprint(userID)
	userType, _ := pkg.ChaincodeQuery("GetUserType", userIDStr)

	var args []string
	args = append(args, userIDStr)

	if userType == "原料供应商" {
		traceCode := pkg.GenerateID()[1:]
		args = append(args, traceCode)
	} else {
		if req.TraceabilityCode == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"message": "非原料供应商必须提供溯源码",
			})
			return
		}
		res, err := pkg.ChaincodeQuery("GetIndustrialProductInfo", req.TraceabilityCode)
		if res == "" || err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"message": "请检查溯源码是否正确!!",
			})
			return
		}
		args = append(args, req.TraceabilityCode)
	}

	args = append(args, fields.Arg1, fields.Arg2, fields.Arg3, fields.Arg4, fields.Arg5)
	args = append(args, "")                    // arg6: 无图片（压缩模式暂不包含图片）
	args = append(args, string(evidenceBytes)) // arg7: 压缩证据

	res, err := pkg.ChaincodeInvoke("Uplink", args)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "uplink failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":                200,
		"message":             "uplink success (compressed)",
		"txid":                res,
		"traceabilityCode":    args[1],
		"compressionEvidence": evidence,
	})
}

// CompressTest 是一个调试接口，接受原始 JSON 数据并返回压缩后的 Base64 及压缩证据。
// POST /compress/test  body: { "data": { ... } }
func CompressTest(c *gin.Context) {
	var req struct {
		Data json.RawMessage `json:"data" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request: " + err.Error(),
		})
		return
	}

	compressor := service.NewCompressor(settings.Cfg.Compression.Algorithm)
	b64, evidence, err := compressor.Compress(req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "compress failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":                200,
		"message":             "compress success",
		"compressedPayload":   b64,
		"compressionEvidence": evidence,
	})
}

// readTraceabilityCode supports both JSON body and form-data.
// Frontend trace query uses JSON; uploading/legacy calls can still use form-data.
func readTraceabilityCode(c *gin.Context) string {
	normalize := func(v string) string {
		v = strings.TrimSpace(v)
		if v == "" {
			return ""
		}
		// keep digits only (frontend也会做类似处理，但后端兜底更稳)
		b := strings.Builder{}
		b.Grow(len(v))
		for _, r := range v {
			if r >= '0' && r <= '9' {
				b.WriteRune(r)
			}
		}
		code := b.String()
		// 兼容历史数据：16-19位（不同环境/历史数据可能不是固定18位）
		if len(code) < 16 || len(code) > 19 {
			return ""
		}
		return code
	}

	// Prefer JSON if Content-Type indicates JSON
	if strings.Contains(strings.ToLower(c.GetHeader("Content-Type")), "application/json") {
		var body struct {
			TraceabilityCode string `json:"traceabilityCode"`
		}
		if err := c.ShouldBindJSON(&body); err == nil {
			return normalize(body.TraceabilityCode)
		}
		// If JSON bind fails, fall back to form/query below
	}
	code := c.PostForm("traceabilityCode")
	if strings.TrimSpace(code) == "" {
		code = c.Query("traceabilityCode")
	}
	return normalize(code)
}
