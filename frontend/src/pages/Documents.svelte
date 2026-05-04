<script lang="ts">
  import { onMount } from 'svelte';
  import { GetDocuments, CreateDocument, DeleteDocument, BrowseForFile, GetDocumentContent } from '../../wailsjs/go/main/App';
  import PageHeader from '../components/PageHeader.svelte';
  import Modal from '../components/Modal.svelte';
  import Badge from '../components/Badge.svelte';
  import Toast from '../components/Toast.svelte';
  import ConfirmDialog from '../components/ConfirmDialog.svelte';
  import EmptyState from '../components/EmptyState.svelte';
  import { formatDate } from '../lib/utils';

  let loading = true;
  let documents: any[] = [];

  // Filter
  let filterCategory = '';

  // Modal state
  let showModal = false;

  // Delete state
  let showConfirm = false;
  let deleteId: number | null = null;
  let deleteLabel = '';

  // Viewer state
  let showViewer = false;
  let viewerDoc: any = null;
  let viewerContent: { base64Data: string; mimeType: string; name: string } | null = null;
  let viewerLoading = false;
  $: viewerDataUrl = viewerContent ? `data:${viewerContent.mimeType};base64,${viewerContent.base64Data}` : '';
  $: viewerIsImage = viewerContent?.mimeType?.startsWith('image/') ?? false;
  $: viewerIsPdf = viewerContent?.mimeType === 'application/pdf' ?? false;
  $: viewerIsText = viewerContent?.mimeType?.startsWith('text/') ?? false;
  $: viewerText = (viewerIsText && viewerContent)
    ? decodeURIComponent(escape(atob(viewerContent.base64Data)))
    : '';

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
      name: '',
      category: 'other',
      filePath: '',
      fileSize: 0,
      mimeType: '',
      uploadDate: todayStr(),
      notes: '',
    };
  }

  let form: any = emptyForm();

  $: canSave = form.name?.trim() && form.filePath;
  const categories = [
    { value: '', label: 'All Categories' },
    { value: 'lease', label: 'Lease' },
    { value: 'inspection', label: 'Inspection' },
    { value: 'receipt', label: 'Receipt' },
    { value: 'insurance', label: 'Insurance' },
    { value: 'tax', label: 'Tax' },
    { value: 'photo', label: 'Photo' },
    { value: 'other', label: 'Other' },
  ];

  const categoryIcons: Record<string, string> = {
    lease: '📄',
    inspection: '🔍',
    receipt: '🧾',
    insurance: '🛡️',
    tax: '💰',
    photo: '📷',
    other: '📎',
  };

  function categoryVariant(c: string): 'success' | 'warning' | 'danger' | 'info' | 'neutral' {
    switch (c) {
      case 'lease': return 'info';
      case 'inspection': return 'warning';
      case 'receipt': return 'success';
      case 'insurance': return 'info';
      case 'tax': return 'danger';
      case 'photo': return 'neutral';
      default: return 'neutral';
    }
  }

  function categoryLabel(c: string): string {
    const cat = categories.find(x => x.value === c);
    return cat ? cat.label : c;
  }

  function formatFileSize(bytes: number): string {
    if (!bytes || bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
  }

  async function loadDocuments() {
    loading = true;
    try {
      documents = (await GetDocuments(filterCategory)) || [];
    } catch {
      documents = [];
    }
    loading = false;
  }

  function openUpload() {
    form = emptyForm();
    showModal = true;
  }

  async function browseFile() {
    try {
      const path = await BrowseForFile();
      if (!path) return; // user cancelled
      form.filePath = path;
      // Auto-fill name from filename if not already set
      const fileName = path.split(/[\\/]/).pop() ?? '';
      if (!form.name) {
        form.name = fileName.replace(/\.[^.]+$/, ''); // strip extension
      }
    } catch (e) {
      toast('Could not open file picker: ' + e, 'error');
    }
  }

  async function openDocumentView(d: any) {
    viewerDoc = d;
    viewerContent = null;
    viewerLoading = true;
    showViewer = true;
    try {
      viewerContent = await GetDocumentContent(d.id);
    } catch (e) {
      showViewer = false;
      toast('Could not load document: ' + e, 'error');
    }
    viewerLoading = false;
  }

  function openDelete(d: any) {
    deleteId = d.id;
    deleteLabel = d.name;
    showConfirm = true;
  }

  async function saveDocument() {
    if (!form.filePath) {
      toast('Please select a file before uploading.', 'error');
      return;
    }
    try {
      const payload = { ...form, fileSize: Number(form.fileSize) };
      await CreateDocument(payload);
      toast('Document uploaded successfully');
      showModal = false;
      await loadDocuments();
    } catch (e) {
      toast('Failed to upload document: ' + e, 'error');
    }
  }

  async function confirmDelete() {
    if (!deleteId) return;
    try {
      await DeleteDocument(deleteId);
      toast('Document deleted');
      await loadDocuments();
    } catch (e) {
      toast('Failed to delete document: ' + e, 'error');
    }
    deleteId = null;
  }

  onMount(loadDocuments);
</script>

<Toast bind:show={showToast} message={toastMessage} type={toastType} />

<PageHeader title="Documents">
  <svelte:fragment slot="actions">
    <button class="btn-primary" on:click={openUpload}>📁 Upload Document</button>
  </svelte:fragment>
</PageHeader>

<!-- Category Filter -->
<div class="flex flex-wrap items-center gap-3 mb-6">
  <div class="flex gap-1 bg-gray-100 rounded-lg p-1">
    {#each categories as cat}
      <button
        class="px-3 py-1.5 rounded-md text-sm font-medium transition-colors {filterCategory === cat.value ? 'bg-white shadow text-gray-900' : 'text-gray-600 hover:text-gray-900'}"
        on:click={() => { filterCategory = cat.value; loadDocuments(); }}
      >{cat.label}</button>
    {/each}
  </div>
</div>

{#if loading}
  <div class="flex items-center justify-center py-20">
    <div class="text-gray-400 text-lg">Loading...</div>
  </div>
{:else if documents.length === 0}
  <div class="card">
    <EmptyState icon="📁" title="No documents" message="Upload your first document to get started.">
      <button class="btn-primary" on:click={openUpload}>📁 Upload Document</button>
    </EmptyState>
  </div>
{:else}
  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
    {#each documents as d}
      <div class="card hover:shadow-md transition-shadow">
        <div class="flex items-start gap-3">
          <span class="text-3xl">{categoryIcons[d.category] || '📎'}</span>
          <div class="flex-1 min-w-0">
            <h3 class="font-medium text-gray-900 truncate" title={d.name}>{d.name}</h3>
            <div class="flex items-center gap-2 mt-1">
              <Badge variant={categoryVariant(d.category)} text={categoryLabel(d.category)} />
            </div>
            <div class="mt-2 text-xs text-gray-500 space-y-0.5">
              <p>Uploaded: {formatDate(d.uploadDate)}</p>
              {#if d.fileSize}
                <p>Size: {formatFileSize(d.fileSize)}</p>
              {/if}
              {#if d.filePath}
                <p class="truncate" title={d.filePath}>Path: {d.filePath}</p>
              {/if}
            </div>
            {#if d.notes}
              <p class="mt-2 text-sm text-gray-600 line-clamp-2">{d.notes}</p>
            {/if}
          </div>
        </div>
        <div class="mt-3 pt-3 border-t border-gray-100 flex justify-end gap-2">
          <button class="btn-secondary btn-sm" on:click={() => openDocumentView(d)}>👁️ View</button>
          <button class="btn-danger btn-sm" on:click={() => openDelete(d)}>🗑️ Delete</button>
        </div>
      </div>
    {/each}
  </div>
{/if}

<!-- Upload Document Modal -->
<Modal bind:show={showModal} title="Upload Document" size="lg">
  <form on:submit|preventDefault={saveDocument} id="document-form">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="col-span-full">
        <label class="label" for="name">Name *</label>
        <input id="name" class="input-field" bind:value={form.name} required />
      </div>
      <div>
        <label class="label" for="category">Category</label>
        <select id="category" class="input-field" bind:value={form.category}>
          <option value="lease">Lease</option>
          <option value="inspection">Inspection</option>
          <option value="receipt">Receipt</option>
          <option value="insurance">Insurance</option>
          <option value="tax">Tax</option>
          <option value="photo">Photo</option>
          <option value="other">Other</option>
        </select>
      </div>
      <div>
        <label class="label" for="uploadDate">Upload Date</label>
        <input id="uploadDate" type="date" class="input-field" bind:value={form.uploadDate} />
      </div>
      <div class="col-span-full">
        <label class="label">File *</label>
        <div class="flex items-center gap-2">
          <button type="button" class="btn-secondary" on:click={browseFile}>📂 Browse…</button>
          {#if form.filePath}
            <span class="text-sm text-gray-700 truncate max-w-xs" title={form.filePath}>
              {form.filePath.split(/[\\/]/).pop()}
            </span>
          {:else}
            <span class="text-sm text-gray-400">No file selected</span>
          {/if}
        </div>
      </div>
      <div class="col-span-full">
        <label class="label" for="notes">Notes</label>
        <textarea id="notes" class="input-field" rows="3" bind:value={form.notes}></textarea>
      </div>
    </div>
  </form>
  <svelte:fragment slot="footer">
    <button class="btn-secondary" on:click={() => showModal = false}>Cancel</button>
    <button class="btn-primary" type="submit" form="document-form" disabled={!canSave}>Upload</button>
  </svelte:fragment>
</Modal>

<!-- Document Viewer -->
<Modal bind:show={showViewer} title={viewerDoc?.name ?? 'Document'} size="xl">
  {#if viewerLoading}
    <div class="flex items-center justify-center h-96 text-gray-400 text-lg">Loading…</div>
  {:else if viewerContent}
    {#if viewerIsImage}
      <div class="flex items-center justify-center overflow-auto max-h-[70vh]">
        <img src={viewerDataUrl} alt={viewerContent.name} class="max-w-full max-h-[70vh] object-contain rounded" />
      </div>
    {:else if viewerIsPdf}
      <iframe
        src={viewerDataUrl}
        title={viewerContent.name}
        class="w-full rounded border border-gray-200"
        style="height: 70vh;"
      ></iframe>
    {:else if viewerIsText}
      <pre class="bg-gray-50 rounded p-4 text-sm text-gray-800 overflow-auto max-h-[70vh] whitespace-pre-wrap break-words">{viewerText}</pre>
    {:else}
      <div class="flex flex-col items-center justify-center gap-4 h-64 text-gray-500">
        <span class="text-5xl">📎</span>
        <p class="text-base font-medium">{viewerContent.name}</p>
        <p class="text-sm text-gray-400">Preview is not available for this file type ({viewerContent.mimeType}).</p>
      </div>
    {/if}
    <p class="mt-3 text-xs text-gray-400 text-right">{viewerContent.mimeType}</p>
  {/if}
  <svelte:fragment slot="footer">
    <button class="btn-secondary" on:click={() => showViewer = false}>Close</button>
  </svelte:fragment>
</Modal>

<!-- Delete Confirmation -->
<ConfirmDialog
  bind:show={showConfirm}
  title="Delete Document"
  message="Are you sure you want to delete &quot;{deleteLabel}&quot;? This action cannot be undone."
  confirmText="Delete"
  danger={true}
  on:confirm={confirmDelete}
/>
