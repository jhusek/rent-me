export function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(amount || 0);
}

export function formatDate(dateStr: string): string {
  if (!dateStr) return '—';
  try {
    return new Date(dateStr + 'T00:00:00').toLocaleDateString('en-US');
  } catch {
    return dateStr;
  }
}

export function getCurrentMonth(): number {
  return new Date().getMonth() + 1;
}

export function getCurrentYear(): number {
  return new Date().getFullYear();
}

export const EXPENSE_CATEGORIES = [
  { value: 'mortgage', label: 'Mortgage', color: 'info' as const },
  { value: 'insurance', label: 'Insurance', color: 'info' as const },
  { value: 'property_tax', label: 'Property Tax', color: 'warning' as const },
  { value: 'repairs', label: 'Repairs', color: 'danger' as const },
  { value: 'maintenance', label: 'Maintenance', color: 'warning' as const },
  { value: 'utilities', label: 'Utilities', color: 'info' as const },
  { value: 'hoa', label: 'HOA', color: 'neutral' as const },
  { value: 'management', label: 'Management', color: 'neutral' as const },
  { value: 'supplies', label: 'Supplies', color: 'neutral' as const },
  { value: 'legal', label: 'Legal', color: 'danger' as const },
  { value: 'advertising', label: 'Advertising', color: 'info' as const },
  { value: 'travel', label: 'Travel', color: 'neutral' as const },
  { value: 'depreciation', label: 'Depreciation', color: 'warning' as const },
  { value: 'other', label: 'Other', color: 'neutral' as const },
];

export const PAYMENT_METHODS = [
  { value: 'cash', label: 'Cash' },
  { value: 'check', label: 'Check' },
  { value: 'bank_transfer', label: 'Bank Transfer' },
  { value: 'venmo', label: 'Venmo' },
  { value: 'zelle', label: 'Zelle' },
  { value: 'paypal', label: 'PayPal' },
  { value: 'other', label: 'Other' },
];

export const MONTHS = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December',
];

export function getCategoryLabel(value: string): string {
  const cat = EXPENSE_CATEGORIES.find(c => c.value === value);
  return cat ? cat.label : value;
}

export function getCategoryColor(value: string): 'success' | 'warning' | 'danger' | 'info' | 'neutral' {
  const cat = EXPENSE_CATEGORIES.find(c => c.value === value);
  return cat ? cat.color : 'neutral';
}

export function getMethodLabel(value: string): string {
  const m = PAYMENT_METHODS.find(p => p.value === value);
  return m ? m.label : value;
}

export function timeAgo(dateStr: string): string {
  if (!dateStr) return '';
  const date = new Date(dateStr);
  const now = new Date();
  const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);
  if (seconds < 0) return date.toLocaleDateString('en-US');
  if (seconds < 60) return 'Just now';
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
  if (seconds < 172800) return 'Yesterday';
  if (seconds < 604800) return `${Math.floor(seconds / 86400)}d ago`;
  return date.toLocaleDateString('en-US');
}
