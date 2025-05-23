package utils

// // メール送信関数
// func SendEmail(to, subject, body string) error {
// 	m := gomail.NewMessage()
// 	//環境変数から送信元emailを読み込む
// 	m.SetHeader("From", "pinoapi@example.com") // 送信元
// 	m.SetHeader("To", to)                      // 送信先
// 	m.SetHeader("Subject", subject)            // 件名
// 	m.SetBody("text/plain", body)              // 本文

// 	d := gomail.NewDialer("smtp.gmail.com", 587, "pinoapi@example.com", "pinoapi-password")

// 	if err := d.DialAndSend(m); err != nil {
// 		log.Println("Failed to send email:", err)
// 		return err
// 	}
// 	return nil
// }

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, body string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "465"
	//環境変数から読み込む

	username := os.Getenv("email")
	password := os.Getenv("email_password")

	from := username
	message := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)

	// TLS 接続を確立
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, &tls.Config{ServerName: smtpHost})
	if err != nil {
		fmt.Println("TLS 接続エラー:", err)
		return err
	}
	defer conn.Close()

	// SMTP クライアント作成
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		fmt.Println("SMTP クライアント作成エラー:", err)
		return err
	}

	// 認証情報を設定
	auth := smtp.PlainAuth("", username, password, smtpHost)
	if err := client.Auth(auth); err != nil {
		fmt.Println("認証エラー:", err)
		return err
	}

	// メール送信
	if err := client.Mail(from); err != nil {
		fmt.Println("MAIL コマンドエラー:", err)
		return err
	}
	if err := client.Rcpt(to); err != nil {
		fmt.Println("RCPT コマンドエラー:", err)
		return err
	}

	// データ送信
	w, err := client.Data()
	if err != nil {
		fmt.Println("DATA コマンドエラー:", err)
		return err
	}
	_, err = w.Write([]byte(message))
	if err != nil {
		fmt.Println("メッセージ送信エラー:", err)
		return err
	}
	w.Close()

	// 接続を終了
	client.Quit()

	fmt.Println("メール送信成功")
	return nil
}
