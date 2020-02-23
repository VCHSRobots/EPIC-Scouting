/*
dash.js handles dashboard settings such as team cookies
*/

function setTeamCookie() {
  var xhttp = new XMLHttpRequest();
  xhttp.open("GET", "/setTeamCookie?team="+document.getElementById("teamselect").value);
  xhttp.send();
}

setTeamCookie();
