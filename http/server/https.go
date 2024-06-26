package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"time"
)

func main() {
	// 创建TLS配置
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// 生成自签名证书和私钥
	certPEM, keyPEM, err := generateSelfSignedCertAndKey()
	if err != nil {
		log.Fatal(err)
	}

	// 将证书和私钥解析为TLS密钥对
	cert, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		log.Fatal(err)
	}

	// 将TLS密钥对添加到TLS配置中
	tlsConfig.Certificates = []tls.Certificate{cert}

	// 创建HTTPS Server
	server := &http.Server{
		Addr:      ":8080",
		TLSConfig: tlsConfig,
	}

	// 设置HTTP请求处理函数
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 打印请求URL
		fmt.Printf("Request URL: %s\n", r.URL.Path)

		// 获取请求查询参数
		queryParams := r.URL.Query()

		// 打印所有查询参数
		for key, values := range queryParams {
			fmt.Printf("%s: %s\n", key, values)
		}

		// 打印请求Header
		fmt.Println("Request Headers:")
		for key, value := range r.Header {
			fmt.Printf("%s: %s\n", key, value)
		}

		// 读取请求Body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		// 打印请求Body
		fmt.Printf("Request Body: %s\n", string(body))

		fmt.Fprint(w, "Hello, world!")
	})

	// 启动HTTPS Server
	log.Fatal(server.ListenAndServeTLS("", ""))
}

// 生成自签名证书和私钥
func generateSelfSignedCertAndKey() (string, string, error) {
	// 生成RSA私钥
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	// 创建自签名证书模板
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "localhost"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		BasicConstraintsValid: true,
	}

	// 使用证书模板和私钥生成自签名证书
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return "", "", err
	}

	// 将证书和私钥写入PEM格式的字符串
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return string(certPEM), string(keyPEM), nil
}
