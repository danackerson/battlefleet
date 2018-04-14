import ga from 'src/battlefleet/analytics'

export default ({ router }) => {
  router.afterEach((to, from) => {
    /*if (sessionId) {
      ga.logPage(to.path, to.name, sessionId)
    } else {*/
      ga.logPage(to.path, to.name)
    //}
  })
}
