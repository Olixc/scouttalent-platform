// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },
  
  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
    '@vueuse/nuxt',
    '@nuxt/image',
  ],

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080',
      authServiceUrl: process.env.NUXT_PUBLIC_AUTH_SERVICE_URL || 'http://localhost:8080',
      profileServiceUrl: process.env.NUXT_PUBLIC_PROFILE_SERVICE_URL || 'http://localhost:8081',
      mediaServiceUrl: process.env.NUXT_PUBLIC_MEDIA_SERVICE_URL || 'http://localhost:8082',
      discoveryServiceUrl: process.env.NUXT_PUBLIC_DISCOVERY_SERVICE_URL || 'http://localhost:8083',
    }
  },

  app: {
    head: {
      title: 'ScoutTalent - Football Talent Marketplace',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Connect players, scouts, and clubs through AI-powered video showcases' },
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }
      ]
    }
  },

  css: ['~/assets/css/main.css'],

  tailwindcss: {
    cssPath: '~/assets/css/main.css',
    configPath: 'tailwind.config.js',
    exposeConfig: false,
    viewer: true,
  },

  typescript: {
    strict: true,
    typeCheck: true
  }
})