<template>
  <footer class="bg-gray-900 text-gray-300">
    <div class="container mx-auto px-4 py-12">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
        <!-- About -->
        <div>
          <h3 class="text-white font-bold text-lg mb-4">ScoutTalent</h3>
          <p class="text-sm">
            Connecting football talent with opportunities through AI-powered video showcases.
          </p>
        </div>

        <!-- Quick Links -->
        <div>
          <h4 class="text-white font-semibold mb-4">Quick Links</h4>
          <ul class="space-y-2 text-sm">
            <li><NuxtLink to="/discover" class="hover:text-white transition-colors">Discover Players</NuxtLink></li>
            <li><NuxtLink to="/search" class="hover:text-white transition-colors">Search</NuxtLink></li>
            <li><NuxtLink to="/feed" class="hover:text-white transition-colors">Video Feed</NuxtLink></li>
          </ul>
        </div>

        <!-- For Users -->
        <div>
          <h4 class="text-white font-semibold mb-4">For Users</h4>
          <ul class="space-y-2 text-sm">
            <li><NuxtLink to="/register" class="hover:text-white transition-colors">Sign Up</NuxtLink></li>
            <li><NuxtLink to="/login" class="hover:text-white transition-colors">Login</NuxtLink></li>
            <li><NuxtLink to="/upload" class="hover:text-white transition-colors">Upload Video</NuxtLink></li>
          </ul>
        </div>

        <!-- Legal -->
        <div>
          <h4 class="text-white font-semibold mb-4">Legal</h4>
          <ul class="space-y-2 text-sm">
            <li><a href="#" class="hover:text-white transition-colors">Privacy Policy</a></li>
            <li><a href="#" class="hover:text-white transition-colors">Terms of Service</a></li>
            <li><a href="#" class="hover:text-white transition-colors">Contact Us</a></li>
          </ul>
        </div>
      </div>

      <div class="border-t border-gray-800 mt-8 pt-8 text-center text-sm">
        <p>&copy; {{ new Date().getFullYear() }} ScoutTalent. All rights reserved.</p>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
</script>