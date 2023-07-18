import React from 'react';
import './App.css';

import Header from './components/Header';
import PollsComponent from './components/Polls';

function App() {
  return (
    <div className="App">
      <Header />
      <PollsComponent />
    </div>
  );
}

export default App;
