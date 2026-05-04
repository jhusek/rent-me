<script lang="ts">
  import { onMount } from 'svelte';
  import { GetTenants, CreateTenant, UpdateTenant, DeleteTenant } from '../../wailsjs/go/main/App';
  import PageHeader from '../components/PageHeader.svelte';
  import Modal from '../components/Modal.svelte';
  import Badge from '../components/Badge.svelte';
  import Toast from '../components/Toast.svelte';
  import ConfirmDialog from '../components/ConfirmDialog.svelte';
  import EmptyState from '../components/EmptyState.svelte';

  let loading = true;
  let tenants: any[] = [];
  let filterStatus = '';

  let showModal = false;
  let modalTitle = 'Add Tenant';
  let editingId: number | null = null;

  let showConfirm = false;
  let deleteId: number | null = null;
  let deleteName = '';

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

  function emptyForm(): any {
    return {
      firstName: '', lastName: '', email: '', phone: '',
      emergencyContactName: '', emergencyContactPhone: '',
      leaseStart: '', leaseEnd: '', monthlyRent: 0, securityDeposit: 0,
      status: 'active', moveInDate: '', moveOutDate: '', notes: '',
    };
  }

  let form: any = emptyForm();

  $: canSave = form.firstName?.trim() && form.lastName?.trim();

  async function loadTenants() {
    try {
      tenants = (await GetTenants(filterStatus)) || [];
    } catch {
      tenants = [];
    } finally {
      loading = false;
    }
  }

  function openAdd() {
    editingId = null;
    modalTitle = 'Add Tenant';
    form = emptyForm();
    showModal = true;
  }

  function openEdit(t: any) {
    editingId = t.id;
    modalTitle = 'Edit Tenant';
    form = { ...t };
    showModal = true;
  }

  function openDelete(t: any) {
    deleteId = t.id;
    deleteName = `${t.firstName} ${t.lastName}`;
    showConfirm = true;
  }

  async function saveTenant() {
    try {
      if (editingId) {
        await UpdateTenant(form);
        toast('Tenant updated successfully');
      } else {
        await CreateTenant(form);
        toast('Tenant added successfully');
      }
      showModal = false;
      await loadTenants();
    } catch (e) {
      toast('Failed to save tenant: ' + e, 'error');
    }
  }

  async function confirmDelete() {
    if (!deleteId) return;
    try {
      await DeleteTenant(deleteId);
      toast('Tenant deleted');
      await loadTenants();
    } catch (e) {
      toast('Failed to delete tenant: ' + e, 'error');
    }
    deleteId = null;
  }

  function setFilter(status: string) {
    filterStatus = status;
    loading = true;
    loadTenants();
  }

  function statusVariant(status: string): 'success' | 'warning' | 'neutral' {
    if (status === 'active') return 'success';
    if (status === 'pending') return 'warning';
    return 'neutral';
  }

  $: filteredTenants = tenants;

  onMount(loadTenants);
</script>

<Toast bind:show={showToast} message={toastMessage} type={toastType} />

<PageHeader title="Tenants">
  <svelte:fragment slot="actions">
    <button class="btn-primary" on:click={openAdd}>+ Add Tenant</button>
  </svelte:fragment>
</PageHeader>

