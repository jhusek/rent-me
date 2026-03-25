<script lang="ts">
  export let show: boolean = false;
  export let title: string = '';
  export let size: 'sm' | 'md' | 'lg' = 'md';

  const sizeClasses = {
    sm: 'max-w-md',
    md: 'max-w-lg',
    lg: 'max-w-2xl',
  };

  function handleBackdrop(e: MouseEvent) {
    if (e.target === e.currentTarget) show = false;
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') show = false;
  }
</script>

<svelte:window on:keydown={handleKeydown} />

{#if show}
  <div 
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
    on:click={handleBackdrop}
  >
    <div class="bg-white rounded-xl shadow-xl w-full {sizeClasses[size]} max-h-[90vh] flex flex-col">
      <div class="flex items-center justify-between p-6 border-b border-gray-200">
        <h2 class="text-lg font-semibold text-gray-900">{title}</h2>
        <button 
          class="text-gray-400 hover:text-gray-600 text-xl"
          on:click={() => show = false}
        >✕</button>
      </div>
      <div class="p-6 overflow-y-auto">
        <slot />
      </div>
      {#if $$slots.footer}
        <div class="p-6 border-t border-gray-200 flex justify-end gap-3">
          <slot name="footer" />
        </div>
      {/if}
    </div>
  </div>
{/if}
