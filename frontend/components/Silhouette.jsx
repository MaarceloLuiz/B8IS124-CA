const Silhouette = ({ imageUrl }) => {
  console.log("Silhouette image URL:", imageUrl);

  return (
    <div className="silhouette-container">
      {imageUrl ? (
        <img
          src={imageUrl}
          alt="Country silhouette"
          className="silhouette-image"
          style={{ maxWidth: '300px', maxHeight: '300px', width: 'auto', height: 'auto', border: 'none' }}
          onError={(e) => {
            console.error("Failed to load image:", imageUrl);
            e.target.onerror = null; // Prevent infinite loop
            e.target.alt = "Country silhouette (failed to load)";
          }}
        />
      ) : (
        <div className="loading-silhouette">Loading...</div>
      )}
    </div>
  );
};

export default Silhouette;