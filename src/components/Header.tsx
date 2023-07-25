import React from 'react';
import './header.css';

function Header() {
  return (
    <header className="header">
      <h1 className="underline text-3xl">Sport Voting</h1>
      <div className="header-buttons">
        <a href="#">Login</a>
        <a href="#">Sign up</a>
      </div>
    </header>
  );
}

export default Header;
