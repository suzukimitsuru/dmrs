// ダイレクトメール返信分離器 dmrs: Direct-Mail Responses Sepalator
// 送信したダイレクトメールを、返信済みと未返信に分離する
package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/sg3des/eml"
)

func replaceExt(filePath, to string) string {
	extLen := len(filepath.Ext(filePath))
	extSep := ""
	if extLen > 0 {
		extSep = ""
	} else {
		extSep = "."
	}
	return filePath[:len(filePath)-extLen] + extSep + to
}

func loggingSettings(filename string) {
	logfile, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetOutput(multiLogFile)
}

func dirwalk(dir string) []string {
	var paths []string
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
		return paths
	}
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}
	return paths
}

func addrToEmails(addrs []eml.Address) []string {
	var emails []string
	for _, addr := range addrs {
		emails = append(emails, addr.Email())
	}
	return emails
}

func include(slice []string, target string) bool {
	for _, val := range slice {
		if val == target {
			return true
		}
	}
	return false
}

func emailsToByte(emails []string) []byte {
	b := []byte{}
	for _, line := range emails {
		b = append(b, []byte(line+";")...)
	}
	return b
}

func main() {
	// 初期設定
	loggingSettings(replaceExt(os.Args[0], "log"))
	log.Printf("%s %s", os.Args[0], "1.0.0")

	// 引数にフォルダ名が無かったら、エラー
	argCount := len(os.Args)
	if argCount < 2 {
		log.Fatalln("メールを入れたフォルダを指定してください")
		os.Exit(1)
	}
	folder := os.Args[argCount-1]
	log.Printf("フォルダ: %s", folder)

	// フォルダ内のメールから、送信済みメールアドレスと受信メールアドレスを抽出
	files := dirwalk(folder)
	sendeds := []string{}
	receives := []string{}
	for _, filename := range files {
		// フォルダ内のファイルが読めなかったら、次のファイルへ移る
		bytes, err := os.ReadFile(filename)
		if err != nil {
			continue
		}
		// メールでは無かったら、次のファイルへ移る
		msg, err := eml.Parse(bytes)
		if err != nil && err.Error() == "unexpected EOF" {
			continue
		}
		// 送信済みメールアドレスと受信メールアドレスを抽出
		bccs := addrToEmails(msg.Bcc)
		sender := addrToEmails([]eml.Address{msg.Sender})
		log.Printf("%s: %s -> %s", filename, sender, bccs)
		sendeds = append(sendeds, bccs...)
		receives = append(receives, sender...)
	}

	// 返信済みアドレスを抽出: 送信済みメールの返信では無いメールが混在していた場合、無視する
	responseds := []string{}
	for _, receive := range receives {
		if include(sendeds, receive) {
			responseds = append(responseds, receive)
		}
	}
	log.Printf("返信済み: %s", responseds)

	// 未返信アドレスを抽出
	unresponses := []string{}
	for _, sended := range sendeds {
		if !include(responseds, sended) {
			unresponses = append(unresponses, sended)
		}
	}
	log.Printf("未返信: %s", unresponses)

	// 返信結果をファイルに保存
	os.WriteFile(filepath.Join(folder, "_返信済.txt"), emailsToByte(responseds), 0666)
	os.WriteFile(filepath.Join(folder, "_未返信.txt"), emailsToByte(unresponses), 0666)
	os.Exit(0)
}
