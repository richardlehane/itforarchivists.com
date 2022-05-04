if (!String.prototype.startsWith) {
  String.prototype.startsWith = function (str){
    return this.lastIndexOf(str, 0) === 0;
  };
}

// https://www.smashingmagazine.com/2018/01/drag-drop-file-uploader-vanilla-js/

function drop_handler(ev) {
  ev.preventDefault();
  // Get the id of the target and add the moved element to the target's DOM
  if (ev.dataTransfer.items) {
    // Use DataTransferItemList interface to access the file(s)
    for (var i = 0; i < ev.dataTransfer.items.length; i++) {
      // If dropped items aren't files, reject them
      if (ev.dataTransfer.items[i].kind === 'file') {
        var file = ev.dataTransfer.items[i].getAsFile();
        alert('... item[' + i + '].name = ' + file.name);
      }
    }
  } else {
    // Use DataTransfer interface to access the file(s)
    for (var i = 0; i < ev.dataTransfer.files.length; i++) {
      alert('... file[' + i + '].name = ' + ev.dataTransfer.files[i].name);
    }
  }
}

/*
$("div#upload").dropzone({
		url: "/siegfried/identify", 
		init: function() {
			this.on("success", function(file, responseText) {
      			$.each(responseText, function(index, value) {
      				for (var key in value) {
                if (key == "ns") {
                  var newline = document.createElement("div");
                  newline.appendChild(document.createTextNode("--------------------"));
                  file.previewTemplate.appendChild(newline);
                }
      					var d = document.createElement("div");
      					d.className = key;
      					var s = document.createElement("strong");
      					s.appendChild(document.createTextNode(key));
      					d.appendChild(s);
                if (key == "id") {
                  d.appendChild(document.createTextNode(": "));
                  if (value[key].startsWith("x-fmt") || value[key].startsWith("fmt")) {
                    var a = document.createElement("a");
                    a.appendChild(document.createTextNode(value[key]));
                    a.href = "http://www.nationalarchives.gov.uk/pronom/".concat(value[key])
                    d.appendChild(a)
                  } else if (value[key].startsWith("fdd")) {
                    var a = document.createElement("a");
                    a.appendChild(document.createTextNode(value[key]));
                    a.href = "https://www.loc.gov/preservation/digital/formats/fdd/".concat(value[key],".shtml")
                    d.appendChild(a)
                  } else {
                    d.appendChild(document.createTextNode(value[key]));
                  }
                } else {
                  d.appendChild(document.createTextNode(": " + value[key]));
                }
      					file.previewTemplate.appendChild(d);
      				}	
      			});
   		})},
		previewTemplate: "<div class=\"dz-preview dz-file-preview\"><div class=\"dz-details\"><div class=\"dz-filename\"><span data-dz-name></span></div><div class=\"dz-size\" data-dz-size></div></div><div class=\"dz-error-message\"><span data-dz-errormessage></span></div></div>"
});
*/