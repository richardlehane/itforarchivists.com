/*
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
    	$( "#results" ).html(text);
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
        $( "#results" ).html(text);
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
    	$( "#results" ).html(text);
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
    	$( "#results" ).html(text);
};

var fn = textFn;

$("#sets-form").submit(function(event) {
    	event.preventDefault();
    	sets = $(this).serializeArray();
    	$.ajax({
		url: "siegfried/sets", 
		method: "POST",
		dataType: "json",
		traditional: true,
		data: sets,
  		success: function( data ) {
  			sets = data;
    		fn();
 		}
		});
});

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
    $("#addlist").filter(":last").clone().show().removeAttr( "id" ).prependTo($("#sets-form"));
}

function del(el) {
    $(el).parent().remove();
}

function cp() {
  	$("#results").select();
	document.execCommand("copy");
}
*/