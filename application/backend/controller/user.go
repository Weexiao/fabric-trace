package controller

import (
	"backend/model"
	"backend/pkg"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	// 将用户信息存入mysql数据库
	var user model.MysqlUser
	user.UserID = pkg.GenerateID()
	user.Username = c.PostForm("username")
	user.Password = c.PostForm("password")
	user.RealInfo = pkg.EncryptByMD5(c.PostForm("username"))
	err := pkg.InsertUser(&user)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "register failed：" + err.Error(),
		})
		return
	}
	// 将用户信息存入区块链
	// userID string, userType string, realInfoHash string
	// 将post请求的参数封装成一个数组args
	var args []string
	args = append(args, user.UserID)
	args = append(args, c.PostForm("userType"))
	args = append(args, user.RealInfo)
	res, err := pkg.ChaincodeInvoke("RegisterUser", args)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "register failed：" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "register success",
		"txid":    res,
		"userID":  user.UserID,
	})
}

func Login(c *gin.Context) {
	var user model.MysqlUser
	user.Username = c.PostForm("username")
	user.Password = c.PostForm("password")
	// 获取用户ID
	var err error
	user.UserID, err = pkg.GetUserID(user.Username)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "没有找到该用户",
		})
		return
	}
	userType, err := GetUserType(user.UserID)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "login failed:" + err.Error(),
		})
		return
	}
	err = pkg.Login(&user)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "login failed:" + err.Error(),
		})
		return
	}

	// 获取用户动态属性
	dynamicAttributes, err := pkg.GetUserDynamicAttributes(user.UserID)
	if err != nil {
		dynamicAttributes = "{}"
	}

	// 生成jwt
	jwt, err := pkg.GenToken(user.UserID, userType, dynamicAttributes)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "login failed:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "login success",
		"jwt":     jwt,
	})
}

func Logout(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":    200,
		"message": "logout success",
	})
}

// 获取用户类型
func GetUserType(userID string) (string, error) {
	userType, err := pkg.ChaincodeQuery("GetUserType", userID)
	if err != nil {
		return "", err
	}
	return userType, nil
}

// 获取用户信息
func GetInfo(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(200, gin.H{
			"message": "get user type failed",
		})
	}

	userType, err := GetUserType(userID.(string))
	if err != nil {
		c.JSON(200, gin.H{
			"message": "get user type failed" + err.Error(),
		})
	}

	username, err := pkg.GetUsername(userID.(string))
	if err != nil {
		c.JSON(200, gin.H{
			"message": "get user name failed" + err.Error(),
		})
	}

	// 获取用户动态属性
	dynamicAttributes, err := pkg.GetUserDynamicAttributes(userID.(string))
	if err != nil {
		dynamicAttributes = "{}"
	}

	c.JSON(200, gin.H{
		"code":              200,
		"message":           "get user type success",
		"userType":          userType,
		"username":          username,
		"dynamicAttributes": dynamicAttributes,
	})
}

// 更新用户动态属性 (管理员操作)
func UpdateUserDynamicAttributes(c *gin.Context) {
	userID := c.PostForm("user_id")
	dynamicAttributes := c.PostForm("dynamic_attributes") // JSON string like {"region":"Sichuan","data_level":"Internal"}

	if userID == "" || dynamicAttributes == "" {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "user_id and dynamic_attributes are required",
		})
		return
	}

	err := pkg.UpdateUserDynamicAttributes(userID, dynamicAttributes)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "update failed:" + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "update success",
	})
}

// 获取所有用户列表 (管理员操作)
func GetAllUsers(c *gin.Context) {
	users, err := pkg.GetAllUsers()
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "get users failed:" + err.Error(),
		})
		return
	}

	// 为每个用户添加用户类型（从区块链获取）
	for _, user := range users {
		userType, err := GetUserType(user.UserID)
		if err != nil {
			// 如果获取失败，设置为空或默认值
			user.UserType = "Unknown"
		} else {
			user.UserType = userType
		}
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "get users success",
		"data":    users,
	})
}
