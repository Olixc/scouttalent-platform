<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Create your account
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Already have an account?
          <NuxtLink to="/login" class="font-medium text-primary-600 hover:text-primary-500">
            Sign in
          </NuxtLink>
        </p>
      </div>

      <form class="mt-8 space-y-6" @submit.prevent="handleRegister">
        <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
          {{ error }}
        </div>

        <div v-if="success" class="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg">
          Account created successfully! Redirecting to login...
        </div>

        <div class="rounded-md shadow-sm space-y-4">
          <div>
            <label for="full-name" class="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
            <input
              id="full-name"
              v-model="form.fullName"
              name="full-name"
              type="text"
              required
              class="input"
              placeholder="John Doe"
            />
          </div>
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-1">Email address</label>
            <input
              id="email"
              v-model="form.email"
              name="email"
              type="email"
              autocomplete="email"
              required
              class="input"
              placeholder="you@example.com"
            />
          </div>
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-1">Password</label>
            <input
              id="password"
              v-model="form.password"
              name="password"
              type="password"
              autocomplete="new-password"
              required
              class="input"
              placeholder="••••••••"
            />
            <p class="mt-1 text-xs text-gray-500">Must be at least 8 characters</p>
          </div>
          <div>
            <label for="confirm-password" class="block text-sm font-medium text-gray-700 mb-1">Confirm Password</label>
            <input
              id="confirm-password"
              v-model="form.confirmPassword"
              name="confirm-password"
              type="password"
              autocomplete="new-password"
              required
              class="input"
              placeholder="••••••••"
            />
          </div>
        </div>

        <div>
          <button
            type="submit"
            :disabled="loading"
            class="w-full btn btn-primary text-lg py-3"
          >
            <span v-if="loading">Creating account...</span>
            <span v-else>Create account</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  layout: false,
})

const authStore = useAuthStore()
const router = useRouter()

const form = ref({
  fullName: '',
  email: '',
  password: '',
  confirmPassword: '',
})

const loading = ref(false)
const error = ref('')
const success = ref(false)

const handleRegister = async () => {
  error.value = ''

  if (form.value.password !== form.value.confirmPassword) {
    error.value = 'Passwords do not match'
    return
  }

  if (form.value.password.length < 8) {
    error.value = 'Password must be at least 8 characters'
    return
  }

  loading.value = true

  try {
    await authStore.register(form.value.email, form.value.password, form.value.fullName)
    success.value = true
    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } catch (err: any) {
    error.value = err.data?.error || 'Failed to create account. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>