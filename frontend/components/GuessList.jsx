const GuessList = ({ guesses, isGameOver }) => {
  const rows = [...guesses, ...Array(6 - guesses.length).fill(null)];

  return (
    <div className="guess-list">
    {rows.map((guess, index) => (
    <GuessRow key={index} guess={guess} isGameOver={isGameOver} />
    ))}
    </div>
  );
};

const GuessRow = ({ guess, isGameOver }) => {
  if (!guess) return (
    <div className="guess-row empty">
    <div className="guess-cell"></div>
    <div className="guess-cell"></div>
    <div className="guess-cell"></div>
    <div className="guess-cell"></div>
    </div>
  );

  // Direction arrow mapping based on the API response format
  const directionArrows = {
    'N': '⬆',
    'NNE': '↗',
    'NE': '↗',
    'ENE': '↗',
    'E': '➡',
    'ESE': '↘',
    'SE': '↘',
    'SSE': '↘',
    'S': '⬇',
    'SSW': '↙',
    'SW': '↙',
    'WSW': '↙',
    'W': '⬅',
    'WNW': '↖',
    'NW': '↖',
    'NNW': '↖'
  };

  // Get the direction arrow based on the direction from the API
  const directionArrow = guess.direction ? directionArrows[guess.direction] || guess.direction : '';

  return (
    <div className="guess-row">
      <div className="guess-cell country-name">{guess.country}</div>
      <div className="guess-cell distance">{guess.distance} km</div>
      <div className="guess-cell direction">
        {directionArrow}
      </div>
      <div className="guess-cell map-link">
        {isGameOver && guess.mapsUrl && (
          <a href={guess.mapsUrl} target="_blank" rel="noopener noreferrer" title="View on Google Maps">
            🌎
          </a>
        )}
      </div>
    </div>
  );
};

export default GuessList;