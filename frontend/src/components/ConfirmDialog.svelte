<script lang="ts">
  import Modal from './Modal.svelte';

  export let show: boolean = false;
  export let title: string = 'Are you sure?';
  export let message: string = 'This action cannot be undone.';
  export let confirmText: string = 'Confirm';
  export let cancelText: string = 'Cancel';
  export let danger: boolean = false;

  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();

  function confirm() {
    dispatch('confirm');
    show = false;
  }
</script>

<Modal bind:show {title} size="sm">
  <p class="text-gray-600">{message}</p>
  <svelte:fragment slot="footer">
    <button class="btn-secondary" on:click={() => show = false}>{cancelText}</button>
    <button class={danger ? 'btn-danger' : 'btn-primary'} on:click={confirm}>{confirmText}</button>
  </svelte:fragment>
</Modal>
