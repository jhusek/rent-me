<script lang="ts">
  import { onMount } from 'svelte';
  import {
    GetPaymentSummary,
    GetExpenseSummary,
    GetProperty,
    GetPayments,
    GetExpenses,
  } from '../../wailsjs/go/main/App';
  import PageHeader from '../components/PageHeader.svelte';
  import Toast from '../components/Toast.svelte';
  import { formatCurrency, formatDate, getCurrentYear, getCategoryLabel, EXPENSE_CATEGORIES } from '../lib/utils';

  let loading = true;
  let selectedYear = getCurrentYear();

  // Data
  let paymentSummary: any = { totalDue: 0, totalPaid: 0, totalLate: 0, totalMissed: 0, totalPartial: 0, lateFees: 0, count: 0 };
  let expenseSummary: any = { totalAmount: 0, byCategory: {}, count: 0, taxDeductible: 0 };
  let property: any = {};
  let payments: any[] = [];
  let expenses: any[] = [];

  // Toast state
  let toastMessage = '';
  let toastType: 'success' | 'error' = 'success';
  let showToast = false;

  function toast(msg: string, type: 'success' | 'error' = 'success') {
    toastMessage = msg;
    toastType = type;
    showToast = true;
  }

  // Computed
  $: totalIncome = paymentSummary.totalPaid || 0;
  $: totalExpenses = expenseSummary.totalAmount || 0;
  $: netIncome = totalIncome - totalExpenses;
  $: incomeBarWidth = totalIncome + totalExpenses > 0 ? (totalIncome / (totalIncome + totalExpenses)) * 100 : 50;

  // Schedule E calculations
  $: byCategory = expenseSummary.byCategory || {};
  $: scheduleE = {
    line3: totalIncome,
    line5: (byCategory['advertising'] as number) || 0,
    line7: ((byCategory['maintenance'] as number) || 0),
    line9: (byCategory['insurance'] as number) || 0,
    line10: (byCategory['legal'] as number) || 0,
    line11: (byCategory['management'] as number) || 0,
    line12: (byCategory['mortgage'] as number) || 0,
    line14: (byCategory['repairs'] as number) || 0,
    line16: (byCategory['property_tax'] as number) || 0,
    line17: (byCategory['utilities'] as number) || 0,
    line19: ((byCategory['hoa'] as number) || 0)
          + ((byCategory['supplies'] as number) || 0)
          + ((byCategory['travel'] as number) || 0)
          + ((byCategory['depreciation'] as number) || 0)
          + ((byCategory['other'] as number) || 0),
  };
  $: scheduleETotal = scheduleE.line5 + scheduleE.line7 + scheduleE.line9 + scheduleE.line10
    + scheduleE.line11 + scheduleE.line12 + scheduleE.line14 + scheduleE.line16
    + scheduleE.line17 + scheduleE.line19;
  $: scheduleENet = scheduleE.line3 - scheduleETotal;

  // Depreciation
  $: purchasePrice = property.purchasePrice || 0;
  $: depPeriod = 27.5;
  $: annualDep = purchasePrice > 0 ? purchasePrice / depPeriod : 0;
  $: yearsOwned = (() => {
    if (!property.purchaseDate) return 0;
    const pd = new Date(property.purchaseDate + 'T00:00:00');
    const now = new Date();
    const diff = (now.getTime() - pd.getTime()) / (1000 * 60 * 60 * 24 * 365.25);
    return Math.max(0, diff);
  })();
  $: accumulatedDep = Math.min(annualDep * yearsOwned, purchasePrice);
  $: remainingBasis = Math.max(0, purchasePrice - accumulatedDep);
  $: depProgress = purchasePrice > 0 ? (accumulatedDep / purchasePrice) * 100 : 0;

  // Expense breakdown for chart
  $: expenseBreakdown = (() => {
    const bc = byCategory as Record<string, number>;
    const entries = Object.entries(bc).filter(([, v]) => v > 0).sort((a, b) => b[1] - a[1]);
    const total = entries.reduce((s, [, v]) => s + v, 0);
    return entries.map(([cat, amount]) => ({
      category: cat,
      label: getCategoryLabel(cat),
      amount,
      pct: total > 0 ? (amount / total) * 100 : 0,
    }));
  })();

  const CATEGORY_COLORS: Record<string, string> = {
    mortgage: 'bg-blue-500',
    insurance: 'bg-cyan-500',
    property_tax: 'bg-yellow-500',
    repairs: 'bg-red-500',
    maintenance: 'bg-orange-500',
    utilities: 'bg-teal-500',
    hoa: 'bg-gray-500',
    management: 'bg-purple-500',
    supplies: 'bg-pink-500',
    legal: 'bg-rose-500',
    advertising: 'bg-indigo-500',
    travel: 'bg-amber-500',
    depreciation: 'bg-slate-500',
    other: 'bg-gray-400',
  };

  async function loadAll() {
    loading = true;
    try {
      const [ps, es, prop, pay, exp] = await Promise.all([
        GetPaymentSummary(selectedYear, 0).catch(() => null),
        GetExpenseSummary(selectedYear).catch(() => null),
        GetProperty().catch(() => null),
        GetPayments(0, selectedYear, 0).catch(() => []),
        GetExpenses('', selectedYear).catch(() => []),
      ]);
      paymentSummary = ps || { totalDue: 0, totalPaid: 0, totalLate: 0, totalMissed: 0, totalPartial: 0, lateFees: 0, count: 0 };
      expenseSummary = es || { totalAmount: 0, byCategory: {}, count: 0, taxDeductible: 0 };
      property = prop || {};
      payments = pay || [];
      expenses = exp || [];
    } catch {
      // defaults already set
    }
    loading = false;
  }

  function changeYear() {
    loadAll();
  }

  // CSV export helpers
  function downloadCSV(filename: string, csvContent: string) {
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    a.click();
    URL.revokeObjectURL(url);
  }

  function exportPaymentsCSV() {
    if (payments.length === 0) {
      toast('No payments to export', 'error');
      return;
    }
    const headers = ['ID', 'Tenant ID', 'Amount', 'Due Date', 'Paid Date', 'Method', 'Status', 'Late Fee', 'Notes'];
    const rows = payments.map((p: any) => [
      p.id, p.tenantId, p.amount, p.dueDate, p.paidDate, p.method, p.status, p.lateFee, `"${(p.notes || '').replace(/"/g, '""')}"`
    ].join(','));
    downloadCSV(`payments_${selectedYear}.csv`, [headers.join(','), ...rows].join('\n'));
    toast('Payments CSV exported');
  }

  function exportExpensesCSV() {
    if (expenses.length === 0) {
      toast('No expenses to export', 'error');
      return;
    }
    const headers = ['ID', 'Category', 'Amount', 'Date', 'Vendor', 'Description', 'Tax Deductible', 'Notes'];
    const rows = expenses.map((e: any) => [
      e.id, e.category, e.amount, e.date, `"${(e.vendor || '').replace(/"/g, '""')}"`,
      `"${(e.description || '').replace(/"/g, '""')}"`, e.taxDeductible ? 'Yes' : 'No',
      `"${(e.notes || '').replace(/"/g, '""')}"`
    ].join(','));
    downloadCSV(`expenses_${selectedYear}.csv`, [headers.join(','), ...rows].join('\n'));
    toast('Expenses CSV exported');
  }

  const yearOptions: number[] = [];
  const currentYear = getCurrentYear();
  for (let y = currentYear - 5; y <= currentYear + 1; y++) {
    yearOptions.push(y);
  }

  const scheduleELines = [
    { line: '3', label: 'Rents received', key: 'line3' },
    { line: '5', label: 'Advertising', key: 'line5' },
    { line: '7', label: 'Cleaning and maintenance', key: 'line7' },
    { line: '9', label: 'Insurance', key: 'line9' },
    { line: '10', label: 'Legal and professional fees', key: 'line10' },
    { line: '11', label: 'Management fees', key: 'line11' },
    { line: '12', label: 'Mortgage interest', key: 'line12' },
    { line: '14', label: 'Repairs', key: 'line14' },
    { line: '16', label: 'Taxes', key: 'line16' },
    { line: '17', label: 'Utilities', key: 'line17' },
    { line: '19', label: 'Other', key: 'line19' },
  ];

  onMount(loadAll);
