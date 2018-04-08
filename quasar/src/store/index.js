import Vue from 'vue'
import Vuex from 'vuex'

import example from './module-example'
import VueNativeSock from 'vue-native-websocket'

Vue.use(Vuex)
Vue.use(require('vue-moment')); // date formatting lib

const store = new Vuex.Store({
  modules: {
    example
  },
  state: {
    count: 0,
    socket: {
      isConnected: false,
      message: 'Space',
      reconnectError: false,
    }
  },
  actions: {
    currentServerTime (context, content) {
      console.log('RCVD currentServerTime: ' + JSON.stringify(content))
      this.state.socket.message = content
      /* interesting code todo something on initial WS connect?
      this.$store.subscribe((mutation, state) => {
        if (mutation.type === 'SOCKET_ONOPEN' ) {
          // your code here
        }
      }*/
    }
  },
  mutations:{
    SOCKET_ONOPEN (state, event)  {
      state.socket.isConnected = true
      console.log('WS Connected')
    },
    SOCKET_ONCLOSE (state, event)  {
      state.socket.isConnected = false
      console.log('WS Disconnected')
    },
    SOCKET_ONERROR (state, event)  {
      console.error('WS ERROR: ', state, event)
    },
    // default handler called for all methods
    SOCKET_ONMESSAGE (state, message)  {
      state.socket.message = message
      console.log('WS MSG: ', message)
    },
    // mutations for reconnect methods
    SOCKET_RECONNECT (state, count) {
      console.info(state, count)
    },
    SOCKET_RECONNECT_ERROR (state) {
      state.socket.reconnectError = true;
    }
  }
})

var hostURL = "wss://"
if (location.protocol == "http:") {
  hostURL = "ws://"
}
hostURL += location.hostname + ':' + process.env.PORT
Vue.use(VueNativeSock, hostURL + '/wsInit', {
  store: store,
  connectManually: true,
  format: 'json',
  reconnection: true,
  reconnectionAttempts: 5,
  reconnectionDelay: 3000,
})

export default store
