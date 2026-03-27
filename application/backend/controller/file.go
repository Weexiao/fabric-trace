package controller

import (
	"backend/pkg"
	"backend/pkg/storage"
	"backend/settings"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// File upload/download controller for off-chain encrypted storage.

// UploadFile handles file upload: encrypt -> IPFS -> chaincode manifest.
func UploadFile(c *gin.Context) {
	userIDRaw, _ := c.Get("userID")
	userID := fmt.Sprint(userIDRaw)
	traceID := c.PostForm("traceabilityCode")
	if traceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "traceabilityCode is required"})
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "file is required"})
		return
	}
	if settings.Cfg.Storage.MaxSizeMB > 0 && fileHeader.Size > settings.Cfg.Storage.MaxSizeMB*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("file exceeds limit %dMB", settings.Cfg.Storage.MaxSizeMB)})
		return
	}
	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// role mapping based on user type from chaincode
	userType, err := pkg.ChaincodeQuery("GetUserType", userID)
	if err != nil || userType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to get user type"})
		return
	}
	role, err := mapUserTypeToRole(userType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	fileID := pkg.GenerateID()
	svc, err := storage.NewService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "init storage failed: " + err.Error()})
		return
	}

	fh, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "open file failed: " + err.Error()})
		return
	}
	defer fh.Close()

	manifest, err := svc.Upload(c.Request.Context(), traceID, fileID, mimeType, fh, fileHeader.Size)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "upload failed: " + err.Error()})
		return
	}
	manifest.Role = role
	manifest.Uploader = userID

	manifestJSON, err := storage.MarshalManifestJSON(manifest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "marshal manifest failed: " + err.Error()})
		return
	}
	if _, err = pkg.ChaincodeInvoke("PutFileManifest", []string{userID, manifestJSON}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "chaincode invoke failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "upload success",
		"data":    manifest,
	})
}

// DownloadFile streams decrypted file back to client with role-based access control.
func DownloadFile(c *gin.Context) {
	userIDRaw, _ := c.Get("userID")
	userID := fmt.Sprint(userIDRaw)
	fileID := c.Param("fileID")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "fileID is required"})
		return
	}

	// fetch manifest from chaincode
	res, err := pkg.ChaincodeQuery("GetFileManifest", fileID)
	if err != nil || res == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "manifest not found"})
		return
	}
	var manifest storage.Manifest
	if err := json.Unmarshal([]byte(res), &manifest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "decode manifest failed: " + err.Error()})
		return
	}

	// role-based access: manufacturer can download all, others only their uploads
	userType, err := pkg.ChaincodeQuery("GetUserType", userID)
	if err != nil || userType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to get user type"})
		return
	}
	if userType != "制造商" && manifest.Uploader != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "not allowed"})
		return
	}

	svc, err := storage.NewService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "init storage failed: " + err.Error()})
		return
	}
	rc, _, err := svc.Download(c.Request.Context(), manifest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "download failed: " + err.Error()})
		return
	}
	defer rc.Close()

	c.Header("Content-Type", manifest.Mime)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", manifest.FileID))
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, rc)
}

// ListManifests lists files under a traceability code.
func ListManifests(c *gin.Context) {
	traceID := c.PostForm("traceabilityCode")
	if traceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "traceabilityCode is required"})
		return
	}
	res, err := pkg.ChaincodeQuery("GetFileManifestsByTrace", traceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": res})
}

func mapUserTypeToRole(userType string) (string, error) {
	switch userType {
	case "原料供应商":
		return "raw_supplier", nil
	case "制造商":
		return "manufacturer", nil
	case "物流承运商":
		return "carrier", nil
	case "经销商":
		return "dealer", nil
	default:
		return "", fmt.Errorf("unsupported user type %s", userType)
	}
}
