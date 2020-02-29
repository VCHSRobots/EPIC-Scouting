/*
data.js manages tables on the match data display page
*/

const urlParams = new URLSearchParams(window.location.search);

//loads match data onto the data page using ajax
function loadData() {
  clearTable();
  if (urlParams.get("display")=="team") {
    sortTeamTableBy("datasort");
  } else if (urlParams.get("display")=="match") {
    sortMatchTableBy("datasort")
  } 
}

function clearTable() {
  while (document.getElementById("teamdata").rows.length > 1) {
    document.getElementById("teamdata").deleteRow(1);
  }
}

//sorts the teams table by the value of the given element
//preforms sort in backend with GET query
function sortTeamTableBy(e) {
  var selector = document.getElementById(e);
  var row = selector.options[selector.selectedIndex].value;
  clearTable();
  var xhttp = new XMLHttpRequest();
  xhttp.responseType = "text";
  xhttp.onreadystatechange = function() {
    if (xhttp.readyState == 4 && xhttp.status == 200) {
      var txt = this.responseText.slice(20, -2);
      var rows = Papa.parse(txt).data;
      var table = document.getElementById("teamdata");
      for (row in rows) {
        tablerow = table.insertRow();
        for (celldata in rows[row]) {
          innerData = rows[row][celldata];
          var cell = tablerow.insertCell();
          cell.innerHTML = innerData;
        }
      }
    }
  }
  xhttp.open("GET", "/teamDataGet?sortby="+row);
  xhttp.send();
}

function sortMatchTableBy(e) {
  var selector = document.getElementById(e);
  var row = selector.options[selector.selectedIndex].value;
  clearTable();
  var xhttp = new XMLHttpRequest();
  xhttp.responseType = "text";
  xhttp.onreadystatechange = function() {
    if (xhttp.readyState == 4 && xhttp.status == 200) {
      var txt = this.responseText.slice(20, -2);
      var rows = Papa.parse(txt).data;
      var table = document.getElementById("teamdata");
      for (row in rows) {
        tablerow = table.insertRow();
        for (celldata in rows[row]) {
          innerData = rows[row][celldata];
          var cell = tablerow.insertCell();
          cell.innerHTML = innerData;
        }
      }
    }
  }
  xhttp.open("GET", "/matchDataGet?sortby="+row);
  xhttp.send();
}

function gotoTeamProfile(team) {
  window.location = "/data?display=teamprofile&team="+String(team);
}