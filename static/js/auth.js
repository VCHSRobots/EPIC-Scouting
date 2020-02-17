/*
auth.js sets and renews user login cookies
*/

function refreshLoginCookie() {

}

function getCookie(name, path="/") {
  var search = name+"=";
  var decodedCookie = decodeURIComponent(document.cookie);
  var splitCookie = decodedCookie.split(";");
  for (var i = 0; i<splitCookie.length; i++) {
    cpart = splitCookie[i];
    cname = cpart.slice(0, search.length+1).trim();
    if (search==cname) {
      return cpart.slice(search.length+1);
    }
  }
  return "";
}

function setCookie(name, value, expiry=0, path="/") {
  var d = new Date();
  d.setTime(d.getTime()+expiry);
  var expires = "expires="+d.toUTCString();
  document.cookie = name+"="+value+";"+expires+";path="+path;
}
