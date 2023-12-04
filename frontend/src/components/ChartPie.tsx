import React, { useEffect, useRef } from 'react';
import * as echarts from 'echarts';

interface PieChartProps {
	data: { value: number; name: string, pollname: string}[];
}

const PieChart: React.FC<PieChartProps> = ({ data }) => {
	const chartContainer = useRef<HTMLDivElement | null>(null);
	const chartInstance = useRef<echarts.ECharts | null>(null);

	useEffect(() => {
		if (chartContainer.current) {
			if (!chartInstance.current) {
				chartInstance.current = echarts.init(chartContainer.current);
			}

			const colors = ['#FF4500', '#008000', '#4169E1', '#FFD700', '#FF1493'];

			const option: echarts.EChartOption = {
				title: {
					text: (data.length > 0 ? data[0].pollname : 'Award race'),
					left: 'center',
				},
				tooltip: {
					trigger: 'item',
					formatter: '{b} ({d}%)',
				},
				series: [
					{
						name: 'Data',
						type: 'pie',
						radius: '55%',
						center: ['50%', '60%'],
						data: data,
						animation: false, // Disable animation for live updates
						label: {
						show: true,
						formatter: '{b} ({d}%)', // Display the name and percentage in the label
						},
						itemStyle: {
								borderColor: '#fff',
								borderWidth: 2,
								shadowColor: 'rgba(0, 0, 0, 0.3)',
								shadowBlur: 10,
								color: (params: any) => colors[params.dataIndex],
						},
					},
				],
				legend: {
						orient: 'vertical',
						left: 'left',
						data: data.map((item) => item.name),
				},
			};

			chartInstance.current.setOption(option);
		}

		// Clean up ECharts instance when component unmounts
		return () => {
			if (chartInstance.current) {
				chartInstance.current.dispose();
				chartInstance.current = null;
			}
		};
	}, [data]);

	return <div ref={chartContainer} style={{ width: '600px', height: '400px' }} />;
};

export default PieChart;
