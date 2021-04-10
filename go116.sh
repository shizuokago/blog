NOW=`pwd`
DIR="app"
HANDLER_DIR="app/handler"

INTERNAL_DIR=${NOW}/${HANDLER_DIR}/internal
LOGIC_DIR=${NOW}/${DIR}/logic

echo ${INTERNAL_DIR}

if [ -e ${INTERNAL_DIR}/static.go116 ]; then

  echo "switch template"
  mv ${INTERNAL_DIR}/template.go ${INTERNAL_DIR}/template.go115
  mv ${INTERNAL_DIR}/template.go116 ${INTERNAL_DIR}/template.go

  echo "switch static"
  mv ${INTERNAL_DIR}/static.go ${INTERNAL_DIR}/static.go115
  mv ${INTERNAL_DIR}/static.go116 ${INTERNAL_DIR}/static.go

  echo "switch logic"
  mv ${LOGIC_DIR}/html.go ${LOGIC_DIR}/html.go115
  mv ${LOGIC_DIR}/html.go116 ${LOGIC_DIR}/html.go

  echo "delete statik"
  rm -r ${INTERNAL_DIR}/statik
  rm -r ${LOGIC_DIR}/statik
  cd ${NOW}

  echo "Success"
else
  echo "Now Version 1.16???"
fi

