<script lang="ts">
  import { onMount } from 'svelte';
  import { GetExpenses, GetExpenseSummary, CreateExpense, UpdateExpense, DeleteExpense } from '../../wailsjs/go/main/App';
  import PageHeader from '../components/PageHeader.svelte';
  import Modal from '../components/Modal.svelte';
  import Badge from '../components/Badge.svelte';
  import Toast from '../components/Toast.svelte';
  import ConfirmDialog from '../components/ConfirmDialog.svelte';
  import EmptyState from '../components/EmptyState.svelte';
  import { formatCurrency, formatDate, getCurrentYear, EXPENSE_CATEGORIES, getCategoryLabel, getCategoryColor } from '../lib/utils';

  let loading = true;
  let expenses: any[] = [];
  let summary: any = { totalAmount: 0, byCategory: {}, count: 0, taxDeductible: 0 };

  // Filters
  let filterCategory = '';
  let filterYear = getCurrentYear();

  // Modal state
  let showModal = false;
  let modalTitle = 'Add Expense';
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
      category: '',
      amount: 0,
      date: todayStr(),
      vendor: '',
      description: '',
      isRecurring: false,
      recurringInterval: 'monthly',
      taxDeductible: true,
      receiptPath: '',
      notes: '',
    };
  }

  let form: any = emptyForm();

  $: canSave = form.category?.trim() && form.amount > 0 && form.date?.trim();

  async function loadExpenses() {
    try {
      expenses = (await GetExpenses(filterCategory, filterYear)) || [];
      expenses.sort((a: any, b: any) => (b.date || '').localeCompare(a.date || ''));
    } catch {
      expenses = [];
    }
  }

  async function loadSummary() {
    try {
      summary = await GetExpenseSummary(filterYear);
      if (!summary) summary = { totalAmount: 0, byCategory: {}, count: 0, taxDeductible: 0 };
    } catch {
      summary = { totalAmount: 0, byCategory: {}, count: 0, taxDeductible: 0 };
    }
  }

  async function loadAll() {
    loading = true;
    await Promise.all([loadExpenses(), loadSummary()]);
    loading = false;
  }

  // Computed values for summary cards
  $: monthlyAvg = summary.count > 0 ? summary.totalAmount / 12 : 0;

  $: topCategory = (() => {
    const bc = summary.byCategory || {};
    let maxCat = '';
    let maxVal = 0;
    for (const [cat, val] of Object.entries(bc)) {
      if ((val as number) > maxVal) {
        maxVal = val as number;
        maxCat = cat;
      }
    }
    return maxCat ? getCategoryLabel(maxCat) : '—';
  })();

  // Current month total (computed from filtered expenses)
  $: currentMonthTotal = (() => {
    const now = new Date();
    const cm = now.getMonth() + 1;
    const cy = now.getFullYear();
    return expenses
      .filter((e: any) => {
        if (!e.date) return false;
        const d = new Date(e.date + 'T00:00:00');
        return d.getMonth() + 1 === cm && d.getFullYear() === cy;
      })
      .reduce((sum: number, e: any) => sum + (e.amount || 0), 0);
  })();

  function openAdd() {
    editingId = null;
    modalTitle = 'Add Expense';
    form = emptyForm();
    showModal = true;
  }

  function openEdit(e: any) {
    editingId = e.id;
    modalTitle = 'Edit Expense';
    form = { ...e };
    showModal = true;
  }

  function openDelete(e: any) {
    deleteId = e.id;
    deleteLabel = `${getCategoryLabel(e.category)} expense of ${formatCurrency(e.amount)}`;
    showConfirm = true;
  }

  async function saveExpense() {
    try {
      const payload = { ...form, amount: Number(form.amount) };
      if (editingId) {
        await UpdateExpense(payload);
        toast('Expense updated successfully');
      } else {
        await CreateExpense(payload);
        toast('Expense added successfully');
      }
      showModal = false;
      await loadAll();
    } catch (e) {
      toast('Failed to save expense: ' + e, 'error');
    }
  }

  async function confirmDelete() {
    if (!deleteId) return;
    try {
      await DeleteExpense(deleteId);
      toast('Expense deleted');
      await loadAll();
    } catch (e) {
      toast('Failed to delete expense: ' + e, 'error');
    }
    deleteId = null;
  }

  function applyFilters() {
    loading = true;
    loadAll();
  }

  const yearOptions: number[] = [];
  const currentYear = getCurrentYear();
  for (let y = currentYear - 5; y <= currentYear + 1; y++) {
    yearOptions.push(y);
  }

  onMount(loadAll);
</script>

<Toast bind:show={showToast} message={toastMessage} type={toastType} />

<PageHeader title="Expenses">
  <svelte:fragment slot="actions">
    <button class="btn-primary" on:click={openAdd}>+ Add Expense</button>
  </svelte:fragment>
</PageHeader>

<!-- Summary Cards -->
<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
  <div class="card !p-4 border-l-4 border-l-red-500">
    <p class="text-xs font-medium text-gray-500 uppercase">This Month</p>
    <p class="text-2xl font-bold text-red-700 mt-1">{formatCurrency(currentMonthTotal)}</p>
  </div>
  <div class="card !p-4 border-l-4 border-l-blue-500">
    <p class="text-xs font-medium text-gray-500 uppercase">Monthly Average</p>
    <p class="text-2xl font-bold text-blue-700 mt-1">{formatCurrency(monthlyAvg)}</p>
  </div>
  <div class="card !p-4 border-l-4 border-l-purple-500">
    <p class="text-xs font-medium text-gray-500 uppercase">YTD Total</p>
    <p class="text-2xl font-bold text-purple-700 mt-1">{formatCurrency(summary.totalAmount)}</p>
  </div>
  <div class="card !p-4 border-l-4 border-l-yellow-500">
    <p class="text-xs font-medium text-gray-500 uppercase">Top Category</p>
    <p class="text-2xl font-bold text-yellow-700 mt-1">{topCategory}</p>
  </div>
