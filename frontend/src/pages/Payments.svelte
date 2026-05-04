<script lang="ts">
  import { onMount } from 'svelte';
  import { GetPayments, GetPaymentSummary, CreatePayment, UpdatePayment, DeletePayment, GetTenants } from '../../wailsjs/go/main/App';
  import PageHeader from '../components/PageHeader.svelte';
  import Modal from '../components/Modal.svelte';
  import Badge from '../components/Badge.svelte';
  import Toast from '../components/Toast.svelte';
  import ConfirmDialog from '../components/ConfirmDialog.svelte';
  import EmptyState from '../components/EmptyState.svelte';
  import { formatCurrency, formatDate, getCurrentMonth, getCurrentYear, PAYMENT_METHODS, MONTHS, getMethodLabel } from '../lib/utils';

  let loading = true;
  let payments: any[] = [];
  let tenants: any[] = [];
  let summary: any = { totalDue: 0, totalPaid: 0, totalLate: 0, totalMissed: 0, totalPartial: 0, lateFees: 0, count: 0 };

  // Filters
  let filterMonth = getCurrentMonth();
  let filterYear = getCurrentYear();
  let filterStatus = '';

  // Modal state
  let showModal = false;
  let modalTitle = 'Record Payment';
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

  function emptyForm(): any {
    const y = filterYear;
    const m = String(filterMonth).padStart(2, '0');
    return {
      tenantId: 0,
      amount: 0,
      dueDate: `${y}-${m}-01`,
      paidDate: '',
      method: 'bank_transfer',
      status: 'due',
      lateFee: 0,
      notes: '',
    };
  }

  let form: any = emptyForm();

  $: canSave = form.tenantId > 0 && form.amount > 0 && form.dueDate?.trim();

  async function loadTenants() {
    try {
      tenants = (await GetTenants('active')) || [];
    } catch {
      tenants = [];
    }
  }

  async function loadPayments() {
    try {
      payments = (await GetPayments(0, filterYear, filterMonth)) || [];
      payments.sort((a: any, b: any) => (b.dueDate || '').localeCompare(a.dueDate || ''));
    } catch {
      payments = [];
    }
  }

  async function loadSummary() {
    try {
      summary = await GetPaymentSummary(filterYear, filterMonth);
      if (!summary) summary = { totalDue: 0, totalPaid: 0, totalLate: 0, totalMissed: 0, totalPartial: 0, lateFees: 0, count: 0 };
    } catch {
      summary = { totalDue: 0, totalPaid: 0, totalLate: 0, totalMissed: 0, totalPartial: 0, lateFees: 0, count: 0 };
    }
  }

  async function loadAll() {
    loading = true;
    await Promise.all([loadPayments(), loadSummary(), loadTenants()]);
    loading = false;
  }

  // Compute YTD income: sum of totalPaid for months 1..12 of the selected year
  let ytdIncome = 0;
  async function loadYTD() {
    try {
      let total = 0;
      for (let m = 1; m <= 12; m++) {
        const s = await GetPaymentSummary(filterYear, m);
        if (s) total += (s.totalPaid || 0);
      }
      ytdIncome = total;
    } catch {
      ytdIncome = 0;
    }
  }

  function tenantName(tenantId: number): string {
    const t = tenants.find((t: any) => t.id === tenantId);
    return t ? `${t.firstName} ${t.lastName}` : `Tenant #${tenantId}`;
  }

  function statusVariant(status: string): 'success' | 'warning' | 'danger' | 'info' | 'neutral' {
    switch (status) {
      case 'paid': return 'success';
      case 'due': return 'warning';
      case 'late': return 'danger';
      case 'partial': return 'info';
      case 'missed': return 'danger';
      default: return 'neutral';
    }
  }

  $: lateCount = payments.filter((p: any) => p.status === 'late').length;

  $: filteredPayments = filterStatus
    ? payments.filter((p: any) => p.status === filterStatus)
    : payments;

  function openAdd() {
    editingId = null;
    modalTitle = 'Record Payment';
    form = emptyForm();
    showModal = true;
  }

  function openEdit(p: any) {
    editingId = p.id;
    modalTitle = 'Edit Payment';
    form = { ...p };
    showModal = true;
  }

  function openDelete(p: any) {
    deleteId = p.id;
    deleteLabel = `payment of ${formatCurrency(p.amount)} for ${tenantName(p.tenantId)}`;
    showConfirm = true;
  }

  function onTenantChange() {
    if (form.tenantId && !editingId) {
      const t = tenants.find((t: any) => t.id === form.tenantId);
      if (t && t.monthlyRent) {
        form.amount = t.monthlyRent;
      }
    }
  }

  // Auto-suggest late fee when paid date > due date
  $: {
    if (form.paidDate && form.dueDate && form.paidDate > form.dueDate && form.lateFee === 0 && !editingId) {
      form.lateFee = 50;
    }
  }

  async function savePayment() {
    try {
      const payload = { ...form, tenantId: Number(form.tenantId), amount: Number(form.amount), lateFee: Number(form.lateFee) };
      if (editingId) {
        await UpdatePayment(payload);
        toast('Payment updated successfully');
      } else {
        await CreatePayment(payload);
        toast('Payment recorded successfully');
      }
      showModal = false;
      await loadAll();
      loadYTD();
    } catch (e) {
      toast('Failed to save payment: ' + e, 'error');
    }
  }

  async function confirmDelete() {
    if (!deleteId) return;
    try {
      await DeletePayment(deleteId);
      toast('Payment deleted');
      await loadAll();
      loadYTD();
    } catch (e) {
      toast('Failed to delete payment: ' + e, 'error');
    }
    deleteId = null;
  }

  function applyFilters() {
    loading = true;
    loadAll();
    loadYTD();
  }

  const yearOptions: number[] = [];
  const currentYear = getCurrentYear();
  for (let y = currentYear - 5; y <= currentYear + 1; y++) {
    yearOptions.push(y);
  }

  const statusFilters = [
    { label: 'All', value: '' },
    { label: 'Due', value: 'due' },
    { label: 'Paid', value: 'paid' },
    { label: 'Late', value: 'late' },
    { label: 'Partial', value: 'partial' },
    { label: 'Missed', value: 'missed' },
  ];

  onMount(() => {
    loadAll();
    loadYTD();
  });
