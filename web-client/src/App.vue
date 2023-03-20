<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <v-card class="h-100">
    <v-layout>
      <v-app-bar title="L'édition libre"></v-app-bar>
      <v-main style="min-height: 300px;">
        <RouterView />
      </v-main>
    </v-layout>
  </v-card>
</template>
