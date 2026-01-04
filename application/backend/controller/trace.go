package controller

import (
	"backend/pkg"
	"fmt"
	"os"

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
	res, err := pkg.ChaincodeQuery("GetIndustrialProductInfo", c.PostForm("traceabilityCode"))
	if err != nil {
		c.JSON(200, gin.H{
			"message": "查询失败：" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})

}

// 获取用户的工业产品ID列表
func GetIndustrialProductList(c *gin.Context) {
	userID, _ := c.Get("userID")
	res, err := pkg.ChaincodeQuery("GetIndustrialProductList", userID.(string))
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
	res, err := pkg.ChaincodeQuery("GetIndustrialProductHistory", c.PostForm("traceabilityCode"))
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

func buildArgs(c *gin.Context, traceCode string) []string {
	var args []string
	userID, _ := c.Get("userID")
	userType, _ := pkg.ChaincodeQuery("GetUserType", userID.(string))
	args = append(args, userID.(string))
	if userType == "原料供应商" {
		args = append(args, traceCode)
	} else {
		res, err := pkg.ChaincodeQuery("GetIndustrialProductInfo", c.PostForm("traceabilityCode"))
		if res == "" || err != nil || len(c.PostForm("traceabilityCode")) != 18 {
			c.JSON(200, gin.H{
				"message": "请检查溯源码是否正确!!",
			})
			return nil
		}
		args = append(args, c.PostForm("traceabilityCode"))
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