</script>

<Toast bind:show={showToast} message={toastMessage} type={toastType} />

<PageHeader title="Reports">
  <svelte:fragment slot="actions">
    <div class="flex items-center gap-3">
      <select class="input-field w-auto" bind:value={selectedYear} on:change={changeYear}>
        {#each yearOptions as y}
          <option value={y}>{y}</option>
        {/each}
      </select>
      <button class="btn-secondary" on:click={exportPaymentsCSV}>📥 Payments CSV</button>
      <button class="btn-secondary" on:click={exportExpensesCSV}>📥 Expenses CSV</button>
    </div>
  </svelte:fragment>
</PageHeader>

{#if loading}
  <div class="flex items-center justify-center py-20">
    <div class="text-gray-400 text-lg">Loading...</div>
  </div>
{:else}
  <div class="space-y-6">

    <!-- Section 1: Annual Income Summary -->
    <div class="card">
      <h2 class="text-lg font-semibold text-gray-900 mb-4">📊 Annual Income Summary — {selectedYear}</h2>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
        <div class="bg-green-50 rounded-lg p-4 border border-green-200">
          <p class="text-xs font-medium text-green-600 uppercase">Total Rental Income</p>
          <p class="text-2xl font-bold text-green-700 mt-1">{formatCurrency(totalIncome)}</p>
        </div>
        <div class="bg-red-50 rounded-lg p-4 border border-red-200">
          <p class="text-xs font-medium text-red-600 uppercase">Total Expenses</p>
          <p class="text-2xl font-bold text-red-700 mt-1">{formatCurrency(totalExpenses)}</p>
        </div>
        <div class="rounded-lg p-4 border {netIncome >= 0 ? 'bg-emerald-50 border-emerald-200' : 'bg-red-50 border-red-200'}">
          <p class="text-xs font-medium uppercase {netIncome >= 0 ? 'text-emerald-600' : 'text-red-600'}">Net Income</p>
          <p class="text-2xl font-bold mt-1 {netIncome >= 0 ? 'text-emerald-700' : 'text-red-700'}">{formatCurrency(netIncome)}</p>
        </div>
      </div>
      <!-- Income vs Expenses bar -->
      <div>
        <div class="flex justify-between text-xs text-gray-500 mb-1">
          <span>Income</span>
          <span>Expenses</span>
        </div>
        <div class="h-4 bg-gray-200 rounded-full overflow-hidden flex">
          <div class="bg-green-500 transition-all duration-500" style="width: {incomeBarWidth}%"></div>
          <div class="bg-red-400 transition-all duration-500" style="width: {100 - incomeBarWidth}%"></div>
        </div>
        <div class="flex justify-between text-xs text-gray-400 mt-1">
          <span>{formatCurrency(totalIncome)}</span>
          <span>{formatCurrency(totalExpenses)}</span>
        </div>
      </div>
    </div>

    <!-- Section 2: Schedule E Helper -->
    <div class="card">
      <h2 class="text-lg font-semibold text-gray-900 mb-1">📝 Schedule E Helper</h2>
      <p class="text-sm text-gray-500 mb-4">Formatted like IRS Schedule E (Supplemental Income and Loss)</p>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b-2 border-gray-300">
              <th class="text-left py-2 px-3 text-xs font-medium text-gray-500 uppercase">Line</th>
              <th class="text-left py-2 px-3 text-xs font-medium text-gray-500 uppercase">Description</th>
              <th class="text-right py-2 px-3 text-xs font-medium text-gray-500 uppercase">Amount</th>
            </tr>
          </thead>
          <tbody>
            {#each scheduleELines as item}
              <tr class="border-b border-gray-100 {item.line === '3' ? 'bg-green-50' : ''}">
                <td class="py-2 px-3 text-sm text-gray-500 font-mono">{item.line}</td>
                <td class="py-2 px-3 text-sm text-gray-700">{item.label}</td>
                <td class="py-2 px-3 text-sm text-right font-medium {item.line === '3' ? 'text-green-700' : 'text-gray-900'}">
                  {formatCurrency(scheduleE[item.key])}
                </td>
              </tr>
            {/each}
            <tr class="border-t-2 border-gray-300 bg-gray-50">
              <td class="py-2 px-3 text-sm font-mono text-gray-500">20</td>
              <td class="py-2 px-3 text-sm font-semibold text-gray-900">Total expenses</td>
              <td class="py-2 px-3 text-sm text-right font-bold text-gray-900">{formatCurrency(scheduleETotal)}</td>
            </tr>
            <tr class="bg-gray-50 {scheduleENet >= 0 ? '' : ''}">
              <td class="py-2 px-3 text-sm font-mono text-gray-500">21</td>
              <td class="py-2 px-3 text-sm font-semibold text-gray-900">Net income (loss)</td>
              <td class="py-2 px-3 text-sm text-right font-bold {scheduleENet >= 0 ? 'text-green-700' : 'text-red-700'}">
                {formatCurrency(scheduleENet)}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Section 3: Depreciation Tracker -->
    <div class="card">
      <h2 class="text-lg font-semibold text-gray-900 mb-4">🏗️ Depreciation Tracker</h2>
      {#if purchasePrice > 0}
        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4 mb-6">
          <div>
            <p class="text-xs text-gray-500 uppercase">Purchase Price</p>
            <p class="text-lg font-bold text-gray-900">{formatCurrency(purchasePrice)}</p>
          </div>
          <div>
            <p class="text-xs text-gray-500 uppercase">Period</p>
            <p class="text-lg font-bold text-gray-900">27.5 years</p>
          </div>
          <div>
            <p class="text-xs text-gray-500 uppercase">Annual Depreciation</p>
            <p class="text-lg font-bold text-blue-700">{formatCurrency(annualDep)}</p>
          </div>
          <div>
            <p class="text-xs text-gray-500 uppercase">Years Owned</p>
            <p class="text-lg font-bold text-gray-900">{yearsOwned.toFixed(1)}</p>
          </div>
          <div>
            <p class="text-xs text-gray-500 uppercase">Accumulated</p>
            <p class="text-lg font-bold text-orange-700">{formatCurrency(accumulatedDep)}</p>
          </div>
          <div>
            <p class="text-xs text-gray-500 uppercase">Remaining Basis</p>
            <p class="text-lg font-bold text-green-700">{formatCurrency(remainingBasis)}</p>
          </div>
        </div>
        <!-- Progress bar -->
        <div>
          <div class="flex justify-between text-xs text-gray-500 mb-1">
            <span>Depreciation Progress</span>
            <span>{depProgress.toFixed(1)}%</span>
          </div>
          <div class="h-3 bg-gray-200 rounded-full overflow-hidden">
            <div class="h-full bg-blue-500 rounded-full transition-all duration-500" style="width: {Math.min(depProgress, 100)}%"></div>
          </div>
          <div class="flex justify-between text-xs text-gray-400 mt-1">
            <span>{formatCurrency(accumulatedDep)} depreciated</span>
            <span>{formatCurrency(remainingBasis)} remaining</span>
          </div>
        </div>
      {:else}
        <p class="text-gray-400 text-sm">Add a purchase price to your property to see depreciation calculations.</p>
      {/if}
    </div>

    <!-- Section 4: Expense Breakdown -->
    <div class="card">
      <h2 class="text-lg font-semibold text-gray-900 mb-4">📊 Expense Breakdown — {selectedYear}</h2>
      {#if expenseBreakdown.length === 0}
        <p class="text-gray-400 text-sm">No expenses recorded for {selectedYear}.</p>
      {:else}
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="border-b-2 border-gray-300">
                <th class="text-left py-2 px-3 text-xs font-medium text-gray-500 uppercase">Category</th>
                <th class="text-right py-2 px-3 text-xs font-medium text-gray-500 uppercase w-28">Amount</th>
                <th class="text-right py-2 px-3 text-xs font-medium text-gray-500 uppercase w-16">%</th>
                <th class="py-2 px-3 text-xs font-medium text-gray-500 uppercase w-48"></th>
              </tr>
            </thead>
            <tbody>
              {#each expenseBreakdown as item}
                <tr class="border-b border-gray-100">
                  <td class="py-2.5 px-3 text-sm text-gray-700">{item.label}</td>
                  <td class="py-2.5 px-3 text-sm text-right font-medium text-gray-900">{formatCurrency(item.amount)}</td>
                  <td class="py-2.5 px-3 text-sm text-right text-gray-500">{item.pct.toFixed(1)}%</td>
                  <td class="py-2.5 px-3">
                    <div class="h-2.5 bg-gray-100 rounded-full overflow-hidden">
                      <div
                        class="h-full rounded-full transition-all duration-500 {CATEGORY_COLORS[item.category] || 'bg-gray-400'}"
                        style="width: {item.pct}%"
                      ></div>
                    </div>
                  </td>
                </tr>
              {/each}
            </tbody>
            <tfoot>
              <tr class="border-t-2 border-gray-300">
                <td class="py-2 px-3 text-sm font-semibold text-gray-900">Total</td>
                <td class="py-2 px-3 text-sm text-right font-bold text-gray-900">{formatCurrency(totalExpenses)}</td>
                <td class="py-2 px-3 text-sm text-right text-gray-500">100%</td>
                <td></td>
              </tr>
            </tfoot>
          </table>
        </div>
      {/if}
    </div>

  </div>
{/if}
