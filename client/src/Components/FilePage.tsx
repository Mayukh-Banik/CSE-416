import React, { useState } from 'react';
import Sidebar from './Sidebar';
import FilePageStyles from '../Stylesheets/FilePageStyle';

const FilePage = () => {
    const classes = FilePageStyles(); // Custom styles
    const [files, setFiles] = useState<File[]>([]); 

    const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const newFiles = event.target.files;

        if (newFiles) {
            // Convert the FileList to an array and merge it with the existing files
            setFiles((prevFiles) => [
                ...prevFiles,
                ...Array.from(newFiles)
            ]);
        }
  };

  return (
    <div >
        <Sidebar/>
      <h1 >Upload Files</h1>
      <input
        type="file"
        multiple
        onChange={handleFileChange}
      />

      {files.length === 0 ? (
        <p>No files uploaded yet.</p>
      ) : (
        <table>
          <thead>
            <tr>
              <th>File Name</th>
            </tr>
          </thead>
          <tbody>
            {files.map((file) => (
              <tr>
                <td>{file.name}</td>
                <td>
                  <a href={URL.createObjectURL(file)} target="_blank">
                    <button>Click here to view</button>
                  </a>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
};

export default FilePage;
