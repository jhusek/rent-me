<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { Chart, BarController, BarElement, CategoryScale, LinearScale, Tooltip, Legend } from 'chart.js';

  Chart.register(BarController, BarElement, CategoryScale, LinearScale, Tooltip, Legend);

  export let labels: string[] = [];
  export let datasets: { label: string; data: number[]; backgroundColor: string }[] = [];

  let canvas: HTMLCanvasElement;
  let chart: Chart;

  onMount(() => {
    chart = new Chart(canvas, {
      type: 'bar',
      data: { labels, datasets },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: { position: 'top' },
          tooltip: {
            callbacks: {
              label: (ctx) => `${ctx.dataset.label}: $${Number(ctx.parsed.y).toLocaleString()}`
            }
          }
        },
        scales: {
          y: {
            beginAtZero: true,
            ticks: {
              callback: (v) => '$' + Number(v).toLocaleString()
            }
          }
        }
      }
    });
  });

  $: if (chart) {
    chart.data.labels = labels;
    chart.data.datasets = datasets;
    chart.update();
  }

  onDestroy(() => { if (chart) chart.destroy(); });
</script>

<div class="relative h-full w-full">
  <canvas bind:this={canvas}></canvas>
</div>
