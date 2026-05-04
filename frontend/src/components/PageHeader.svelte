<script lang="ts">
  export let title: string;
</script>

<div class="flex items-center justify-between mb-6">
  <h1 class="text-2xl font-bold text-gray-900">{title}</h1>
  <slot name="actions" />
</div>
