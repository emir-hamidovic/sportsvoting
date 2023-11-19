import { useState, useCallback, useEffect } from "react";
import { useParams } from "react-router-dom";
import PieChart from "./ChartPie";
import axiosInstance from "../utils/axios-instance";

interface Votes {
    value: number,
    name: string,
    pollname: string
  }

const Results = () => {
    const { pollId } = useParams();
    const [data, setData] = useState<Votes[]>([]);
    const fetchData = useCallback(async () => {
        try {
            const response = await axiosInstance.get<Votes[]>(`/playervotes/${pollId}`);
            if (response.data === null) {
                setData([])
            } else {
                setData(response.data);
            }
        } catch (error) {
            console.error('Error fetching data:', error);
        }
    }, [pollId]);

    useEffect(() => {
        fetchData();
        const intervalId = setInterval(() => {
            fetchData();
        }, 2000);

        return () => clearInterval(intervalId);
    }, [fetchData]);

    return (
    <div>
        <PieChart data={data}/>
    </div>
    )
}

export default Results
