# GoogleAppEngine Blog Engine


## Using

	"golang.org/x/tools/present"

	"github.com/gorilla/mux"
	"github.com/knightso/base"
	"github.com/pborman/uuid"

	"github.com/gopherjs/gopherjs"

	"github.com/nfnt/resize"
	"github.com/robfig/graphics-go/graphics"


## local serve
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


### if you change design
```bash
    go run gopherjs/deploy.go
```


## Sample

http://shizuoka-go.appspot.com/
