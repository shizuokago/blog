NOW=`pwd`
DIR="app"
HANDLER="app/handler"

INTERNAL_DIR=${NOW}/${HANDLER_DIR}/internal
LOGIC_DIR=${NOW}/${DIR}/logic

if [ -e ${INTERNAL_DIR}/static.go115 ]; then

  echo "switch template"
  mv ${INTERNAL_DIR}/template.go ${INTERNAL_DIR}/template.go116
  mv ${INTERNAL_DIR}/template.go115 ${INTERNAL_DIR}/template.go

  echo "switch static"
  mv ${INTERNAL_DIR}/static.go ${INTERNAL_DIR}/static.go116
  mv ${INTERNAL_DIR}/static.go115 ${INTERNAL_DIR}/static.go

  echo "generate internal statik"
  cd ${INTERNAL_DIR}
  statik -src _assets

  echo "switch logic"
  mv ${LOGIC_DIR}/html.go ${LOGIC_DIR}/html.go116
  mv ${LOGIC_DIR}/html.go115 ${LOGIC_DIR}/html.go

  echo "generate logic statik"
  cd ${LOGIC_DIR}
  statik -src _entry

  cd ${NOW}

  echo "Success"
else
  echo "Now Version 1.15???"
fi

