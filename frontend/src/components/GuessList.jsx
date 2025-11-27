import NIcon from '../assets/icons/N.png';
import NEIcon from '../assets/icons/NE.png';
import EIcon from '../assets/icons/E.png';
import SEIcon from '../assets/icons/SE.png';
import SIcon from '../assets/icons/S.png';
import SWIcon from '../assets/icons/SW.png';
import WIcon from '../assets/icons/W.png';
import NWIcon from '../assets/icons/NW.png';
import WorldIcon from '../assets/icons/world.png';

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

 // Direction mapping to simplify the number of icons needed
 const directionIcons = {
    'N': NIcon,
    'NNE': NEIcon,
    'NE': NEIcon,
    'ENE': NEIcon,
    'E': EIcon,
    'ESE': SEIcon,
    'SE': SEIcon,
    'SSE': SEIcon,
    'S': SIcon,
    'SSW': SWIcon,
    'SW': SWIcon,
    'WSW': SWIcon,
    'W': WIcon,
    'WNW': NWIcon,
    'NW': NWIcon,
    'NNW': NWIcon
  };

  // Get the icon name based on the direction from the API
const iconSrc = guess.direction ? directionIcons[guess.direction] : '';

  return (
    <div className="guess-row">
      <div className="guess-cell country-name">{guess.country}</div>
      <div className="guess-cell distance">{guess.distance} km</div>
      <div className="guess-cell direction">
        {iconSrc && (
          <img 
            src={iconSrc} 
            alt={guess.direction} 
            style={{ width: '24px', height: '24px' }} 
          />
        )}
      </div>
      <div className="guess-cell map-link">
        {isGameOver && guess.mapsUrl && (
          <a href={guess.mapsUrl} target="_blank" rel="noopener noreferrer" title="View on Google Maps">
            <img 
              src={WorldIcon} 
              alt="View on Google Maps" 
              style={{ width: '24px', height: '24px' }} 
            />
          </a>
        )}
      </div>
    </div>
  );
};

export default GuessList;