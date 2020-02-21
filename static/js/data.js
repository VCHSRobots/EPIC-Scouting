/*
data.js manages tables on the match data display page
*/

//loads match data onto the data page using ajax
function loadData() {
  clearTable();
  sortTableBy("match");
}

function clearTable() {
  while (document.getElementById("teamdata").rows.length > 1) {
    document.getElementById("teamdata").deleteRow(1);
  }
}

//sorts the matches table by the given element
//preforms sort in backend with GET query
function sortTableBy(row) {
  clearTable();
  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    txt = this.responseText;
    rows = papa.Parse(txt);
    table = document.getElementById("teamdata");
    for (row in rows) {
      row = table.addRow();
      for (celldata in row) {
        cell = row.addCell();
        cell.InnerHTML = celldata;
      }
    }
  }
  xhttp.open("GET", "/matchdataGET?sortrow="+row);
  xhttp.send();
}