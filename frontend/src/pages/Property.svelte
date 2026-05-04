<script lang="ts">
  import { onMount } from 'svelte';
  import { GetProperty, UpdateProperty } from '../../wailsjs/go/main/App';
  import PageHeader from '../components/PageHeader.svelte';
  import Toast from '../components/Toast.svelte';
  import EmptyState from '../components/EmptyState.svelte';

  let loading = true;
  let editing = false;
  let hasProperty = false;

  let property: any = {};
  let form: any = {};

  let toastMessage = '';
  let toastType: 'success' | 'error' = 'success';
  let showToast = false;

  function toast(msg: string, type: 'success' | 'error' = 'success') {
    toastMessage = msg;
    toastType = type;
    showToast = true;
  }

  function formatCurrency(amount: number): string {
    return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(amount || 0);
  }

  function formatDate(dateStr: string): string {
    if (!dateStr) return '—';
    return new Date(dateStr).toLocaleDateString('en-US');
  }

  async function loadProperty() {
    try {
      const p = await GetProperty();
      property = p || {};
      hasProperty = !!(p && p.id);
    } catch {
      property = {};
      hasProperty = false;
    } finally {
      loading = false;
    }
  }

  function startEditing() {
    form = { ...property };
    editing = true;
  }

  function cancelEditing() {
    editing = false;
  }

  async function saveProperty() {
    if (!canSave) return;
    saving = true;
    try {
      await UpdateProperty(form);
      property = { ...form };
      hasProperty = true;
      editing = false;
      toast('Property updated successfully');
      await loadProperty();
    } catch (e) {
      toast('Failed to save property: ' + e, 'error');
    } finally {
      saving = false;
    }
  }

  let saving = false;

  $: canSave = (form.name?.trim() || form.address?.trim()) && !saving;

  onMount(loadProperty);

  const propertyTypes = ['Single Family', 'Condo', 'Townhouse', 'Duplex', 'Multi-Family', 'Other'];
</script>

<Toast bind:show={showToast} message={toastMessage} type={toastType} />

