/** @type Object */ export const Bar: any;
/** @type Object */ export const Bubble: any;
/** @type Object */ export const Doughnut: any;
/** @type Object */ export const Line: any;
/** @type Object */ export const Pie: any;
/** @type Object */ export const PolarArea: any;
/** @type Object */ export const Radar: any;
/** @type Object */ export const Scatter: any;
export function generateChart(chartId: any, chartType: any, chartController: any): {
    props: {
        chartData: {
            type: ObjectConstructor;
            required: boolean;
        };
        chartOptions: {
            type: ObjectConstructor;
            default: () => void;
        };
        datasetIdKey: {
            type: StringConstructor;
            default: string;
        };
        chartId: {
            type: StringConstructor;
            default: any;
        };
        width: {
            type: NumberConstructor;
            default: number;
        };
        height: {
            type: NumberConstructor;
            default: number;
        };
        cssClasses: {
            type: StringConstructor;
            default: string;
        };
        styles: {
            type: ObjectConstructor;
            default: () => void;
        };
        plugins: {
            type: ArrayConstructor;
            default: () => any[];
        };
    };
    data(): {
        _chart: any;
        _id: string;
    };
    computed: {
        hasAnnotationPlugin(): boolean;
    };
    created(): void;
    mounted(): void;
    watch: {
        chartData: {
            handler: (newValue: any, oldValue: any) => void;
            deep: boolean;
        };
        chartOptions: {
            handler: (newValue: any) => void;
            deep: boolean;
        };
    };
    methods: {
        renderChart(data: any, options: any): void;
        chartDataHandler(newValue: any, oldValue: any): void;
        chartOptionsHandler(options: any): void;
        updateChart(): void;
        getCurrentChart(): any;
        setCurrentChart(chart: any): void;
    };
    beforeDestroy(): void;
    render: (createElement: any) => any;
};
//# sourceMappingURL=index.d.ts.map