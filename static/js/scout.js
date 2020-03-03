/*
scout.js middle man that posts data to the backend server and will show the qr codes as failsafes if unable to do so
*/

var connected;

function submitMatchData(form) {
  //Parse data to CSV
  //TODO: Add session id verification?
  var data = [form.match.value, form.team.value, form.alliance.value, form.autoLineCross.value, form.autoLowBalls.value, form.autoHighBalls.value, form.autoBackBalls.value, form.autoShots.value, form.autoPickups.value, form.shotQuantity.value, form.lowFuel.value, form.highFuel.value, form.backFuel.value, form.stageOneComplete.value, form.stageOneTime.value, form.stageTwoComplete.value, form.stageTwoTime.value, form.fouls.value, form.techFouls.value, form.cards.value, form.climbed.value, form.balanced.value, form.climbTime.value, form.comments.value];
  var csvstring = data.join();
  //Try to post the data to the server
  checkConnection();
  if (connected) {
    var jsonstring = JSON.stringify(Papa.parse(csvstring));
    xhttp = new XMLHttpRequest();
    xhttp.open("POST", "/matchPOST", true);
    xhttp.send(jsonstring);
  } //else {
    //If that fails, prepare QR Code
    makeQrCode(csvstring);
  //}
}

function submitPitData(form) {
  //Try to post the data to the server
  //If that fails, prepare QR Code
}