</script>

<Toast bind:show={showToast} message={toastMessage} type={toastType} />

<PageHeader title="Payments">
  <svelte:fragment slot="actions">
    <button class="btn-primary" on:click={openAdd}>+ Record Payment</button>
  </svelte:fragment>
</PageHeader>

<!-- Summary Cards -->
<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
  <div class="card !p-4 border-l-4 border-l-green-500">
    <p class="text-xs font-medium text-gray-500 uppercase">Collected This Month</p>
    <p class="text-2xl font-bold text-green-700 mt-1">{formatCurrency(summary.totalPaid)}</p>
  </div>
  <div class="card !p-4 border-l-4 border-l-yellow-500">
    <p class="text-xs font-medium text-gray-500 uppercase">Due This Month</p>
    <p class="text-2xl font-bold text-yellow-700 mt-1">{formatCurrency(summary.totalDue)}</p>
  </div>
  <div class="card !p-4 border-l-4 border-l-red-500">
    <p class="text-xs font-medium text-gray-500 uppercase">Late Payments</p>
    <p class="text-2xl font-bold text-red-700 mt-1">{lateCount}</p>
  </div>
  <div class="card !p-4 border-l-4 border-l-blue-500">
    <p class="text-xs font-medium text-gray-500 uppercase">YTD Income</p>
    <p class="text-2xl font-bold text-blue-700 mt-1">{formatCurrency(ytdIncome)}</p>
  </div>
</div>

