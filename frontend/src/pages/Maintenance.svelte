<script lang="ts">
  import { onMount } from 'svelte';
  import { GetMaintenanceRequests, CreateMaintenanceRequest, UpdateMaintenanceRequest, DeleteMaintenanceRequest } from '../../wailsjs/go/main/App';
  import PageHeader from '../components/PageHeader.svelte';
  import Modal from '../components/Modal.svelte';
  import Badge from '../components/Badge.svelte';
  import Toast from '../components/Toast.svelte';
  import ConfirmDialog from '../components/ConfirmDialog.svelte';
  import EmptyState from '../components/EmptyState.svelte';
  import { formatCurrency, formatDate } from '../lib/utils';

  let loading = true;
  let requests: any[] = [];

  // Filters
  let filterStatus = '';
  let filterPriority = '';

  // Modal state
  let showModal = false;
  let modalTitle = 'New Request';
  let editingId: number | null = null;

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
      priority: 'medium',
      status: 'open',
      reportedDate: todayStr(),
      completedDate: '',
      cost: 0,
      vendor: '',
      notes: '',
    };
  }

  let form: any = emptyForm();

  $: canSave = form.title?.trim();
  async function loadRequests() {
    loading = true;
    try {
      requests = (await GetMaintenanceRequests('')) || [];
    } catch {
      requests = [];
    }
    loading = false;
  }

  // Summary computations
  $: openCount = requests.filter((r: any) => r.status === 'open').length;
  $: inProgressCount = requests.filter((r: any) => r.status === 'in_progress').length;

  $: completedThisMonth = (() => {
    const now = new Date();
    const cm = now.getMonth() + 1;
    const cy = now.getFullYear();
    return requests.filter((r: any) => {
      if (r.status !== 'completed' || !r.completedDate) return false;
      const d = new Date(r.completedDate + 'T00:00:00');
      return d.getMonth() + 1 === cm && d.getFullYear() === cy;
    }).length;
  })();

  $: totalCostThisMonth = (() => {
    const now = new Date();
    const cm = now.getMonth() + 1;
    const cy = now.getFullYear();
    return requests
      .filter((r: any) => {
        if (r.status !== 'completed' || !r.completedDate) return false;
        const d = new Date(r.completedDate + 'T00:00:00');
        return d.getMonth() + 1 === cm && d.getFullYear() === cy;
      })
      .reduce((sum: number, r: any) => sum + (r.cost || 0), 0);
  })();

  $: filteredRequests = requests.filter((r: any) => {
    if (filterStatus && r.status !== filterStatus) return false;
    if (filterPriority && r.priority !== filterPriority) return false;
    return true;
  });

  const statusFilters = [
    { label: 'All', value: '' },
    { label: 'Open', value: 'open' },
    { label: 'In Progress', value: 'in_progress' },
    { label: 'Completed', value: 'completed' },
    { label: 'Cancelled', value: 'cancelled' },
  ];

  const priorityOptions = [
    { label: 'All Priorities', value: '' },
    { label: 'Emergency', value: 'emergency' },
    { label: 'High', value: 'high' },
    { label: 'Medium', value: 'medium' },
    { label: 'Low', value: 'low' },
  ];

  function priorityVariant(p: string): 'success' | 'warning' | 'danger' | 'info' | 'neutral' {
    switch (p) {
      case 'emergency': return 'danger';
      case 'high': return 'warning';
      case 'medium': return 'info';
      case 'low': return 'success';
      default: return 'neutral';
    }
  }

  function priorityLabel(p: string): string {
    switch (p) {
      case 'emergency': return 'Emergency';
      case 'high': return 'High';
      case 'medium': return 'Medium';
      case 'low': return 'Low';
      default: return p;
    }
  }

  function statusVariant(s: string): 'success' | 'warning' | 'danger' | 'info' | 'neutral' {
    switch (s) {
      case 'open': return 'warning';
      case 'in_progress': return 'info';
      case 'completed': return 'success';
      case 'cancelled': return 'neutral';
      default: return 'neutral';
    }
  }

  function statusLabel(s: string): string {
    switch (s) {
      case 'open': return 'Open';
      case 'in_progress': return 'In Progress';
      case 'completed': return 'Completed';
      case 'cancelled': return 'Cancelled';
      default: return s;
    }
  }

  function openAdd() {
    editingId = null;
    modalTitle = 'New Request';
    form = emptyForm();
    showModal = true;
  }

  function openEdit(r: any) {
    editingId = r.id;
    modalTitle = 'Edit Request';
    form = { ...r };
    showModal = true;
  }

  function openDelete(r: any) {
    deleteId = r.id;
    deleteLabel = r.title;
    showConfirm = true;
  }

  // Auto-fill completed date when status changes to completed
  $: {
    if (form.status === 'completed' && !form.completedDate) {
      form.completedDate = todayStr();
    }
  }

  async function saveRequest() {
    try {
      const payload = { ...form, cost: Number(form.cost) };
      if (editingId) {
        await UpdateMaintenanceRequest(payload);
        toast('Request updated successfully');
      } else {
        await CreateMaintenanceRequest(payload);
        toast('Request created successfully');
      }
      showModal = false;
      await loadRequests();
    } catch (e) {
      toast('Failed to save request: ' + e, 'error');
    }
  }

  async function confirmDelete() {
    if (!deleteId) return;
    try {
      await DeleteMaintenanceRequest(deleteId);
      toast('Request deleted');
      await loadRequests();
    } catch (e) {
      toast('Failed to delete request: ' + e, 'error');
    }
    deleteId = null;
  }

  onMount(loadRequests);
