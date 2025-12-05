import { useState, useEffect, useRef } from 'react';
import Silhouette from './components/Silhouette';
import GuessInput from './components/GuessInput';
import GuessList from './components/GuessList';
import EndGame from './components/EndGame';
import config from './config';
import './App.css';

function App() {
  const [guesses, setGuesses] = useState([]);
  const [sessionId, setSessionId] = useState('');
  const [silhouette, setSilhouette] = useState('');
  const [territories, setTerritories] = useState([]);
  const [answer, setAnswer] = useState(null);
  const [answerMapUrl, setAnswerMapUrl] = useState(null);
  const [isGameOver, setIsGameOver] = useState(false);
  
  const initialized = useRef(false);

  useEffect(() => {
    // Prevent double initialization in React StrictMode
    if (initialized.current) return;
    initialized.current = true;

    const initializeGame = async () => {
      try {
        // Get existing session from localStorage (if any)
        let session = localStorage.getItem('worldleSession');

        // ALWAYS call /api/newgame (pass session if exists)
        const url = session 
          ? `${config.API_URL}/api/newgame?sessionId=${session}`
          : `${config.API_URL}/api/newgame`;
        
        const sessionRes = await fetch(url);
        const sessionData = await sessionRes.json();
        
        // Store the session returned by backend
        session = sessionData.sessionId;
        localStorage.setItem('worldleSession', session);
        setSessionId(session);

        // Now fetch silhouette (game is initialized)
        const silhouetteRes = await fetch(`${config.API_URL}/api/silhouette`, {
          headers: { 'Accept': 'image/png,image/*' }
        });

        const blob = await silhouetteRes.blob();
        const imageUrl = URL.createObjectURL(blob);
        setSilhouette(imageUrl);

        // Get territories
        const territoriesRes = await fetch(`${config.API_URL}/api/territories`);
        setTerritories(await territoriesRes.json());
      } catch (error) {
        console.error('Initialization error:', error);
      }
    };

    initializeGame();
  }, []);

  const handleGuess = async (country) => {
    if (isGameOver) return;

    try {
      const res = await fetch(`${config.API_URL}/api/guess`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ sessionId, guess: country })
      });

      // Check if response is OK
      if (!res.ok) {
        const errorText = await res.text();
        console.error(`Server error (${res.status}):`, errorText);
        alert(`Failed to process guess for "${country}". The server couldn't find coordinates for this location.`);
        return;
      }

      // Try to parse JSON
      const result = await res.json();
      console.log("API response:", result);

      const newGuess = {
        country,
        distance: result.distance,
        direction: result.direction,
        mapsUrl: result.url
      };

      setGuesses(prev => [...prev, newGuess]);

      if (result.isCorrect || guesses.length >= 5) {
        // Game is over, fetch the answer
        const answerRes = await fetch(`${config.API_URL}/api/answer`, {
          method: 'GET'
        });

        const answerData = await answerRes.json();
        setAnswer(answerData.answer);
        setAnswerMapUrl(answerData.url);
        setIsGameOver(true);

        if (result.isCorrect) {
          localStorage.removeItem('worldleSession');
        }
      }
    } catch (error) {
      console.error('Guess error:', error);
      alert(`Failed to submit guess. Error: ${error.message}`);
    }
  };

  return (
    <div className="container">
      <h1 className="title">Worldle</h1>
      <Silhouette imageUrl={silhouette} />
      <GuessInput territories={territories} onSubmit={handleGuess} disabled={isGameOver} />
      <GuessList guesses={guesses} isGameOver={isGameOver} />
      {isGameOver && <EndGame answer={answer} answerMapUrl={answerMapUrl} />}
    </div>
  );
}

export default App;