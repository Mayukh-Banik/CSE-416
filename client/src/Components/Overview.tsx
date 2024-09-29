import React, { useEffect, useRef } from 'react';
import Sidebar from './Sidebar'; 
import { Chart, BarController, BarElement, CategoryScale, LinearScale, Title, Tooltip, Legend } from 'chart.js';
import { transactions } from './Transactions'; 

export interface Transaction {
    dateTime: string;
    transactionId: string;
    sender: string;
    receiver: string;
    fileName: string;
    fileSize: string;
    status: "Complete" | "Pending" | "Failed";
    fee: string;
}

const aggregateDataByDate = (transactions: Transaction[]) => {
    const aggregatedData: Record<string, number> = {};

    transactions.forEach(transaction => {
        const date = new Date(transaction.dateTime).toLocaleDateString(); 
        const fee = parseFloat(transaction.fee); 

        if (!aggregatedData[date]) {
            aggregatedData[date] = 0; 
        }
        aggregatedData[date] += fee; 
    });

    return aggregatedData;
};

const prepareChartData = (aggregatedData: Record<string, number>) => {
    const labels = Object.keys(aggregatedData).sort(); 
    const data = labels.map(label => aggregatedData[label]); 

    return { labels, data };
};

const OverviewPage: React.FC = () => {
    const chartRef = useRef<HTMLCanvasElement | null>(null);

    Chart.register(BarController, BarElement, CategoryScale, LinearScale, Title, Tooltip, Legend);

    useEffect(() => {
        if (chartRef.current) {
            const ctx = chartRef.current.getContext('2d');
            if (ctx) {
                const aggregatedData = aggregateDataByDate(transactions);
                const { labels, data } = prepareChartData(aggregatedData);

                const chart = new Chart(ctx, {
                    type: 'bar',
                    data: {
                        labels: labels,
                        datasets: [{
                            data: data, 
                            backgroundColor: 'rgba(75, 192, 192, 0.6)',
                        }],
                    },
                    options: {
                        scales: {
                            x: {
                                title: {
                                    display: true,
                                    text: 'Time', 
                                },
                            },
                            y: {
                                title: {
                                    display: true,
                                    text: 'Fees($)',
                                },
                            },
                        }
                    }
                });
                return () => {
                    chart.destroy();
                };
                
            }
        }
    }, [transactions]); 

    return (
        <div style={{ display: 'flex', height: '100vh', justifyContent: 'center', alignItems: 'center' }}>
            <Sidebar />
            <div className="chart-container" style={{ width: '100%', height: '40%' }}> 
                <canvas
                    ref={chartRef}
                    style={{ width: '100%', height: '100%' }} 
                />
            </div>
        </div>
    );
    
};

export default OverviewPage;
