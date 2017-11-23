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

function disconnectServer() {
  userDisconnect = true;
  wsApp.$data.ws.close();
  wsApp.$data.connectionState = 'CLOSING';
}

function connectServer() {
  if (wsApp.$data.ws == null || wsApp.$data.ws.readyState == WebSocket.CLOSED) {
    var scheme = "wss:";
    if (location.protocol == "http:") {
      scheme = "ws:";
    }
    wsApp.$data.ws = new WebSocket(scheme + '//' + window.location.host + '/wsInit');

    if (wsApp.$data.ws.readyState != WebSocket.CLOSED) {
      wsApp.$data.ws.addEventListener('message', function(e) {
          wsApp.serverTime = e.data;
          wsApp.$data.ws.send("ACK");
      });

      wsApp.$data.ws.onopen = function (evt) {
        wsApp.$data.connectionState = 'OPEN';
        wsApp.$data.ws.send("OPEN");
        console.log('Socket open: Waiting for data');
        userDisconnect = false;
      }

      wsApp.$data.ws.onerror = function(evt) {
        wsApp.$data.connectionState = 'ERROR';
        wsApp.$data.ws.send("ERR");
        console.error('Socket encountered error: ', evt.message, 'Closing socket');
        wsApp.$data.ws.close();
      }

      wsApp.$data.ws.onclose = function (evt) {
        wsApp.$data.connectionState = 'CLOSED';
        if (!userDisconnect) {
          console.log('Socket is closed.', evt.reason, ' Reconnect will be attempted.');
          setTimeout(function() {
            connectServer();
          }, 4400); // wait 4.4secs for reconnect
        }
      }
    } else {
      wsApp.$data.ws = null;
    }
  } else {
    console.log('Socket already exists.')
  }
}

window.onload = function(){
  wsApp = new Vue({
    el: '#app',
    delimiters: ['${', '}'],
    data: {
      ws: null,
      connectionState: 'INITIAL',
      serverTime: '',
    }
  });
}
