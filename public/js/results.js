
   
    var columnTitles = ["filename","id","mime","size","modified","basis"];
    var dataSet = [
    [ "Agency in the Archive -Post Edit - GR.pdf", "fmt/20", "text/plain", "5421", "2011/04/25", "extension match" ],
    [ "Always Becoming 2015 Visitor Claim Form.doc", "fmt/934", "text/plain", "8422", "2011/07/25", "extension match" ],
    [ "D14 23426  Digital Records Operations recruitment (BN 14-3129) - Tab A.pdf", "fmt/1", "text/plain", "1562", "2009/01/12", "extension match" ],
    [ "dig.txt", "fmt/934", "text/plain", "6224", "2012/03/29", "extension match" ],
    [ "ESS.txt", "fmt/934", "text/plain", "5407", "2008/11/28", "extension match" ],
    [ "GH SAG talk.docx", "fmt/40", "text/plain", "4804", "2012/12/02", "extension match" ],
    [ "glenpdp.txt", "fmt/934", "text/plain", "9608", "2012/08/06", "extension match" ],
    [ "gravatar.png", "fmt/40", "text/plain", "6200", "2010/10/14", "extension match" ],
    [ "ICMS_EA_Management_Presso_090916_v03.pptx", "fmt/934", "text/plain", "2360", "2009/09/15", "extension match" ],
    [ "Increment_ Jolliffe.pdf", "fmt/40", "text/plain", "1667", "2008/12/13", "extension match" ],
    [ "Increment_Humphries.pdf", "fmt/934", "text/plain", "3814", "2008/12/19", "extension match" ],
    [ "isilon.txt", "fmt/40", "text/plain", "9497", "2013/03/03", "extension match" ],
    [ "jerry maguire.txt", "fmt/934", "text/plain", "6741", "2008/10/16", "extension match" ],
    [ "Kate_Sweetapple_Data_poetics.pdf", "fmt/40", "text/plain", "3597", "2012/12/18", "extension match" ],
    [ "memex.txt", "fmt/934", "text/plain", "1965", "2010/03/17", "extension match" ],
    [ "Phase 1.docx", "fmt/40", "text/plain", "1581", "2012/11/27", "extension match" ],
    [ "Project Brief 2015-03C - TAB B.docx", "fmt/934", "text/plain", "3059", "2010/06/09", "extension match" ],
    [ "Recordkeeping Procurement brief memo.docx", "fmt/40", "text/plain", "1721", "2009/04/10", "extension match" ],
    [ "RFP 2015-03C-2 - DA Web Design Services.docx", "fmt/20", "text/plain", "2558", "2012/10/13", "extension match" ],
    [ "RIMPA.pptx", "Personnel Lead", "text/plain", "2290", "2012/09/26", "extension match" ],
    [ "SLNSW - ICT Strategic Plan - 2013-2017 Draft 30112012 (Combined Doc).pdf", "fmt/40", "text/plain", "1937", "2011/09/03", "extension match" ],
    [ "SRA Final Strategic Planning Workshop Agenda and PreReading - WshopUpdate.docx", "x-fmt/12", "text/plain", "6154", "2009/06/25", "extension match" ],
    [ "SRA_EA_Report_v1.1.pdf", "fmt/40", "New York", "8330", "2011/12/12", "extension match" ],
    [ "State Records Authority - Digital Records Operations recruitment.docx", "fmt/1", "text/plain", "3023", "2010/09/20", "extension match" ],
    [ "Thumbs.db", "fmt/40", "text/plain", "5797", "2009/10/09", "extension match" ],
    [ "WoG Initiatives.docx", "x-fmt/12", "text/plain", "8822", "2010/12/22", "extension match" ],
    [ "XML schemas.lnk", "fmt/40", "text/plain", "9239", "2010/11/14", "extension match" ]
    ];
    var dataCount = [['fmt/40', 11],
    ['fmt/20', 2],
    ['x-fmt/12', 2],
    ['fmt/1', 2],
    ['fmt/934', 9]
    ];

    function drawChart() {
      var dt = google.visualization.arrayToDataTable(dataCount, true);
      var options = {
        pieHole: 0.4,
        pieSliceText: 'label',
      };
      var chart = new google.visualization.PieChart(document.getElementById('piechart'));
      var searchVal = ""
      function selectHandler(e) {
        var selection = chart.getSelection();
        if (selection.length > 0) {
          searchVal = dt.getValue(selection[0].row, 0);
        }
        var table = $('#example').DataTable();
        table.search(searchVal).draw();
      }
      google.visualization.events.addListener(chart, 'select', selectHandler);
      chart.draw(dt, options);
    }

    google.charts.load('current', {'packages':['corechart']});
    google.charts.setOnLoadCallback(function() {
        $(function() {
            drawChart();
            $('#example').DataTable( {
                data: dataSet,
                columns: columnTitles.map(function(col){return {title: col}}),
                dom: 'Bfrtip',
                buttons: [
                'copy', 'csv', 'excel', 'pdf', 'print'
                ]   
            });
        })
    });