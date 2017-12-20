var wsApp;

var userDisconnect = false;

function loginAuth0() {
  var webAuth = new auth0.WebAuth({
    domain: AUTH0_DOMAIN,
    clientID: AUTH0_CLIENT_ID,
    redirectUri: AUTH0_CALLBACK_URL,
    audience: `https://${AUTH0_DOMAIN}/userinfo`,
    responseType: 'code',
    scope: 'openid profile'
  });

  webAuth.authorize();
}

function logoutAuth0(save_games) {
  if (save_games > 3) {
    alert("You can only save 3 games. Please delete " + (save_games - 3) + "...");
  } else {
    window.location.href = "/account/?action=logout";
  }
}

function confirmAccountDeletion(cmdrName) {
  var confirm = prompt("Permanently DELETE your account and all games?", "Retype your Commander Name to confirm...");
  if (confirm == cmdrName) {
    window.location.href = "/account/?action=delete";
  } else if (confirm != null){
    alert(confirm + " is NOT " + cmdrName + ". Try again!");
  }
}

function confirmGameDeletion(gameID) {
  var confirmed = confirm("Permanently DELETE your game?");
  if (confirmed) {
    window.location.href = "/games/" + gameID + "?action=delete";
  }
}
