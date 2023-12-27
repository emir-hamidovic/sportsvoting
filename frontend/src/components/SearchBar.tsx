import React, { useState } from 'react';
import '../css/SearchBar.css';

const SearchBar = () => {
  const [searchQuery, setSearchQuery] = useState('');

  const handleSearchInputChange = (event: any) => {
    setSearchQuery(event.target.value);
  }

  const handleSearchSubmit = (event: any) => {
    event.preventDefault();
    console.log(`Searching for: ${searchQuery}`);
    // Add your search functionality here
  }

  return (
    <form onSubmit={handleSearchSubmit}>
      <input 
        type="text" 
        placeholder="Search..."
      />
      <button type="submit">Search</button>
    </form>
  );
}

export default SearchBar;