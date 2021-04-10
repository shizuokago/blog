# GoogleAppEngine Blog Engine

Shizuoka.goでは、GoogleAppEngine上でブログを動作させています。
以前までGopherJSを利用していましたが、WASMに変更をかけました。

## Datastore

Datastoreを利用する為、開発環境では以下を行います

```bash
    gcloud beta emulators datastore start
```

で起動しておきます。

実環境の判定はmet.OnGCE()で行い、ProjectIDも同一でのみ動作します。
開発環境はProjectID=blogで動作します。

## Web

```bash
    go run _cmd/main.go
```

で動作確認を行えます。

## Deploy

```bash
   gcloud app deploy --project=[Project ID] app.yaml index.yaml
```

# デザイン変更について

デザインを変更するには以下のファイルを変更する必要があります。

  一覧を表示するテンプレートはこちら

      handler/internal/_assets/index.tmpl

  記事のテンプレートはこちら

      logic/_entry/entry.tmpl
      logic/_entry/action.tmpl

## wasm生成

    GOOS=js GOARCH=wasm -o editor.wasm _cmd/editor/editor.go 

    実際に使用するWASMはgzipを行って実行しています。

    app/_cmd/editor/wasm.sh を参考にしてください。

## Authentication

/admin/ にアクセスするとGoogle認証が入ります。
Blog設定で指定したアドレスのみ認証可能です。

※初回は誰でもログイン可能であるため、設定してください。

app/handler/internal/_assets/templates/authentication.tmpl

にプロジェクトのAPIのClientIDが設定されていますので、
そこを変更する必要があります。