</script>

<Toast bind:show={showToast} message={toastMessage} type={toastType} />

<PageHeader title="Maintenance">
  <svelte:fragment slot="actions">
    <button class="btn-primary" on:click={openAdd}>+ New Request</button>
  </svelte:fragment>
</PageHeader>

<!-- Summary Cards -->
<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
  <div class="card !p-4 border-l-4 border-l-yellow-500">
    <p class="text-xs font-medium text-gray-500 uppercase">Open Requests</p>
    <p class="text-2xl font-bold text-yellow-700 mt-1">{openCount}</p>
  </div>
  <div class="card !p-4 border-l-4 border-l-blue-500">
    <p class="text-xs font-medium text-gray-500 uppercase">In Progress</p>
    <p class="text-2xl font-bold text-blue-700 mt-1">{inProgressCount}</p>
  </div>
  <div class="card !p-4 border-l-4 border-l-green-500">
    <p class="text-xs font-medium text-gray-500 uppercase">Completed This Month</p>
    <p class="text-2xl font-bold text-green-700 mt-1">{completedThisMonth}</p>
  </div>
  <div class="card !p-4 border-l-4 border-l-red-500">
    <p class="text-xs font-medium text-gray-500 uppercase">Total Cost This Month</p>
    <p class="text-2xl font-bold text-red-700 mt-1">{formatCurrency(totalCostThisMonth)}</p>
  </div>
</div>

