import MainLayout from 'layouts/default'
import MainPage from 'pages/index'
import LeftPanel from 'pages/left'
import RightPanel from 'pages/right'

const routes = [
  {
    path: '/',
    component: MainLayout,
    children: [
      { path: '',
        components: {
          main: MainPage,
          left: LeftPanel,
          right: RightPanel
        }
      }
    ]
  }
]

// Always leave this as last one
routes.push({ path: '*', component: () => import('pages/404') })

export default routes
