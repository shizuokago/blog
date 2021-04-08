# GoogleAppEngine Blog Engine

Shizuoka.goでは、GoogleAppEngine上でブログを動作させています。
GopherJSを利用していましたが、WASMに変更をかけました。

## Datastore

Datastoreを利用する為、開発環境では以下を行います

```bash
    gcloud beta emulators datastore start
```

で起動しておきます。
開発環境はProjectID=blogで動作。実環境ではmetaより取得してきます。

開発環境の判定は動作位置が「/srv」かどうかで判定しています。
※そのため開発環境でも動作位置が/srvの場合動作しません

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
      handler/internal/_assets/index.tmpl

  記事のテンプレート
      logic/_entry/entry.tmpl
      logic/_entry/action.tmpl

## wasm生成

      cmd/editor/editor.go 

## Authentication

/admin/ にアクセスするとGoogle認証が入ります。
Blog設定で指定したアドレスのみ認証可能です。

