package main

import (
	"crypto/tls"
	"gobankid/soap"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	cert, err := tls.LoadX509KeyPair("cert.crt", "key.key")
	if err != nil {
		log.Fatal(err)
	}

	s := soap.NewClient("https://appapi.test.bankid.com/rp/v4?wsdl", cert)

	router.GET("/authenticate", authenticateHandler(s))
	router.GET("/collect", collectHandler(s))

	router.Run()
}

func authenticateHandler(s *soap.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := &soap.EndUserInfo{
			UserInfoType: "IP_ADDR",
			Value:        "192.168.0.1",
		}
		authResp, err := s.Authenticate("190102030400", u)
		if err != nil {
			log.Fatalln(err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, authResp.OrderRef)
	}
}

func collectHandler(s *soap.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderRef := c.Query("orderref")
		if orderRef == "" {
			log.Println("No orderRef")
			return
		}

		collResp, err := s.Collect(orderRef)
		if err != nil {
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		log.Println(collResp.Status)

		if collResp.Status == soap.StatusComplete {
			c.JSON(http.StatusOK, gin.H{"status": collResp.Status, "userInfo": collResp.UserInfo})
		}
		c.JSON(http.StatusOK, gin.H{"status": collResp.Status})
	}
}
