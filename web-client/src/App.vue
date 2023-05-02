<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { useAuthentication } from "./services/authentication"
const { logout } = useAuthentication()
</script>

<template>
  <v-card class="h-100">
    <v-layout>
      <v-app-bar title="L'édition libre">
        <v-btn @click="logout">Logout</v-btn>
      </v-app-bar>
      <v-main style="min-height: 300px;">
        <RouterView />
      </v-main>
    </v-layout>
  </v-card>
</template>
