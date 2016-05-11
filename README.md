# GoogleAppEngine Blog Engine


## Using

	"golang.org/x/tools/present"

	"github.com/gorilla/mux"
	"github.com/knightso/base"
	"github.com/pborman/uuid"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"

	"github.com/PuerkitoBio/goquery"

	"github.com/nfnt/resize"
	"github.com/robfig/graphics-go/graphics"


## AppEngine Deploy

- production
```bash
   goapp deploy app
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
    go run cmd/editor.go
```

woking directory(in "go run cmd/editor.js")
```bash
    go run cmd/deploy.go
```

watch design(in "go run cmd/deploy.js")
### if you change design(watch)
```bash
    go run cmd/watch.go
```





## Sample

http://shizuoka-go.appspot.com/
