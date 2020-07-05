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


//publish
	//jQuery(PUBLISH).On(jquery.CLICK, func(e jquery.Event) {
		//ajax("publish")
	//})

//save
	//jQuery(SAVE).On(jquery.CLICK, func(e jquery.Event) {
		//ajax("save")
	//})


// view
	//jQuery("button#viewBtn").On(jquery.CLICK, func(e jquery.Event) {
		//url := "/entry/" + jQuery(ARTICLE_ID).Val()
		//js.Global.Call("open", url, "_blank")
	//})

//private 
	//jQuery("button#private").On(jquery.CLICK, func(e jquery.Event) {
		//ajax("private")
	//})

//delete
	//jQuery("button#delete").On(jquery.CLICK, func(e jquery.Event) {
		//url := "/admin/article/delete/" + jQuery(ARTICLE_ID).Val()
		//l := js.Global.Get("location")
		//l.Set("href", url)
	//})
