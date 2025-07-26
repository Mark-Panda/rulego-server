<script lang="js" setup>
import MainLayout from '@src/layout/components/default/main.vue';
</script>

<template>
  <main-layout>
    <router-view></router-view>
 
    <div class="copyright mt-4 text-center text-xs text-gray-500 dark:text-gray-400">
      Copyright © 2025 
      <a href="https://github.com/rulego/rulego-server" target="_blank" class="hover:text-primary transition-colors">
        RuleGo-Server
      </a>. All rights reserved.
      <a href="https://rulego.cc" target="_blank" class="ml-1 hover:text-primary transition-colors">
        Powered by RuleGo
      </a>
    </div>
  </main-layout>
</template>
