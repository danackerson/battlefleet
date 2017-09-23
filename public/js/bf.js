window.onload = function(){
  var wsApp = new Vue({
    el: '#app',
    delimiters: ['${', '}'],
    data: {
      ws: null,
      connectionState: 'INITIAL',
      serverTime: '',
    },
    created () {
      var self = this;
      this.ws = new WebSocket('ws://' + window.location.host + '/wsInit');
      this.ws.addEventListener('message', function(e) {
          self.serverTime = e.data;
      });
      this.ws.onopen = function (evt) {
        self.connectionState = 'OPEN';
      }
      this.ws.onerror = function(evt) {
        self.connectionState = 'ERROR';
      }
      this.ws.onclose = function (evt) {
        self.connectionState = 'CLOSED';
      }
    }
  });


}