{#if loading}
  <div class="flex items-center justify-center py-20">
    <div class="text-gray-400 text-lg">Loading...</div>
  </div>
{:else if !hasProperty && !editing}
  <div class="card">
    <EmptyState icon="🏠" title="Welcome to RentMe!" message="Set up your rental property details to get started.">
      <button class="btn-primary" on:click={startEditing}>🏠 Set Up Property</button>
    </EmptyState>
  </div>
{:else if editing}
  <PageHeader title={hasProperty ? 'Edit Property' : 'Set Up Property'} />

  <form on:submit|preventDefault={saveProperty} class="space-y-6">
    <!-- Address -->
    <div class="card">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">📍 Address</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="md:col-span-2">
          <label class="label" for="name">Property Name <span class="text-red-500">*</span></label>
          <input id="name" class="input-field" bind:value={form.name} placeholder="e.g. Main Street Rental" />
        </div>
        <div class="md:col-span-2">
          <label class="label" for="address">Street Address <span class="text-red-500">*</span></label>
          <input id="address" class="input-field" bind:value={form.address} placeholder="123 Main St" />
        </div>
        <div>
          <label class="label" for="city">City</label>
          <input id="city" class="input-field" bind:value={form.city} />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="label" for="state">State</label>
            <input id="state" class="input-field" bind:value={form.state} maxlength="2" placeholder="CA" />
          </div>
          <div>
            <label class="label" for="zip">ZIP</label>
            <input id="zip" class="input-field" bind:value={form.zip} placeholder="90210" />
          </div>
        </div>
      </div>
    </div>

    <!-- Details -->
    <div class="card">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">🏗️ Property Details</h3>
      <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
        <div>
          <label class="label" for="beds">Beds</label>
          <input id="beds" type="number" class="input-field" bind:value={form.beds} min="0" />
        </div>
        <div>
          <label class="label" for="baths">Baths</label>
          <input id="baths" type="number" class="input-field" bind:value={form.baths} min="0" step="0.5" />
        </div>
        <div>
          <label class="label" for="sqft">Sq Ft</label>
          <input id="sqft" type="number" class="input-field" bind:value={form.sqft} min="0" />
        </div>
        <div>
          <label class="label" for="yearBuilt">Year Built</label>
          <input id="yearBuilt" type="number" class="input-field" bind:value={form.yearBuilt} />
        </div>
        <div>
          <label class="label" for="lotSize">Lot Size</label>
          <input id="lotSize" class="input-field" bind:value={form.lotSize} placeholder="0.25 acres" />
        </div>
        <div>
          <label class="label" for="propertyType">Property Type</label>
          <select id="propertyType" class="input-field" bind:value={form.propertyType}>
            <option value="">Select...</option>
            {#each propertyTypes as pt}
              <option value={pt}>{pt}</option>
            {/each}
          </select>
        </div>
      </div>
    </div>

    <!-- Financial -->
    <div class="card">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">💰 Financial Details</h3>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label class="label" for="purchasePrice">Purchase Price</label>
          <input id="purchasePrice" type="number" class="input-field" bind:value={form.purchasePrice} min="0" step="0.01" />
        </div>
        <div>
          <label class="label" for="purchaseDate">Purchase Date</label>
          <input id="purchaseDate" type="date" class="input-field" bind:value={form.purchaseDate} />
        </div>
        <div>
          <label class="label" for="mortgagePayment">Monthly Mortgage</label>
          <input id="mortgagePayment" type="number" class="input-field" bind:value={form.mortgagePayment} min="0" step="0.01" />
        </div>
        <div>
          <label class="label" for="mortgageRate">Rate (%)</label>
          <input id="mortgageRate" type="number" class="input-field" bind:value={form.mortgageRate} min="0" step="0.01" />
        </div>
        <div>
          <label class="label" for="mortgageTermYears">Term (years)</label>
          <input id="mortgageTermYears" type="number" class="input-field" bind:value={form.mortgageTermYears} min="0" />
        </div>
        <div>
          <label class="label" for="mortgageStart">Mortgage Start</label>
          <input id="mortgageStart" type="date" class="input-field" bind:value={form.mortgageStart} />
        </div>
        <div>
          <label class="label" for="insuranceMonthly">Insurance (monthly)</label>
          <input id="insuranceMonthly" type="number" class="input-field" bind:value={form.insuranceMonthly} min="0" step="0.01" />
        </div>
        <div>
          <label class="label" for="propertyTaxAnnual">Property Tax (annual)</label>
          <input id="propertyTaxAnnual" type="number" class="input-field" bind:value={form.propertyTaxAnnual} min="0" step="0.01" />
        </div>
        <div>
          <label class="label" for="hoaMonthly">HOA (monthly)</label>
          <input id="hoaMonthly" type="number" class="input-field" bind:value={form.hoaMonthly} min="0" step="0.01" />
        </div>
      </div>
    </div>

    <!-- Notes -->
    <div class="card">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">📝 Notes</h3>
      <textarea class="input-field" rows="4" bind:value={form.notes} placeholder="Any notes about the property..."></textarea>
    </div>

    <!-- Actions -->
    <div class="flex justify-end gap-3">
      {#if hasProperty}
        <button type="button" class="btn-secondary" on:click={cancelEditing}>Cancel</button>
      {/if}
      <button type="submit" class="btn-primary" disabled={!canSave}>💾 Save Property</button>
    </div>
  </form>
{:else}
  <PageHeader title="Property">
    <svelte:fragment slot="actions">
      <button class="btn-primary" on:click={startEditing}>✏️ Edit Property</button>
    </svelte:fragment>
  </PageHeader>

  <div class="space-y-6">
    <!-- Address Card -->
    <div class="card">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">📍 Address</h3>
      <div class="space-y-2">
        {#if property.name}
          <p class="text-xl font-medium text-gray-900">{property.name}</p>
        {/if}
        <p class="text-gray-600">
          {property.address || '—'}
        </p>
        <p class="text-gray-600">
          {property.city || ''}{property.city && property.state ? ', ' : ''}{property.state || ''} {property.zip || ''}
        </p>
      </div>
    </div>

    <!-- Details Card -->
    <div class="card">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">🏗️ Property Details</h3>
      <div class="grid grid-cols-2 md:grid-cols-3 gap-6">
        <div>
          <p class="text-sm text-gray-500">Beds</p>
          <p class="text-lg font-medium">{property.beds || '—'}</p>
        </div>
        <div>
          <p class="text-sm text-gray-500">Baths</p>
          <p class="text-lg font-medium">{property.baths || '—'}</p>
        </div>
        <div>
          <p class="text-sm text-gray-500">Sq Ft</p>
          <p class="text-lg font-medium">{property.sqft ? property.sqft.toLocaleString() : '—'}</p>
        </div>
        <div>
          <p class="text-sm text-gray-500">Year Built</p>
          <p class="text-lg font-medium">{property.yearBuilt || '—'}</p>
        </div>
        <div>
          <p class="text-sm text-gray-500">Lot Size</p>
          <p class="text-lg font-medium">{property.lotSize || '—'}</p>
        </div>
        <div>
          <p class="text-sm text-gray-500">Type</p>
          <p class="text-lg font-medium">{property.propertyType || '—'}</p>
        </div>
      </div>
    </div>

    <!-- Financial Card -->
    <div class="card">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">💰 Financial Details</h3>
      <div class="grid grid-cols-2 md:grid-cols-3 gap-6">
        <div>
          <p class="text-sm text-gray-500">Purchase Price</p>
          <p class="text-lg font-medium">{formatCurrency(property.purchasePrice)}</p>
        </div>
        <div>
          <p class="text-sm text-gray-500">Purchase Date</p>
          <p class="text-lg font-medium">{formatDate(property.purchaseDate)}</p>
        </div>
        <div>
          <p class="text-sm text-gray-500">Monthly Mortgage</p>
          <p class="text-lg font-medium">{formatCurrency(property.mortgagePayment)}</p>
        </div>
        <div>
          <p class="text-sm text-gray-500">Rate</p>
          <p class="text-lg font-medium">{property.mortgageRate ? property.mortgageRate + '%' : '—'}</p>
        </div>
        <div>
          <p class="text-sm text-gray-500">Term</p>
          <p class="text-lg font-medium">{property.mortgageTermYears ? property.mortgageTermYears + ' years' : '—'}</p>
        </div>
        <div>
          <p class="text-sm text-gray-500">Mortgage Start</p>
          <p class="text-lg font-medium">{formatDate(property.mortgageStart)}</p>
        </div>
        <div>
          <p class="text-sm text-gray-500">Insurance (monthly)</p>
          <p class="text-lg font-medium">{formatCurrency(property.insuranceMonthly)}</p>
        </div>
        <div>
          <p class="text-sm text-gray-500">Property Tax (annual)</p>
          <p class="text-lg font-medium">{formatCurrency(property.propertyTaxAnnual)}</p>
        </div>
        <div>
          <p class="text-sm text-gray-500">HOA (monthly)</p>
          <p class="text-lg font-medium">{formatCurrency(property.hoaMonthly)}</p>
        </div>
      </div>
    </div>

    <!-- Notes -->
    {#if property.notes}
      <div class="card">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">📝 Notes</h3>
        <p class="text-gray-600 whitespace-pre-wrap">{property.notes}</p>
      </div>
    {/if}
  </div>
{/if}
