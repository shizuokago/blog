window.addEventListener("DOMContentLoaded", async () => {
  const memory = new WebAssembly.Memory({
    initial: 1024,
    maximum: 1024
  });
  var mem = { memory:memory };

  const go = new Go();
  const url = "/admin/editor.wasm.gz";
  const pako = window.pako;
  let wasm = pako.ungzip(await (await fetch(url)).arrayBuffer());
  if (wasm[0] === 0x1f && wasm[1] === 0x8b) {
      wasm = pako.ungzip(wasm);
  }

  go.importObject["env"] = mem;
  console.log(go.importObject);
  const result = await WebAssembly.instantiate(wasm, go.importObject).catch((err) => {
      console.error(err);
  });
  go.run(result.instance);

  function resize() {
    var inn = window.innerHeight;
    var editor = document.querySelector("#editor");
    var left   = document.querySelector("#left");
    var right  = document.querySelector("#right");
    var out    = document.querySelector("#result");

    var h = inn - 230;
    left.height = h + "px";
    right.height = h + 10;
    out.height = h + 10;
    editor.style.height = h + "px";
  }

  resize();
  window.addEventListener("resize",function() {
     resize();
  },false);

  document.querySelector("button#publish").addEventListener("click",function(e) {
      confirmDialog("Publish?","Would you like to publish the current article?","YES!",function() {
        var d = waitDialog();
        var id = document.querySelector("input#ID");
        var url = "/admin/article/publish/" + id.value;
        var params = createArticleParam();

        request(url,params,function(){
          d.close();
        },function() {
          d.close();
        });
      });
  });

  function createArticleParam()  {
    var params = {};
    var title = document.querySelector("input#Title");
    var tags = document.querySelector("input#Tags");
    var md = document.querySelector("textarea#editor");
    params["Title"] = title.value;
    params["Tags"] = tags.value;
    params["Markdown"] = md.value;
    return params;
  }

  var md = document.querySelector("textarea#editor");
  var saveVal = md.value;

  function save(d) {

    if ( md.value === saveVal ) {
      if ( d !== undefined ) {
        d.close();
      } else {
        toast("auto saved.");
      }
      return;
    }

    var id = document.querySelector("input#ID");
    var url = "/admin/article/save/" + id.value;
    var params = createArticleParam();

    request(url,params,function(){
      saveVal = md.value;
      if ( d !== undefined ) {
        d.close();
      } else {
        toast("auto saved.");
      }
    },function() {
      if ( d !== undefined ) {
        d.close();
      } else {
        toast("auto save error");
      }
    });
  }

  var auto = document.querySelector("input#AutoSave").value;

  if ( auto === "on" ) {
    setInterval(function() {
        save();
    },60000);
  }

  document.querySelector("button#save").addEventListener("click",function(e) {
    var d = waitDialog();
    save(d);
  });

  document.querySelector("button#viewBtn").addEventListener("click",function(e) {
      var id = document.querySelector("input#ID");
      var url = "/entry/" + id.value;
      window.open(url,"_blank");
  });

  document.querySelector("button#private").addEventListener("click",function(e) {
      confirmDialog("Private?","Do you want to delete the publish data and make it private?","Yes!",function() {
        var d = waitDialog();
        var id = document.querySelector("input#ID");
        var url = "/admin/article/private/" + id.value;
        var params = createArticleParam();

        request(url,params,function(){
          d.close();
        },function() {
          d.close();
        });
      });
  });

  document.querySelector("button#delete").addEventListener("click",function(e) {
      confirmDialog("De","if you delete it,published data will also be deleted.","DELETE",function() {
        waitDialog();
        var id = document.querySelector("input#ID");
	    var url = "/admin/article/delete/" + id.value;
        location.href=url;
      });
  });

});
