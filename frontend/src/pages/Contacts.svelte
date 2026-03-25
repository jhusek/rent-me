<script lang="ts">
  import { onMount } from 'svelte';
  import { GetContacts, CreateContact, UpdateContact, DeleteContact } from '../../wailsjs/go/main/App';
  import PageHeader from '../components/PageHeader.svelte';
  import Modal from '../components/Modal.svelte';
  import Badge from '../components/Badge.svelte';
  import Toast from '../components/Toast.svelte';
  import ConfirmDialog from '../components/ConfirmDialog.svelte';
  import EmptyState from '../components/EmptyState.svelte';

  let loading = true;
  let contacts: any[] = [];

  let showModal = false;
  let modalTitle = 'Add Contact';
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

  const roles = [
    'plumber', 'electrician', 'handyman', 'contractor',
    'insurance_agent', 'property_manager', 'realtor',
    'attorney', 'accountant', 'hoa', 'other',
  ];

  function roleLabel(role: string): string {
    return role.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase());
  }

  function roleBadgeVariant(role: string): 'info' | 'success' | 'warning' | 'neutral' {
    if (['plumber', 'electrician', 'handyman', 'contractor'].includes(role)) return 'info';
    if (['insurance_agent', 'property_manager', 'realtor'].includes(role)) return 'success';
    if (['attorney', 'accountant'].includes(role)) return 'warning';
    return 'neutral';
  }

  function emptyForm(): any {
    return { name: '', role: '', phone: '', email: '', company: '', notes: '' };
  }

  let form: any = emptyForm();

  $: canSave = form.name?.trim();

  async function loadContacts() {
    try {
      contacts = (await GetContacts()) || [];
    } catch {
      contacts = [];
    } finally {
      loading = false;
    }
  }

  function openAdd() {
    editingId = null;
    modalTitle = 'Add Contact';
    form = emptyForm();
    showModal = true;
  }

  function openEdit(c: any) {
    editingId = c.id;
    modalTitle = 'Edit Contact';
    form = { ...c };
    showModal = true;
  }

  function openDelete(c: any) {
    deleteId = c.id;
    deleteName = c.name;
    showConfirm = true;
  }

  async function saveContact() {
    try {
      if (editingId) {
        await UpdateContact(form);
        toast('Contact updated successfully');
      } else {
        await CreateContact(form);
        toast('Contact added successfully');
      }
      showModal = false;
      await loadContacts();
    } catch (e) {
      toast('Failed to save contact: ' + e, 'error');
    }
  }

  async function confirmDelete() {
    if (!deleteId) return;
    try {
      await DeleteContact(deleteId);
      toast('Contact deleted');
      await loadContacts();
    } catch (e) {
      toast('Failed to delete contact: ' + e, 'error');
    }
    deleteId = null;
  }

  onMount(loadContacts);
</script>

<Toast bind:show={showToast} message={toastMessage} type={toastType} />

<PageHeader title="Contacts">
  <svelte:fragment slot="actions">
    <button class="btn-primary" on:click={openAdd}>+ Add Contact</button>
  </svelte:fragment>
</PageHeader>

{#if loading}
  <div class="flex items-center justify-center py-20">
    <div class="text-gray-400 text-lg">Loading...</div>
  </div>
{:else if contacts.length === 0}
  <div class="card">
    <EmptyState icon="📇" title="No contacts yet" message="Add vendors, agents, and service providers.">
      <button class="btn-primary" on:click={openAdd}>+ Add Contact</button>
    </EmptyState>
  </div>
{:else}
  <div class="card overflow-hidden !p-0">
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b border-gray-200 bg-gray-50">
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Name</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Role</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Phone</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Email</th>
            <th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Company</th>
            <th class="text-right px-6 py-3 text-xs font-medium text-gray-500 uppercase">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          {#each contacts as c}
            <tr class="hover:bg-gray-50">
              <td class="px-6 py-4 font-medium text-gray-900">{c.name}</td>
              <td class="px-6 py-4">
                {#if c.role}
                  <Badge variant={roleBadgeVariant(c.role)} text={roleLabel(c.role)} />
                {:else}
                  <span class="text-gray-400">—</span>
                {/if}
              </td>
              <td class="px-6 py-4 text-sm text-gray-600">{c.phone || '—'}</td>
              <td class="px-6 py-4 text-sm text-gray-600">{c.email || '—'}</td>
              <td class="px-6 py-4 text-sm text-gray-600">{c.company || '—'}</td>
              <td class="px-6 py-4 text-right">
                <button class="btn-secondary btn-sm" on:click={() => openEdit(c)}>✏️</button>
                <button class="btn-danger btn-sm ml-1" on:click={() => openDelete(c)}>🗑️</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
{/if}

<!-- Add/Edit Modal -->
<Modal bind:show={showModal} title={modalTitle}>
  <form on:submit|preventDefault={saveContact} id="contact-form">
    <div class="space-y-4">
      <div>
        <label class="label" for="name">Name *</label>
        <input id="name" class="input-field" bind:value={form.name} required />
      </div>
      <div>
        <label class="label" for="role">Role</label>
        <select id="role" class="input-field" bind:value={form.role}>
          <option value="">Select role...</option>
          {#each roles as r}
            <option value={r}>{roleLabel(r)}</option>
          {/each}
        </select>
      </div>
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="label" for="cphone">Phone</label>
          <input id="cphone" class="input-field" bind:value={form.phone} />
        </div>
        <div>
          <label class="label" for="cemail">Email</label>
          <input id="cemail" type="email" class="input-field" bind:value={form.email} />
        </div>
      </div>
      <div>
        <label class="label" for="company">Company</label>
        <input id="company" class="input-field" bind:value={form.company} />
      </div>
      <div>
        <label class="label" for="cnotes">Notes</label>
        <textarea id="cnotes" class="input-field" rows="3" bind:value={form.notes}></textarea>
      </div>
    </div>
  </form>
  <svelte:fragment slot="footer">
    <button class="btn-secondary" on:click={() => showModal = false}>Cancel</button>
    <button class="btn-primary" type="submit" form="contact-form" disabled={!canSave}>
      {editingId ? 'Save Changes' : 'Add Contact'}
    </button>
  </svelte:fragment>
</Modal>

<!-- Delete Confirmation -->
<ConfirmDialog
  bind:show={showConfirm}
  title="Delete Contact"
  message="Are you sure you want to delete {deleteName}? This action cannot be undone."
  confirmText="Delete"
  danger={true}
  on:confirm={confirmDelete}
/>
