// SearchBar.js

import React, { useState } from 'react';
import '../css/SearchBar.css';

const SearchBar = () => {
  const [searchQuery, setSearchQuery] = useState('');

  const handleSearchInputChange = (event: any) => {
    setSearchQuery(event.target.value);
  }

  const handleKeyPress = (event: any) => {
    if (event.key === 'Enter') {
      event.preventDefault();
      handleSearchSubmit();
    }
  }

  const handleSearchSubmit = () => {
    console.log(`Searching for: ${searchQuery}`);
    // Add your search functionality here

    // Optionally, you can reset the searchQuery state after the search
    setSearchQuery('');
  }

  return (
    <form>
      <input 
        type="text" 
        placeholder="🔍 Search..."
        value={searchQuery}
        onChange={handleSearchInputChange}
        onKeyPress={handleKeyPress}
      />
    </form>
  );
}

export default SearchBar;
