<script lang="ts">
  import { onMount } from 'svelte';

  const STORAGE_KEY = 'rentme_settings';

  interface LateFeeSettings {
    gracePeriod: number;
    feeType: 'flat' | 'percentage';
    feeAmount: number;
  }

  let settings: LateFeeSettings = {
    gracePeriod: 5,
    feeType: 'flat',
    feeAmount: 50,
  };

  let saved = false;

  function loadSettings() {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (raw) {
        const parsed = JSON.parse(raw);
        settings = { ...settings, ...parsed };
      }
    } catch {
      // use defaults
    }
  }

  function saveSettings() {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(settings));
    saved = true;
    setTimeout(() => (saved = false), 2000);
  }

  onMount(loadSettings);
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <h1 class="text-2xl font-bold text-gray-900">⚙️ Settings</h1>
  </div>

  <!-- Late Fee Configuration -->
  <div class="card">
    <h2 class="text-lg font-semibold mb-4">💰 Late Fee Configuration</h2>
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div>
        <label for="grace-period" class="block text-sm font-medium text-gray-700 mb-1">Grace Period (days)</label>
        <input
          id="grace-period"
          type="number"
          min="0"
          max="30"
          bind:value={settings.gracePeriod}
          class="input-field"
          placeholder="5"
        />
        <p class="text-xs text-gray-400 mt-1">Days after due date before late fee applies</p>
      </div>
      <div>
        <label for="fee-type" class="block text-sm font-medium text-gray-700 mb-1">Late Fee Type</label>
        <select id="fee-type" bind:value={settings.feeType} class="input-field">
          <option value="flat">Flat Amount ($)</option>
          <option value="percentage">Percentage (%)</option>
        </select>
      </div>
      <div>
        <label for="fee-amount" class="block text-sm font-medium text-gray-700 mb-1">
          Late Fee {settings.feeType === 'flat' ? 'Amount ($)' : 'Rate (%)'}
        </label>
        <input
          id="fee-amount"
          type="number"
          min="0"
          step={settings.feeType === 'flat' ? '1' : '0.5'}
          bind:value={settings.feeAmount}
          class="input-field"
          placeholder={settings.feeType === 'flat' ? '50' : '5'}
        />
        <p class="text-xs text-gray-400 mt-1">
          {settings.feeType === 'flat'
            ? `$${settings.feeAmount} flat fee`
            : `${settings.feeAmount}% of monthly rent`}
        </p>
      </div>
    </div>
    <div class="mt-4 flex items-center gap-3">
      <button class="btn-primary" on:click={saveSettings}>💾 Save Settings</button>
      {#if saved}
        <span class="text-sm text-green-600 font-medium">✓ Saved</span>
      {/if}
    </div>
  </div>

  <!-- Display Info -->
  <div class="card">
    <h2 class="text-lg font-semibold mb-4">🌐 Display Preferences</h2>
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div>
        <p class="text-sm font-medium text-gray-700 mb-1">Currency</p>
        <p class="text-gray-900 font-medium">USD ($)</p>
        <p class="text-xs text-gray-400 mt-1">Default currency for all amounts</p>
      </div>
      <div>
        <p class="text-sm font-medium text-gray-700 mb-1">Date Format</p>
        <p class="text-gray-900 font-medium">MM/DD/YYYY</p>
        <p class="text-xs text-gray-400 mt-1">Used across the application</p>
      </div>
    </div>
  </div>

  <!-- Data Management -->
  <div class="card">
    <h2 class="text-lg font-semibold mb-4">🗄️ Data Management</h2>
    <div class="space-y-4">
      <div>
        <p class="text-sm font-medium text-gray-700 mb-1">Database Location</p>
        <div class="flex items-center gap-2">
          <code class="px-3 py-1.5 bg-gray-100 text-gray-700 rounded text-sm font-mono">./rentme.db</code>
          <span class="text-xs text-gray-400">SQLite database</span>
        </div>
      </div>
      <div>
        <p class="text-sm font-medium text-gray-700 mb-1">Backup</p>
        <p class="text-sm text-gray-600">
          To back up your data, copy the <code class="px-1 py-0.5 bg-gray-100 rounded text-xs font-mono">rentme.db</code> file from the application directory to a safe location.
        </p>
      </div>
      <div>
        <p class="text-sm font-medium text-gray-700 mb-1">Documents Folder</p>
        <div class="flex items-center gap-2">
          <code class="px-3 py-1.5 bg-gray-100 text-gray-700 rounded text-sm font-mono">./documents/</code>
          <span class="text-xs text-gray-400">Uploaded document storage</span>
        </div>
      </div>
    </div>
  </div>

  <!-- About -->
  <div class="card">
    <h2 class="text-lg font-semibold mb-4">ℹ️ About</h2>
    <div class="space-y-2">
      <div class="flex items-center gap-3">
        <span class="text-3xl">🏡</span>
        <div>
          <h3 class="text-xl font-bold text-gray-900">RentMe</h3>
          <p class="text-sm text-gray-500">Version 0.1.0</p>
        </div>
      </div>
      <p class="text-gray-600">A simple tool for managing your rental property.</p>
      <p class="text-sm text-gray-400">Built with Go + Wails + Svelte</p>
    </div>
  </div>
</div>
