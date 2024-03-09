package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

var greenColor = color.New(color.FgGreen).SprintFunc()
var blueColor = color.New(color.FgBlue).SprintFunc()
var redColor = color.New(color.FgRed).SprintFunc()
var redCyan = color.New(color.FgCyan).SprintFunc()

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With,user-agent")
		c.Header("Access-Control-Allow-Methods", "POST,OPTIONS")
		c.Header("Content-Type", "application/json")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func procReq(c *gin.Context) {
	dtSend := time.Now()
	fmt.Println(redCyan(">>>>---------------------------------------------------------" + dtSend.Format("01-02-2006 15:04:05.000000")))

	corpoRequest, erro := io.ReadAll(io.Reader(c.Request.Body))
	if erro != nil {
		fmt.Println("erro", erro)
		c.JSON(http.StatusBadRequest, string(`{"error":true}`))
		return
	}
	defer c.Request.Body.Close()

	// fmt.Println(c.Request.Header)
	// fmt.Println("headers inicio ##")
	// for k, m := range c.Request.Header {
	// 	fmt.Println(k, m)
	// }
	// fmt.Println("headers fim ##")
	// fmt.Println("header " + c.Request.Header.Get("Authorization"))

	qrTo, ison := c.GetQuery("to")
	if !ison {
		fmt.Println("falta 'to'")
		c.JSON(http.StatusBadRequest, string(`{"error":true}`))
		return
	}

	qrSerial, ison := c.GetQuery("serial")
	if !ison {
		fmt.Println("falta 'Serial'")
		c.JSON(http.StatusBadRequest, string(`{"error":true}`))
		return
	}

	fmt.Println(">To: ", qrTo)
	fmt.Println(">From: " + redColor(qrSerial) + " | " + c.Request.Header.Get("Authorization"))
	fmt.Println(">Sent: " + greenColor(string(corpoRequest)))
	fmt.Println(redCyan("----------------------------------------------------------------------------------->>>>"))

	req, err := http.NewRequest(http.MethodPost, qrTo, bytes.NewBuffer(corpoRequest))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", c.Request.Header.Get("Authorization"))
	req.Header.Add("traceId", "fake_traceId")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Go dev Proxy >")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	corpoResp, erro := io.ReadAll(io.Reader(res.Body))
	if erro != nil {
		fmt.Println("erro", erro)
		return
	}
	outbody := string(corpoResp)

	dtRec := time.Now()
	fmt.Println(blueColor("<<<<---------------------------------------------------------" + dtRec.Format("01-02-2006 15:04:05.000000")))
	fmt.Println("<Resp to: " + redColor(qrSerial) + " code: " + res.Status)
	fmt.Println("<Rec: ", greenColor(outbody))
	fmt.Println(blueColor("-----------------------------------------------------------------------------------<<<<\n"))

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("User-Agent", "Go dev Proxy <")
	// c.Header("Content-Type", "application/json")

	c.Data(http.StatusOK, "application/json", corpoResp)
}

func main() {
	fmt.Println(greenColor("Starting Proxy"))

	gin.SetMode(gin.ReleaseMode)
	// router := gin.Default()
	router := gin.New()

	router.Use(CORSMiddleware())
	router.SetTrustedProxies([]string{"http://localhost"})
	router.POST("/", procReq)

	fmt.Println(greenColor("Running at port: 1248"))
	router.Run("localhost:1248")
}
