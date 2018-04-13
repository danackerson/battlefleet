import ga from 'src/battlefleet/analytics'

export default ({ router }) => {
  router.afterEach((to, from) => {
    var sessionId = $cookies.get('battlefleetID')
    if (sessionId) {
      ga.logPage(to.path, to.name, sessionId)
    } else {
      ga.logPage(to.path, to.name)
    }
  })
}
