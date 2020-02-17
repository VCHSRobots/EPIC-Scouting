function showGraph(graphName, getPath) {
  xhttp = new XMLHttpRequest();
  xhttp.open("GET", getPath, true);
  xhttp.responseType = "arraybuffer";
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 & this.status == 200) {
      var arr = new Uint8Array(this.response);
      var rawstr = String.fromCharCode.apply(null, arr);
      var b64=btoa(rawstr);
      dataurl = "data:image/jpeg;base64,"+b64;
      document.getElementById(graphName).src = dataurl;
    }
    
  }
  xhttp.send();
}