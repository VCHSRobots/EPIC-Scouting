/*
db.js calls for backend database querys
*/

function listCampaigns() {
  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (this.readyStae == 4 && this.status == 200) {
      document.getElementById("campaigns").innerHTML = this.responseText;
    }
  };
  xhttp.open("GET", "http://localhost:443/listcampaigns")
}

function getDatabaseSize() {
  //currently sends a dummy value for testing purposes
  
}
