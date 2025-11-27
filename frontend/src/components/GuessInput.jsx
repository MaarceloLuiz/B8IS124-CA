import { useState, useEffect, useRef, useCallback } from 'react';

const GuessInput = ({ territories, onSubmit, disabled }) => {
  const [input, setInput] = useState('');
  const [suggestions, setSuggestions] = useState([]);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const inputRef = useRef(null);
  const suggestionsRef = useRef(null);
  const selectedItemRef = useRef(null);

  // FIX: Wrapped in useCallback to stabilize the function reference
  const getSuggestions = useCallback((value) => {
    const query = value.toLowerCase();
    // If territories is not an array or is empty, return empty array
    if (!Array.isArray(territories) || territories.length === 0) {
      console.log("Territories not available:", territories);
      return [];
    }

    // Check if territories are strings or objects
    const isStringArray = typeof territories[0] === 'string';
    
    // If input is empty, return ALL territories
    if (!query) {
      return territories;
    }
    
    // Filter based on whether we have strings or objects
    if (isStringArray) {
      return territories
        .filter(t => t && t.toLowerCase().includes(query));
    } else {
      return territories
        .filter(t => t && t.name)
        .filter(t =>
          t.name.toLowerCase().includes(query) ||
          (t.code && t.code.toLowerCase() === query)
        );
    }
  }, [territories]); // Dependency ensures it only updates when territories change

  const handleSubmit = () => {
    if (!input.trim() || disabled) return;
    
    // Try to find a matching territory
    const matchingTerritory = suggestions.find(territory => {
      const territoryName = typeof territory === 'string' ? territory : territory.name;
      return territoryName.toLowerCase() === input.toLowerCase();
    });
    
    if (matchingTerritory) {
      // If we found a match, submit it
      const territoryValue = typeof matchingTerritory === 'string' ? matchingTerritory : matchingTerritory.name;
      onSubmit(territoryValue);
      setInput('');
      setSuggestions([]);
      setShowSuggestions(false);
      setSelectedIndex(-1);
    } else if (suggestions.length > 0) {
      // If no exact match but we have suggestions, use the first one
      const territoryValue = typeof suggestions[0] === 'string' ? suggestions[0] : suggestions[0].name;
      onSubmit(territoryValue);
      setInput('');
      setSuggestions([]);
      setShowSuggestions(false);
      setSelectedIndex(-1);
    } else {
      // If no suggestions, submit the raw input
      onSubmit(input);
      setInput('');
      setSuggestions([]);
      setShowSuggestions(false);
      setSelectedIndex(-1);
    }
  };

  // Function to handle selecting an item from the list
  const handleSelectItem = (territory) => {
    // Get the display text for the territory
    const displayText = typeof territory === 'string' ? territory : territory.name;
    
    // Set the input value to the selected territory
    setInput(displayText);
    
    // Hide the suggestions
    setShowSuggestions(false);
    setSelectedIndex(-1);
    
    // Focus the input field
    inputRef.current.focus();
  };

  useEffect(() => {
    const timer = setTimeout(() => {
      try {
        const newSuggestions = getSuggestions(input);
        setSuggestions(newSuggestions);
        setShowSuggestions(newSuggestions.length > 0);
        // Reset selected index when suggestions change
        setSelectedIndex(-1);
      } catch (error) {
        console.error("Error getting suggestions:", error);
        setSuggestions([]);
        setShowSuggestions(false);
        setSelectedIndex(-1);
      }
    }, 200);
    return () => clearTimeout(timer);
  }, [input, getSuggestions]); // FIX: Added getSuggestions to dependencies

  // Scroll selected item into view
  useEffect(() => {
    if (selectedItemRef.current && suggestionsRef.current) {
      selectedItemRef.current.scrollIntoView({
        behavior: 'smooth',
        block: 'nearest'
      });
    }
  }, [selectedIndex]);

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      if (showSuggestions && selectedIndex >= 0 && selectedIndex < suggestions.length) {
        // If suggestions are shown and an item is selected, just fill the input
        const selected = suggestions[selectedIndex];
        const displayText = typeof selected === 'string' ? selected : selected.name;
        setInput(displayText);
        setShowSuggestions(false);
        setSelectedIndex(-1);
      } else {
        // Otherwise, try to submit the current input
        handleSubmit();
      }
    }
    else if (e.key === 'ArrowDown') {
      if (suggestions.length > 0) {
        // Update the selected index
        setSelectedIndex(prev => Math.min(prev + 1, suggestions.length - 1));
        // Ensure suggestions are shown
        setShowSuggestions(true);
      }
    }
    else if (e.key === 'ArrowUp') {
      if (selectedIndex > -1) {
        // Update the selected index
        setSelectedIndex(prev => Math.max(prev - 1, 0));
      }
    }
    else if (e.key === 'Escape') {
      // Close suggestions on Escape
      setShowSuggestions(false);
      setSelectedIndex(-1);
    }
  };

  // Handle clicks outside the component to close suggestions
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (
        inputRef.current && 
        !inputRef.current.contains(event.target) && 
        suggestionsRef.current && 
        !suggestionsRef.current.contains(event.target)
      ) {
        setShowSuggestions(false);
        setSelectedIndex(-1); // Reset selected index when closing suggestions
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  return (
    <div className="guess-input-container">
      <input
        ref={inputRef}
        type="text"
        value={input}
        onChange={(e) => setInput(e.target.value)}
        onKeyDown={handleKeyDown}
        onFocus={() => {
          // Show all territories when focused with empty input
          const newSuggestions = getSuggestions(input);
          setSuggestions(newSuggestions);
          setShowSuggestions(true); // Always show on focus
        }}
        placeholder="Start typing country name..."
        aria-label="Country guess input"
        disabled={disabled}
        autoFocus
      />
      {showSuggestions && suggestions.length > 0 && (
        <ul 
          ref={suggestionsRef} 
          className="suggestions-list" 
          style={{ maxHeight: '200px', overflowY: 'auto' }}
        >
          {suggestions.map((territory, index) => {
            // Handle both string and object territories
            const displayText = typeof territory === 'string' ? territory : territory.name;
            const isSelected = index === selectedIndex;
            
            return (
              <li
                key={index}
                ref={isSelected ? selectedItemRef : null}
                className={`suggestion-item ${isSelected ? 'selected' : ''}`}
                onClick={() => handleSelectItem(territory)}
              >
                {displayText}
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
};

export default GuessInput;