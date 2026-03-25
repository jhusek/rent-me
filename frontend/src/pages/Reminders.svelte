<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import {
    GetReminders,
    GetUpcomingReminders,
    CreateReminder,
    DismissReminder,
    DeleteReminder,
  } from '../../wailsjs/go/main/App';
  import PageHeader from '../components/PageHeader.svelte';
  import Modal from '../components/Modal.svelte';
  import Badge from '../components/Badge.svelte';
  import Toast from '../components/Toast.svelte';
  import ConfirmDialog from '../components/ConfirmDialog.svelte';
  import EmptyState from '../components/EmptyState.svelte';
  import { formatDate } from '../lib/utils';

  type Tab = 'upcoming' | 'all' | 'dismissed';
  let activeTab: Tab = 'upcoming';
  let loading = true;
  let reminders: any[] = [];

  // Modal state
  let showModal = false;

  // Delete state
  let showConfirm = false;
  let deleteId: number | null = null;
  let deleteLabel = '';

  // Toast state
  let toastMessage = '';
  let toastType: 'success' | 'error' = 'success';
  let showToast = false;

  function toast(msg: string, type: 'success' | 'error' = 'success') {
    toastMessage = msg;
    toastType = type;
    showToast = true;
  }

  function todayStr(): string {
    const d = new Date();
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
  }

  function emptyForm(): any {
    return {
      title: '',
      description: '',
      dueDate: todayStr(),
      type: 'custom',
      isRecurring: false,
      recurringInterval: 'monthly',
      notes: '',
    };
  }

  let form: any = emptyForm();

  $: canSave = form.title?.trim() && form.dueDate?.trim();
  const TYPE_ICONS: Record<string, string> = {
    lease_expiry: '📋',
    rent_due: '💰',
    maintenance: '🔧',
    insurance: '🛡️',
    tax: '💰',
    custom: '📌',
  };

  const TYPE_LABELS: Record<string, string> = {
    lease_expiry: 'Lease Expiry',
    rent_due: 'Rent Due',
    maintenance: 'Maintenance',
    insurance: 'Insurance',
    tax: 'Tax',
    custom: 'Custom',
  };

  const TYPE_COLORS: Record<string, 'success' | 'warning' | 'danger' | 'info' | 'neutral'> = {
    lease_expiry: 'warning',
    rent_due: 'success',
    maintenance: 'info',
    insurance: 'neutral',
    tax: 'danger',
    custom: 'neutral',
  };

  function daysUntil(dateStr: string): number {
    if (!dateStr) return 0;
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    const due = new Date(dateStr + 'T00:00:00');
    return Math.ceil((due.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
  }

  function dueLabel(dateStr: string): string {
    const days = daysUntil(dateStr);
    if (days < 0) return `${Math.abs(days)} day${Math.abs(days) !== 1 ? 's' : ''} overdue`;
    if (days === 0) return 'Due today';
    if (days === 1) return 'Due tomorrow';
    return `in ${days} day${days !== 1 ? 's' : ''}`;
  }

  function isDueOverdue(dateStr: string): boolean {
    return daysUntil(dateStr) < 0;
  }

  function isDueSoon(dateStr: string): boolean {
    const days = daysUntil(dateStr);
    return days >= 0 && days <= 3;
  }

  async function loadReminders() {
    loading = true;
    try {
      if (activeTab === 'upcoming') {
        reminders = (await GetUpcomingReminders(30)) || [];
      } else if (activeTab === 'dismissed') {
        const all = (await GetReminders(true)) || [];
        reminders = all.filter((r: any) => r.isDismissed);
      } else {
        reminders = (await GetReminders(false)) || [];
      }
    } catch {
      reminders = [];
    }
    loading = false;
  }

  function switchTab(tab: Tab) {
    activeTab = tab;
    loadReminders();
  }

  function openAdd() {
    form = emptyForm();
    showModal = true;
  }

  function openDelete(r: any) {
    deleteId = r.id;
    deleteLabel = r.title;
    showConfirm = true;
  }

  async function saveReminder() {
    try {
      await CreateReminder(form);
      toast('Reminder created successfully');
      showModal = false;
      await loadReminders();
    } catch (e) {
      toast('Failed to create reminder: ' + e, 'error');
    }
  }

  async function dismiss(id: number) {
    try {
      await DismissReminder(id);
      toast('Reminder dismissed');
      await loadReminders();
    } catch (e) {
      toast('Failed to dismiss reminder: ' + e, 'error');
    }
  }

  async function confirmDelete() {
    if (!deleteId) return;
    try {
      await DeleteReminder(deleteId);
      toast('Reminder deleted');
      await loadReminders();
    } catch (e) {
      toast('Failed to delete reminder: ' + e, 'error');
    }
    deleteId = null;
  }

  onMount(loadReminders);
</script>

<Toast bind:show={showToast} message={toastMessage} type={toastType} />

<PageHeader title="Reminders">
  <svelte:fragment slot="actions">
    <button class="btn-primary" on:click={openAdd}>+ Add Reminder</button>
  </svelte:fragment>
</PageHeader>

<!-- Tabs -->
<div class="flex gap-1 mb-6 bg-gray-100 rounded-lg p-1 w-fit">
  <button
    class="px-4 py-2 text-sm font-medium rounded-md transition-colors {activeTab === 'upcoming' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-500 hover:text-gray-700'}"
    on:click={() => switchTab('upcoming')}
  >
    Upcoming
  </button>
  <button
    class="px-4 py-2 text-sm font-medium rounded-md transition-colors {activeTab === 'all' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-500 hover:text-gray-700'}"
    on:click={() => switchTab('all')}
  >
    All
  </button>
  <button
    class="px-4 py-2 text-sm font-medium rounded-md transition-colors {activeTab === 'dismissed' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-500 hover:text-gray-700'}"
    on:click={() => switchTab('dismissed')}
  >
    Dismissed
  </button>
</div>

{#if loading}
  <div class="flex items-center justify-center py-20">
    <div class="text-gray-400 text-lg">Loading...</div>
  </div>
{:else if reminders.length === 0}
  <div class="card">
    <EmptyState
      icon="🔔"
      title="No reminders"
      message={activeTab === 'upcoming' ? 'No upcoming reminders in the next 30 days.' : activeTab === 'dismissed' ? 'No dismissed reminders.' : 'No active reminders found.'}
    >
      {#if activeTab !== 'dismissed'}
        <button class="btn-primary" on:click={openAdd}>+ Add Reminder</button>
      {/if}
    </EmptyState>
  </div>
{:else}
  <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
    {#each reminders as r}
      <div class="card !p-0 overflow-hidden border-l-4 {isDueOverdue(r.dueDate) ? 'border-l-red-500' : isDueSoon(r.dueDate) ? 'border-l-yellow-500' : 'border-l-blue-500'}">
        <div class="p-5">
          <!-- Header -->
          <div class="flex items-start justify-between mb-3">
            <div class="flex items-center gap-2 min-w-0">
              <span class="text-xl flex-shrink-0">{TYPE_ICONS[r.type] || '📌'}</span>
              <h3 class="font-semibold text-gray-900 truncate">{r.title}</h3>
            </div>
            <Badge variant={TYPE_COLORS[r.type] || 'neutral'} text={TYPE_LABELS[r.type] || r.type} />
          </div>

          <!-- Description -->
          {#if r.description}
            <p class="text-sm text-gray-500 mb-3 line-clamp-2">{r.description}</p>
          {/if}

          <!-- Due date -->
          <div class="flex items-center gap-2 text-sm mb-3">
            <span class="text-gray-400">📅</span>
            <span class="text-gray-600">{formatDate(r.dueDate)}</span>
            <span class="font-medium {isDueOverdue(r.dueDate) ? 'text-red-600' : isDueSoon(r.dueDate) ? 'text-yellow-600' : 'text-blue-600'}">
              — {dueLabel(r.dueDate)}
            </span>
          </div>

          <!-- Recurring indicator -->
          {#if r.isRecurring}
            <div class="flex items-center gap-1 text-sm text-gray-500 mb-3">
              <span>🔄</span>
              <span>Recurring ({r.recurringInterval || 'monthly'})</span>
            </div>
          {/if}

          <!-- Dismissed indicator -->
          {#if r.isDismissed}
            <div class="mb-3">
              <Badge variant="neutral" text="Dismissed" />
            </div>
          {/if}
        </div>

        <!-- Actions -->
        <div class="px-5 py-3 bg-gray-50 border-t border-gray-100 flex justify-end gap-2">
          {#if !r.isDismissed}
            <button
              class="btn-secondary btn-sm"
              on:click={() => dismiss(r.id)}
              title="Dismiss"
            >✓ Dismiss</button>
          {/if}
          <button
            class="btn-danger btn-sm"
            on:click={() => openDelete(r)}
            title="Delete"
          >🗑️ Delete</button>
        </div>
      </div>
    {/each}
  </div>
{/if}

<!-- Add Reminder Modal -->
<Modal bind:show={showModal} title="Add Reminder" size="lg">
  <form on:submit|preventDefault={saveReminder} id="reminder-form">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="col-span-full">
        <label class="label" for="rem-title">Title *</label>
        <input id="rem-title" class="input-field" bind:value={form.title} required placeholder="e.g. Rent due, Lease renewal..." />
      </div>
      <div class="col-span-full">
        <label class="label" for="rem-desc">Description</label>
        <textarea id="rem-desc" class="input-field" rows="2" bind:value={form.description} placeholder="Optional details..."></textarea>
      </div>
      <div>
        <label class="label" for="rem-due">Due Date *</label>
        <input id="rem-due" type="date" class="input-field" bind:value={form.dueDate} required />
      </div>
      <div>
        <label class="label" for="rem-type">Type</label>
        <select id="rem-type" class="input-field" bind:value={form.type}>
          <option value="custom">📌 Custom</option>
          <option value="lease_expiry">📋 Lease Expiry</option>
          <option value="rent_due">💰 Rent Due</option>
          <option value="maintenance">🔧 Maintenance</option>
          <option value="insurance">🛡️ Insurance</option>
          <option value="tax">💰 Tax</option>
        </select>
      </div>

      <hr class="col-span-full border-gray-200" />

      <div class="flex items-center gap-6 col-span-full">
        <label class="flex items-center gap-2 cursor-pointer">
          <input type="checkbox" class="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500" bind:checked={form.isRecurring} />
          <span class="text-sm font-medium text-gray-700">Recurring Reminder</span>
        </label>
      </div>

      {#if form.isRecurring}
        <div>
          <label class="label" for="rem-interval">Recurring Interval</label>
          <select id="rem-interval" class="input-field" bind:value={form.recurringInterval}>
            <option value="monthly">Monthly</option>
            <option value="quarterly">Quarterly</option>
            <option value="annually">Annually</option>
          </select>
        </div>
      {/if}

      <div class="col-span-full">
        <label class="label" for="rem-notes">Notes</label>
        <textarea id="rem-notes" class="input-field" rows="2" bind:value={form.notes}></textarea>
      </div>
    </div>
  </form>
  <svelte:fragment slot="footer">
    <button class="btn-secondary" on:click={() => showModal = false}>Cancel</button>
    <button class="btn-primary" type="submit" form="reminder-form" disabled={!canSave}>Add Reminder</button>
  </svelte:fragment>
</Modal>

<!-- Delete Confirmation -->
<ConfirmDialog
  bind:show={showConfirm}
  title="Delete Reminder"
  message="Are you sure you want to delete &quot;{deleteLabel}&quot;? This action cannot be undone."
  confirmText="Delete"
  danger={true}
  on:confirm={confirmDelete}
/>
