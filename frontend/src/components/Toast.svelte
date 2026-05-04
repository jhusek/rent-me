<script lang="ts">
  import { fade } from 'svelte/transition';

  export let message: string = '';
  export let type: 'success' | 'error' | 'info' = 'success';
  export let show: boolean = false;

  const typeClasses = {
    success: 'bg-green-600',
    error: 'bg-red-600',
    info: 'bg-blue-600',
  };

  const icons = {
    success: '✓',
    error: '✕',
    info: 'ℹ',
  };

  $: if (show) {
    setTimeout(() => { show = false; }, 3000);
  }
</script>

{#if show}
  <div 
    class="fixed top-4 right-4 z-[100] flex items-center gap-3 px-4 py-3 rounded-lg text-white shadow-lg {typeClasses[type]}"
    transition:fade={{ duration: 200 }}
  >
    <span class="font-bold">{icons[type]}</span>
    <span>{message}</span>
  </div>
{/if}
