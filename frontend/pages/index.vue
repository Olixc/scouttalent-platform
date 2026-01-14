<template>
  <div>
    <!-- Hero Section -->
    <section class="bg-gradient-to-br from-primary-600 to-primary-800 text-white py-20">
      <div class="container mx-auto px-4">
        <div class="max-w-4xl mx-auto text-center">
          <h1 class="text-5xl md:text-6xl font-bold mb-6">
            Showcase Your Football Talent to the World
          </h1>
          <p class="text-xl md:text-2xl mb-8 text-primary-100">
            Connect with scouts, clubs, and opportunities through AI-powered video showcases
          </p>
          <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <NuxtLink to="/register" class="btn bg-white text-primary-600 hover:bg-gray-100 text-lg px-8 py-3">
              Get Started
            </NuxtLink>
            <NuxtLink to="/discover" class="btn border-2 border-white text-white hover:bg-white hover:text-primary-600 text-lg px-8 py-3">
              Discover Talent
            </NuxtLink>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="py-20 bg-gray-50">
      <div class="container mx-auto px-4">
        <h2 class="text-4xl font-bold text-center mb-12">Why ScoutTalent?</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div class="card text-center">
            <div class="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </div>
            <h3 class="text-xl font-bold mb-2">Video Showcases</h3>
            <p class="text-gray-600">Upload and share your best football moments with professional scouts and clubs</p>
          </div>

          <div class="card text-center">
            <div class="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
              </svg>
            </div>
            <h3 class="text-xl font-bold mb-2">AI-Powered Matching</h3>
            <p class="text-gray-600">Get discovered by the right scouts based on your skills and playing style</p>
          </div>

          <div class="card text-center">
            <div class="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
              </svg>
            </div>
            <h3 class="text-xl font-bold mb-2">Global Network</h3>
            <p class="text-gray-600">Connect with scouts, clubs, and players from around the world</p>
          </div>
        </div>
      </div>
    </section>

    <!-- How It Works -->
    <section class="py-20">
      <div class="container mx-auto px-4">
        <h2 class="text-4xl font-bold text-center mb-12">How It Works</h2>
        <div class="max-w-4xl mx-auto">
          <div class="space-y-8">
            <div class="flex items-start gap-6">
              <div class="flex-shrink-0 w-12 h-12 bg-primary-600 text-white rounded-full flex items-center justify-center font-bold text-xl">
                1
              </div>
              <div>
                <h3 class="text-2xl font-bold mb-2">Create Your Profile</h3>
                <p class="text-gray-600">Sign up and build your player profile with your details, position, and stats</p>
              </div>
            </div>

            <div class="flex items-start gap-6">
              <div class="flex-shrink-0 w-12 h-12 bg-primary-600 text-white rounded-full flex items-center justify-center font-bold text-xl">
                2
              </div>
              <div>
                <h3 class="text-2xl font-bold mb-2">Upload Your Videos</h3>
                <p class="text-gray-600">Share your best football moments - skills, goals, assists, and match highlights</p>
              </div>
            </div>

            <div class="flex items-start gap-6">
              <div class="flex-shrink-0 w-12 h-12 bg-primary-600 text-white rounded-full flex items-center justify-center font-bold text-xl">
                3
              </div>
              <div>
                <h3 class="text-2xl font-bold mb-2">Get Discovered</h3>
                <p class="text-gray-600">Our AI matches you with relevant scouts and clubs based on your profile and videos</p>
              </div>
            </div>

            <div class="flex items-start gap-6">
              <div class="flex-shrink-0 w-12 h-12 bg-primary-600 text-white rounded-full flex items-center justify-center font-bold text-xl">
                4
              </div>
              <div>
                <h3 class="text-2xl font-bold mb-2">Connect & Grow</h3>
                <p class="text-gray-600">Receive opportunities, connect with scouts, and take your career to the next level</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="bg-primary-600 text-white py-16">
      <div class="container mx-auto px-4 text-center">
        <h2 class="text-4xl font-bold mb-4">Ready to Start Your Journey?</h2>
        <p class="text-xl mb-8 text-primary-100">Join thousands of players already showcasing their talent</p>
        <NuxtLink to="/register" class="btn bg-white text-primary-600 hover:bg-gray-100 text-lg px-8 py-3 inline-block">
          Sign Up Now
        </NuxtLink>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: 'Home',
})
</script>