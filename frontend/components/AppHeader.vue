<template>
  <header class="bg-white border-b border-gray-200 sticky top-0 z-50">
    <nav class="container mx-auto px-4 py-4">
      <div class="flex items-center justify-between">
        <!-- Logo -->
        <NuxtLink to="/" class="flex items-center space-x-2">
          <div class="w-10 h-10 bg-primary-600 rounded-lg flex items-center justify-center">
            <span class="text-white font-bold text-xl">ST</span>
          </div>
          <span class="text-xl font-bold text-gray-900">ScoutTalent</span>
        </NuxtLink>

        <!-- Desktop Navigation -->
        <div class="hidden md:flex items-center space-x-8">
          <NuxtLink to="/discover" class="text-gray-600 hover:text-primary-600 font-medium transition-colors">
            Discover
          </NuxtLink>
          <NuxtLink to="/search" class="text-gray-600 hover:text-primary-600 font-medium transition-colors">
            Search
          </NuxtLink>
          <NuxtLink to="/feed" class="text-gray-600 hover:text-primary-600 font-medium transition-colors">
            Feed
          </NuxtLink>
        </div>

        <!-- Auth Buttons -->
        <div class="hidden md:flex items-center space-x-4">
          <template v-if="authStore.isAuthenticated">
            <NuxtLink to="/profile" class="text-gray-600 hover:text-primary-600 font-medium transition-colors">
              Profile
            </NuxtLink>
            <NuxtLink to="/upload" class="btn btn-primary">
              Upload Video
            </NuxtLink>
            <button @click="handleLogout" class="text-gray-600 hover:text-red-600 font-medium transition-colors">
              Logout
            </button>
          </template>
          <template v-else>
            <NuxtLink to="/login" class="text-gray-600 hover:text-primary-600 font-medium transition-colors">
              Login
            </NuxtLink>
            <NuxtLink to="/register" class="btn btn-primary">
              Sign Up
            </NuxtLink>
          </template>
        </div>

        <!-- Mobile Menu Button -->
        <button @click="mobileMenuOpen = !mobileMenuOpen" class="md:hidden">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path v-if="!mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Mobile Menu -->
      <div v-if="mobileMenuOpen" class="md:hidden mt-4 pb-4 border-t border-gray-200 pt-4">
        <div class="flex flex-col space-y-4">
          <NuxtLink to="/discover" class="text-gray-600 hover:text-primary-600 font-medium">
            Discover
          </NuxtLink>
          <NuxtLink to="/search" class="text-gray-600 hover:text-primary-600 font-medium">
            Search
          </NuxtLink>
          <NuxtLink to="/feed" class="text-gray-600 hover:text-primary-600 font-medium">
            Feed
          </NuxtLink>
          <template v-if="authStore.isAuthenticated">
            <NuxtLink to="/profile" class="text-gray-600 hover:text-primary-600 font-medium">
              Profile
            </NuxtLink>
            <NuxtLink to="/upload" class="btn btn-primary">
              Upload Video
            </NuxtLink>
            <button @click="handleLogout" class="text-left text-gray-600 hover:text-red-600 font-medium">
              Logout
            </button>
          </template>
          <template v-else>
            <NuxtLink to="/login" class="text-gray-600 hover:text-primary-600 font-medium">
              Login
            </NuxtLink>
            <NuxtLink to="/register" class="btn btn-primary">
              Sign Up
            </NuxtLink>
          </template>
        </div>
      </div>
    </nav>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '~/stores/auth'

const authStore = useAuthStore()
const mobileMenuOpen = ref(false)

const handleLogout = () => {
  authStore.logout()
  navigateTo('/login')
}
</script>