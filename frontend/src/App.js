import { useState, useEffect } from 'react';
import Silhouette from './components/Silhouette';
import GuessInput from './components/GuessInput';
import GuessList from './components/GuessList';
import EndGame from './components/EndGame';
import './App.css';

function App() {
  const [guesses, setGuesses] = useState([]);
  const [sessionId, setSessionId] = useState('');
  const [silhouette, setSilhouette] = useState('');
  const [territories, setTerritories] = useState([]);
  const [answer, setAnswer] = useState(null);
  const [answerMapUrl, setAnswerMapUrl] = useState(null);
  const [isGameOver, setIsGameOver] = useState(false);

  useEffect(() => {
    const initializeGame = async () => {
      try {
        let session = localStorage.getItem('worldleSession');

        if (!session) {
          // Create a new game session
          const sessionRes = await fetch('http://localhost:8080/api/newgame');
          const sessionText = await sessionRes.text(); // Get response as text first

          try {
            // Try to parse as JSON
            const sessionData = JSON.parse(sessionText);
            session = sessionData.sessionId;
            localStorage.setItem('worldleSession', session);
            setSessionId(session);
          } catch (e) {
            // If not valid JSON, extract session ID from text
            console.log("Response is not JSON:", sessionText);
            // Assuming the response contains the session ID somewhere
            const match = sessionText.match(/session[:\s]+([a-zA-Z0-9-]+)/i);
            if (match && match[1]) {
              session = match[1];
              localStorage.setItem('worldleSession', session);
              setSessionId(session);
            } else {
              console.error("Could not extract session ID from response");
            }
          }
        } else {
          setSessionId(session);
        }

        // Get the silhouette image directly as a blob
        const silhouetteRes = await fetch(`http://localhost:8080/api/silhouette`, {
          headers: {
            'Accept': 'image/png,image/*'
          }
        });

        // Check if the response is an image
        const contentType = silhouetteRes.headers.get('content-type');
        if (contentType && contentType.includes('image')) {
          // Create a blob URL from the image data
          const blob = await silhouetteRes.blob();
          const imageUrl = URL.createObjectURL(blob);
          setSilhouette(imageUrl);
        } else {
          console.error("Silhouette response is not an image:", contentType);
          try {
            const data = await silhouetteRes.json();
            if (data.imageUrl) {
              setSilhouette(`http://localhost:8080${data.imageUrl}`);
            }
          } catch (e) {
            console.error("Could not parse silhouette response:", e);
          }
        }

        // Get territories list
        const territoriesRes = await fetch('http://localhost:8080/api/territories');
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
      const res = await fetch('http://localhost:8080/api/guess', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ sessionId, guess: country })
      });

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
        // Game is over, fetch the answer from the new endpoint
        const answerRes = await fetch('http://localhost:8080/api/answer', {
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
