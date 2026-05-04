<script lang="ts">
  import { onMount } from 'svelte';
  import { GetDashboardData } from '../../wailsjs/go/main/App';
  import type { services } from '../../wailsjs/go/models';
  import BarChart from '../components/BarChart.svelte';
  import DoughnutChart from '../components/DoughnutChart.svelte';
  import Badge from '../components/Badge.svelte';
  import { formatCurrency, formatDate, getCategoryLabel, timeAgo } from '../lib/utils';

  let loading = true;
  let error = '';
  let data: services.DashboardData | null = null;

  const CATEGORY_COLORS: Record<string, string> = {
    mortgage: '#3b82f6',
    insurance: '#8b5cf6',
    property_tax: '#f59e0b',
    repairs: '#ef4444',
    maintenance: '#f97316',
    utilities: '#06b6d4',
    hoa: '#6b7280',
    management: '#84cc16',
    supplies: '#14b8a6',
    legal: '#dc2626',
    advertising: '#a855f7',
    travel: '#ec4899',
    depreciation: '#eab308',
    other: '#9ca3af',
  };

  const REMINDER_ICONS: Record<string, string> = {
    lease: '📋',
    payment: '💰',
    maintenance: '🔧',
    insurance: '🛡️',
    tax: '🏛️',
    inspection: '🔍',
    other: '📌',
  };

  const ACTION_ICONS: Record<string, string> = {
    created: '✨',
    updated: '✏️',
    deleted: '🗑️',
    paid: '💵',
    completed: '✅',
  };

  function formatMonthLabel(monthStr: string): string {
    if (!monthStr) return '';
    const [year, month] = monthStr.split('-');
    const shortMonths = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec'];
    const idx = parseInt(month, 10) - 1;
    return idx >= 0 && idx < 12 ? shortMonths[idx] : monthStr;
  }

  // Rent status helpers
  function getRentStatusText(rs: services.DashboardData['rentStatus']): string {
    if (!rs) return '—';
    const late = (rs.late || 0) + (rs.missed || 0);
    const due = (rs.due || 0) + (rs.partial || 0);
    if (late > 0) return `${late} Late`;
    if (due > 0) return `${due} Due`;
    if (rs.paid > 0) return 'Paid ✓';
    return 'No Data';
  }

  function getRentStatusColor(rs: services.DashboardData['rentStatus']): string {
    if (!rs) return 'text-gray-400';
    const late = (rs.late || 0) + (rs.missed || 0);
    const due = (rs.due || 0) + (rs.partial || 0);
    if (late > 0) return 'text-red-600';
    if (due > 0) return 'text-yellow-600';
    if (rs.paid > 0) return 'text-green-600';
    return 'text-gray-400';
  }

  function getRentBorderColor(rs: services.DashboardData['rentStatus']): string {
    if (!rs) return 'border-l-gray-300';
    const late = (rs.late || 0) + (rs.missed || 0);
    const due = (rs.due || 0) + (rs.partial || 0);
    if (late > 0) return 'border-l-red-500';
    if (due > 0) return 'border-l-yellow-500';
    if (rs.paid > 0) return 'border-l-green-500';
    return 'border-l-gray-300';
  }

  // Derived chart data
  $: barLabels = data?.incomeByMonth?.map(m => formatMonthLabel(m.month)) || [];
  $: barDatasets = [
    {
      label: 'Income',
      data: data?.incomeByMonth?.map(m => m.amount) || [],
      backgroundColor: 'rgba(59, 130, 246, 0.8)',
    },
    {
      label: 'Expenses',
      data: data?.expensesByMonth?.map(m => m.amount) || [],
      backgroundColor: 'rgba(239, 68, 68, 0.8)',
    },
  ];
  $: hasBarData = (data?.incomeByMonth?.some(m => m.amount > 0) || data?.expensesByMonth?.some(m => m.amount > 0)) ?? false;

  $: doughnutLabels = data?.expensesByCategory?.map(c => getCategoryLabel(c.category)) || [];
  $: doughnutData = data?.expensesByCategory?.map(c => c.amount) || [];
  $: doughnutColors = data?.expensesByCategory?.map(c => CATEGORY_COLORS[c.category] || CATEGORY_COLORS.other) || [];
  $: hasDoughnutData = (data?.expensesByCategory?.length ?? 0) > 0;

  $: isFirstRun = data && !data.property?.address && !data.currentTenant && data.ytdIncome === 0 && data.ytdExpenses === 0;

  $: netSign = (data?.netIncome ?? 0) >= 0 ? '+' : '';

  $: propertyAddress = data?.property
    ? [data.property.address, data.property.city, data.property.state, data.property.zip].filter(Boolean).join(', ')
    : '';

  $: tenantName = data?.currentTenant
    ? `${data.currentTenant.firstName} ${data.currentTenant.lastName}`
    : 'No active tenant';

  async function loadDashboard() {
    loading = true;
    error = '';
    try {
      data = await GetDashboardData();
    } catch (e) {
      error = String(e);
    } finally {
      loading = false;
    }
  }

  onMount(loadDashboard);
