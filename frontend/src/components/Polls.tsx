import { useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import useAuth from '../hooks/use-auth';
import axiosInstance from '../utils/axios-instance';
import '../css/Polls.css';

const Polls = () => {
    const navigate = useNavigate();
    const { auth } = useAuth();
    const [polls, setPolls] = useState([]);
    const fetchData = useCallback(async () => {
        try {
            const response = await axiosInstance.get('/polls/get');
            setPolls(response.data);
        } catch (error) {
            console.error('Error fetching data:', error);
        }
    }, []);

    useEffect(() => {
        // Fetch data from Go server
        fetchData();
    }, [fetchData]);

    return (
		<div className="polls">
  <h1>Editors' Picks</h1>
  <div className="polls-grid">
    {polls &&
      polls.map((poll) => (
        <div key={poll['id']} className="poll">
          <div className="poll-details">
            <h2>{poll['name']}</h2>
            <div className="image-description-row">
              <div className="poll-image">
                <img src={`${poll['image']}`} alt={poll['name']} />
              </div>
              <div className="poll-description">
                <p>{poll['description']}</p>
                <p>Season: {poll['season']}</p>
              </div>
            </div>
			<div className="poll-actions">
            {auth.user && (
                <button
                  className="vote-button"
                  onClick={() => navigate(`/poll/${poll['id']}`)}
                >
                  Vote
                </button>
			)}
                <button
                  className="results-button"
                  onClick={() => navigate(`/results/${poll['id']}`)}
                >
                  Check Results
                </button>
              </div>
          </div>
        </div>
      ))}
  </div>
</div>

  
    );
};

export default Polls;
