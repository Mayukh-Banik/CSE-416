import React from "react";
import { useParams } from "react-router-dom";
import { Box, Typography } from "@mui/material";
import Sidebar from "./Sidebar";
import { useTheme } from '@mui/material/styles';

interface File {
  id: string;
  name: string;
  size: number; // size in bytes
  uploadedAt: string;
  rating: number;
}

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const files: File[] = [
  { id: "tx001", name: "file001.pdf", size: 102400, uploadedAt: "2024-10-01", rating: 4 },
  { id: "tx002", name: "file009.docx", size: 204800, uploadedAt: "2024-10-05", rating: 5 }
  // Add more files as needed...
];

const FileViewPage: React.FC = () => {
  const theme = useTheme();
  const { fileId } = useParams<{ fileId: string }>();
  const file = files.find((file) => file.id === fileId);

  if (!fileId || !file) {
    return (
      <Box sx={{ padding: 3 }}>
        <Typography variant="h5">File not found</Typography>
      </Box>
    );
  }

  return (
    <Box
      sx={{
        padding: 2,
        marginTop: '70px',
        marginLeft: `${drawerWidth}px`, 
        transition: 'margin-left 0.3s ease',
        [theme.breakpoints.down('sm')]: {
          marginLeft: `${collapsedDrawerWidth}px`,
        },
      }}
    >
      <Sidebar />
      <Box sx={{ flexGrow: 1, padding: 3 }}>
        <Typography variant="h4" gutterBottom>
          File Details
        </Typography>
        <Typography variant="body1">File Name: {file.name}</Typography>
        <Typography variant="body1">
          Size: {(file.size / 1024).toFixed(2)} KB {/* Optionally, handle larger files */}
        </Typography>
        <Typography variant="body1">Uploaded At: {file.uploadedAt}</Typography>
        <Typography variant="body1">Rating: {file.rating}/5</Typography>

        {/* Future section: Account's Hosting */}
        <Box sx={{ mt: 4 }}>
          <Typography variant="h5">Account's Hosting</Typography>
          {/* Additional hosting-related details go here */}
        </Box>
      </Box>
    </Box>
  );
};

export default FileViewPage;