<!-- Filters Row -->
<div class="flex flex-wrap items-center gap-3 mb-6">
  <div class="flex gap-1 bg-gray-100 rounded-lg p-1">
    {#each statusFilters as sf}
      <button
        class="px-3 py-1.5 rounded-md text-sm font-medium transition-colors {filterStatus === sf.value ? 'bg-white shadow text-gray-900' : 'text-gray-600 hover:text-gray-900'}"
        on:click={() => { filterStatus = sf.value; }}
      >{sf.label}</button>
    {/each}
  </div>
  <select class="input-field w-auto" bind:value={filterPriority}>
    {#each priorityOptions as po}
      <option value={po.value}>{po.label}</option>
    {/each}
  </select>
</div>

{#if loading}
  <div class="flex items-center justify-center py-20">
    <div class="text-gray-400 text-lg">Loading...</div>
  </div>
{:else if filteredRequests.length === 0}
  <div class="card">
    <EmptyState icon="🔧" title="No maintenance requests" message="No requests match your current filters.">
      <button class="btn-primary" on:click={openAdd}>+ New Request</button>
    </EmptyState>
  </div>
{:else}
  <div class="card overflow-hidden !p-0">
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b border-gray-200 bg-gray-50">
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Title</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Priority</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Status</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Reported</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Vendor</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Cost</th>
            <th class="text-right px-6 py-3 text-xs font-medium text-gray-500 uppercase">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          {#each filteredRequests as r}
            <tr class="hover:bg-gray-50">
              <td class="px-6 py-4 font-medium text-gray-900">{r.title}</td>
              <td class="px-6 py-4">
                <Badge variant={priorityVariant(r.priority)} text={priorityLabel(r.priority)} />
              </td>
              <td class="px-6 py-4">
                <Badge variant={statusVariant(r.status)} text={statusLabel(r.status)} />
              </td>
              <td class="px-6 py-4 text-sm text-gray-600">{formatDate(r.reportedDate)}</td>
              <td class="px-6 py-4 text-sm text-gray-600">{r.vendor || '—'}</td>
              <td class="px-6 py-4 text-sm text-gray-900">{r.cost ? formatCurrency(r.cost) : '—'}</td>
              <td class="px-6 py-4 text-right">
                <button class="btn-secondary btn-sm" on:click={() => openEdit(r)}>✏️</button>
                <button class="btn-danger btn-sm ml-1" on:click={() => openDelete(r)}>🗑️</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
{/if}

<!-- Add/Edit Request Modal -->
<Modal bind:show={showModal} title={modalTitle} size="lg">
  <form on:submit|preventDefault={saveRequest} id="maintenance-form">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="col-span-full">
        <label class="label" for="title">Title *</label>
        <input id="title" class="input-field" bind:value={form.title} required />
      </div>
      <div class="col-span-full">
        <label class="label" for="description">Description</label>
        <textarea id="description" class="input-field" rows="3" bind:value={form.description}></textarea>
      </div>
      <div>
        <label class="label" for="priority">Priority</label>
        <select id="priority" class="input-field" bind:value={form.priority}>
          <option value="low">Low</option>
          <option value="medium">Medium</option>
          <option value="high">High</option>
          <option value="emergency">Emergency</option>
        </select>
      </div>
      <div>
        <label class="label" for="status">Status</label>
        <select id="status" class="input-field" bind:value={form.status}>
          <option value="open">Open</option>
          <option value="in_progress">In Progress</option>
          <option value="completed">Completed</option>
          <option value="cancelled">Cancelled</option>
        </select>
      </div>
      <div>
        <label class="label" for="reportedDate">Reported Date</label>
        <input id="reportedDate" type="date" class="input-field" bind:value={form.reportedDate} />
      </div>
      <div>
        <label class="label" for="vendor">Vendor</label>
        <input id="vendor" class="input-field" bind:value={form.vendor} />
      </div>

      {#if form.status === 'completed'}
        <div>
          <label class="label" for="completedDate">Completed Date *</label>
          <input id="completedDate" type="date" class="input-field" bind:value={form.completedDate} required />
        </div>
        <div>
          <label class="label" for="cost">Cost *</label>
          <input id="cost" type="number" class="input-field" bind:value={form.cost} min="0" step="0.01" required />
        </div>
        <div class="col-span-full">
          <p class="text-sm text-amber-600 bg-amber-50 rounded-lg p-3">
            💡 Completing this request will log the cost as a maintenance expense.
          </p>
        </div>
      {/if}

      <div class="col-span-full">
        <label class="label" for="notes">Notes</label>
        <textarea id="notes" class="input-field" rows="3" bind:value={form.notes}></textarea>
      </div>
    </div>
  </form>
  <svelte:fragment slot="footer">
    <button class="btn-secondary" on:click={() => showModal = false}>Cancel</button>
    <button class="btn-primary" type="submit" form="maintenance-form" disabled={!canSave}>
      {editingId ? 'Save Changes' : 'Create Request'}
    </button>
  </svelte:fragment>
</Modal>

<!-- Delete Confirmation -->
<ConfirmDialog
  bind:show={showConfirm}
  title="Delete Request"
  message="Are you sure you want to delete &quot;{deleteLabel}&quot;? This action cannot be undone."
  confirmText="Delete"
  danger={true}
  on:confirm={confirmDelete}
/>
