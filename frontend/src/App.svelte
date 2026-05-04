<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import Router, { link, location } from 'svelte-spa-router';
  import { GetUpcomingReminders } from '../wailsjs/go/main/App';
  
  import Dashboard from './pages/Dashboard.svelte';
  import Property from './pages/Property.svelte';
  import Tenants from './pages/Tenants.svelte';
  import Payments from './pages/Payments.svelte';
  import Expenses from './pages/Expenses.svelte';
  import Maintenance from './pages/Maintenance.svelte';
  import Documents from './pages/Documents.svelte';
  import Reminders from './pages/Reminders.svelte';
  import Reports from './pages/Reports.svelte';
  import Contacts from './pages/Contacts.svelte';
  import Settings from './pages/Settings.svelte';

  const routes = {
    '/': Dashboard,
    '/property': Property,
    '/tenants': Tenants,
    '/payments': Payments,
    '/expenses': Expenses,
    '/maintenance': Maintenance,
    '/documents': Documents,
    '/reminders': Reminders,
    '/reports': Reports,
    '/contacts': Contacts,
    '/settings': Settings,
  };

  const navItems = [
    { path: '/', label: 'Dashboard', icon: '📊' },
    { path: '/property', label: 'Property', icon: '🏠' },
    { path: '/tenants', label: 'Tenants', icon: '👤' },
    { path: '/payments', label: 'Payments', icon: '💰' },
    { path: '/expenses', label: 'Expenses', icon: '🧾' },
    { path: '/maintenance', label: 'Maintenance', icon: '🔧' },
    { path: '/documents', label: 'Documents', icon: '📁' },
    { path: '/reminders', label: 'Reminders', icon: '🔔' },
    { path: '/reports', label: 'Reports', icon: '📈' },
    { path: '/contacts', label: 'Contacts', icon: '📞' },
  ];

  $: currentPath = $location;

  let reminderCount = 0;
  let refreshInterval: ReturnType<typeof setInterval>;

  async function refreshReminderCount() {
    try {
      const upcoming = await GetUpcomingReminders(7);
      reminderCount = (upcoming || []).length;
    } catch {
      reminderCount = 0;
    }
  }

  onMount(() => {
    refreshReminderCount();
    refreshInterval = setInterval(refreshReminderCount, 60000);
  });

  onDestroy(() => {
    if (refreshInterval) clearInterval(refreshInterval);
  });

  // Refresh count when navigating away from reminders page
  $: if (currentPath) {
    refreshReminderCount();
  }
</script>

<div class="flex h-screen overflow-hidden">
  <!-- Sidebar -->
  <aside class="w-64 bg-sidebar flex flex-col flex-shrink-0">
    <!-- Logo -->
    <div class="p-6 border-b border-gray-700">
      <h1 class="text-xl font-bold text-white flex items-center gap-2">
        🏡 RentMe
      </h1>
      <p class="text-xs text-gray-400 mt-1">Property Manager</p>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 py-4 overflow-y-auto">
      {#each navItems as item}
        <a
          href="#{item.path}"
          use:link
          class="flex items-center gap-3 px-6 py-2.5 text-sm transition-colors duration-150
            {currentPath === item.path 
              ? 'bg-sidebar-active text-white border-r-2 border-primary-500' 
              : 'text-gray-300 hover:bg-sidebar-hover hover:text-white'}"
        >
          <span class="text-lg">{item.icon}</span>
          <span>{item.label}</span>
          {#if item.path === '/reminders' && reminderCount > 0}
            <span class="ml-auto inline-flex items-center justify-center w-5 h-5 text-xs font-bold text-white bg-red-500 rounded-full">
              {reminderCount > 9 ? '9+' : reminderCount}
            </span>
          {/if}
        </a>
      {/each}
    </nav>

    <!-- Footer -->
    <div class="p-4 border-t border-gray-700">
      <a
        href="#/settings"
        use:link
        class="flex items-center gap-3 px-2 py-2 text-sm text-gray-400 hover:text-white transition-colors rounded-lg hover:bg-sidebar-hover"
      >
        <span class="text-lg">⚙️</span>
        <span>Settings</span>
      </a>
      <div class="text-xs text-gray-600 mt-2 px-2">v0.1.0</div>
    </div>
  </aside>

  <!-- Main Content -->
  <main class="flex-1 overflow-y-auto bg-gray-50">
    <div class="p-8">
      <Router {routes} />
    </div>
  </main>
</div>
