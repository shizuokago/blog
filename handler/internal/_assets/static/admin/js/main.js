function waitDialog(title) {
    if(title === undefined) title = "Wait...";

    var dialog = document.querySelector('dialog.waitDialog');
    dialog.querySelector('h3.mdl-dialog__title').textContent = title;

    dialog.showModal();
    return dialog;
}

function confirmDialog(title,msg,btn,func) {
    var dialog = document.querySelector('dialog.confirmDialog');

    if(title === undefined) title = "Confirm";
    if(msg === undefined) msg = "OK?";
    if(btn === undefined) btn = "OK";

    dialog.querySelector('h4.mdl-dialog__title').textContent = title;
    dialog.querySelector('p.confirm_msg').textContent = msg;
    dialog.querySelector('button.ok').textContent = btn;

    dialog.showModal();
    dialog.querySelector('button.ok').addEventListener('click', func);
    dialog.querySelector('button.cancel').addEventListener('click', function() {
        dialog.close();
    });
    return dialog;
}

function toast(msg) {
   var snackbarContainer = document.querySelector('#TOAST');
   var data = {message: msg};
   snackbarContainer.MaterialSnackbar.showSnackbar(data);
}

function parsePostValue(params) {
  var value = "";
  Object.keys(params).forEach(function (k) {
    if ( value != "" ) {
      value += "&"
    }
    value += k + "=" + params[k];
  });
  return value;
}

function request(url,params,successFunc,errorFunc) {

  var xhr = new XMLHttpRequest();

  xhr.open('POST',url);
  xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
  xhr.responseType = 'json';

  xhr.onload = function() {
    var resp = xhr.response;
    if ( successFunc !== undefined ) {
      successFunc(resp);
    }
  };

  xhr.onerror = function(prog) {
    var resp = xhr.response;
    var msg = "connection refused";
    if ( resp != null ) {
        msg = resp;
    }
    console.log(msg);
    console.log(resp);

    alert("Error",msg,function() {
      if ( errorFunc !== undefined ) {
        errorFunc(resp);
      }
    });
  };

  xhr.send(parsePostValue(params));
  return;
}

