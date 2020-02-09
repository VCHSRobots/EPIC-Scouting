/*
scout.js middle man that posts data to the backend server and will show the qr codes as failsafes if unable to do so
*/

var connected;

function submitMatchData(form) {
  //This doesn't seem to return right
  alert(form.team.value);
  //Try to post the data to the server
  //If that fails, prepare QR Code
}

function submitPitData(form) {
  //Try to post the data to the server
  //If that fails, prepare QR Code
}

//displays qr code based on whether the post succeded or not
function handleData() {

}
