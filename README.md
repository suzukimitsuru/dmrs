# ダイレクトメール返信分離器 dmrs: Direct-Mail Responses Sepalator

送信したダイレクトメールを、返信済みと未返信に分離する

## ビルド手順

``` zsh
brew install go # go1.19.4 darwin/amd64
brew install gox
go get github.com/sg3des/eml
gox --osarch "windows/amd64 windows/386 darwin/amd64"
```
