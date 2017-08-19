var PIES = {
  options: {
    pieHole: 0.2,
    pieSliceText: 'label'
  },
  hide: "mimechart",
  show: "fmtchart", 
  redacted: false
};

if (!String.prototype.startsWith) {
  String.prototype.startsWith = function (str){
    return this.lastIndexOf(str, 0) === 0;
  };
}

function reveal(target) {
  if (PIES.show === target) {
    return;
  }
  PIES.hide = PIES.show;
  PIES.show = target;
  $(".chart").toggle();
  $("#charts button").toggleClass("pure-button-active");
}

function chart(dv, location, searchNm) {
  var chart = new google.visualization.PieChart(document.getElementById(location));
  var searchVal = "";
  function selectHandler(e) {
    var selection = chart.getSelection();
    if (selection.length > 0) {
      searchVal = dv.getValue(selection[0].row, 0);
    }
    find(searchVal, searchNm, "no MIME");
  }
  google.visualization.events.addListener(chart, 'select', selectHandler);
  return chart;
}

function find(val, name, subs) {
  if (PIES.lastSearch) {
    PIES.dt.column(PIES.lastSearch+':name').search('');
  } 
  if (val === subs) {
    val = "";
  } 
  PIES.dt.column(name+':name').search('^' + val + '$', true, false).draw();
  PIES.lastSearch = name;
}

function datatable(typ) {
    var a = 0;
    var data = new google.visualization.DataTable();
    data.addColumn('string');
    data.addColumn('number');
    PIES[typ+"_idxs"] = [];
    for (i = 0; i < RESULTS.results.length; i++) {
      PIES[typ+"_idxs"].push([]);
      for (j = 0; j < RESULTS.results[i][typ].length; j++) {
          data.addRow(RESULTS.results[i][typ][j]);
          PIES[typ+"_idxs"][i].push(a);
          a++;
      }
    }
    return data;
}

function table(cols, data) {
  if (PIES.dt) {
    PIES.dt.destroy();
    $('#table').empty();
  }
  PIES.dt = $('#table').DataTable( {
      data: data,
      columns: cols.map(function(col, idx){
        var hid = true;
        if (idx >= cols.length - 2) {
            hid = false;
        }
        return {
            title: col,
            name: col,
            visible: hid,
            render: function ( data, type, row, meta ) {
              if (type !== 'display') {
                return data;
              } 
              if (col === "id") {
                if (data.startsWith("fmt") || data.startsWith("x-fmt")) {
                  return '<a href="http://www.nationalarchives.gov.uk/pronom/'+data+'">' + data + '</a>';
                }
                if (data.startsWith("fdd")) {
                  return '<a href="https://www.loc.gov/preservation/digital/formats/fdd/'+data+'.shtml">' + data + '</a>';
                }
              } else if (idx > 5 && data.length > 80) {
                return '<span title="'+data+'">'+data.substr( 0, 78 )+'...</span>';
              }
              return data;
          }
      }}),
      dom: 'Bfrtip',
      buttons: [
        'copy', 'csv'
      ]
  });
}

function load(num) {
  PIES.lastSearch = false;
  // set text
  $("#errNo span").text(RESULTS.results[num].errors);
  $("#warnNo span").text(RESULTS.results[num].warnings);
  $("#unkNo span").text(RESULTS.results[num].unknowns);
  $("#multiNo span").text(RESULTS.results[num].multipleIDs);
  $("#dupesNo span").text(RESULTS.results[num].duplicates);

  // load charts
  $("#"+PIES.hide).show();
  PIES.fmtView.setRows(PIES["fmtCounts_idxs"][num]);
  PIES.fmtChart.draw(PIES.fmtView, PIES.options);
  PIES.mimeView.setRows(PIES["mimeCounts_idxs"][num]);
  PIES.mimeChart.draw(PIES.mimeView, PIES.options);
  $("#"+PIES.hide).hide();
  // load table
  table(RESULTS.results[num].titles, RESULTS.results[num].rows);
}

function initialize() {
  PIES.fmtView = new google.visualization.DataView(datatable("fmtCounts"));
  PIES.mimeView = new google.visualization.DataView(datatable("mimeCounts"));
  PIES.fmtChart = chart(PIES.fmtView, "fmtchart", "id");
  PIES.mimeChart = chart(PIES.mimeView, "mimechart", "mime");
  load(0);

  $("#errNo").click(function() { find('.+', 'errors'); return false; });
  $("#warnNo").click(function() { find('.+', 'warning'); return false; });
  $("#unkNo").click(function() { find('UNKNOWN', 'id'); return false; });
  $("#multiNo").click(function() { find('true', 'hasMultiID'); return false; });
  $("#dupesNo").click(function() { find('true', 'isDuplicate'); return false; });
  $("#reset").click(function() { find('.+', 'id'); return false; });
  
  $("#redactNow").click(function() {
        if (PIES.redacted) {
          return false;
        }
        $(".publish-button").toggle();
        var formData = new FormData();
        var blob = new Blob([JSON.stringify(RESULTS)], { type: "application/json"});
        formData.append("results", blob);
        $.ajax({
          url: "redact", 
          method: "POST",
          processData: false,
          contentType: false,
          data: formData,
          success: function(data, textStatus) {
              RESULTS = data;
              PIES.redacted = true;
              load(0);
              $(".publish-button").toggle();
          },
          error: function (jqXHR, textStatus, errorThrown) {
            $("#share-form").append("<div><strong>" +textStatus + "</strong>: " + errorThrown+ "</div>");
          }
        });
        return false;
    });  

  $("#share-form").submit(function(event) {
        $(".publish-button").toggle();
        var formData = new FormData();
        var blob = new Blob([JSON.stringify(RESULTS)], { type: "application/json"});
        formData.append("results", blob);
        $.each($('#share-form').serializeArray(), function(i, field) {
          formData.append(field.name, field.value);
        });
        if (PIES.redacted || formData.redact !== "true") {
          formData.redact = "false";
        }
        event.preventDefault();
        $.ajax({
          url: "share", 
          method: "POST",
          processData: false,
          contentType: false,
          data: formData,
          success: function(data, textStatus) {
            if (data.success) {
              window.location.replace(data.success);
            }
          },
          error: function (jqXHR, textStatus, errorThrown) {
            $("#share-form").append("<div><strong>" +textStatus + "</strong>: " + errorThrown+ "</div>");
          }
        });
    });
}

google.charts.load('current', {'packages':['corechart']});

google.charts.setOnLoadCallback(function() {
  $(initialize())
});