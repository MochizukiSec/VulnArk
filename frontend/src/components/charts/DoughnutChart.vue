<template>
  <div class="chart-wrapper">
    <canvas ref="doughnutChart"></canvas>
  </div>
</template>

<script>
import { onMounted, watch, ref } from 'vue';
import { Chart, DoughnutController, ArcElement, Tooltip, Legend } from 'chart.js';

// 注册Chart.js组件
Chart.register(DoughnutController, ArcElement, Tooltip, Legend);

export default {
  name: 'DoughnutChart',
  props: {
    chartData: {
      type: Object,
      required: true
    },
    options: {
      type: Object,
      default: () => ({})
    }
  },
  setup(props) {
    const doughnutChart = ref(null);
    let chartInstance = null;

    const renderChart = () => {
      if (chartInstance) {
        chartInstance.destroy();
      }

      const ctx = doughnutChart.value.getContext('2d');
      chartInstance = new Chart(ctx, {
        type: 'doughnut',
        data: props.chartData,
        options: props.options
      });
    };

    onMounted(() => {
      renderChart();
    });

    watch(() => props.chartData, () => {
      renderChart();
    }, { deep: true });

    return {
      doughnutChart
    };
  }
};
</script>

<style scoped>
.chart-wrapper {
  width: 100%;
  height: 100%;
}
</style>
