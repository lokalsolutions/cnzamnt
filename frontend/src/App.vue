<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { createArtwork, getArtworks, getMe } from './api';
import type { Artwork, User } from './types';

const user = ref<User | null>(null);
const artworks = ref<Artwork[]>([]);
const loading = ref(true);
const saving = ref(false);
const error = ref('');
const showCreate = ref(false);

const form = reactive({
  title: '',
  image_url: '',
  caption: '',
});

const balance = computed(() => {
  if (!user.value) return '...';
  return `${formatCNZ(user.value.cnz_balance)} CNZ`;
});

async function loadApp() {
  loading.value = true;
  error.value = '';
  try {
    const [me, feed] = await Promise.all([getMe(), getArtworks()]);
    user.value = me;
    artworks.value = feed;
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Unable to load CnzAMnt';
  } finally {
    loading.value = false;
  }
}

async function submitArtwork() {
  if (!form.title.trim() || !form.image_url.trim()) {
    error.value = 'Title and image URL are required.';
    return;
  }

  saving.value = true;
  error.value = '';
  try {
    const artwork = await createArtwork({
      title: form.title.trim(),
      image_url: form.image_url.trim(),
      caption: form.caption.trim(),
    });
    artworks.value = [artwork, ...artworks.value];
    form.title = '';
    form.image_url = '';
    form.caption = '';
    showCreate.value = false;
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Unable to create artwork';
  } finally {
    saving.value = false;
  }
}

function commentCount(artwork: Artwork) {
  return artwork.comment_count ?? artwork.comments?.length ?? 0;
}

function formatCNZ(value: number) {
  return new Intl.NumberFormat(undefined, {
    maximumFractionDigits: 2,
  }).format(value);
}

onMounted(loadApp);
</script>

<template>
  <main class="app-shell">
    <header class="topbar">
      <div>
        <p class="eyebrow">Art feedback</p>
        <h1>CnzAMnt</h1>
      </div>
      <div class="balance-pill">
        <span>Balance</span>
        <strong>{{ balance }}</strong>
      </div>
    </header>

    <section class="composer-band">
      <button class="primary-action" type="button" @click="showCreate = !showCreate">
        {{ showCreate ? 'Close' : 'Post artwork' }}
      </button>
      <button class="ghost-action" type="button" @click="loadApp">Refresh</button>
    </section>

    <form v-if="showCreate" class="create-panel" @submit.prevent="submitArtwork">
      <label>
        <span>Title</span>
        <input v-model="form.title" autocomplete="off" placeholder="Sketchbook study" />
      </label>
      <label>
        <span>Image URL</span>
        <input v-model="form.image_url" inputmode="url" placeholder="https://..." />
      </label>
      <label>
        <span>Caption</span>
        <textarea v-model="form.caption" rows="3" placeholder="What were you trying to make?" />
      </label>
      <button class="submit-action" type="submit" :disabled="saving">
        {{ saving ? 'Posting...' : 'Post' }}
      </button>
    </form>

    <section v-if="loading" class="state-panel">
      <div class="loader"></div>
      <p>Loading art feed...</p>
    </section>

    <section v-else-if="error" class="state-panel error-state">
      <p>{{ error }}</p>
      <p class="state-note">For local testing, the backend expects a dev user with id 1.</p>
    </section>

    <section v-else-if="artworks.length === 0" class="state-panel">
      <p>No artwork yet.</p>
      <p class="state-note">Post the first piece and let the feed breathe a little.</p>
    </section>

    <section v-else class="feed" aria-label="Artwork feed">
      <article v-for="artwork in artworks" :key="artwork.id" class="art-card">
        <div class="image-frame">
          <img :src="artwork.image_url" :alt="artwork.title" loading="lazy" />
        </div>
        <div class="art-copy">
          <div class="art-heading">
            <div>
              <h2>{{ artwork.title }}</h2>
              <p>@{{ artwork.artist?.handle ?? 'artist' }}</p>
            </div>
            <div class="comment-count">
              <strong>{{ commentCount(artwork) }}</strong>
              <span>comments</span>
            </div>
          </div>
          <p class="caption">{{ artwork.caption || 'No caption yet.' }}</p>
        </div>
      </article>
    </section>
  </main>
</template>
