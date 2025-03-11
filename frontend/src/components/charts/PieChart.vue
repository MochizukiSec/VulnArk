<template>
  <div class="chart-wrapper">
    <canvas ref="pieChart"></canvas>
  </div>
</template>

<script>
import { onMounted, watch, ref } from 'vue';
import { Chart, PieController, ArcElement, Tooltip, Legend } from 'chart.js';

// 注册Chart.js组件
Chart.register(PieController, ArcElement, Tooltip, Legend);

export default {
  name: 'PieChart',
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
    const pieChart = ref(null);
    let chartInstance = null;

    const renderChart = () => {
      if (chartInstance) {
        chartInstance.destroy();
      }

      const ctx = pieChart.value.getContext('2d');
      chartInstance = new Chart(ctx, {
        type: 'pie',
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
      pieChart
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