/*
data.js manages tables on the match data display page
*/

//loads match data onto the data page using ajax
function loadData() {
  var urlParams = new URLSearchParams(window.location.search);
  if (urlParams.get("display")=="team") {
    sortTeamTableBy("datasort");
  } else if (urlParams.get("display")=="match") {
    sortMatchTableBy("datasort");
  }
}

function clearTable(team) {
  while (document.getElementById(team).rows.length > 1) {
    document.getElementById(team).deleteRow(1);
  }
}

//sorts the teams table by the value of the given element
//preforms sort in backend with GET query
function sortTeamTableBy(e) {
  var selector = document.getElementById(e);
  var row = selector.options[selector.selectedIndex].value;
  clearTable("teamdata");
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
  var urlParams = new URLSearchParams(window.location.search);
  var selector = document.getElementById(e);
  var selected = encodeURIComponent(selector.options[selector.selectedIndex].value);
  var team = encodeURIComponent(urlParams.get("team"));
  clearTable("matchdata");
  var xhttp = new XMLHttpRequest();
  xhttp.responseType = "text";
  xhttp.onreadystatechange = function() {
    if (xhttp.readyState == 4 && xhttp.status == 200) {
      var txt = this.responseText.slice(20, -1);
      var rows = Papa.parse(txt).data;
      var table = document.getElementById("matchdata");
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
  xhttp.open("GET", "/matchDataGet?sortby="+selected+"&team="+team);
  xhttp.send();
}

//Shows graph based on settings inputted on the page
function showGraph() {
  var xhttp = new XMLHttpRequest();
  var select = document.getElementById("graphselect");
  var selection = encodeURIComponent(select.options[select.selectedIndex].value);
  var urlParams = new URLSearchParams(window.location.search);
  var profileTeam = urlParams.get("team");
  var url = "/getGraph?subject="+selection+"&team="+profileTeam;
  xhttp.open("GET", url, true);
  xhttp.responseType = "arraybuffer";
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 & this.status == 200) {
      var arr = new Uint8Array(this.response);
      var rawstr = String.fromCharCode.apply(null, arr);
      var b64=btoa(rawstr);
      dataurl = "data:image/jpeg;base64,"+b64;
      document.getElementById("graph").src = dataurl;
    }
  }
  xhttp.send();
}

function gotoTeamProfile(team) {
  window.location = "/data?display=teamprofile&team="+String(team);
}