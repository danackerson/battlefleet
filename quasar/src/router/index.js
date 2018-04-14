import Vue from 'vue'
import VueRouter from 'vue-router'
import auth from 'src/battlefleet/auth'
import store from 'src/store'
import routes from './routes'

Vue.use(auth)
Vue.use(VueRouter)

const Router = new VueRouter({
  /*
   * NOTE! Change Vue Router mode from quasar.conf.js -> build -> vueRouterMode
   *
   * If you decide to go with "history" mode, please also set "build.publicPath"
   * to something other than an empty string.
   * Example: '/' instead of ''
   */

  // Leave as is and change from quasar.conf.js instead!
  mode: process.env.VUE_ROUTER_MODE,
  base: process.env.VUE_ROUTER_BASE,
  scrollBehavior: () => ({ y: 0 }),
  routes
})

Router.beforeEach((to, from, next) => {
  /*
  if (to.meta) {
    store.commit('module-example/index', to.meta)
  }
  */

  /*
  // Inform Google Analytics
  if (ga !== undefined) {
    ga('set', 'page', to.path)
    ga('send', 'pageview')
  }
  */

  next()
})

export default Router
