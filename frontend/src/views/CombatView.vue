<script setup lang="ts">
import { useRouter } from 'vue-router'
const router = useRouter()
</script>

<template>
  <div style="min-height: 100vh; background: var(--bg-dark); padding: 32px;">
    <button class="btn btn-ghost" @click="router.back()">← Back</button>
    <h1 style="font-family: var(--font-display); color: var(--accent-gold); margin-top: 24px;">
      Combat Tracker
    </h1>
    <p style="color: var(--text-secondary); font-family: var(--font-body); margin-top: 8px;">
      Coming in Module 10!
    </p>
  </div>
</template>
