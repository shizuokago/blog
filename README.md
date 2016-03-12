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

    blog design
    ./app/templates/index.tmpl
    ./app/templates/entry/*

    backgroung image
    ./app/static/images/backgound


### if you change design
```bash
    go run gopherjs/deploy.go
```

