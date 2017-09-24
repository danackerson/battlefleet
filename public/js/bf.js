var wsApp;

function connectServer() {
  if (wsApp.$data.ws == null || wsApp.$data.ws.readyState == WebSocket.CLOSED) {
    wsApp.$data.ws = new WebSocket('ws://' + window.location.host + '/wsInit');

    if (wsApp.$data.ws.readyState != WebSocket.CLOSED) {
      wsApp.$data.ws.addEventListener('message', function(e) {
          wsApp.serverTime = e.data;
      });

      wsApp.$data.ws.onopen = function (evt) {
        wsApp.$data.connectionState = 'OPEN';
        console.log('Socket open: Waiting for data');
      }

      wsApp.$data.ws.onerror = function(evt) {
        wsApp.$data.connectionState = 'ERROR';
        console.error('Socket encountered error: ', evt.message, 'Closing socket');
        wsApp.$data.ws.close();
      }

      wsApp.$data.ws.onclose = function (evt) {
        wsApp.$data.connectionState = 'CLOSED';
        console.log('Socket is closed.', evt.reason, ' Reconnect will be attempted.');
        setTimeout(function() {
          connectServer();
        }, 4400); // wait 4.4secs for reconnect
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
