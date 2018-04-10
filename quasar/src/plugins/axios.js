import axios from 'axios'

export default ({ Vue }) => {
  axios.defaults.withCredentials = true
  Vue.prototype.$axios = axios
}