<!-- Filter Tabs -->
<div class="flex gap-1 mb-6 bg-gray-100 rounded-lg p-1 w-fit">
  {#each [{ label: 'All', value: '' }, { label: 'Active', value: 'active' }, { label: 'Past', value: 'past' }] as tab}
    <button
      class="px-4 py-1.5 rounded-md text-sm font-medium transition-colors {filterStatus === tab.value ? 'bg-white shadow text-gray-900' : 'text-gray-600 hover:text-gray-900'}"
      on:click={() => setFilter(tab.value)}
    >{tab.label}</button>
  {/each}
</div>

{#if loading}
  <div class="flex items-center justify-center py-20">
    <div class="text-gray-400 text-lg">Loading...</div>
  </div>
{:else if filteredTenants.length === 0}
  <div class="card">
    <EmptyState icon="👤" title="No tenants yet" message="Add your first tenant to start tracking leases and payments.">
      <button class="btn-primary" on:click={openAdd}>+ Add Tenant</button>
    </EmptyState>
  </div>
{:else}
  <div class="card overflow-hidden !p-0">
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b border-gray-200 bg-gray-50">
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Name</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Status</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Lease Period</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Monthly Rent</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Phone</th>
            <th class="text-right px-6 py-3 text-xs font-medium text-gray-500 uppercase">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          {#each filteredTenants as t}
            <tr class="hover:bg-gray-50">
              <td class="px-6 py-4 font-medium text-gray-900">{t.firstName} {t.lastName}</td>
              <td class="px-6 py-4">
                <Badge variant={statusVariant(t.status)} text={t.status} />
              </td>
              <td class="px-6 py-4 text-sm text-gray-600">
                {formatDate(t.leaseStart)} – {formatDate(t.leaseEnd)}
              </td>
              <td class="px-6 py-4 text-sm text-gray-900">{formatCurrency(t.monthlyRent)}</td>
              <td class="px-6 py-4 text-sm text-gray-600">{t.phone || '—'}</td>
              <td class="px-6 py-4 text-right">
                <button class="btn-secondary btn-sm" on:click={() => openEdit(t)}>✏️</button>
                <button class="btn-danger btn-sm ml-1" on:click={() => openDelete(t)}>🗑️</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
{/if}

<!-- Add/Edit Modal -->
<Modal bind:show={showModal} title={modalTitle} size="lg">
  <form on:submit|preventDefault={saveTenant} id="tenant-form">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <label class="label" for="firstName">First Name *</label>
        <input id="firstName" class="input-field" bind:value={form.firstName} required />
      </div>
      <div>
        <label class="label" for="lastName">Last Name *</label>
        <input id="lastName" class="input-field" bind:value={form.lastName} required />
      </div>
      <div>
        <label class="label" for="email">Email</label>
        <input id="email" type="email" class="input-field" bind:value={form.email} />
      </div>
      <div>
        <label class="label" for="phone">Phone</label>
        <input id="phone" class="input-field" bind:value={form.phone} />
      </div>
      <div>
        <label class="label" for="emergencyName">Emergency Contact</label>
        <input id="emergencyName" class="input-field" bind:value={form.emergencyContactName} placeholder="Name" />
      </div>
      <div>
        <label class="label" for="emergencyPhone">Emergency Phone</label>
        <input id="emergencyPhone" class="input-field" bind:value={form.emergencyContactPhone} />
      </div>

      <hr class="col-span-full border-gray-200" />

      <div>
        <label class="label" for="leaseStart">Lease Start</label>
        <input id="leaseStart" type="date" class="input-field" bind:value={form.leaseStart} />
      </div>
      <div>
        <label class="label" for="leaseEnd">Lease End</label>
        <input id="leaseEnd" type="date" class="input-field" bind:value={form.leaseEnd} />
      </div>
      <div>
        <label class="label" for="monthlyRent">Monthly Rent</label>
        <input id="monthlyRent" type="number" class="input-field" bind:value={form.monthlyRent} min="0" step="0.01" />
      </div>
      <div>
        <label class="label" for="securityDeposit">Security Deposit</label>
        <input id="securityDeposit" type="number" class="input-field" bind:value={form.securityDeposit} min="0" step="0.01" />
      </div>
      <div>
        <label class="label" for="status">Status</label>
        <select id="status" class="input-field" bind:value={form.status}>
          <option value="active">Active</option>
          <option value="past">Past</option>
          <option value="pending">Pending</option>
        </select>
      </div>
      <div><!-- spacer --></div>
      <div>
        <label class="label" for="moveInDate">Move-in Date</label>
        <input id="moveInDate" type="date" class="input-field" bind:value={form.moveInDate} />
      </div>
      <div>
        <label class="label" for="moveOutDate">Move-out Date</label>
        <input id="moveOutDate" type="date" class="input-field" bind:value={form.moveOutDate} />
      </div>
      <div class="col-span-full">
        <label class="label" for="notes">Notes</label>
        <textarea id="notes" class="input-field" rows="3" bind:value={form.notes}></textarea>
      </div>
    </div>
  </form>
  <svelte:fragment slot="footer">
    <button class="btn-secondary" on:click={() => showModal = false}>Cancel</button>
    <button class="btn-primary" type="submit" form="tenant-form" disabled={!canSave}>
      {editingId ? 'Save Changes' : 'Add Tenant'}
    </button>
  </svelte:fragment>
</Modal>

<!-- Delete Confirmation -->
<ConfirmDialog
  bind:show={showConfirm}
  title="Delete Tenant"
  message="Are you sure you want to delete {deleteName}? This action cannot be undone."
  confirmText="Delete"
  danger={true}
  on:confirm={confirmDelete}
/>
