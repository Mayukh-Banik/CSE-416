import { createTheme } from "@mui/material";

const GeneralTheme = createTheme({
    palette: {
      primary: {
        main: "#A7C7E7", // Pastel blue
      },
      secondary: {
        main: "#f50057", // Pink (can be used later for accents if needed)
      },
      background: {
        default: "#ffffff", // White
      },
    },
    typography: {
      fontFamily: "'Arial', sans-serif",
      h2: {
        fontSize: "3rem",  // Increased font size for h2 (Welcome message)
        fontWeight: 700,   // Bolder weight for prominence
      },
      body1: {
        fontSize: "1.2rem", // Slightly larger body text
        color: "#666",      // Subtle grey for body text
      },
    },
  });
  
  export default GeneralTheme;