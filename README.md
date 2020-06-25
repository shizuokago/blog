# GoogleAppEngine Blog Engine

## Install 


```bash




```

## AppEngine Deploy

- production
```bash
   gcloud app deploy --project=[Project ID] app.yaml index.yaml
```

## Original

  index HTML

    ./app/templates/index.tmpl

  article HTML

    ./app/templates/entry/entry.tmpl
    ./app/templates/entry/action.tmpl

  stylesheet

    ./app/static/css/style.css

  backgroung image

    ./app/static/images/bg_1024.jpg
    ./app/static/images/bg_2048.jpg
    ./app/static/images/bg_2880.jpg

## Command

### if you change design

create editor js
```bash
    editor
```

woking directory(in "go run cmd/editor.js")
```bash
    deploy
```

watch design(in "go run cmd/deploy.js")
### if you change design(watch)
```bash
    watcher
```

## Sample

http://shizuoka-go.appspot.com/
