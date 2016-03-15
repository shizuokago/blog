# GoogleAppEngine Blog Engine



## Using

	"github.com/gopherjs/gopherjs"

	"golang.org/x/tools/present"
	"github.com/gorilla/mux"
	"github.com/knightso/base"
	"github.com/pborman/uuid"

## local serve
## AppEngine Deploy

- local
```bash
   goapp serve app
```

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


### if you change design
```bash
    go run gopherjs/deploy.go
```

