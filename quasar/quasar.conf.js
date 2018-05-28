// Configuration for your app

module.exports = function (ctx) {
  console.log(ctx)

  return {
    // app plugins (/src/plugins)
    plugins: [
      'axios',
      'google-analytics'
    ],
    css: [
      'app.styl'
    ],
    extras: [
      //ctx.theme.mat ? 'roboto-font' : null,
      'material-icons'
      // 'ionicons',
      // 'mdi',
      // 'fontawesome'
    ],
    supportIE: false,
    build: {
      env: ctx.dev ? {
        PORT: 8083, // Websocket & Axios port to Golang server
        AUTH0_CLIENT_ID: '"kvOcNm3klMGSxTfzD5mvg23C7vgcYvij"',
        AUTH0_CALLBACK_URL: '"http://localhost:8443/callback"'
      }
      :{
        PORT: 443,
        AUTH0_CLIENT_ID: '"' + process.env.AUTH0_CLIENT_ID + '"',
        AUTH0_CALLBACK_URL: '"' + process.env.AUTH0_CALLBACK_URL + '"'
      },
      scopeHoisting: true,
      vueRouterMode: 'history',
      publicPath: '/',
      // gzip: true,
      // analyze: true,
      // extractCSS: false,
      // useNotifier: false,
      extendWebpack (cfg) {
        cfg.module.rules.push({
          enforce: 'pre',
          test: /\.(js|vue)$/,
          loader: 'eslint-loader',
          exclude: /(node_modules|quasar)/
        })
      }
    },
    devServer: {
      https: false,
      port: 8443,
      open: true // opens browser window automatically
    },
    //devtool: '#eval-source-map',
    // framework: 'all' --- includes everything; for dev only!
    framework: {
      components: [
        'QCollapsible',
        'QField',
        'QInput',
        'QLayout',
        'QLayoutHeader',
        'QLayoutFooter',
        'QLayoutDrawer',
        'QList',
        'QPage',
        'QPageContainer',
        'QToolbar',
        'QToolbarTitle',
        'QBtn',
        'QIcon'
      ],
      directives: [
        'Ripple'
      ],
      // Quasar plugins
      plugins: [
        'Notify'
      ]
    },
    // animations: 'all' --- includes all animations
    animations: [
    ],
    pwa: {
      cacheExt: 'js,html,css,ttf,eot,otf,woff,woff2,json,svg,gif,jpg,jpeg,png,wav,ogg,webm,flac,aac,mp4,mp3',
      manifest: {
        // name: 'Quasar App',
        // short_name: 'Quasar-PWA',
        // description: 'Best PWA App in town!',
        display: 'standalone',
        orientation: 'portrait',
        background_color: '#ffffff',
        theme_color: '#027be3',
        icons: [
          {
            'src': 'statics/icons/icon-128x128.png',
            'sizes': '128x128',
            'type': 'image/png'
          },
          {
            'src': 'statics/icons/icon-192x192.png',
            'sizes': '192x192',
            'type': 'image/png'
          },
          {
            'src': 'statics/icons/icon-256x256.png',
            'sizes': '256x256',
            'type': 'image/png'
          },
          {
            'src': 'statics/icons/icon-384x384.png',
            'sizes': '384x384',
            'type': 'image/png'
          },
          {
            'src': 'statics/icons/icon-512x512.png',
            'sizes': '512x512',
            'type': 'image/png'
          }
        ]
      }
    },
    cordova: {
      // id: 'org.cordova.quasar.app'
    },
    electron: {
      extendWebpack (cfg) {
        // do something with cfg
      },
      packager: {
        // OS X / Mac App Store
        // appBundleId: '',
        // appCategoryType: '',
        // osxSign: '',
        // protocol: 'myapp://path',

        // Window only
        // win32metadata: { ... }
      }
    },

    // leave this here for Quasar CLI
    starterKit: '1.0.2'
  }
}
