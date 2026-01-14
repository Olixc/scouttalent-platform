import { defineStore } from 'pinia'

interface User {
  id: string
  email: string
  full_name: string
}

interface AuthState {
  user: User | null
  token: string | null
  profileId: string | null
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    token: null,
    profileId: null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    currentUser: (state) => state.user,
  },

  actions: {
    async login(email: string, password: string) {
      const config = useRuntimeConfig()
      const response = await $fetch(`${config.public.authServiceUrl}/api/v1/auth/login`, {
        method: 'POST',
        body: { email, password },
      })

      this.token = response.access_token
      this.user = response.user
      this.profileId = response.profile_id

      // Store in localStorage
      if (process.client) {
        localStorage.setItem('token', response.access_token)
        localStorage.setItem('user', JSON.stringify(response.user))
        if (response.profile_id) {
          localStorage.setItem('profileId', response.profile_id)
        }
      }

      return response
    },

    async register(email: string, password: string, fullName: string) {
      const config = useRuntimeConfig()
      const response = await $fetch(`${config.public.authServiceUrl}/api/v1/auth/register`, {
        method: 'POST',
        body: {
          email,
          password,
          full_name: fullName,
        },
      })

      return response
    },

    async fetchCurrentUser() {
      if (!this.token) return

      const config = useRuntimeConfig()
      try {
        const response = await $fetch(`${config.public.authServiceUrl}/api/v1/auth/me`, {
          headers: {
            Authorization: `Bearer ${this.token}`,
          },
        })

        this.user = response.user
        this.profileId = response.profile_id
      } catch (error) {
        console.error('Failed to fetch current user:', error)
        this.logout()
      }
    },

    logout() {
      this.user = null
      this.token = null
      this.profileId = null

      if (process.client) {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        localStorage.removeItem('profileId')
      }
    },

    initializeFromStorage() {
      if (process.client) {
        const token = localStorage.getItem('token')
        const userStr = localStorage.getItem('user')
        const profileId = localStorage.getItem('profileId')

        if (token) {
          this.token = token
          if (userStr) {
            this.user = JSON.parse(userStr)
          }
          if (profileId) {
            this.profileId = profileId
          }
        }
      }
    },
  },
})