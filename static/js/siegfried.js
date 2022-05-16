if (!String.prototype.startsWith) {
  String.prototype.startsWith = function (str) {
    return this.lastIndexOf(str, 0) === 0;
  };
}

// https://www.smashingmagazine.com/2018/01/drag-drop-file-uploader-vanilla-js/

function drop_handler(ev) {
  // Get the id of the target and add the moved element to the target's DOM
  if (ev.dataTransfer.items) {
    // Use DataTransferItemList interface to access the file(s)
    for (var i = 0; i < ev.dataTransfer.items.length; i++) {
      // If dropped items aren't files, reject them
      if (ev.dataTransfer.items[i].kind === 'file') {
        var file = ev.dataTransfer.items[i].getAsFile();
        upload(file);
      }
    }
  } else {
    // Use DataTransfer interface to access the file(s)
    for (var i = 0; i < ev.dataTransfer.files.length; i++) {
      upload(ev.dataTransfer.files[i]);
    }
  }
  ev.preventDefault();
}

function upload(file) {
  let url = '/siegfried/identify';
  let line = "--------------------";
  let formData = new FormData();
  formData.append('file', file)

  fetch(url, {
    method: 'POST',
    body: formData
  })
    .then((response) => {
      return response.json();
    })
    .then((data) => {
      var fileline = document.createElement("div");
      fileline.appendChild(document.createTextNode(line));
      var filename = document.createElement("div");
      filename.appendChild(document.createTextNode(file.name));
      document.getElementById('upload').appendChild(fileline);
      document.getElementById('upload').appendChild(filename);
      for (var idx in data) {
        let value = data[idx]
        for (var key in value) {
          if (key == "ns") {
            var newline = document.createElement("div");
            newline.appendChild(document.createTextNode(line));
            document.getElementById('upload').appendChild(newline);
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
              a.href = "https://www.loc.gov/preservation/digital/formats/fdd/".concat(value[key], ".shtml")
              d.appendChild(a)
            } else {
              d.appendChild(document.createTextNode(value[key]));
            }
          } else {
            d.appendChild(document.createTextNode(": " + value[key]));
          }
          document.getElementById('upload').appendChild(d);
        }
      }
    })
    .catch(function (error) {
      console.log(error);
    });
}

