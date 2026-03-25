<script lang="ts">
  export let icon: string = '📭';
  export let title: string = 'Nothing here yet';
  export let message: string = '';
</script>

<div class="text-center py-12">
  <div class="text-5xl mb-4">{icon}</div>
  <h3 class="text-lg font-medium text-gray-900 mb-2">{title}</h3>
  {#if message}
    <p class="text-gray-500 mb-4">{message}</p>
  {/if}
  <slot />
</div>
