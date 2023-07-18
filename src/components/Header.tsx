import React from 'react';
import './header.css';

function Header() {
  return (
    <header className="header">
      <div className="logo">Sport Voting</div>
      <div className="header-buttons">
        <a href="#">Login</a>
        <a href="#">Sign up</a>
      </div>
    </header>
  );
}

export default Header;
