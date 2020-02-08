/*
qr.js backup system for data transfer if the client cannot connect to the server
*/

//this is the trouble line
//import {QRious} from "js/qrious.js";

var connected = false;
var parser = new DOMParser();
var serializer = new XMLSerializer();

//creates a qr of the scouting fields should the connection be lost
function makeQrCode(dataString) {
  //creates qr code
  var qr = new QRious();
  alert(dataString);
  qr.set({
    element: document.querySelector("img"),
    value: dataString
  });
  alert("here");
  qr.canvas.parentNode.appendChild(qr.image);
}

//checks if the client can connect to the server
function checkConnection() {
  connected = false;
  var xhttp = new XMLHttpRequest(); 
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
      connected = true;
    }
  }
  //TODO: put the real domain name here
  xhttp.open("GET", "http://localhost:443/", true);
  xhttp.send();
  //TODO: setText is the function that handles events after we know if we're connected or not.
  setTimeout(setText, 100, xhttp);
}

function setText(xhttp) {
  xhttp.abort();
  if (connected) {
    document.getElementById("test").textContent = "connected";
  } else {
    document.getElementById("test").textContent = "not connected";
  }
}