<!-- Filters Row -->
<div class="flex flex-wrap items-center gap-3 mb-6">
  <select class="input-field w-auto" bind:value={filterMonth} on:change={applyFilters}>
    {#each MONTHS as month, i}
      <option value={i + 1}>{month}</option>
    {/each}
  </select>
  <select class="input-field w-auto" bind:value={filterYear} on:change={applyFilters}>
    {#each yearOptions as y}
      <option value={y}>{y}</option>
    {/each}
  </select>
  <div class="flex gap-1 bg-gray-100 rounded-lg p-1">
    {#each statusFilters as sf}
      <button
        class="px-3 py-1.5 rounded-md text-sm font-medium transition-colors {filterStatus === sf.value ? 'bg-white shadow text-gray-900' : 'text-gray-600 hover:text-gray-900'}"
        on:click={() => { filterStatus = sf.value; }}
      >{sf.label}</button>
    {/each}
  </div>
</div>

{#if loading}
  <div class="flex items-center justify-center py-20">
    <div class="text-gray-400 text-lg">Loading...</div>
  </div>
{:else if filteredPayments.length === 0}
  <div class="card">
    <EmptyState icon="💳" title="No payments found" message="No payment records match your current filters.">
      <button class="btn-primary" on:click={openAdd}>+ Record Payment</button>
    </EmptyState>
  </div>
{:else}
  <div class="card overflow-hidden !p-0">
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b border-gray-200 bg-gray-50">
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Due Date</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Tenant</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Amount</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Status</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Paid Date</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Method</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Late Fee</th>
            <th class="text-right px-6 py-3 text-xs font-medium text-gray-500 uppercase">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          {#each filteredPayments as p}
            <tr class="hover:bg-gray-50">
              <td class="px-6 py-4 text-sm text-gray-900">{formatDate(p.dueDate)}</td>
              <td class="px-6 py-4 font-medium text-gray-900">{tenantName(p.tenantId)}</td>
              <td class="px-6 py-4 text-sm text-gray-900">{formatCurrency(p.amount)}</td>
              <td class="px-6 py-4">
                <Badge variant={statusVariant(p.status)} text={p.status} />
              </td>
              <td class="px-6 py-4 text-sm text-gray-600">{formatDate(p.paidDate)}</td>
              <td class="px-6 py-4 text-sm text-gray-600">{p.method ? getMethodLabel(p.method) : '—'}</td>
              <td class="px-6 py-4 text-sm text-gray-600">{p.lateFee ? formatCurrency(p.lateFee) : '—'}</td>
              <td class="px-6 py-4 text-right">
                <button class="btn-secondary btn-sm" on:click={() => openEdit(p)}>✏️</button>
                <button class="btn-danger btn-sm ml-1" on:click={() => openDelete(p)}>🗑️</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
{/if}

<!-- Record/Edit Payment Modal -->
<Modal bind:show={showModal} title={modalTitle} size="lg">
  <form on:submit|preventDefault={savePayment} id="payment-form">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <label class="label" for="tenantId">Tenant *</label>
        <select id="tenantId" class="input-field" bind:value={form.tenantId} on:change={onTenantChange} required>
          <option value={0} disabled>Select a tenant</option>
          {#each tenants as t}
            <option value={t.id}>{t.firstName} {t.lastName}</option>
          {/each}
        </select>
      </div>
      <div>
        <label class="label" for="amount">Amount *</label>
        <input id="amount" type="number" class="input-field" bind:value={form.amount} min="0" step="0.01" required />
      </div>
      <div>
        <label class="label" for="dueDate">Due Date *</label>
        <input id="dueDate" type="date" class="input-field" bind:value={form.dueDate} required />
      </div>
      <div>
        <label class="label" for="paidDate">Paid Date</label>
        <input id="paidDate" type="date" class="input-field" bind:value={form.paidDate} />
      </div>
      <div>
        <label class="label" for="method">Payment Method</label>
        <select id="method" class="input-field" bind:value={form.method}>
          {#each PAYMENT_METHODS as m}
            <option value={m.value}>{m.label}</option>
          {/each}
        </select>
      </div>
      <div>
        <label class="label" for="status">Status *</label>
        <select id="status" class="input-field" bind:value={form.status} required>
          <option value="due">Due</option>
          <option value="paid">Paid</option>
          <option value="partial">Partial</option>
          <option value="late">Late</option>
          <option value="missed">Missed</option>
        </select>
      </div>
      <div>
        <label class="label" for="lateFee">Late Fee</label>
        <input id="lateFee" type="number" class="input-field" bind:value={form.lateFee} min="0" step="0.01" />
      </div>
      <div><!-- spacer --></div>
      <div class="col-span-full">
        <label class="label" for="notes">Notes</label>
        <textarea id="notes" class="input-field" rows="3" bind:value={form.notes}></textarea>
      </div>
    </div>
  </form>
  <svelte:fragment slot="footer">
    <button class="btn-secondary" on:click={() => showModal = false}>Cancel</button>
    <button class="btn-primary" type="submit" form="payment-form" disabled={!canSave}>
      {editingId ? 'Save Changes' : 'Record Payment'}
    </button>
  </svelte:fragment>
</Modal>

<!-- Delete Confirmation -->
<ConfirmDialog
  bind:show={showConfirm}
  title="Delete Payment"
  message="Are you sure you want to delete this {deleteLabel}? This action cannot be undone."
  confirmText="Delete"
  danger={true}
  on:confirm={confirmDelete}
/>