</script>

{#if loading}
  <div class="flex items-center justify-center py-20">
    <div class="flex flex-col items-center gap-3">
      <div class="w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
      <p class="text-gray-500">Loading dashboard...</p>
    </div>
  </div>
{:else if error}
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
    </div>
    <div class="card text-center py-12">
      <p class="text-red-500 text-lg mb-2">Failed to load dashboard</p>
      <p class="text-gray-500 text-sm mb-4">{error}</p>
      <button class="btn-primary" on:click={loadDashboard}>Try Again</button>
    </div>
  </div>
{:else if data && isFirstRun}
  <div class="space-y-6">
    <div class="card text-center py-16">
      <p class="text-5xl mb-4">🏡</p>
      <h1 class="text-3xl font-bold text-gray-900 mb-2">Welcome to RentMe!</h1>
      <p class="text-gray-500 mb-8 max-w-md mx-auto">Get started by setting up your rental property in a few quick steps.</p>
      <div class="max-w-sm mx-auto space-y-3 text-left">
        <a href="#/property" class="flex items-center gap-4 p-4 rounded-lg bg-blue-50 hover:bg-blue-100 transition-colors group">
          <span class="flex items-center justify-center w-8 h-8 rounded-full bg-blue-600 text-white text-sm font-bold">1</span>
          <div>
            <p class="font-semibold text-gray-900 group-hover:text-blue-700">Set up your property</p>
            <p class="text-sm text-gray-500">Add address and details</p>
          </div>
          <span class="ml-auto text-gray-400 group-hover:text-blue-600">→</span>
        </a>
        <a href="#/tenants" class="flex items-center gap-4 p-4 rounded-lg bg-green-50 hover:bg-green-100 transition-colors group">
          <span class="flex items-center justify-center w-8 h-8 rounded-full bg-green-600 text-white text-sm font-bold">2</span>
          <div>
            <p class="font-semibold text-gray-900 group-hover:text-green-700">Add your tenant</p>
            <p class="text-sm text-gray-500">Lease and contact info</p>
          </div>
          <span class="ml-auto text-gray-400 group-hover:text-green-600">→</span>
        </a>
        <a href="#/payments" class="flex items-center gap-4 p-4 rounded-lg bg-yellow-50 hover:bg-yellow-100 transition-colors group">
          <span class="flex items-center justify-center w-8 h-8 rounded-full bg-yellow-600 text-white text-sm font-bold">3</span>
          <div>
            <p class="font-semibold text-gray-900 group-hover:text-yellow-700">Record your first payment</p>
            <p class="text-sm text-gray-500">Track rent and income</p>
          </div>
          <span class="ml-auto text-gray-400 group-hover:text-yellow-600">→</span>
        </a>
      </div>
    </div>
  </div>
{:else if data}
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p class="text-sm text-gray-500 mt-1">Overview of your rental property</p>
      </div>
      <button class="btn-secondary text-sm" on:click={loadDashboard}>↻ Refresh</button>
    </div>

    <!-- Row 1: Summary Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Monthly Income -->
      <div class="card !p-4 border-l-4 border-l-green-500">
        <p class="text-xs font-medium text-gray-500 uppercase tracking-wide">Monthly Income</p>
        <p class="text-2xl font-bold text-green-700 mt-1">{formatCurrency(data.monthlyIncome)}</p>
        <p class="text-xs text-gray-400 mt-1">This month collected</p>
      </div>
      <!-- Monthly Expenses -->
      <div class="card !p-4 border-l-4 border-l-red-500">
        <p class="text-xs font-medium text-gray-500 uppercase tracking-wide">Monthly Expenses</p>
        <p class="text-2xl font-bold text-red-700 mt-1">{formatCurrency(data.monthlyExpenses)}</p>
        <p class="text-xs text-gray-400 mt-1">This month spent</p>
      </div>
      <!-- Net Income (YTD) -->
      <div class="card !p-4 border-l-4 border-l-blue-500">
        <p class="text-xs font-medium text-gray-500 uppercase tracking-wide">Net Income (YTD)</p>
        <p class="text-2xl font-bold mt-1" class:text-green-700={data.netIncome >= 0} class:text-red-700={data.netIncome < 0}>
          {netSign}{formatCurrency(Math.abs(data.netIncome))}
        </p>
        <p class="text-xs text-gray-400 mt-1">{formatCurrency(data.ytdIncome)} income — {formatCurrency(data.ytdExpenses)} expenses</p>
      </div>
      <!-- Rent Status -->
      <div class="card !p-4 border-l-4 {getRentBorderColor(data.rentStatus)}">
        <p class="text-xs font-medium text-gray-500 uppercase tracking-wide">Rent Status</p>
        <p class="text-2xl font-bold mt-1 {getRentStatusColor(data.rentStatus)}">
          {getRentStatusText(data.rentStatus)}
        </p>
        {#if data.rentStatus}
          <p class="text-xs text-gray-400 mt-1">
            {data.rentStatus.paid || 0} paid · {(data.rentStatus.due || 0) + (data.rentStatus.partial || 0)} pending · {(data.rentStatus.late || 0) + (data.rentStatus.missed || 0)} overdue
          </p>
        {/if}
      </div>
    </div>

    <!-- Row 2: Charts -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Income vs Expenses Bar Chart -->
      <div class="card">
        <h2 class="text-lg font-semibold mb-4">📊 Income vs Expenses</h2>
        {#if hasBarData}
          <div class="h-72">
            <BarChart labels={barLabels} datasets={barDatasets} />
          </div>
        {:else}
          <div class="h-72 flex items-center justify-center text-gray-400">
            <div class="text-center">
              <p class="text-4xl mb-2">📈</p>
              <p>No income or expense data yet</p>
              <p class="text-sm mt-1">Charts will appear once you record transactions</p>
            </div>
          </div>
        {/if}
      </div>
      <!-- Expense Breakdown Doughnut -->
      <div class="card">
        <h2 class="text-lg font-semibold mb-4">🍩 Expense Breakdown</h2>
        {#if hasDoughnutData}
          <div class="h-72">
            <DoughnutChart labels={doughnutLabels} data={doughnutData} colors={doughnutColors} />
          </div>
        {:else}
          <div class="h-72 flex items-center justify-center text-gray-400">
            <div class="text-center">
              <p class="text-4xl mb-2">📊</p>
              <p>No expense categories to display</p>
              <p class="text-sm mt-1">Add expenses to see the breakdown</p>
            </div>
          </div>
        {/if}
      </div>
    </div>

    <!-- Row 3: Activity & Alerts -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Upcoming Reminders -->
      <div class="card">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold">🔔 Upcoming Reminders</h2>
          <a href="#/reminders" class="text-sm text-blue-600 hover:text-blue-800 font-medium">View All →</a>
        </div>
        {#if data.upcomingReminders && data.upcomingReminders.length > 0}
          <div class="space-y-3">
            {#each data.upcomingReminders.slice(0, 5) as reminder}
              <div class="flex items-start gap-3 p-3 rounded-lg bg-gray-50 hover:bg-gray-100 transition-colors">
                <span class="text-xl flex-shrink-0 mt-0.5">{REMINDER_ICONS[reminder.type] || REMINDER_ICONS.other}</span>
                <div class="flex-1 min-w-0">
                  <p class="font-medium text-gray-900 truncate">{reminder.title}</p>
                  <p class="text-sm text-gray-500">Due {formatDate(reminder.dueDate)}</p>
                </div>
                <Badge variant={reminder.type === 'payment' ? 'warning' : reminder.type === 'maintenance' ? 'danger' : 'info'} text={reminder.type} />
              </div>
            {/each}
          </div>
        {:else}
          <div class="flex flex-col items-center justify-center py-8 text-gray-400">
            <p class="text-3xl mb-2">🎉</p>
            <p class="text-sm">No upcoming reminders</p>
          </div>
        {/if}
      </div>

      <!-- Recent Activity -->
      <div class="card">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold">📋 Recent Activity</h2>
        </div>
        {#if data.recentActivity && data.recentActivity.length > 0}
          <div class="space-y-3">
            {#each data.recentActivity.slice(0, 10) as activity}
              <div class="flex items-start gap-3 p-2 rounded-lg hover:bg-gray-50 transition-colors">
                <span class="text-lg flex-shrink-0 mt-0.5">{ACTION_ICONS[activity.action] || '📝'}</span>
                <div class="flex-1 min-w-0">
                  <p class="text-sm text-gray-900">{activity.description}</p>
                  <p class="text-xs text-gray-400">{timeAgo(activity.createdAt)}</p>
                </div>
              </div>
            {/each}
          </div>
        {:else}
          <div class="flex flex-col items-center justify-center py-8 text-gray-400">
            <p class="text-3xl mb-2">📭</p>
            <p class="text-sm">No recent activity</p>
          </div>
        {/if}
      </div>
    </div>

    <!-- Row 4: Quick Status -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Property Card -->
      <div class="card">
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-lg font-semibold">🏠 Property</h2>
          <a href="#/property" class="text-sm text-blue-600 hover:text-blue-800 font-medium">View Details →</a>
        </div>
        <div class="space-y-2">
          <div>
            <p class="text-sm text-gray-500">Address</p>
            <p class="font-medium text-gray-900">{propertyAddress || 'Not configured'}</p>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <p class="text-sm text-gray-500">Current Tenant</p>
              <p class="font-medium text-gray-900">{tenantName}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">Lease Ends</p>
              <p class="font-medium text-gray-900">{data.currentTenant?.leaseEnd ? formatDate(data.currentTenant.leaseEnd) : '—'}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Open Maintenance -->
      <div class="card">
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-lg font-semibold">🔧 Maintenance</h2>
          <a href="#/maintenance" class="text-sm text-blue-600 hover:text-blue-800 font-medium">View All →</a>
        </div>
        <div class="flex items-center gap-4">
          <div class="flex items-center justify-center w-16 h-16 rounded-full {data.openMaintenance > 0 ? 'bg-yellow-100' : 'bg-green-100'}">
            <span class="text-2xl font-bold {data.openMaintenance > 0 ? 'text-yellow-700' : 'text-green-700'}">{data.openMaintenance}</span>
          </div>
          <div>
            <p class="font-medium text-gray-900">
              {#if data.openMaintenance === 0}
                All clear!
              {:else if data.openMaintenance === 1}
                1 open request
              {:else}
                {data.openMaintenance} open requests
              {/if}
            </p>
            <p class="text-sm text-gray-500">
              {#if data.openMaintenance === 0}
                No maintenance issues pending
              {:else}
                Needs attention
              {/if}
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
{/if}
