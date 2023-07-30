import React, { useCallback, useEffect, useState } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';

const PollsComponent = () => {
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
    <div className="App-polls">
        {polls.map(poll => (
          <div key={poll['id']} className="App-poll">
            <div className="App-poll-image">
              <img src={poll['image']} alt={poll['name']} />
            </div>
            <div className="App-poll-details">
              <h2>{poll['name']}</h2>
              <p>{poll['description']}</p>
            </div>
            <div className="App-poll-actions">
              <Link to={`/${poll['endpoint']}`}>
                <button className="vote-button">Vote</button>
              </Link>
              <button className="results-button">Check Results</button>
            </div>
          </div>
        ))}
    </div>
  );
};

export default PollsComponent;
