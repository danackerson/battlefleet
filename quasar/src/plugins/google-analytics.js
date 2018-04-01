import ga from 'src/battlefleet/analytics'

export default ({ router }) => {
  router.afterEach((to, from) => {
    if (typeof sessionId == 'undefined') {
      ga.logPage(to.path, to.name)
    } else {
      ga.logPage(to.path, to.name, sessionId)
    }
  })
}
