import React from 'react';
import './App.css';

import Header from './components/Header';

function App() {
  const polls = [
    {
      id: 1,
      name: 'LeBron James vs. Michael Jordan',
      description: 'Who is the greatest basketball player of all time?',
      image: 'https://cdn.nba.com/headshots/nba/latest/1040x760/2544.png'
    },
    {
      id: 2,
      name: 'Kobe Bryant vs. Shaquille O\'Neal',
      description: 'Who was the better player on the Lakers\' championship teams?',
      image: 'http://t1.gstatic.com/licensed-image?q=tbn:ANd9GcTdh5LY0yPLQzIfpwCrj-UVjKeUz5c9RfXBW42rMKUW13vFtDhqEYcEVJkn9TGuF1BGLCgh_gSTMyXXIn8'
    },
    {
      id: 3,
      name: '2016 NBA Finals Game 7',
      description: 'Was this the greatest NBA Finals game of all time?',
      image: 'https://images2.minutemediacdn.com/image/fetch/w_2000,h_2000,c_fit/https%3A%2F%2Fhoopshabit.com%2Ffiles%2F2016%2F06%2Fstephen-curry-lebron-james-nba-finals-cleveland-cavaliers-golden-state-warriors.jpg'
    },
    {
      id: 4,
      name: 'Michael Jordan vs. Kobe Bryant',
      description: 'Who was the more dominant scorer?',
      image: 'https://d1si3tbndbzwz9.cloudfront.net/basketball/player/3/headshot.png'
    }
  ];

  return (
    <div className="App">
      <Header />
      <div className="App-polls">
        {polls.map(poll => (
          <div key={poll.id} className="App-poll">
            <div className="App-poll-image">
              <img src={poll.image} alt={poll.name} />
            </div>
            <div className="App-poll-details">
              <h2>{poll.name}</h2>
              <p>{poll.description}</p>
            </div>
            <div className="App-poll-actions">
              <button className="vote-button">Vote</button>
              <button className="results-button">Check Results</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default App;
