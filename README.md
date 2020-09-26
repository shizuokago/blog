# GoogleAppEngine Blog Engine

Shizuoka.goでは、GoogleAppEngine上でブログを動作させています。

編集画面でGopherJSを利用し、エディタを実現しています。

## Datastore

Datastoreを利用する為、

```bash
    gcloud beta emulators datastore start
```

で起動しておきます。

## Web

```bash
    go run cmd/main.go
```

で動作確認を行えます。

## Deploy

```bash
   gcloud app deploy --project=[Project ID] app.yaml index.yaml
```

# デザイン変更について

デザインを変更するには以下のファイルを変更する必要があります。

  一覧を表示するテンプレート
      cmd/templates/index.tmpl

  記事のテンプレート
      cmd/templates/entry/entry.tmpl
      cmd/templates/entry/action.tmpl

## JS処理

  エディタ部分のJS処理

      cmd/editor/editor.go 

  JSを出力する処理

      cmd/editor/deploy.go 

  JSを出力する処理

      cmd/editor/deploy.go 

## GopherJS

GopherJSの動作が、1.12以外では動作しない為、
テンプレートの更新にはそれらが必要になります。
※基本的に下書き用に利用している為、Generateする部分は別になります。

環境変数 GOPHERJS_GOROOT を利用する必要があります。

またGopherJS自体がWindowsでは動作しませんのでWSLなどで出力しましょう。


## Authentication

/admin/ にアクセスするとGoogle認証が入ります。

Blog設定で指定したアドレスのみ認証可能です。

