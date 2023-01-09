# ダイレクトメール返信分離器 dmrs: Direct-Mail Responses Sepalator

送信したダイレクトメールを、返信済みと未返信に分離する

## インストール手順

実行ファイルのショートカットをデスクトップに作る

| 実行ファイル             | OS      |
|------------------------|---------|
| dmrs_windows_amd64.exe | Windows |
| dmrs_darwin_amd64      | Mac OSX |

## 操作手順

- 1.Thunderbird での操作
  - a) ダイレクトメールの送信: 宛先は BCC に指定する事
  - b) 返信の受信
  - c) 送信済みダイレクトメールと返信メールを、フォルダにドラッグ＆ドロップ
- 2.dmrs での操作
  - a) メールを入れたフォルダを、dmrs ショートカットにドラッグ＆ドロップ
  - b) 返信済みと未返信のメールアドレスを、メールを入れたフォルダに作成
    - _返信済.txt: ダイレクトメールに返信したメールアドレス
    - _未返信.txt: ダイレクトメールに返信していないメールアドレス
- 3.Thunderbird での操作
  - 返信済みのメールアドレスに、詳しい内容を送る
  - 未返信のメールアドレスに、定期的に再度ダイレクトメールを送る

## ビルド手順

``` zsh
brew install go # go1.19.4 darwin/amd64
brew install gox
go get github.com/sg3des/eml
gox --osarch "windows/amd64 darwin/amd64"
```
