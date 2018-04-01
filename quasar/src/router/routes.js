import MainLayout from 'layouts/default'
import MainPage from 'pages/index'
import LeftPanel from 'pages/left'
import RightPanel from 'pages/right'
import Callback from 'components/callback'

const routes = [
  {
    path: '/',
    component: MainLayout,
    children: [
      { path: '',
        name: 'home',
        components: {
          main: MainPage,
          left: LeftPanel,
          right: RightPanel
        }
      }
    ]
  },
  {
    path: '/callback',
    name: 'callback',
    component: Callback
  }
]

// Always leave this as last one
routes.push({ path: '*', component: () => import('pages/404') })

export default routes
