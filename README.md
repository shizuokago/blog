# GoogleAppEngine Blog Engine

Shizuoka.goでは、GoogleAppEngine上でブログを動作させています。

編集画面でGopherJSを利用し、エディタを実現しています。
デザインを変更するには以下のファイルを変更する必要があります。

  一覧を表示するテンプレート
      cmd/templates/index.tmpl

  記事のテンプレート
      cmd/templates/entry/entry.tmpl
      cmd/templates/entry/action.tmpl

## GopherJS

GopherJSの動作が、1.12以外では動作しない為、
テンプレートの更新にはそれらが必要になります。
※基本的に下書き用に利用している為、Generateする部分は別になります。

GOPHERJS_GOROOT を利用する必要があります。

## Install 


## Deploy

- production

```bash
   gcloud app deploy --project=[Project ID] app.yaml index.yaml
```

## Command


## Sample

    http://shizuoka-go.appspot.com/
