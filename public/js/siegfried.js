if (!String.prototype.startsWith) {
  String.prototype.startsWith = function (str){
    return this.lastIndexOf(str, 0) === 0;
  };
}

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
