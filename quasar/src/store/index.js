import Vue from 'vue'
import Vuex from 'vuex'

import example from './module-example'
import VueNativeSock from 'vue-native-websocket'

Vue.use(Vuex)
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

const store = new Vuex.Store({
  modules: {
    example
  },
  state: {
    socket: {
      isConnected: false,
      message: '',
      reconnectError: false,
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
      state.message = message
      console.log('WS MSG: ', message)
    },
    // mutations for reconnect methods
    SOCKET_RECONNECT(state, count) {
      console.info('WS Reconnect...', state, count)
    },
    SOCKET_RECONNECT_ERROR(state) {
      state.socket.reconnectError = true;
    },
  }
})

export default store
