<template>
  <div class="chart-wrapper">
    <canvas ref="lineChart"></canvas>
  </div>
</template>

<script>
import { onMounted, watch, ref } from 'vue';
import { 
  Chart, 
  LineController, 
  LineElement, 
  PointElement, 
  LinearScale, 
  CategoryScale, 
  Tooltip, 
  Legend 
} from 'chart.js';

// 注册Chart.js组件
Chart.register(
  LineController, 
  LineElement, 
  PointElement, 
  LinearScale, 
  CategoryScale, 
  Tooltip, 
  Legend
);

export default {
  name: 'LineChart',
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
    const lineChart = ref(null);
    let chartInstance = null;

    const renderChart = () => {
      if (chartInstance) {
        chartInstance.destroy();
      }

      const ctx = lineChart.value.getContext('2d');
      chartInstance = new Chart(ctx, {
        type: 'line',
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
      lineChart
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