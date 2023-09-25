import { useCallback, useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import useAuth from '../hooks/use-auth';

const Polls = () => {
  const navigate = useNavigate();
  const {auth} = useAuth();
  const [polls, setPolls] = useState([]);
  const fetchData = useCallback(async () => {
    try {
      const response = await axios.get('http://localhost:8080/getpolls');
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
        {polls.map(poll => (
          <div key={poll['id']} className="poll">
            <div className="poll-image">
              <img src={`${poll['image']}`} alt={poll['name']} />
            </div>
            <div className="poll-details">
              <h2>{poll['name']}</h2>
              <p>{poll['description']}</p>
            </div>
            <div className="poll-actions">
              {auth.user ? 
              <button className="vote-button" onClick={() => navigate(`/${poll['endpoint']}/${poll['id']}`)}>Vote</button> : 
              ''
              }
              <button className="results-button" onClick={() => navigate(`/results/${poll['id']}`)}>Check Results</button>
            </div>
          </div>
        ))}
    </div>
  );
};

export default Polls;