</div>

<!-- Filters Row -->
<div class="flex flex-wrap items-center gap-3 mb-6">
  <select class="input-field w-auto" bind:value={filterCategory} on:change={applyFilters}>
    <option value="">All Categories</option>
    {#each EXPENSE_CATEGORIES as cat}
      <option value={cat.value}>{cat.label}</option>
    {/each}
  </select>
  <select class="input-field w-auto" bind:value={filterYear} on:change={applyFilters}>
    {#each yearOptions as y}
      <option value={y}>{y}</option>
    {/each}
  </select>
</div>

{#if loading}
  <div class="flex items-center justify-center py-20">
    <div class="text-gray-400 text-lg">Loading...</div>
  </div>
{:else if expenses.length === 0}
  <div class="card">
    <EmptyState icon="📊" title="No expenses found" message="No expense records match your current filters.">
      <button class="btn-primary" on:click={openAdd}>+ Add Expense</button>
    </EmptyState>
  </div>
{:else}
  <div class="card overflow-hidden !p-0">
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b border-gray-200 bg-gray-50">
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Date</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Category</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Amount</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Vendor</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Description</th>
            <th class="text-center px-6 py-3 text-xs font-medium text-gray-500 uppercase">Recurring</th>
            <th class="text-center px-6 py-3 text-xs font-medium text-gray-500 uppercase">Tax Ded.</th>
            <th class="text-right px-6 py-3 text-xs font-medium text-gray-500 uppercase">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          {#each expenses as e}
            <tr class="hover:bg-gray-50">
              <td class="px-6 py-4 text-sm text-gray-900">{formatDate(e.date)}</td>
              <td class="px-6 py-4">
                <Badge variant={getCategoryColor(e.category)} text={getCategoryLabel(e.category)} />
              </td>
              <td class="px-6 py-4 text-sm font-medium text-gray-900">{formatCurrency(e.amount)}</td>
              <td class="px-6 py-4 text-sm text-gray-600">{e.vendor || '—'}</td>
              <td class="px-6 py-4 text-sm text-gray-600 max-w-xs truncate">{e.description || '—'}</td>
              <td class="px-6 py-4 text-center text-sm">{e.isRecurring ? '🔄' : ''}</td>
              <td class="px-6 py-4 text-center text-sm">{e.taxDeductible ? '✓' : '✗'}</td>
              <td class="px-6 py-4 text-right">
                <button class="btn-secondary btn-sm" on:click={() => openEdit(e)}>✏️</button>
                <button class="btn-danger btn-sm ml-1" on:click={() => openDelete(e)}>🗑️</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
{/if}

<!-- Add/Edit Expense Modal -->
<Modal bind:show={showModal} title={modalTitle} size="lg">
  <form on:submit|preventDefault={saveExpense} id="expense-form">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <label class="label" for="category">Category *</label>
        <select id="category" class="input-field" bind:value={form.category} required>
          <option value="" disabled>Select a category</option>
          {#each EXPENSE_CATEGORIES as cat}
            <option value={cat.value}>{cat.label}</option>
          {/each}
        </select>
      </div>
      <div>
        <label class="label" for="amount">Amount *</label>
        <input id="amount" type="number" class="input-field" bind:value={form.amount} min="0" step="0.01" required />
      </div>
      <div>
        <label class="label" for="date">Date *</label>
        <input id="date" type="date" class="input-field" bind:value={form.date} required />
      </div>
      <div>
        <label class="label" for="vendor">Vendor</label>
        <input id="vendor" class="input-field" bind:value={form.vendor} />
      </div>
      <div class="col-span-full">
        <label class="label" for="description">Description</label>
        <input id="description" class="input-field" bind:value={form.description} />
      </div>

      <hr class="col-span-full border-gray-200" />

      <div class="flex items-center gap-6 col-span-full">
        <label class="flex items-center gap-2 cursor-pointer">
          <input type="checkbox" class="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500" bind:checked={form.isRecurring} />
          <span class="text-sm font-medium text-gray-700">Recurring Expense</span>
        </label>
        <label class="flex items-center gap-2 cursor-pointer">
          <input type="checkbox" class="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500" bind:checked={form.taxDeductible} />
          <span class="text-sm font-medium text-gray-700">Tax Deductible</span>
        </label>
      </div>

      {#if form.isRecurring}
        <div>
          <label class="label" for="recurringInterval">Recurring Interval</label>
          <select id="recurringInterval" class="input-field" bind:value={form.recurringInterval}>
            <option value="monthly">Monthly</option>
            <option value="quarterly">Quarterly</option>
            <option value="annually">Annually</option>
          </select>
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
    <button class="btn-primary" type="submit" form="expense-form" disabled={!canSave}>
      {editingId ? 'Save Changes' : 'Add Expense'}
    </button>
  </svelte:fragment>
</Modal>

<!-- Delete Confirmation -->
<ConfirmDialog
  bind:show={showConfirm}
  title="Delete Expense"
  message="Are you sure you want to delete this {deleteLabel}? This action cannot be undone."
  confirmText="Delete"
  danger={true}
  on:confirm={confirmDelete}
/>
