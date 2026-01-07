package controller

import (
	"backend/pkg"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

	// Query IPFS file manifests (CID/Hash) associated with this traceability code
	manifestsRaw, mErr := pkg.ChaincodeQuery("GetFileManifestsByTrace", traceCode)
	fileHashesByRole := map[string][]string{}
	// init keys to keep response stable
	fileHashesByRole["raw_supplier"] = []string{}
	fileHashesByRole["manufacturer"] = []string{}
	fileHashesByRole["carrier"] = []string{}
	fileHashesByRole["dealer"] = []string{}

	if mErr == nil && strings.TrimSpace(manifestsRaw) != "" {
		trim := strings.TrimSpace(manifestsRaw)
		var arr []any
		if err := json.Unmarshal([]byte(trim), &arr); err == nil {
			seen := map[string]map[string]bool{}
			for k := range fileHashesByRole {
				seen[k] = map[string]bool{}
			}
			for _, item := range arr {
				m, ok := item.(map[string]any)
				if !ok {
					continue
				}
				role, _ := m["role"].(string)
				hash, _ := m["hash"].(string)
				role = strings.TrimSpace(role)
				hash = strings.TrimSpace(hash)
				if role == "" || hash == "" {
					continue
				}
				if _, exists := fileHashesByRole[role]; !exists {
					// ignore unknown roles to avoid leaking
					continue
				}
				if seen[role][hash] {
					continue
				}
				seen[role][hash] = true
				fileHashesByRole[role] = append(fileHashesByRole[role], hash)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "query success",
		"data": gin.H{
			"product":          productEnriched,
			"fileHashesByRole": fileHashesByRole,
		},
	})
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
	args = append(args, c.PostForm("arg1"), c.PostForm("arg2"), c.PostForm("arg3"), c.PostForm("arg4"), c.PostForm("arg5"))
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
	}
	args = append(args, "")
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
