<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const email = ref('')
const password = ref('')
const error = ref<string | null>(null)
const loading = ref(false)

async function handleLogin() {
  error.value = null
  loading.value = true
  try {
    await auth.login(email.value, password.value)
    await router.push(auth.isDM ? { name: 'dm-dashboard' } : { name: 'characters' })
  } catch (e: any) {
    error.value = e.message ?? 'Invalid email or password'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <!-- Background texture -->
    <div class="bg-layer" />

    <div class="login-container fade-in">
      <!-- Logo / Title -->
      <div class="login-brand">
        <div class="brand-icon">⚔</div>
        <h1 class="brand-name">Behold DnD</h1>
        <p class="brand-tagline">Your campaign, at your fingertips</p>
      </div>

      <div class="ornament" />

      <!-- Form -->
      <form class="login-form" @submit.prevent="handleLogin">
        <div class="field">
          <label class="field-label">Email</label>
          <input
            v-model="email"
            class="input"
            type="email"
            placeholder="adventurer@realm.com"
            required
            autocomplete="email"
          />
        </div>

        <div class="field">
          <label class="field-label">Password</label>
          <input
            v-model="password"
            class="input"
            type="password"
            placeholder="••••••••"
            required
            autocomplete="current-password"
          />
        </div>

        <p v-if="error" class="login-error">{{ error }}</p>

        <button class="btn btn-primary login-btn" type="submit" :disabled="loading">
          <span v-if="loading">Entering the realm...</span>
          <span v-else>Enter the Realm</span>
        </button>
      </form>

      <div class="ornament" />
      <p class="login-footer">Contact your Dungeon Master for access</p>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: var(--bg-dark);
}

/* Subtle radial glow behind the card */
.bg-layer {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(ellipse 60% 60% at 50% 40%, rgba(201,168,71,0.06) 0%, transparent 70%),
    radial-gradient(ellipse 80% 50% at 50% 80%, rgba(139,26,26,0.08) 0%, transparent 70%);
  pointer-events: none;
}

.login-container {
  width: 100%;
  max-width: 400px;
  padding: 48px 40px;
  background: var(--bg-card);
  border: 1px solid var(--border-light);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-card), 0 0 60px rgba(0,0,0,0.6);
  position: relative;
  z-index: 1;
}

/* Gold corner accents */
.login-container::before,
.login-container::after {
  content: '';
  position: absolute;
  width: 20px;
  height: 20px;
  border-color: var(--accent-gold-dim);
  border-style: solid;
}
.login-container::before { top: 12px; left: 12px; border-width: 1px 0 0 1px; }
.login-container::after  { bottom: 12px; right: 12px; border-width: 0 1px 1px 0; }

.login-brand {
  text-align: center;
  margin-bottom: 24px;
}

.brand-icon {
  font-size: 36px;
  color: var(--accent-gold);
  margin-bottom: 12px;
  display: block;
  filter: drop-shadow(0 0 12px rgba(201,168,71,0.4));
}

.brand-name {
  font-family: var(--font-deco);
  font-size: 28px;
  font-weight: 700;
  letter-spacing: 0.1em;
  background: linear-gradient(180deg, #e8dcc8 0%, var(--accent-gold) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 6px;
}

.brand-tagline {
  font-family: var(--font-body);
  font-style: italic;
  color: var(--text-secondary);
  font-size: 15px;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin: 24px 0;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-label {
  font-family: var(--font-display);
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--text-muted);
}

.login-error {
  font-family: var(--font-body);
  font-size: 14px;
  color: var(--accent-red-bright);
  text-align: center;
  padding: 8px;
  background: rgba(139,26,26,0.1);
  border: 1px solid var(--accent-red);
  border-radius: var(--radius-sm);
}

.login-btn {
  width: 100%;
  justify-content: center;
  padding: 12px;
  font-size: 14px;
  margin-top: 4px;
}

.login-footer {
  text-align: center;
  font-family: var(--font-body);
  font-style: italic;
  font-size: 13px;
  color: var(--text-muted);
  margin-top: 16px;
}
</style>
