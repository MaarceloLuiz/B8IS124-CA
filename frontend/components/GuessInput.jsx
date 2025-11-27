import { useState, useEffect, useRef } from 'react';

const GuessInput = ({ territories, onSubmit, disabled }) => {
  const [input, setInput] = useState('');
  const [suggestions, setSuggestions] = useState([]);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const inputRef = useRef(null);

  const getSuggestions = (value) => {
    if (!value) return [];
    const query = value.toLowerCase();

    // Filter out territories without a name property and then filter by name
    return territories
      .filter(t => t && t.name) // Make sure territory and territory.name exist
      .filter(t =>
        t.name.toLowerCase().includes(query) ||
        (t.code && t.code.toLowerCase() === query)
      )
      .slice(0, 8);
  };

  const handleSubmit = (country) => {
    if (!country || disabled) return;
    onSubmit(country);
    setInput('');
    setSuggestions([]);
    setShowSuggestions(false);
  };

  useEffect(() => {
    const timer = setTimeout(() => {
      try {
        const newSuggestions = getSuggestions(input);
        setSuggestions(newSuggestions);
        setShowSuggestions(newSuggestions.length > 0);
      } catch (error) {
        console.error("Error getting suggestions:", error);
        setSuggestions([]);
        setShowSuggestions(false);
      }
    }, 200);
    return () => clearTimeout(timer);
  }, [input, territories]);

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      const selected = suggestions[selectedIndex];
      handleSubmit(selected?.name || input);
    }
    if (e.key === 'ArrowDown') {
      setSelectedIndex(prev => Math.min(prev + 1, suggestions.length - 1));
    }
    if (e.key === 'ArrowUp') {
      setSelectedIndex(prev => Math.max(prev - 1, -1));
    }
  };

  return (
    <div className="guess-input-container">
      <input
        ref={inputRef}
        type="text"
        value={input}
        onChange={(e) => setInput(e.target.value)}
        onKeyDown={handleKeyDown}
        onFocus={() => setShowSuggestions(suggestions.length > 0)}
        placeholder="Start typing country name..."
        aria-label="Country guess input"
        disabled={disabled}
        autoFocus
      />
      {showSuggestions && (
        <ul className="suggestions-list">
          {suggestions.map((country, index) => (
            <li
              key={country.code || index}
              className={`suggestion-item ${index === selectedIndex ? 'selected' : ''}`}
              onClick={() => handleSubmit(country.name)}
            >
              {country.name}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default GuessInput;