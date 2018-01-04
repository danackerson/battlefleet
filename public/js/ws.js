var wsApp;

var userDisconnect = false;

// see https://github.com/mpalmerlee/HexagonTools
// https://www.redblobgames.com/grids/hexagons/
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
          wsApp.game = JSON.parse(e.data); // MUST set wsApp.game for vue actions in game.tmpl
          bootstrapGameData(wsApp.game);
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

function bootstrapGameData(game) {
  //alert(game["LastTurn"]);
  // TODO: render game["Ships"]
  renderShips(game["Ships"]);
}

window.onload = function(){
  wsApp = new Vue({
    el: '#app',
    delimiters: ['${', '}'],
    data: {
      ws: null,
      connectionState: 'INITIAL',
      game: 'Loading...',
    }
  });
}

var formatter = {
    date: function (value, format) {
        if (value) {
            return moment(String(value)).format(format || 'DD.MM.YY hh:mm:ss')
        }
    }
};

Vue.component('format', {
    template: `<span>{{ formatter[fn](value, format) }}</span>`,
    props: ['value', 'fn', 'format'],
    computed: {
        formatter() {
            return formatter;
        }
    }
});
