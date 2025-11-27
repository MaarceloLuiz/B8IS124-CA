const EndGame = ({ answer, answerMapUrl }) => {
  console.log("EndGame props:", { answer, answerMapUrl }); // Debug log

  return (
    <div className="end-game">
      <div className="correct-answer">
        Correct answer: <span className="answer-text">{answer || 'Unknown'}</span>
      </div>
      {answerMapUrl && (
        <a href={answerMapUrl} target="_blank" rel="noopener noreferrer">
          <button className="map-button large">View on Google Maps</button>
        </a>
      )}
    </div>
  );
};

export default EndGame;