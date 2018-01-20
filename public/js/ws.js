var wsApp;

var userDisconnect = false;

function disconnectServer(reason) {
  userDisconnect = true;
  wsApp.$data.ws.close();
  wsApp.$data.connectionState = 'CLOSING';
  if (reason == "idle") {
    console.log('Idle disconnect from client-side (' + idleTimeout + ' secs with no action exceeded)...');
  }
}

var idleTimeout = 60000; // 60 secs
var idleTimeoutHandle = resetIdleTimeout();
function resetIdleTimeout() {
  window.clearTimeout(idleTimeoutHandle);

  idleTimeoutHandle = window.setTimeout(function() {
    disconnectServer("idle");
  }, idleTimeout);  // wait idle seconds to auto-disconnect on idle game

  return idleTimeoutHandle;
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
        if (e.data != null) {
          wsApp.game = JSON.parse(e.data);
          renderShips(wsApp.game["Ships"]);
          resetIdleTimeout();
        }
      });

      wsApp.$data.ws.onopen = function (evt) {
        wsApp.$data.connectionState = 'OPEN';
        sendWebSocketMessage(JSON.stringify({'cmd': 'OPEN'}));
        console.log('Socket open: Waiting for data');

        userDisconnect = false;
        resetIdleTimeout();
      }

      wsApp.$data.ws.onerror = function(evt) {
        wsApp.$data.connectionState = 'ERROR';
        sendWebSocketMessage(JSON.stringify({'cmd': 'ERR'}));
        console.error('Socket encountered error: ', evt.message, 'Closing socket');
        wsApp.$data.ws.close();

        window.clearTimeout(idleTimeoutHandle);
      }

      wsApp.$data.ws.onclose = function (evt) {
        wsApp.$data.connectionState = 'CLOSED';
        clearBoard();
        window.clearTimeout(idleTimeoutHandle);

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

// Hexagon Helpers
// https://www.redblobgames.com/grids/hexagons/
// https://github.com/mpalmerlee/HexagonTools

var sprites = [];
// cleanup board in prep for next reconnect/rendering
function clearBoard() {
  for (i=0; i < sprites.length; i++) {
    sprites[i].disable();
    sprites.splice(i,1);
  }

  movingPiece = null;
  //console.log("board cleared...");
}

function renderShips(ships) {
  // populate the board
  for (i=0; i < ships.length; i++) {
    var spriteConfig = {
      container: board.group,
      url: '/images/ships/'+ships[i]['Class']+'_'+ships[i]['Type']+'.png',
      scale: 10,
      heightOffset: 6
    };

    // render ship image
    sprites[i] = new Sprite(spriteConfig);
    sprites[i].activate();
    sprites[i].uniqueId = ships[i]['ID'];
    // determine position of ship on
    q = ships[i]['Position']['X'];
    r = ships[i]['Position']['Y'];
    s = -q - r;
    shipGridCoordinates = q + '.' + r + '.' + s;
    cell = board.grid.cells[shipGridCoordinates];
    tile = board.getTileAtCell(cell);

    // place ship on proper Tile
    board.setEntityOnTile(sprites[i], tile);

    // help out the Newbies
    if (shipGridCoordinates == "0.0.0") {
      //alert("Select your ship and move to another hexagon");
    	scene.add(arrow);
    } else {
      scene.remove(arrow);
    }
  }
}

// keep track of states
var movingPiece = null;
function moveEntityToCell(movingPiece, tile) {
  origin = "(" + movingPiece.tile.cell.q + "," + movingPiece.tile.cell.r + ")";
  destination = "(" + tile.cell.q + "," + tile.cell.r + ")";

  var ship = new Object();
  ship.ID = movingPiece.uniqueId;
  ship.origin = {Q: movingPiece.tile.cell.q, R: movingPiece.tile.cell.r};
  ship.destination = {Q: tile.cell.q, R: tile.cell.r};

  var msg = new Object();
  msg.cmd = "MOV";
  msg.payload = ship;

  sendWebSocketMessage(JSON.stringify(msg));

  board.setEntityOnTile(movingPiece, tile);

  clearBoard(); // response from server will redraw the board
}

function sendWebSocketMessage(jsonString) {
  wsApp.$data.ws.send(jsonString);
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
  connectServer();
}

// formatter/format components for rendering pretty DATETIME
// via <format :value="game.LastTurn" fn="date" /> in game.tmpl file
// moment is a JS library for pretty DATETIMEs: https://momentjs.com/
var formatter = {
    date: function (value, format) {
        if (value) {
            return moment(String(value)).format(format || 'DD.MM.YY HH:mm:ss')
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
