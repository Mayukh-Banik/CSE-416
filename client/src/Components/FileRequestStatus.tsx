import React, { useEffect, useRef, useState } from 'react';

interface FileRequestProp {
  targetID: string;
  fileHash: string;
}

const FileRequestStatus: React.FC<FileRequestProp> = ({ targetID, fileHash }) => {
  const [status, setStatus] = useState('');
  const ws = useRef<WebSocket | null>(null); // Use a ref to avoid re-render on WebSocket changes

  useEffect(() => {
    // Instantiate the WebSocket connection
    ws.current = new WebSocket('ws://localhost:8081'); // Replace with your actual WebSocket server URL

    ws.current.onopen = () => console.log("WebSocket connection established");
    
    ws.current.onmessage = (event) => {
      const message = JSON.parse(event.data);

      // Check if the message corresponds to this specific fileHash
      if (message.fileHash === fileHash) {
        if (message.status === 'accepted') {
          setStatus('Download accepted!');
          // Call your download function here
          handleFileDownload(fileHash);
        } else if (message.status === 'declined') {
          setStatus('Download declined!');
        }
      }
    };

    ws.current.onerror = (error) => console.error("WebSocket error:", error);
    
    // Clean up WebSocket on component unmount
    return () => {
      if (ws.current) {
        ws.current.close();
      }
    };
  }, [fileHash]);

  // Placeholder for download handling logic
  const handleFileDownload = (fileHash: string) => {
    // Implement your download logic here, e.g., make an API request to download the file
    console.log(`Initiating download for file with hash: ${fileHash}`);
  };

  return (
    <div>
      <h4>Download Status for {fileHash}:</h4>
      <p>{status}</p>
    </div>
  );
};

export default FileRequestStatus;
