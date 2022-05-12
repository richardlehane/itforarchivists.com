var sets;

var textFn = function(){
    	if(!sets) {
    		return;
    	};
    	var text = sets[0];
    	for (i = 1; i < sets.length; i++) { 
    		if (i > 0) {
    			text += ", ";
    		}
    		text += sets[i];
		}
    	document.getElementById("results").value = text;
};

var textnlFn = function(){
        if(!sets) {
            return;
        };
        var text = sets[0];
        for (i = 1; i < sets.length; i++) { 
            if (i > 0) {
                text += "\n";
            }
            text += sets[i];
        }
        document.getElementById("results").value = text;
};

var goFn = function(){
    	if(!sets) {
    		return
    	};
    	var text = "func IsType(puid string) bool {\n    switch puid {\n    case ";
    	for (i = 0; i < sets.length; i++) { 
    		if (i > 0) {
    			text += ", ";
    		}
    		text += '"' + sets[i] + '"';
		}
		text += ":\n        return true\n    }\n    return false\n}";
    	document.getElementById("results").value = text;
};

var pyFn = function(){
    	if(!sets) {
    		return;
    	};
    	var text = "def IsType(puid):";
    	line = "\n    return puid in [";
    	for (i = 0; i < sets.length; i++) { 
    		if (i > 0) {
    			line += ", ";
    		}
    		if (line.length > 40) {
    			text += line.substring(0, line.length - 1);
    			line = "\n                    ";
    		}
    		line += "'" + sets[i] + "'";
		}
		text += line +  "]\n";
    	document.getElementById("results").value = text;
};

var fn = textFn;

document.getElementById("sets-form").addEventListener("submit", (e) => {
	e.preventDefault();
	const form = new FormData(document.getElementById("sets-form"));
	fetch("siegfried/sets", {
		method: 'POST',
		body: form
	})
	.then(response => {
		return response.json()
	})
	.then(val => {
		sets = val;
		fn();
	})
})

function changeMode(m) {
   	if(m.value == "text") {
   		fn = textFn;
   	} else if (m.value == "text-nl") {
        fn = textnlFn; 
    } else if (m.value == "golang") {
  		fn = goFn;
   	} else if (m.value == "python") {
   		fn = pyFn;
   	}
   	fn();
};

function add() {
    let newinput = document.getElementById("addlist").cloneNode(true);
    newinput.removeAttribute("id");
	newinput.removeAttribute("style");
	document.getElementById("sets-form").prepend(newinput);
}

function del(el) {
    el.parentElement.remove();
}

function cp() {
	let text = document.getElementById("results").value;
	navigator.clipboard.writeText(text);
}
